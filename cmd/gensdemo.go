package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tal-tech/rigger/common"
	"github.com/tal-tech/rigger/internal"
)

var Gensdemo = &cobra.Command{
	Use:           "gensdemo yourservice",
	Short:         "生成server端proto代码",
	Run:           gensdemo,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func gensdemo(c *cobra.Command, args []string) {
	if len(args) < 1 {
		fmt.Fprintln(os.Stdout, "请使用rigger gensdemo yourservice")
		return
	}

	serviceName := args[0]

	curDir, _ := common.GetCurPath()

	interfaceFile := curDir + "/app/serviceInterface/interface.go"

	parseResult, err := internal.GenParseResult(interfaceFile)
	if err != nil {
		fmt.Fprintf(os.Stdout, "解析interface.go文件失败:%v，请在项目根目录下执行rigger gensproto", err)
		return
	}

	common.CreateDir(curDir + "/proto")
	buffer := internal.GenServerProto(parseResult, serviceName)
	serviceFile := curDir + "/proto/demo.go"
	common.WriteToFile(buffer, serviceFile, true)
	return

}
