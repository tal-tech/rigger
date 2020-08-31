package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tal-tech/rigger/internal"
)

var Genfn = &cobra.Command{
	Use:           "genfn [zhongtai|irc|oa|tiku] apiName funcName path",
	Short:         "生成xesSDK Func代码",
	Run:           genfn,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func genfn(c *cobra.Command, args []string) {
	if len(args) < 4 {
		fmt.Fprintln(os.Stdout, "请使用rigger genfn [zhongtai|irc|oa|tiku] apiName funcName path")
		return
	}

	tplType := args[0]

	apiName := args[1]

	funcName := args[2]

	path := args[3]

	buffer, _ := internal.GenXesSDKFunc(tplType, apiName, funcName, path)

	fmt.Printf("%s\n", buffer)

	return

}
