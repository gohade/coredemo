package command

import (
	"fmt"
	"github.com/gohade/hade/framework"
	"github.com/pkg/errors"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/gohade/hade/framework/cobra"
	"github.com/gohade/hade/framework/contract"
	"github.com/gohade/hade/framework/util"

	"github.com/fsnotify/fsnotify"
)

const (
	frontendPrefix = "/dist"
)

var (
	refreshTime int32
	proxyURL    string
)

// devConfig 代表调试模式的配置信息
type devConfig struct {
	Port    string   // 调试模式最终监听的端口，默认为8070
	Backend struct { // 后端调试模式配置
		RefreshTime   int    // 调试模式后端更新时间，如果文件变更，等待3s才进行一次更新，能让频繁保存变更更为顺畅, 默认1s
		Port          string // 后端监听端口， 默认 8072
		MonitorFolder string // 监听文件夹，默认为AppFolder
	}
	Frontend struct { // 前端调试模式配置
		Port string // 前端启动端口, 默认8071
	}
}

// 初始化配置文件
func initDevConfig(c framework.Container) *devConfig {
	devConfig := &devConfig{
		Port: "8087",
		Backend: struct {
			RefreshTime   int
			Port          string
			MonitorFolder string
		}{
			1,
			"8072",
			"",
		},
		Frontend: struct {
			Port string
		}{
			"8071",
		},
	}
	configer := c.MustMake(contract.ConfigKey).(contract.Config)
	if configer.IsExist("app.dev.port") {
		devConfig.Port = configer.GetString("app.dev.port")
	}
	if configer.IsExist("app.dev.backend.refresh_time") {
		devConfig.Backend.RefreshTime = configer.GetInt("app.dev.backend.refresh_time")
	}
	if configer.IsExist("app.dev.backend.port") {
		devConfig.Backend.Port = configer.GetString("app.dev.backend.port")
	}
	monitorFolder := configer.GetString("app.dev.backend.monitor_folder")
	if monitorFolder == "" {
		appService := c.MustMake(contract.AppKey).(contract.App)
		devConfig.Backend.MonitorFolder = appService.AppFolder()
	}

	if configer.IsExist("app.dev.frontend.port") {
		devConfig.Frontend.Port = configer.GetString("app.dev.frontend.port")
	}
	return devConfig
}

// Proxy 代表serve启动的backend的服务器代理
type Proxy struct {
	devConfig   *devConfig    // 配置文件
	proxyServer *http.Server  // proxy的服务
	backendPid  int           // 当前的backend服务的pid
	frontendPid int           // 当前的frontend服务的pid
	lock        *sync.RWMutex // 锁机制
}

// NewProxy 初始化一个Proxy
func NewProxy(c framework.Container) *Proxy {
	devConfig := initDevConfig(c)
	return &Proxy{
		devConfig: devConfig,
		lock:      &sync.RWMutex{},
	}
}

// 重新启动一个proxy网关
func (p *Proxy) newProxyReverseProxy(frontend, backend *url.URL) *httputil.ReverseProxy {
	if p.frontendPid == 0 {
		if p.backendPid != 0 {
			director := func(req *http.Request) {
				req.URL.Scheme = backend.Scheme
				req.URL.Host = backend.Host
			}
			return &httputil.ReverseProxy{Director: director}
		}
	}

	if p.backendPid == 0 {
		if p.frontendPid != 0 {
			director := func(req *http.Request) {
				req.URL.Scheme = frontend.Scheme
				req.URL.Host = frontend.Host
			}
			return &httputil.ReverseProxy{Director: director}
		}
	}

	// 两个都有进程
	director := func(req *http.Request) {
		if req.URL.Path == "/" || req.URL.Path == "/app.js" {
			req.URL.Scheme = frontend.Scheme
			req.URL.Host = frontend.Host
		} else {
			req.URL.Scheme = backend.Scheme
			req.URL.Host = backend.Host
		}
	}

	NotFoundErr := errors.New("response is 404, need to redirect")
	return &httputil.ReverseProxy{Director: director,
		ModifyResponse: func(response *http.Response) error {
			if response.StatusCode == 404 {
				return NotFoundErr
			}
			return nil
		},
		ErrorHandler: func(writer http.ResponseWriter, request *http.Request, err error) {
			if errors.Is(err, NotFoundErr) {
				httputil.NewSingleHostReverseProxy(backend).ServeHTTP(writer, request)
			}
		}}
}

// rebuildBackend 重新编译后端
func (p *Proxy) rebuildBackend() error {
	// 重新编译hade
	cmdBuild := exec.Command("./hade", "build", "backend")
	cmdBuild.Stdout = os.Stdout
	cmdBuild.Stderr = os.Stderr
	if err := cmdBuild.Start(); err == nil {
		err = cmdBuild.Wait()
		if err != nil {
			return err
		}
	}
	return nil
}

// restartBackend 启动后端服务
func (p *Proxy) restartBackend() error {

	// 杀死之前的进程
	if p.backendPid != 0 {
		syscall.Kill(p.backendPid, syscall.SIGKILL)
		p.backendPid = 0
	}

	// 设置随机端口，真实后端的端口
	port := p.devConfig.Backend.Port
	hadeAddress := fmt.Sprintf(":" + port)
	// 使用命令行启动后端进程
	cmd := exec.Command("./hade", "app", "start", "--address="+hadeAddress)
	cmd.Stdout = os.NewFile(0, os.DevNull)
	cmd.Stderr = os.Stderr
	fmt.Println("启动后端服务: ", "http://127.0.0.1:"+port)
	err := cmd.Start()
	if err != nil {
		fmt.Println(err)
	}
	p.backendPid = cmd.Process.Pid
	fmt.Println("后端服务pid:", p.backendPid)
	return nil
}

