package internal

import (
	"bytes"
	"fmt"
	"strings"
)

func GenClientProto(parse *ParseResult, serviceName string) *bytes.Buffer {
	//imports
	buffer := bytes.NewBufferString("package proto\n\n")
	buffer.WriteString(`type SayHelloRequest struct {
	Greeting string
}
type SayHelloResponse struct {
	Reply string
}

type UserInfoRequest struct {
	Id int
}
type UserInfoResponse struct {
	Name string
	Age  int
	City string
}

type AddUserRequest struct {
	Name string
	Age  int
	City string
}
type AddUserResponse struct {
	Id int
}

type UpdateUserRequest struct {
	Id   int
	Name string
	Age  int
	City string
}

type UpdateUserResponse struct {
}

`)

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

func GenApiRouter(parse *ParseResult, serviceName string) *bytes.Buffer {
	//imports
	buffer := bytes.NewBufferString("package router\n\n")
	buffer.WriteString(fmt.Sprintf(`import (
	"%s/app/controller/demo"
	"github.com/gin-gonic/gin"
)

//The routing method is exactly the same as Gin                                                                     
func RegisterRouter(router *gin.Engine) {
	entry := router.Group("/demo")
	entry.GET("/test", demo.%sDemo)

`, serviceName, strings.Title(serviceName)))

	//struct
	for _, fn := range parse.Fns {
		buffer.WriteString(fmt.Sprintf("	entry.GET(\"/%v\", demo.%v)\n", strings.ToLower(fn.Name), fn.Name))
	}
	buffer.WriteString(`}`)
	return buffer
}

func GenApiController(parse *ParseResult, serviceName string) *bytes.Buffer {
	//imports
	buffer := bytes.NewBufferString("package demo\n\n")
	buffer.WriteString(fmt.Sprintf(`import (
	"net/http"

	"%s/app/service/demo"
	"%s/utils"

	"github.com/gin-gonic/gin"
)

`, serviceName, serviceName))

	//struct
	for _, fn := range parse.Fns {
		buffer.WriteString(fmt.Sprintf(`func %s(ctx *gin.Context) {
	goCtx := utils.TransferToContext(ctx)
	ret, err := demo.%s(goCtx)
	if err != nil {
		resp := utils.Error(err)
		ctx.JSON(http.StatusOK, resp)
	} else {
		resp := utils.Success(ret)
		ctx.JSON(http.StatusOK, resp)
	}
}`, fn.Name, fn.Name))
		buffer.WriteString("\n\n")
	}
	return buffer
}

func GenApiService(parse *ParseResult, serviceName string) *bytes.Buffer {
	//imports
	buffer := bytes.NewBufferString("package demo\n\n")
	buffer.WriteString(fmt.Sprintf(`import (
	"context"
	"%s/proto"
	"%s/rpc"

	logger "github.com/tal-tech/loggerX"
)

`, serviceName, serviceName))

	//struct
	for _, fn := range parse.Fns {
		buffer.WriteString(fmt.Sprintf(`func %s(ctx context.Context) (interface{}, error) {
	ins, _ := rpc.NewClient()
	req := proto.%sRequest{}
	resp := proto.%sResponse{}
	err := ins.%s(ctx, &req, &resp)
	if err != nil {
		logger.Ex(ctx, "service error:%%v", err)
	}
	return resp, err
}`, fn.Name, fn.Name, fn.Name, fn.Name))
		buffer.WriteString("\n\n")
	}
	return buffer
}
