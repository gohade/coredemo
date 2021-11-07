package command

import (
	"context"
	"fmt"
	"github.com/gohade/hade/framework"
	"github.com/gohade/hade/framework/cobra"
	"github.com/gohade/hade/framework/contract"
	"github.com/gohade/hade/framework/provider/ssh"
	"github.com/gohade/hade/framework/util"
	"github.com/pkg/errors"
	"github.com/pkg/sftp"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func initDeployCommand() *cobra.Command {
	deployCommand.AddCommand(deployFrontendCommand)
	deployCommand.AddCommand(deployBackendCommand)
	deployCommand.AddCommand(deployAllCommand)
	return deployCommand
}

var deployCommand = &cobra.Command{
	Use:   "deploy",
	Short: "部署相关命令",
	RunE: func(c *cobra.Command, args []string) error {
		if len(args) == 0 {
			c.Help()
		}
		return nil
	},
}

// 创建部署的folder
func createDeployFolder(c framework.Container) (string, error) {
	appService := c.MustMake(contract.AppKey).(contract.App)
	deployFolder := appService.DeployFolder()

	deployVersion := time.Now().Format("20060102150405")
	versionFolder := filepath.Join(deployFolder, deployVersion)
	if !util.Exists(versionFolder) {
		return versionFolder, os.Mkdir(versionFolder, os.ModePerm)
	}
	return versionFolder, nil
}

// deployFrontendCommand 部署前端
var deployFrontendCommand = &cobra.Command{
	Use:   "frontend",
	Short: "部署前端",
	RunE: func(c *cobra.Command, args []string) error {
		container := c.GetContainer()

		deployFolder, err := createDeployFolder(container)
		if err != nil {
			return err
		}

		if err := deployBuildFrontend(c, deployFolder); err != nil {
			return err
		}

		return deployUploadAction(deployFolder, container, "frontend")
	},
}

// 部署后端
var deployBackendCommand = &cobra.Command{
	Use:   "backend",
	Short: "部署后端",
	RunE: func(c *cobra.Command, args []string) error {
		container := c.GetContainer()

		deployFolder, err := createDeployFolder(container)
		if err != nil {
			return err
		}

		if err := deployBuildBackend(c, deployFolder); err != nil {
			return err
		}

		return deployUploadAction(deployFolder, container, "backend")
	},
}

func deployBuildBackend(c *cobra.Command, deployFolder string) error {
	container := c.GetContainer()
	configService := container.MustMake(contract.ConfigKey).(contract.Config)
	appService := container.MustMake(contract.AppKey).(contract.App)
	envService := container.MustMake(contract.EnvKey).(contract.Env)
	logger := container.MustMake(contract.LogKey).(contract.Log)

	env := envService.AppEnv()

	binFile := "hade"

	// 编译前端
	path, err := exec.LookPath("go")
	if err != nil {
		log.Fatalln("hade go: 请在Path路径中先安装go")
	}

	deployBinFile := filepath.Join(deployFolder, binFile)
	cmd := exec.Command(path, "build", "-o", deployBinFile, "./")
	cmd.Env = os.Environ()
	if configService.GetString("deploy.backend.goos") != "" {
		cmd.Env = append(cmd.Env, "GOOS="+configService.GetString("deploy.backend.goos"))
	}
	if configService.GetString("deploy.backend.goarch") != "" {
		cmd.Env = append(cmd.Env, "GOARCH="+configService.GetString("deploy.backend.goarch"))
	}

	ctx := context.Background()
	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error(ctx, "go build err", map[string]interface{}{
			"err": err,
			"out": string(out),
		})
		return err
	}
	logger.Info(ctx, "编译成功", nil)

	// 复制.env文件
	if util.Exists(filepath.Join(appService.BaseFolder(), ".env")) {
		if err := util.CopyFile(filepath.Join(appService.BaseFolder(), ".env"), filepath.Join(deployFolder, ".env")); err != nil {
			return err
		}
	}

	// 复制config文件
	deployConfigFolder := filepath.Join(deployFolder, "config", env)
	if !util.Exists(deployConfigFolder) {
		if err := os.MkdirAll(deployConfigFolder, os.ModePerm); err != nil {
			return err
		}
	}
	if err := util.CopyFolder(filepath.Join(appService.ConfigFolder(), env), deployConfigFolder); err != nil {
		return err
	}

	logger.Info(ctx, "build local ok", nil)
	return nil
}