// 启动前端服务
func (p *Proxy) restartFrontend() error {
	// 启动前端调试模式
	// 如果已经开启了npm run serve， 什么都不做
	if p.frontendPid != 0 {
		return nil
	}

	// 否则开启npm run serve
	port := p.devConfig.Frontend.Port
	path, err := exec.LookPath("npm")
	if err != nil {
		return err
	}
	cmd := exec.Command(path, "run", "dev")
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("%s%s", "PORT=", port))
	cmd.Stdout = os.NewFile(0, os.DevNull)
	cmd.Stderr = os.Stderr

	// 因为npm run serve 是控制台挂起模式，所以这里使用go routine启动
	err = cmd.Start()
	fmt.Println("启动前端服务: ", "http://127.0.0.1:"+port)
	if err != nil {
		fmt.Println(err)
	}
	p.frontendPid = cmd.Process.Pid
	fmt.Println("前端服务pid:", p.frontendPid)

	return nil
}

// 重启后端服务, 如果frontend为nil，则没有包含后端
func (p *Proxy) restartProxy(startFrontend, startBackend bool) error {
	var backendURL, frontendURL *url.URL
	var err error

	// 启动后端
	if startBackend {
		if err := p.restartBackend(); err != nil {
			return err
		}
	}
	// 启动前端
	if startFrontend {
		if err := p.restartFrontend(); err != nil {
			return err
		}
	}

	// 如果已经启动过proxy了，就不要进行设置了
	if p.proxyServer != nil {
		return nil
	}

	if frontendURL, err = url.Parse(fmt.Sprintf("%s%s", "http://127.0.0.1:", p.devConfig.Frontend.Port)); err != nil {
		return err
	}

	if backendURL, err = url.Parse(fmt.Sprintf("%s%s", "http://127.0.0.1:", p.devConfig.Backend.Port)); err != nil {
		return err
	}

	// 设置反向代理
	proxyReverse := p.newProxyReverseProxy(frontendURL, backendURL)
	p.proxyServer = &http.Server{
		Addr:    "127.0.0.1:" + p.devConfig.Port,
		Handler: proxyReverse,
	}

	fmt.Println("代理服务启动:", "http://"+p.proxyServer.Addr)
	// 启动proxy服务
	err = p.proxyServer.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

// monitorBackend 监听应用文件
func (p *Proxy) monitorBackend() error {
	// 监听
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	appFolder := p.devConfig.Backend.MonitorFolder
	fmt.Println("监控文件夹：", appFolder)
	filepath.Walk(appFolder, func(path string, info os.FileInfo, err error) error {
		if info != nil && !info.IsDir() {
			return nil
		}
		if util.IsHiddenDirectory(path) {
			return nil
		}
		return watcher.Add(path)
	})

	refreshTime := p.devConfig.Backend.RefreshTime
	t := time.NewTimer(time.Duration(refreshTime) * time.Second)
	t.Stop()
	for {
		select {
		case <-t.C:
			fmt.Println("...检测到文件更新，重启服务开始...")
			if err := p.rebuildBackend(); err != nil {
				fmt.Println("重新编译失败：", err.Error())
			} else {
				if err := p.restartBackend(); err != nil {
					fmt.Println("重新启动失败：", err.Error())
				}
			}
			fmt.Println("...检测到文件更新，重启服务结束...")
			t.Stop()
		case _, ok := <-watcher.Events:
			if !ok {
				continue
			}
			t.Reset(time.Duration(refreshTime) * time.Second)
		case err, ok := <-watcher.Errors:
			if !ok {
				continue
			}
			fmt.Println("监听文件夹错误：", err.Error())
			t.Reset(time.Duration(refreshTime) * time.Second)
		}
	}
	return nil
}

// 初始化Dev命令
func initDevCommand() *cobra.Command {
	devCommand.AddCommand(devBackendCommand)
	devCommand.AddCommand(devFrontendCommand)
	devCommand.AddCommand(devAllCommand)
	return devCommand
}

// devCommand 为调试模式的一级命令
var devCommand = &cobra.Command{
	Use:   "dev",
	Short: "调试模式",
	RunE: func(c *cobra.Command, args []string) error {
		c.Help()
		return nil
	},
}

// devBackendCommand 启动后端调试模式
var devBackendCommand = &cobra.Command{
	Use:   "backend",
	Short: "启动后端调试模式",
	RunE: func(c *cobra.Command, args []string) error {
		proxy := NewProxy(c.GetContainer())
		go proxy.monitorBackend()
		if err := proxy.restartProxy(false, true); err != nil {
			return err
		}
		return nil
	},
}

// devFrontendCommand 启动前端调试模式
var devFrontendCommand = &cobra.Command{
	Use:   "frontend",
	Short: "前端调试模式",
	RunE: func(c *cobra.Command, args []string) error {

		// 启动前端服务
		proxy := NewProxy(c.GetContainer())
		return proxy.restartProxy(true, false)

	},
}

var devAllCommand = &cobra.Command{
	Use:   "all",
	Short: "同时启动前端和后端调试",
	RunE: func(c *cobra.Command, args []string) error {
		proxy := NewProxy(c.GetContainer())
		go proxy.monitorBackend()
		if err := proxy.restartProxy(true, true); err != nil {
			return err
		}
		return nil
	},
}
