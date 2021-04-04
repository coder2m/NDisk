package main

import (
	"fmt"
	"os"

	"github.com/coder2z/g-saber/xcolor"
	"github.com/coder2z/g-saber/xflag"
	"github.com/coder2z/g-saber/xsignals"

	"github.com/coder2z/ndisk/cmd/getway/app"
)

func main() {
	xflag.NewRootCommand(&xflag.CommandNode{
		Name: "NGetWay",
		Command: &xflag.Command{
			Use:                "NGetWay",
			DisableSuggestions: false,
		},
	})
	xflag.Register(
		xflag.CommandNode{
			Name: "run",
			Command: &xflag.Command{
				Use:   "run",
				Short: "run your app",
				Run: func(cmd *xflag.Command, args []string) {
					err := app.Run(xsignals.SetupSignalHandler())
					if err != nil {
						fmt.Println(xcolor.Red(err.Error()))
						os.Exit(1)
					}
				},
			},
			Flags: func(c *xflag.Command) {
				c.Flags().StringP("xcfg", "c", "", "配置文件")
				_ = c.MarkFlagRequired("xcfg")
			},
		},
	)
	_ = xflag.Parse()
}
