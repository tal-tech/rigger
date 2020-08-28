package cmd

import (
	"bytes"
	"fmt"
	"os"
	osexec "os/exec"

	"github.com/spf13/cobra"
	"github.com/tal-tech/rigger/common"
)

var Clean = &cobra.Command{
	Use:           "clean",
	Short:         "清理编译产生的文件",
	Run:           clean,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func clean(c *cobra.Command, args []string) {
	exists, err := common.PathExists("Makefile")
	if err != nil {
		fmt.Fprintln(os.Stdout, err)
		return
	}
	if !exists {
		fmt.Fprintln(os.Stdout, "没有发现makefile文件")
		return
	}

	cmd := osexec.Command("make", "clean")

	var buffer bytes.Buffer
	cmd.Stderr = &buffer

	output, err := cmd.Output()

	handlerCmdOutput(output, err, buffer)
	return
}
