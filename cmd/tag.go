package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	osexec "os/exec"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var Tag = &cobra.Command{
	Use:   "tag",
	Short: "使用git tag给项目打标签",
	//Run:   printHelpInfo,
	RunE: func(c *cobra.Command, args []string) error {
		return c.Usage()
	},
}

var TagNow = &cobra.Command{
	Use:           "now",
	Short:         "展示当前tag",
	Run:           showVersionNow,
	SilenceUsage:  true,
	SilenceErrors: true,
}

var TagUp = &cobra.Command{
	Use:           "up",
	Short:         "升级tag,使用up x或up y或up z",
	Run:           upgradeVersion,
	SilenceUsage:  true,
	SilenceErrors: true,
}

var TagPush = &cobra.Command{
	Use:           "push",
	Short:         "推送到远端",
	Run:           pushToRemote,
	SilenceUsage:  true,
	SilenceErrors: true,
}

var TagInit = &cobra.Command{
	Use:           "init",
	Short:         "初始化一个tag",
	Run:           initTag,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	Tag.AddCommand(TagNow)
	Tag.AddCommand(TagUp)
	Tag.AddCommand(TagPush)
	Tag.AddCommand(TagInit)
}

/*
var (
	prompt  = "rigger-tag > "
	tagCmds = map[string]*Command{
		"now":  &Command{"now", "展示当前tag", showVersionNow},
		"up":   &Command{"up", "升级tag,使用up x或up y或up z", upgradeVersion},
		"push": &Command{"push", "推送到远端", pushToRemote},
		"exit": &Command{"exit", "退出命令行模式", exitCli},
		"init": &Command{"init", "初始化一个tag", initTag},
	}
)

func tag(args []string) {
	printHelpInfo()
	r, err := readline.New(prompt)

	if err != nil {
		fmt.Fprint(os.Stdout, err)
		os.Exit(1)
	}

	defer r.Close()

	if err := fetch(); err != nil {
		os.Exit(1)
	}

	for {

		args, err := r.Readline()
		if err != nil {
			fmt.Fprint(os.Stdout, err)

			return
		}

		args = strings.TrimSpace(args)

		// skip no args
		if len(args) == 0 {
			continue
		}

		parts := strings.Split(args, " ")
		if len(parts) == 0 {
			continue
		}

		name := parts[0]

		// get alias
		if n, ok := alias[name]; ok {
			name = n
		}

		if cmd, ok := tagCmds[name]; ok {
			cmd.exec(parts[1:])
		} else {
			helpTagCmd(parts[1:])
		}
	}
}
*/

func printHelpInfo(c *cobra.Command, args []string) {
	fmt.Fprintln(os.Stdout, "=====================================================")
	fmt.Fprintln(os.Stdout, "      version format v1.0.0")
	fmt.Fprintln(os.Stdout, "                      x.y.z")
	fmt.Fprintln(os.Stdout, "                      | | |--->修复bug或者添加小功能")
	fmt.Fprintln(os.Stdout, "                      | |----->添加较大功能")
	fmt.Fprintln(os.Stdout, "                      |------->重构项目或项目有重大更新")
	fmt.Fprintln(os.Stdout, "=====================================================")
}

/*
func helpTagCmd(args []string) {
	w := tabwriter.NewWriter(os.Stdout, 0, 8, 1, '\t', 0)

	fmt.Fprintln(os.Stdout, "Commands(用于管理tag号):")

	var keys []string
	for k, _ := range tagCmds {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		cmd := tagCmds[k]
		fmt.Fprintln(w, "\t", cmd.name, "\t\t", cmd.usage)
	}

	w.Flush()
}
*/
// version v1.0.0
func showNextVersion(next string) {
	fmt.Fprintf(os.Stdout, "即将更新tag到%s\n", next)
}

