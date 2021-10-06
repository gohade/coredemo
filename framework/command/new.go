package command

import (
	"bytes"
	"context"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/google/go-github/v39/github"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/gohade/hade/framework/cobra"
	"github.com/gohade/hade/framework/util"
)

// new相关的名称
func initNewCommand() *cobra.Command {
	return newCommand
}

// 创建一个新应用
var newCommand = &cobra.Command{
	Use:     "new",
	Aliases: []string{"create", "init"},
	Short:   "创建一个新的应用",
	RunE: func(c *cobra.Command, args []string) error {
		currentPath := util.GetExecDirectory()

		var name string
		var folder string
		var mod string
		var version string
		var release *github.RepositoryRelease
		{
			prompt := &survey.Input{
				Message: "请输入目录名称：",
			}
			err := survey.AskOne(prompt, &name)
			if err != nil {
				return err
			}

			folder = filepath.Join(currentPath, name)
			if util.Exists(folder) {
				isForce := false
				prompt2 := &survey.Confirm{
					Message: "目录" + name + "已经存在,是否删除重新创建？",
					Default: false,
				}
				err := survey.AskOne(prompt2, &isForce)
				if err != nil {
					return err
				}

				if isForce {
					os.RemoveAll(folder)
				} else {
					fmt.Println("目录已存在，创建应用失败")
					return nil
				}
			}
		}
		{
			prompt := &survey.Input{
				Message: "请输入模块名称(go.mod中的module, 默认为文件夹名称)：",
			}
			err := survey.AskOne(prompt, &mod)
			if err != nil {
				return err
			}
			if mod == "" {
				mod = name
			}
		}
		{
			// 获取hade的版本
			client := github.NewClient(nil)
			opt := &github.ListOptions{Page: 0, PerPage: 10}
			releases, _, err := client.Repositories.ListReleases(context.Background(), "gohade", "hade", opt)
			if err != nil {
				return err
			}

			prompt := &survey.Input{
				Message: "请输入版本名称(参考 https://github.com/gohade/hade/releases，默认为最新版本)：",
			}
			err = survey.AskOne(prompt, &version)
			if err != nil {
				return err
			}
			if version != "" {
				// 确认版本是否正确
				release, _, err = client.Repositories.GetReleaseByTag(context.Background(), "gohade", "hade", version)
				if err != nil {
					return err
				}
				if release == nil {
					fmt.Println("版本不存在，创建应用失败，请参考 https://github.com/gohade/hade/releases")
					return nil
				}
			}
			if version == "" {
				version = releases[0].GetTagName()
				release = releases[0]
			}
		}

		// 拷贝template项目
		url := release.GetZipballURL()
		err := util.DownloadFile("template-main.zip", url)
		if err != nil {
			return err
		}

		_, err = util.Unzip("template-main.zip", currentPath)
		if err != nil {
			return err
		}

		if err := os.Rename(filepath.Join(currentPath, "/template-main"), folder); err != nil {
			return err
		}

		if err := os.Remove("template-main.zip"); err != nil {
			return err
		}
		fmt.Println("remove " + path.Join(folder, ".git"))
		os.RemoveAll(path.Join(folder, ".git"))

		// 删除framework 目录
		os.RemoveAll(path.Join(folder, "framework"))

		filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
			fmt.Println("read file:" + path)
			if info.IsDir() {
				return nil
			}

			c, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			if path == filepath.Join(folder, "go.mod") {
				fmt.Println("更新文件:" + path)
				c = bytes.ReplaceAll(c, []byte("module github.com/gohade/hade"), []byte("module "+mod))
			}

			isContain := bytes.Contains(c, []byte("github.com/gohade/hade/app"))
			if isContain {
				fmt.Println("update file:" + path)
				c = bytes.ReplaceAll(c, []byte("github.com/gohade/hade/app"), []byte(mod+"/app"))
				err = ioutil.WriteFile(path, c, 0644)
				if err != nil {
					return err
				}
			}

			return nil
		})
		return nil
	},
}
