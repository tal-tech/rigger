package cmd

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tal-tech/rigger/common"
	"github.com/tal-tech/rigger/internal"
)

var BasePath string

var ProjectPath string

var ImportPath string

var Genc = &cobra.Command{
	Use:           "genc [go|php] yourservice",
	Short:         "生成sdk代码",
	Run:           genc,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func genc(c *cobra.Command, args []string) {
	if len(args) < 2 {
		fmt.Fprintln(os.Stdout, "请使用rigger genc php yourservice 或rigger genc go yourservice 命令来生成客户端代码")
		return
	}

	code := args[0]

	serviceName := args[1]

	curDir, _ := common.GetCurPath()

	newrpcfile := "service/service.go"

	oldrpcfile := "rpc/" + serviceName + ".go"

	rpcFile := os.Getenv("GOPATH") + "/src/" + serviceName + "/app/"
	if ProjectPath != "" {
		rpcFile = strings.TrimRight(ProjectPath, serviceName) + "/" + serviceName + "/app/"
	}
	if ok, _ := common.PathExists(rpcFile + oldrpcfile); ok {
		rpcFile = rpcFile + oldrpcfile
	} else {
		// Support search current path
		if ok, _ := common.PathExists(curDir + "/app/" + newrpcfile); ok {
			rpcFile = curDir + "/app/" + newrpcfile
		} else {
			rpcFile = rpcFile + newrpcfile
		}
	}
	ImportPath = strings.TrimRight(ImportPath, serviceName)

	var buffer *bytes.Buffer

	var outputFile string

	var err error

	if code == "php" {
		common.CreateDir(curDir + "/http/")
		outputFile = curDir + "/http/" + serviceName + ".go"
		buffer, err = internal.GenPHPHttpClient(rpcFile)
		outputFile = strings.TrimRight(outputFile, ".go")
		outputFile = outputFile + ".php"
	} else if code == "go" {
		common.CreateDir(curDir + "/rpc/")
		outputFile = curDir + "/rpc/" + serviceName + ".go"
		buffer, err = internal.GenGoRpcClient(rpcFile, BasePath, ImportPath)
	}

	if err != nil {
		fmt.Fprintf(os.Stdout, "生成代码失败 %v\n", err)
	}

	err = common.WriteToFile(buffer, outputFile, false)

	if err != nil {
		fmt.Fprintf(os.Stdout, "生成代码失败 %v\n", err)
	}

	fmt.Fprintln(os.Stdout, "客户端代码生成成功:"+outputFile)
	return

}
