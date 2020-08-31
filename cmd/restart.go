package cmd

import (
	"time"

	"github.com/spf13/cobra"
)

var Restart = &cobra.Command{
	Use:           "restart",
	Short:         "重启项目",
	Run:           restart,
	SilenceUsage:  true,
	SilenceErrors: true,
}

var Rebuild bool

func restart(c *cobra.Command, args []string) {
	stop(c, args)
	time.Sleep(time.Millisecond * 500)
	if Rebuild {
		args = []string{}
		build(c, args)
	}
	time.Sleep(time.Millisecond * 500)
	start(c, args)
}
