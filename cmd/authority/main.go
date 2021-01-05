package main

import (
	"fmt"
	"github.com/myxy99/component/pkg/xcolor"
	"github.com/myxy99/component/pkg/xflag"
	"github.com/myxy99/component/pkg/xsignals"
	"github.com/myxy99/ndisk/cmd/authority/app"
	"os"
)

func main() {
	xflag.NewRootCommand(&xflag.Command{
		Use:                "NAuth",
		DisableSuggestions: false,
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
