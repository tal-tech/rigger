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

var Gens = &cobra.Command{
	Use:           "gens yourservice",
	Short:         "生成service代码",
	Run:           gens,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func gens(c *cobra.Command, args []string) {
	if len(args) < 1 {
		fmt.Fprintln(os.Stdout, "请使用rigger gens yourservice")
		return
	}

	serviceName := args[0]

	curDir, _ := common.GetCurPath()

	interfaceFile := curDir + "/app/serviceInterface/interface.go"

	parseResult, err := internal.GenParseResult(interfaceFile)
	if err != nil {
		fmt.Fprintf(os.Stdout, "解析interface.go文件失败:%v，请在项目根目录下执行rigger gens", err)
		return
	}

	common.CreateDir(curDir + "/app/service/")
	buffer := internal.GenService(parseResult, serviceName)
	serviceFile := curDir + "/app/service/service.go"
	common.WriteToFile(buffer, serviceFile, true)
	buffer = internal.GenServiceBridge(parseResult)
	serviceBridgeFile := curDir + "/app/service/serviceBridge.go"
	common.WriteToFile(buffer, serviceBridgeFile, true)
	buffer = internal.GenServiceInit(parseResult, serviceName)
	serviceInitFile := curDir + "/app/serviceInit.go"
	common.WriteToFile(buffer, serviceInitFile, true)

	//Impl File
	fileBuffer := make(map[string]*bytes.Buffer, 0)
	for _, fn := range parseResult.Fns {
		buffer, ok := fileBuffer[fn.Comment]
		if !ok {
			buffer = internal.GenImplFile(parseResult.Imports, fn)
			fileBuffer[fn.Comment] = buffer
		}
		buffer.WriteString(internal.GenImplFunc(fn))
	}
	for k, b := range fileBuffer {
		filename := strings.TrimLeft(k, "//")
		filename = strings.Replace(filename, ".", "/", -1)
		last := strings.LastIndex(filename, "/")
		if last > 0 {
			fmt.Println(last, filename[:last])
			common.CreateDir(curDir + "/app/serviceImpl/" + filename[:last])
		}
		filename = filename + ".go"
		common.WriteToFile(b, curDir+"/app/serviceImpl/"+filename, false)
	}
	/*
		rpcDir := os.Getenv("GOPATH") + "/src/" + serviceName + "/app/rpc"
		if ProjectPath != "" {
			rpcDir = strings.TrimRight(ProjectPath, serviceName) + "/" + serviceName + "/app/rpc"
		}
		ImportPath = strings.TrimRight(ImportPath, serviceName)

		dirList, err := ioutil.ReadDir(rpcDir)
		if err != nil {
			fmt.Fprintln(os.Stdout, err)
			return
		}

		var buffer *bytes.Buffer

		for _, v := range dirList {
			var outputFile string
			rpcFile := rpcDir + "/" + v.Name()

			if code == "php" {
				common.CreateDir(curDir + "/http/")
				outputFile = curDir + "/http/" + v.Name()
				buffer, err = internal.GenPHPHttpClient(rpcFile)
				outputFile = strings.TrimRight(outputFile, ".go")
				outputFile = outputFile + ".php"
			} else if code == "go" {
				common.CreateDir(curDir + "/rpc/")
				outputFile = curDir + "/rpc/" + v.Name()
				buffer, err = internal.GenGoRpcClient(rpcFile, BasePath, ImportPath)
			}

			if err != nil {
				fmt.Fprintf(os.Stdout, "生成代码失败 %v\n", err)
				break
			}

			err = common.WriteToFile(buffer, outputFile)

			if err != nil {
				fmt.Fprintf(os.Stdout, "生成代码失败 %v\n", err)
				continue
			}

			fmt.Fprintln(os.Stdout, "客户端代码生成成功:"+outputFile)
		}
	*/
	return

}