func upgradeVersion(c *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stdout, errors.New("请使用up x,up y或up z"))
		return
	}

	output, err := getVersionNow()

	if err != nil {
		fmt.Fprint(os.Stdout, err)
		return
	}

	if err := checkVersionFormat(output); err != nil {
		fmt.Fprintln(os.Stdout, err)
		return
	}

	branch, err := getBranch()
	if err != nil {
		fmt.Fprintf(os.Stdout, "获取当前分支失败 %v\n", err)
		return
	}
	if strings.Trim(string(branch), "\n") != "master" {
		fmt.Fprintln(os.Stdout, errors.New("必须在master分支更新tag"))
		return
	}
	next, err := getNextVersion(output, args[0])
	if err != nil {
		fmt.Fprintln(os.Stdout, err)
		return
	}
	showNextVersion(next)

	arg := `git tag -a ` + next + ` -m "更新版本到"` + next
	cmd := osexec.Command(syscmd, "-c", arg)
	var buffer bytes.Buffer
	cmd.Stderr = &buffer

	output, err = cmd.Output()

	handlerCmdOutput(output, err, buffer)

	nowversion, _ := getVersionNow()

	fmt.Fprintf(os.Stdout, "当前tag已更新到%s\n", nowversion)
}

func getNextVersion(now []byte, arg string) (string, error) {
	var next string
	//去掉换行符
	now = bytes.Trim(now, "\n")
	//去掉v
	version := fmt.Sprintf("%s", now[1:])
	//分割
	vers := strings.Split(version, ".")

	if len(vers) != 3 {
		return "", errors.New("tag格式错误")
	}

	var err error
	var x, y, z int
	if arg == "x" {
		x, err = strconv.Atoi(vers[0])
		x++
		vers[0] = strconv.Itoa(x)
	} else if arg == "y" {
		y, err = strconv.Atoi(vers[1])
		y++
		vers[1] = strconv.Itoa(y)
	} else if arg == "z" {
		z, err = strconv.Atoi(vers[2])
		z++
		vers[2] = strconv.Itoa(z)
	}

	if err != nil {
		return "", err
	}

	next = "v" + vers[0] + "." + vers[1] + "." + vers[2]

	return next, nil
}

func pushToRemote(c *cobra.Command, args []string) {
	arg := `git push origin --tags`
	cmd := osexec.Command(syscmd, "-c", arg)

	var buffer bytes.Buffer
	cmd.Stderr = &buffer

	output, err := cmd.Output()

	handlerCmdOutput(output, err, buffer)

	fmt.Fprintln(os.Stdout, "已推送所有tag到远端仓库")
}

func exitCli(args []string) {
	os.Exit(0)
}

func showVersionNow(c *cobra.Command, args []string) {
	output, err := getVersionNow()

	handlerCmdOutput(output, err, bytes.Buffer{})
}

func initTag(c *cobra.Command, args []string) {
	output, err := getVersionNow()
	if err != nil {
		fmt.Fprintln(os.Stdout, err)
		return
	}

	if len(output) > 0 {
		fmt.Fprintf(os.Stdout, "tag已经存在,无需初始化\n")
		return
	}

	arg := `git tag -a v0.0.1 -m "初始化tag"`
	cmd := osexec.Command(syscmd, "-c", arg)

	var buffer bytes.Buffer
	cmd.Stderr = &buffer

	output, err = cmd.Output()

	handlerCmdOutput(output, err, buffer)

	fmt.Fprintln(os.Stdout, "tag已经初始完成")
}

func getVersionNow() ([]byte, error) {
	arg := "git tag -l | sort -rV|head -n 1"
	//arg := "git tag -l --sort=-v:refname |head -n 1"
	cmd := osexec.Command(syscmd, "-c", arg)
	output, err := cmd.Output()

	return output, err
}

func handlerCmdOutput(output []byte, err error, buffer bytes.Buffer) {

	if err != nil {
		fmt.Fprintln(os.Stdout, err)
		fmt.Fprintln(os.Stdout, buffer.String())
		return
	}

	fmt.Fprint(os.Stdout, string(output))
}

func checkVersionFormat(version []byte) error {
	if len(version) == 0 {
		return errors.New("当前还没有tag")
	}

	if string(version[0]) != "v" {
		return errors.New("tag格式非法")
	}

	return nil
}

func fetch() error {
	arg := "git fetch -p"
	cmd := osexec.Command(syscmd, "-c", arg)
	var buffer bytes.Buffer
	cmd.Stderr = &buffer

	output, err := cmd.Output()

	handlerCmdOutput(output, err, buffer)

	if err != nil {
		return err
	}

	return nil
}

func getBranch() ([]byte, error) {
	arg := `git branch |grep \* |awk '{print $2}'`
	//arg := "git symbolic-ref --short -q HEAD"

	cmd := osexec.Command(syscmd, "-c", arg)
	output, err := cmd.Output()

	return output, err
}
