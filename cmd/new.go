package cmd

import (
	"bytes"
	"fmt"
	"os"
	osexec "os/exec"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tal-tech/rigger/common"
	"github.com/tal-tech/rigger/config"
)

var New = &cobra.Command{
	Use:           "new [micro|api|async|proxy|job|custom] [servicename]",
	Short:         "根据工程模板创建项目(job会在当前目录下生成)",
	Long:          "请使用rigger new templatename(micro/api/async/proxy/job/custom) yourservicename",
	Run:           new,
	SilenceUsage:  true,
	SilenceErrors: true,
}

var DefaultName string

var TitleName string

var GitUrl string

var DefaultReplacer *strings.Replacer

func init() {
	DefaultReplacer = strings.NewReplacer("\t", "", "\r", "", "\n", "", ".", "", "-", "")
}

func Filter(msg string) string {
	replacer := DefaultReplacer
	return replacer.Replace(msg)
}

var newProjectInfo config.NewProjectInfo

func new(c *cobra.Command, args []string) {
	if len(args) <= 1 {
		fmt.Fprintln(os.Stdout, "请使用rigger new templatename(micro/api/async/proxy/custom) yourservicename")
		return
	}
	serviceName := Filter(args[1])
	templatename := args[0]

	switch templatename {
	case "micro":
		newProjectInfo = config.NewOdinInfo
	case "api":
		newProjectInfo = config.NewGaeaInfo
	case "async":
		newProjectInfo = config.NewTritonInfo
	case "proxy":
		newProjectInfo = config.NewPanInfo
	case "job":
		newProjectInfo = config.NewJobInfo
		newJobPro(serviceName)
		return
	case "custom":
		if GitUrl == "" {
			fmt.Fprintln(os.Stdout, "请通过-g指定模板项目git地址")
			return
		}
		newProjectInfo.TplRepo = GitUrl
		if DefaultName != "" {
			newProjectInfo.ReplaceContent = append(newProjectInfo.ReplaceContent, config.ReplaceContentItem{DefaultName, config.DefaultReplaceName})
		}
		if TitleName != "" {
			newProjectInfo.ReplaceContent = append(newProjectInfo.ReplaceContent, config.ReplaceContentItem{TitleName, config.TitleReplaceName})
		}

	default:
		fmt.Fprintln(os.Stdout, "请使用rigger new templatename(micro/api/async/proxy/custom) yourservicename")
		return
	}

	exists, _ := common.PathExists(getServiceDir(serviceName))

	if exists {
		fmt.Fprintf(os.Stdout, "项目(%s)已存在\n", getServiceDir(serviceName))
	}

	if err := cloneTpl(serviceName); err != nil {
		fmt.Println(err)
		return
	}

	if err := cleanGitFile(serviceName); err != nil {
		fmt.Println(err)
		return
	}

	if err := replaceContent(serviceName); err != nil {
		fmt.Println(err)
		return
	}

	if err := replaceDir(serviceName); err != nil {
		fmt.Println(err)
		return
	}

	replaceFile(serviceName)

	fmt.Fprintln(os.Stdout, serviceName+"项目已创建完成, 使用:\n cd "+getServiceDir(serviceName)+" && rigger build \n开始你的微服务之旅！")
	return
}

func cloneTpl(serviceName string) error {
	//todo 放入gopath
	arg := "git clone " + newProjectInfo.TplRepo + " " + getServiceDir(serviceName)
	cmd := osexec.Command("/bin/sh", "-c", arg)

	var buffer bytes.Buffer
	cmd.Stderr = &buffer

	output, err := cmd.Output()
	handlerCmdOutput(output, err, buffer)

	return err
}

func cleanGitFile(serviceName string) error {
	arg := "pushd " + getServiceDir(serviceName) + "&& rm -rf .git && popd"
	cmd := osexec.Command("/bin/sh", "-c", arg)

	var buffer bytes.Buffer
	cmd.Stderr = &buffer

	_, err := cmd.Output()

	handlerCmdOutput([]byte{}, err, buffer)

	return err
}

func replaceContent(serviceName string) error {
	for _, item := range newProjectInfo.ReplaceContent {
		arg := `grep '` + item.Key + `' -rl ` + getServiceDir(serviceName) + `|xargs ` + sedI() + ` 's/` + item.Key + `/` + item.Fn(serviceName) + `/g'`

		cmd := osexec.Command("/bin/sh", "-c", arg)

		var buffer bytes.Buffer
		cmd.Stderr = &buffer

		output, err := cmd.Output()

		handlerCmdOutput(output, err, buffer)

		if err != nil {
			return err
		}
	}
	return nil
}

func replaceFile(serviceName string) error {
	for key, _ := range newProjectInfo.ReplaceFile {
		arg := `pushd ` + getServiceDir(serviceName) +
			`&&find . -name '` + key + `.go' |awk -F "` + key + `.go" '{print $1}' |` +
			`xargs -I'{}' mv {}` + key + `.go {}` + serviceName + `.go&&popd`

		cmd := osexec.Command("/bin/sh", "-c", arg)

		var buffer bytes.Buffer
		cmd.Stderr = &buffer

		_, err := cmd.Output()

		handlerCmdOutput([]byte{}, err, buffer)
	}
	return nil
}

func replaceDir(serviceName string) error {
	for key, _ := range newProjectInfo.ReplaceDir {
		arg := `pushd ` + getServiceDir(serviceName) +
			`&&find . -name '` + key + `' -type d |awk -F "` + key + `" '{print $1}'| ` +
			`xargs -I'{}' mv {}` + key + ` {}` + serviceName + `&&popd`

		cmd := osexec.Command("/bin/sh", "-c", arg)

		var buffer bytes.Buffer
		cmd.Stderr = &buffer

		_, err := cmd.Output()

		handlerCmdOutput([]byte{}, err, buffer)
	}
	return nil
}

func sedI() string {
	if runtime.GOOS == config.MacOS {
		return `sed -i ""`
	} else {
		return `sed -i`
	}
}

func getServiceDir(serviceName string) string {
	return os.Getenv("GOPATH") + "/src/" + serviceName
}

func newJobPro(serviceName string) {
	//仅克隆jobWorker
	exist, _ := pathExists(serviceName)
	if exist {
		handlerCmdOutput([]byte("["+serviceName+"] 此目录已存在!\n"), nil, bytes.Buffer{})
		return
	}

	tempDir := "/tmp/" + serviceName
	arg := "git clone " + newProjectInfo.TplRepo + " " + tempDir
	cmd := osexec.Command("/bin/sh", "-c", arg)

	var buffer bytes.Buffer
	cmd.Stderr = &buffer
	output, err := cmd.Output()
	handlerCmdOutput(output, err, buffer)
	if err != nil {
		return
	}

	buffer.Reset()
	shellCmd := `mv ` + tempDir + `/clijob/jobWorker ` + serviceName + " && rm -rf " + tempDir
	cmd = osexec.Command("/bin/sh", "-c", shellCmd)
	output, err = cmd.Output()
	handlerCmdOutput(output, err, buffer)
	if err != nil {
		return
	}

	fmt.Fprintln(os.Stdout, serviceName+" Job server 已生成!")
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
