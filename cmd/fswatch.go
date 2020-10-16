package cmd

import (
	"bytes"
	"fmt"
	"os"
	osexec "os/exec"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"github.com/tal-tech/rigger/common"
)

var Fswatch = &cobra.Command{
	Use:           "fswatch",
	Short:         "启动项目并watch",
	Run:           fswatch,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func fswatch(c *cobra.Command, args []string) {
	var buffer bytes.Buffer
	cmd := osexec.Command("which", "fswatch")
	cmd.Stdout = &buffer
	cmd.Stderr = os.Stdout
	err := cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stdout, "未找到fswatch，请先安装fswatch(go get github.com/codeskyblue/fswatch)\n")
		return
	}
	fsPath := buffer.String()
	if fsPath == "" {
		fmt.Fprintf(os.Stdout, "请先安装fswatch(go get github.com/codeskyblue/fswatch)\n")
		return
	}
	curdir, err := common.GetCurPath()

	if err != nil {
		fmt.Fprintln(os.Stdout, err)
	}

	serviceName, err := common.GetServiceName()
	if err != nil {
		fmt.Fprintln(os.Stdout, err)
	}

	detail, exist := processExistByName(serviceName)
	if exist {
		fmt.Fprintf(os.Stdout, "服务:%s 已启动,详细信息如下\n===============================\n%s",
			serviceName, detail)
		return
	}

	cmd = osexec.Command("fswatch")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	err = cmd.Start()
	if err != nil {
		fmt.Fprintf(os.Stdout, "启动失败 %v\n", err)
		return
	}
	pid := strconv.Itoa(cmd.Process.Pid)

	if Foreground {
		fmt.Fprintf(os.Stdout, "启动成功pid:%s\n", pid)
		cmd.Wait()
	}

	pinfo, _ := getProcessByPid(pid)
	if pinfo != "" {
		fmt.Fprintf(os.Stdout, "启动失败\n")
		return
	}

	pidFile, _ := getPidFile()

	err = common.CreateDir(curdir + "/run/")

	if err != nil {
		fmt.Fprintf(os.Stdout, "创建run目录失败 %v\n", err)
		return
	}

	f, err := os.OpenFile(pidFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		fmt.Fprintf(os.Stdout, "创建pid文件失败 %v\n", err)
		return
	}
	defer f.Close()

	f.WriteString(pid)
	time.Sleep(time.Second)
	fmt.Fprintf(os.Stdout, "启动成功pid:%s\n", pid)
}
