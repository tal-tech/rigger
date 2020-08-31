package cmd

import (
	"bytes"
	osexec "os/exec"
)

func shell(args []string) {
	cmdargs := ""
	for _, arg := range args {
		cmdargs = cmdargs + " " + arg
	}
	cmd := osexec.Command("/bin/bash", "-c", cmdargs)

	var buffer bytes.Buffer
	cmd.Stderr = &buffer

	output, err := cmd.Output()

	handlerCmdOutput(output, err, buffer)
}
