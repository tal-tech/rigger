package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	osexec "os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tal-tech/rigger/common"
)

var Stop = &cobra.Command{
	Use:           "stop",
	Short:         "停止项目",
	Run:           stop,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func stop(c *cobra.Command, args []string) {
	pidFile, err := getPidFile()
	if err != nil {
		fmt.Fprintf(os.Stdout, "获取pid文件失败 %v\n", err)
		return
	}

	pid, err := readPid(pidFile)
	if err != nil {
		fmt.Fprintf(os.Stdout, "读取pid失败 %v\n", err)
		return
	}

	cmd := osexec.Command("kill", pid)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout

	err = cmd.Start()

	if err != nil {
		fmt.Fprintf(os.Stdout, "停止失败 %v\n", err)
		return
	}
	err = cmd.Wait()
	if err != nil {
		fmt.Fprintf(os.Stdout, "停止失败 %v\n", err)
		return
	}

	serviceName, _ := common.GetServiceName()
	fmt.Fprintln(os.Stdout, "进程"+serviceName+"(pid:"+string(pid)+")已停止")
	return
}

func getPidFile() (string, error) {
	curdir, err := common.GetCurPath()

	if err != nil {
		return "", err
	}

	serviceName, err := common.GetServiceName()
	if err != nil {
		return "", err
	}

	pidFile := curdir + "/run/" + serviceName + ".pid"

	return pidFile, nil
}

func readPid(pidFile string) (string, error) {
	f, err := os.OpenFile(pidFile, os.O_RDONLY, 0600)
	if err != nil {
		//fmt.Fprintf(os.Stdout, "读取pid文件失败 %v\n", err)
		return "", err
	}

	defer f.Close()

	pid, err := ioutil.ReadAll(f)
	if err != nil {
		//fmt.Fprintf(os.Stdout, "获取pid失败 %v", err)
		return "", err
	}

	return strings.TrimSpace(string(pid)), nil
}
