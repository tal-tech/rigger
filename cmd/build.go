package cmd

import (
	"bytes"
	"fmt"
	"os"
	osexec "os/exec"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/tal-tech/rigger/common"
)

var Build = &cobra.Command{
	Use:           "build",
	Short:         "使用makefile编译项目",
	Run:           build,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func build(c *cobra.Command, args []string) {
	exists, err := common.PathExists("Makefile")
	if err != nil {
		fmt.Fprintln(os.Stdout, err)
		return
	}
	if !exists {
		fmt.Fprintln(os.Stdout, "没有发现makefile文件")
		return
	}

	if runtime.GOOS == "windows" {
		args = []string{"build-windows"}
	}

	cmd := osexec.Command("make", args...)

	var buffer bytes.Buffer
	cmd.Stderr = &buffer

	output, err := cmd.Output()

	handlerCmdOutput(output, err, buffer)
	return
}
