package cmd

import (
	"bytes"
	osexec "os/exec"

	"github.com/spf13/cobra"
	"github.com/tal-tech/rigger/common"
)

var Example = &cobra.Command{
	Use:           "example [tags]",
	Short:         "运行example目录中的调用示例",
	Run:           example,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func example(c *cobra.Command, args []string) {

	var tag string

	if len(args) > 0 {
		tag = args[0]
	}

	serviceName, _ := common.GetServiceName()

	arg := `go run -tags "` + tag + `" ` + getServiceDir(serviceName) + `/examples/main.go` +
		` -p=` + getServiceDir(serviceName) + ` -c=conf/conf.ini`

	cmd := osexec.Command(syscmd, "-c", arg)

	var buffer bytes.Buffer
	cmd.Stderr = &buffer

	output, err := cmd.Output()

	handlerCmdOutput(output, err, buffer)
	return
}
