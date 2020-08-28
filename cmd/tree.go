package cmd

import (
	"bytes"
	"fmt"
	"os"
	osexec "os/exec"

	"github.com/spf13/cobra"
	"github.com/tal-tech/rigger/config"
)

var Tree = &cobra.Command{
	Use:           "tree",
	Short:         "查看golang生态组件",
	Run:           tree,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func tree(c *cobra.Command, args []string) {
	if len(args) == 0 {
		printAll()
		fmt.Fprintln(os.Stdout, "-----------------------------------")
		fmt.Fprintln(os.Stdout, "请使用rigger tree name查看详情")
	} else {
		name := args[0]
		var module config.ModuleInfo
		for _, m := range config.Modules {
			if m.Name == name {
				module = m
				break
			}
		}
		if module.Name != "" {
			fmt.Printf("%s:%s\n", module.Name, module.Describe)
			if module.MainPage != "" {
				fmt.Printf("文档地址:%s\n", module.MainPage)
			}
			fmt.Printf("git地址:%s\n", module.GitPath)
			if module.Summary != "" {
				fmt.Printf("%s\n", module.Summary)
			}
		}
		if module.GoGet != "" {
			fmt.Fprintln(os.Stdout, "------------------------")
			fmt.Fprintln(os.Stdout, "确认是否下载(go get)当前组件  y/n ?")
			var input string
			fmt.Scanln(&input)

			if input == "y" {
				cmd := osexec.Command("go", "get", module.GoGet)

				var buffer bytes.Buffer
				cmd.Stderr = &buffer

				output, err := cmd.Output()

				handlerCmdOutput(output, err, buffer)
			}
		}
	}
	return
}

func printAll() {
	for _, m := range config.Modules {
		fmt.Printf("%s: %s\n", m.Name, m.Describe)
	}
}
