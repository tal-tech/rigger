package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tal-tech/rigger/common"
	"github.com/tal-tech/rigger/internal"
)

var Gencdemo = &cobra.Command{
	Use:           "gencdemo yourservice",
	Short:         "生成client端proto及api调用代码",
	Run:           gencdemo,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func gencdemo(c *cobra.Command, args []string) {
	if len(args) < 1 {
		fmt.Fprintln(os.Stdout, "请使用rigger gencdemo yourservice")
		return
	}

	serviceName := args[0]

	curDir, _ := common.GetCurPath()

	interfaceFile := curDir + "/app/serviceInterface/interface.go"
	if ProjectPath != "" {
		interfaceFile = ProjectPath + "/app/serviceInterface/interface.go"
	}

	parseResult, err := internal.GenParseResult(interfaceFile)
	if err != nil {
		fmt.Fprintf(os.Stdout, "解析interface.go文件失败:%v，请在项目根目录下执行rigger gensproto", err)
		return
	}

	common.CreateDir(curDir + "/proto")
	buffer := internal.GenClientProto(parseResult, serviceName)
	serviceFile := curDir + "/proto/rpc.go"
	common.WriteToFile(buffer, serviceFile, true)

	buffer = internal.GenApiRouter(parseResult, serviceName)
	serviceFile = curDir + "/app/router/router.go"
	common.WriteToFile(buffer, serviceFile, true)

	buffer = internal.GenApiController(parseResult, serviceName)
	serviceFile = curDir + "/app/controller/demo/rpc.go"
	common.WriteToFile(buffer, serviceFile, true)

	buffer = internal.GenApiService(parseResult, serviceName)
	serviceFile = curDir + "/app/service/demo/rpc.go"
	common.WriteToFile(buffer, serviceFile, true)
	return

}