func deployUploadAction(deployFolder string, container framework.Container, end string) error {
	configService := container.MustMake(contract.ConfigKey).(contract.Config)
	sshService := container.MustMake(contract.SSHKey).(contract.SSHService)
	logger := container.MustMake(contract.LogKey).(contract.Log)

	// 遍历所有deploy的服务器
	deployNodes := configService.GetStringSlice("deploy.connections")
	if len(deployNodes) == 0 {
		return errors.New("deploy connections len is zero")
	}
	remoteFolder := configService.GetString("deploy.remote_folder")
	if remoteFolder == "" {
		return errors.New("remote folder is empty")
	}

	preActions := []string{}
	postActions := []string{}

	if end == "frontend" || end == "both" {
		preActions = append(preActions, configService.GetStringSlice("deploy.frontend.pre_action")...)
		postActions = append(postActions, configService.GetStringSlice("deploy.frontend.post_action")...)
	}
	if end == "backend" || end == "both" {
		preActions = append(preActions, configService.GetStringSlice("deploy.backend.pre_action")...)
		postActions = append(postActions, configService.GetStringSlice("deploy.backend.post_action")...)
	}

	for _, node := range deployNodes {
		sshClient, err := sshService.GetClient(ssh.WithConfigPath(node))
		if err != nil {
			return err
		}
		client, err := sftp.NewClient(sshClient)
		if err != nil {
			return err
		}

		for _, action := range preActions {
			session, err := sshClient.NewSession()
			if err != nil {
				return err
			}
			defer session.Close()
			bts, err := session.CombinedOutput(action)
			if err != nil {
				return err
			}
			logger.Info(context.Background(), "execute pre action", map[string]interface{}{
				"cmd":        action,
				"connection": node,
				"out":        strings.ReplaceAll(string(bts), "\n", ""),
			})
		}

		if err := uploadFolderToSFTP(container, deployFolder, remoteFolder, client); err != nil {
			logger.Info(context.Background(), "upload folder failed", map[string]interface{}{
				"err": err,
			})
			return err
		}
		logger.Info(context.Background(), "upload folder success", nil)

		for _, action := range postActions {
			session, err := sshClient.NewSession()
			if err != nil {
				return err
			}
			defer session.Close()
			bts, err := session.CombinedOutput(action)
			if err != nil {
				return err
			}
			logger.Info(context.Background(), "execute post action", map[string]interface{}{
				"cmd":        action,
				"connection": node,
				"out":        strings.ReplaceAll(string(bts), "\n", ""),
			})
		}
	}
	return nil
}

var deployAllCommand = &cobra.Command{
	Use:   "all",
	Short: "全部部署",
	RunE: func(c *cobra.Command, args []string) error {
		container := c.GetContainer()

		deployFolder, err := createDeployFolder(container)
		if err != nil {
			return err
		}

		// 编译前端
		if err := deployBuildFrontend(c, deployFolder); err != nil {
			return err
		}

		// 编译后端
		if err := deployBuildBackend(c, deployFolder); err != nil {
			return err
		}

		// 上传前端+后端
		return deployUploadAction(deployFolder, container, "both")
	},
}

func write(w io.WriteCloser, command string) error {
	_, err := w.Write([]byte(command + "\n"))
	return err
}

