package cmd

import (
	"fmt"
	"os"
	osexec "os/exec"

	"github.com/spf13/cobra"
)

var Status = &cobra.Command{
	Use:           "status",
	Short:         "当前服务的运行状态",
	Run:           status,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func status(c *cobra.Command, args []string) {
	pidFile, _ := getPidFile()

	pid, err := readPid(pidFile)

	if err != nil {
		fmt.Fprint(os.Stdout, "读取pid失败 %v\n", err)
		return
	}

	arg := `ps aux |grep ` + pid + ` |grep -v grep`

	cmd := osexec.Command("/bin/sh", "-c", arg)

	output, err := cmd.Output()

	if err != nil {
		fmt.Fprintln(os.Stdout, "服务未启动，可以使用rigger start 启动你的服务")
		return
	}
	fmt.Fprint(os.Stdout, string(output))
	return
}
