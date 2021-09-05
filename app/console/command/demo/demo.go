package demo

import (
	"github.com/gohade/hade/framework/cobra"
	"log"
)

var Demo2Command = &cobra.Command{
	Use:   "demo2",
	Short: "demo for app",
	RunE: func(c *cobra.Command, args []string) error {
		container := c.GetContainer()
		log.Println(container)
		return nil
	},
}
