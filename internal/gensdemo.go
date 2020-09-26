package internal

import (
	"bytes"
	"fmt"
)

func GenServerProto(parse *ParseResult, serviceName string) *bytes.Buffer {
	//imports
	buffer := bytes.NewBufferString("package proto\n\n")

	//struct
	for _, fn := range parse.Fns {
		paths := parseComment(fn.Comment)
		if len(paths) <= 0 {
			continue
		}
		if paths[0] != "DemoService" {
			continue
		}
		buffer.WriteString(fmt.Sprintf("type %sRequest struct {\n}\n", fn.Name))
		buffer.WriteString(fmt.Sprintf("type %sResponse struct {\n	Reply []interface{}\n}\n\n", fn.Name))
	}

	return buffer
}
