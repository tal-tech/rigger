package cmd

import (
	"bytes"
	"fmt"
	"os"
	osexec "os/exec"

	"github.com/spf13/cobra"
	"github.com/tal-tech/rigger/common"
)

var Frame = &cobra.Command{
	Use:           "frame",
	Short:         "管理框架插件",
	Long:          "请使用rigger frame [Plugin|Middleware] name(pprof/perf/expvar/maxfd|perf/trace)",
	Run:           frame,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func frame(c *cobra.Command, args []string) {
	if len(args) <= 1 {
		fmt.Fprintln(os.Stdout, "请使用rigger frame [Plugin|Middleware] name(pprof/perf/expvar/maxfd|perf/trace)")
		return
	}

	serviceName, _ := common.GetServiceName()

	addType := args[0]

	name := args[1]

	var replaceContent string
	switch addType {
	case "Plugin":
		replaceContent = getPlugin(name)
	case "Middleware":
		replaceContent = getMiddleware(name)
	default:
		fmt.Fprintln(os.Stdout, "请使用rigger frame [Plugin|Middleware] name(pprof/perf/expvar/maxfd|perf/trace)")
		return
	}

	if replaceContent == "" {
		fmt.Fprintln(os.Stdout, "请使用rigger frame [Plugin|Middleware] name(pprof/perf/expvar/maxfd|perf/trace)")
		return
	}

	command := sedI() + ` '\/\/Optional ` + addType + `/a\\t` + replaceContent + `' cmd/` + serviceName + `/main.go`

	cmd := osexec.Command("/bin/sh", "-c", command)

	var buffer bytes.Buffer
	cmd.Stderr = &buffer

	output, err := cmd.Output()

	handlerCmdOutput(output, err, buffer)
	return
}

func getPlugin(name string) string {
	switch name {
	case "pprof":
		return `s.AddBeforeServerStartFunc(bs.InitPprof())`
	case "perf":
		return `s.AddBeforeServerStartFunc(bs.InitPerfutil())`
	case "expvar":
		return `s.AddBeforeServerStartFunc(bs.InitExpvar())`
	case "maxfd":
		return `s.AddBeforeServerStartFunc(bs.InitMaxFd())`
	default:
		return ""
	}
}

func getMiddleware(name string) string {
	switch name {
	case "trace":
		return `engine.Use(middleware.TraceMiddleware())`
	case "perf":
		return `engine.Use(middleware.PerfMiddleware())`
	default:
		return ""
	}
}
