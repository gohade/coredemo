package cobra

import (
	"github.com/gohade/hade/framework"
	"github.com/robfig/cron/v3"
	"log"
)

// SetContainer 设置服务容器
func (c *Command) SetContainer(container framework.Container) {
	c.container = container
}

// GetContainer 获取容器
func (c *Command) GetContainer() framework.Container {
	return c.Root().container
}

// CommandSpec add by yejianfeng
type CommandSpec struct {
	Cmd  *Command
	Args []string
	Spec string
}

// AddCronCommand add by yejianfeng
// add Cron command for cron execute
func (c *Command) AddCronCommand(spec string, cmd *Command, args ...string) {
	root := c.Root()
	if root.Cron == nil {
		root.Cron = cron.New()
		root.CronSepcs = []CommandSpec{}
	}
	root.CronSepcs = append(root.CronSepcs, CommandSpec{
		Cmd:  cmd,
		Args: args,
		Spec: spec,
	})
	root.Cron.AddFunc(spec, func() {
		var cronCmd Command
		ctx := root.Context()
		cronCmd = *cmd
		cronCmd.SetParentNull()
		cronCmd.args = []string{}
		err := cronCmd.ExecuteContext(ctx)
		if err != nil {
			log.Println(err)
		}
	})
}

// by yejianfeng, set parent
func (c *Command) SetParentNull() {
	c.parent = nil
}
