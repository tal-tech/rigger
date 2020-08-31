package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tal-tech/rigger/xorm"
)

var (
	ReverseSingle   bool
	ReverseTmplPath string
)

var Reverse = &cobra.Command{
	Use:           "reverse",
	Short:         "一键生成MySQL表对象实体文件",
	Long:          "reverse [-s] [-t tmplPath] driverName datasourceName [generatedPath] [tableFilterReg]",
	Run:           runReverse,
	SilenceUsage:  true,
	SilenceErrors: true,
	//Args:          cobra.MinimumNArgs(2),
}

func runReverse(cmd *cobra.Command, args []string) {
	if len(args) < 2 {
		fmt.Fprintln(os.Stdout, "请使用rigger reverse [-s] [-t tmplPath] driverName datasourceName [generatedPath] [tableFilterReg]")
		return
	}
	c := xorm.NewReverse(
		xorm.SetReverseSingle(ReverseSingle),
		xorm.SetReverseTmplPath(ReverseTmplPath),
	)
	c.Run(cmd, args)
}