func readUntil(r io.Reader, matchingByte []byte) (*string, error) {
	var buf [64 * 1024]byte
	var t int
	for {
		fmt.Println("start read")
		n, err := r.Read(buf[t:])
		fmt.Println("start read one line")
		if err != nil {
			return nil, err
		}
		t += n
		if isMatch(buf[:t], t, matchingByte) {
			stringResult := string(buf[:t])
			return &stringResult, nil
		}
		fmt.Println(string(buf[:t]))
	}
}

func isMatch(bytes []byte, t int, matchingBytes []byte) bool {
	if t >= len(matchingBytes) {
		for i := 0; i < len(matchingBytes); i++ {
			if bytes[t-len(matchingBytes)+i] != matchingBytes[i] {
				return false
			}
		}
		return true
	}
	return false
}

func uploadFolderToSFTP(container framework.Container, localFolder, remoteFolder string, client *sftp.Client) error {
	logger := container.MustMake(contract.LogKey).(contract.Log)
	return filepath.Walk(localFolder, func(path string, info os.FileInfo, err error) error {
		relPath := strings.Replace(path, localFolder, "", 1)
		if relPath == "" {
			return nil
		}
		if info.IsDir() {
			logger.Info(context.Background(), "mkdir: "+filepath.Join(remoteFolder, relPath), nil)
			return client.MkdirAll(filepath.Join(remoteFolder, relPath))
		}

		rf, err := os.Open(filepath.Join(localFolder, relPath))
		if err != nil {
			return errors.New("read file " + filepath.Join(localFolder, relPath) + " error:" + err.Error())
		}
		rfStat, err := rf.Stat()
		if err != nil {
			return err
		}
		f, err := client.Create(filepath.Join(remoteFolder, relPath))
		if err != nil {
			return errors.New("create file " + filepath.Join(remoteFolder, relPath) + " error:" + err.Error())
		}
		// 大于2M的文件显示进度
		if rfStat.Size() > 2*1024*1024 {
			logger.Info(context.Background(), "upload local file: "+filepath.Join(localFolder, relPath)+
				" to remote file: "+filepath.Join(remoteFolder, relPath)+" start", nil)
			go func(localFile, remoteFile string) {
				ticker := time.NewTicker(10 * time.Second)

				for range ticker.C {
					remoteFileInfo, err := client.Stat(remoteFile)
					if err != nil {
						logger.Error(context.Background(), "stat error", map[string]interface{}{
							"err":         err,
							"remote_file": remoteFile,
						})
						continue
					}
					size := remoteFileInfo.Size()
					if size >= rfStat.Size() {
						break
					}
					percent := int(size * 100 / rfStat.Size())
					logger.Info(context.Background(), "upload local file: "+filepath.Join(localFolder, relPath)+
						" to remote file: "+filepath.Join(remoteFolder, relPath)+fmt.Sprintf(" %v%% %v/%v", percent, size, rfStat.Size()), nil)
				}
			}(filepath.Join(localFolder, relPath), filepath.Join(remoteFolder, relPath))
		}
		if _, err := f.ReadFromWithConcurrency(rf, 10); err != nil {
			return errors.New("Write file " + filepath.Join(remoteFolder, relPath) + " error:" + err.Error())
		}
		logger.Info(context.Background(), "upload local file: "+filepath.Join(localFolder, relPath)+
			" to remote file: "+filepath.Join(remoteFolder, relPath)+" finish", nil)
		return nil
	})
}

func deployBuildFrontend(c *cobra.Command, deployFolder string) error {
	container := c.GetContainer()
	appService := container.MustMake(contract.AppKey).(contract.App)

	// 编译前端
	if err := buildFrontendCommand.RunE(c, []string{}); err != nil {
		return err
	}

	// 复制前端文件到deploy文件夹
	frontendFolder := filepath.Join(deployFolder, "dist")
	if err := os.Mkdir(frontendFolder, os.ModePerm); err != nil {
		return err
	}

	buildFolder := filepath.Join(appService.BaseFolder(), "dist")
	if err := util.CopyFolder(buildFolder, frontendFolder); err != nil {
		return err
	}
	return nil
}
