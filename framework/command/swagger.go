package command

import (
	"fmt"
	"github.com/gohade/hade/framework/cobra"
	"github.com/gohade/hade/framework/contract"
	"github.com/swaggo/swag/gen"
	"path/filepath"
)

func initSwaggerCommand() *cobra.Command {
	swaggerCommand.AddCommand(swaggerGenCommand)
	return swaggerCommand
}

var swaggerCommand = &cobra.Command{
	Use:   "swagger",
	Short: "swagger对应命令",
	RunE: func(c *cobra.Command, args []string) error {
		if len(args) == 0 {
			c.Help()
		}
		return nil
	},
}

// envCommand show current envionment
var swaggerGenCommand = &cobra.Command{
	Use:   "gen",
	Short: "生成对应的swagger文件, contain swagger.yaml, doc.go",
	Run: func(c *cobra.Command, args []string) {
		container := c.GetContainer()
		appService := container.MustMake(contract.AppKey).(contract.App)

		outputDir := filepath.Join(appService.AppFolder(), "http", "swagger")

		conf := &gen.Config{
			SearchDir:          "./app/http/",
			Excludes:           "",
			OutputDir:          outputDir,
			MainAPIFile:        "swagger.go",
			PropNamingStrategy: "",
			ParseVendor:        false,
			ParseDependency:    false,
			ParseInternal:      false,
			MarkdownFilesDir:   "",
			GeneratedTime:      false,
		}
		err := gen.New().Build(conf)
		if err != nil {
			fmt.Println(err)
		}
	},
}
