package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	osexec "os/exec"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"github.com/tal-tech/rigger/common"
)

var Start = &cobra.Command{
	Use:           "start",
	Short:         "启动项目",
	Run:           start,
	SilenceUsage:  true,
	SilenceErrors: true,
}

var Foreground bool

func start(c *cobra.Command, args []string) {
	curdir, err := common.GetCurPath()

	if err != nil {
		fmt.Fprintln(os.Stdout, err)
	}

	serviceName, err := common.GetServiceName()
	if err != nil {
		fmt.Fprintln(os.Stdout, err)
	}

	pfile, _ := getPidFile()
	pidstr, err := ioutil.ReadFile(pfile)
	if err == nil {
		pid := string(pidstr)
		detail, err := getProcessByPid(pid)
		if detail != "" && err == nil {
			fmt.Fprintf(os.Stdout, "服务:%s 已启动,详细信息如下\n===============================\n%s",
				serviceName, detail)
			return
		}
	}

	arg := curdir + `/bin/` + serviceName

	if len(args) == 0 {
		args = append(args, "-p="+curdir, "-c=conf/conf.ini")
	} else {
		arg = args[0]
		args = args[1:]
	}

	cmd := osexec.Command(arg, args...)
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
	if pinfo == "" {
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

func getProcessByPid(pid string) (string, error) {
	arg := `ps aux |grep ` + pid + ` |grep -v grep`

	cmd := osexec.Command("/bin/sh", "-c", arg)

	output, err := cmd.Output()

	return string(output), err
}

func processExistByName(name string) (string, bool) {
	arg := `ps aux |grep ` + name + ` |grep -v grep`

	cmd := osexec.Command("/bin/sh", "-c", arg)

	output, err := cmd.Output()

	return string(output), err == nil
}
