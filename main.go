package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tal-tech/rigger/cmd"
)

var version bool

const VERSION = "v1.0.1"

var RootCmd = &cobra.Command{
	Use: "rigger",
	RunE: func(c *cobra.Command, args []string) error {
		if version {
			fmt.Println(VERSION)
			return nil
		}
		return c.Usage()
	},
}

func init() {
	RootCmd.AddCommand(cmd.Build)
	RootCmd.AddCommand(cmd.Clean)
	RootCmd.AddCommand(cmd.Example)
	RootCmd.AddCommand(cmd.Genc)
	RootCmd.AddCommand(cmd.Gens)
	RootCmd.AddCommand(cmd.Gensdemo)
	RootCmd.AddCommand(cmd.Gencdemo)
	RootCmd.AddCommand(cmd.New)
	//RootCmd.AddCommand(cmd.Help)
	RootCmd.AddCommand(cmd.Start)
	RootCmd.AddCommand(cmd.Restart)
	RootCmd.AddCommand(cmd.Status)
	RootCmd.AddCommand(cmd.Stop)
	RootCmd.AddCommand(cmd.Tag)
	RootCmd.AddCommand(cmd.Fswatch)
	RootCmd.AddCommand(cmd.Frame)
	RootCmd.AddCommand(cmd.Tree)
	RootCmd.AddCommand(cmd.Reverse)
	cmd.Genc.Flags().StringVarP(&cmd.BasePath, "basepath", "b", "tal_tech", "service BasePath")
	cmd.Genc.Flags().StringVarP(&cmd.ProjectPath, "projectpath", "p", "", "your rpc project path")
	cmd.Genc.Flags().StringVarP(&cmd.ImportPath, "importpath", "i", "", "service proto import path")
	cmd.Gencdemo.Flags().StringVarP(&cmd.ProjectPath, "projectpath", "p", "", "your rpc project path")
	cmd.Restart.Flags().BoolVarP(&cmd.Rebuild, "rebuild", "r", false, "build before start")
	cmd.Start.Flags().BoolVarP(&cmd.Foreground, "foreground", "f", false, "run at foreground")
	cmd.Fswatch.Flags().BoolVarP(&cmd.Foreground, "foreground", "f", false, "run at foreground")
	cmd.New.Flags().StringVarP(&cmd.DefaultName, "defaultnamereplace", "d", "", "template project name to replace(gaea->myproject)")
	cmd.New.Flags().StringVarP(&cmd.TitleName, "titlenamereplace", "t", "", "template project name to replace(Gaea->Myproject)")
	cmd.New.Flags().StringVarP(&cmd.GitUrl, "giturl", "g", "", "template project git url(git@github.com/tal-tech/gaea.git)")
	RootCmd.Flags().BoolVarP(&version, "version", "v", false, "show version")
	cmd.Reverse.Flags().BoolVarP(&cmd.ReverseSingle, "single", "s", false, "Generated one go file for every table")
	cmd.Reverse.Flags().StringVarP(&cmd.ReverseTmplPath, "tmplpath", "t", "", "Template dir for generated. the default templates dir has provide 1 template")
}

func main() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	//cmd.Run()
}
