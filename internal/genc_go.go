package internal

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"

	"github.com/tal-tech/rigger/common"
)

type Fn struct {
	Args []Arg
	Name string
}

type Arg struct {
	Star string
	Name string
	X    string
	Sel  string
}

func GenGoRpcClient(input string, basePath string, importPath string) (*bytes.Buffer, error) {
	exists, err := common.PathExists(input)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("file " + input + " not exist")
	}

	fset := token.NewFileSet()

	fs, err := parser.ParseFile(fset, input, nil, parser.ParseComments)

	if err != nil {
		return nil, err
	}

	buffer := bytes.NewBufferString(genImport(importPath, strings.ToLower(getServiceName(fs))) + "\n")
	buffer.WriteString(genVar(getServiceName(fs), basePath) + "\n" + genStruct() + "\n")
	buffer.WriteString(genNewFunc() + "\n")

	for _, decl := range fs.Decls {
		fn := parseFunc(decl)

		if fn == nil {
			continue

		}
		if len(fn.Args) != 3 {
			continue
		}

		buffer.WriteString(genFunc(fn.Args[0], fn.Name, fn.Args[1], fn.Args[2]) + "\n")
	}

	return buffer, nil
}

func GenEtcdDiscovery() *bytes.Buffer {
	tpl := `//+build etcd

package rpc

import (
	rpcxClient "github.com/smallnest/rpcx/client"
)

func NewDiscovery(basePath string, serviceName string, addrs []string) rpcxClient.ServiceDiscovery {

	return rpcxClient.NewEtcdDiscovery(basePath, serviceName, addrs, nil)
		
}`

	return bytes.NewBufferString(tpl)
}

func GenZkDiscovery() *bytes.Buffer {
	tpl := `//+build zookeeper

package rpc

import (
	rpcxClient "github.com/smallnest/rpcx/client"
)

func NewDiscovery(basePath string, serviceName string, addrs []string) rpcxClient.ServiceDiscovery {

	return rpcxClient.NewZookeeperDiscovery(basePath, serviceName, addrs, nil)
		
}`

	return bytes.NewBufferString(tpl)
}

func getFieldType(field interface{}) string {
	return fmt.Sprintf("%T", field)
}

func genImport(importPath, serviceName string) string {
	importTpl := `package rpc

import (
	"%s/proto"
	"github.com/tal-tech/xtools/rpcxutil"
	"context"
	"sync"
)`

	return fmt.Sprintf(importTpl, importPath)
}

func genVar(serviceName string, basePath string) string {
	varTpl := `
var (
	client *Client
	once   sync.Once
)

const (
	ServiceName = "%s"
	BasePath    = "/%s"
)`

	return fmt.Sprintf(varTpl, serviceName, basePath)
}

func genStruct() string {
	structTpl := `

type Client struct {
	wrapclient *rpcxutil.WrapClient

	//使用本地服务发现的client，仅当注册中心故障时启用
	localWrapclient *rpcxutil.WrapClient
	lock sync.RWMutex
}`

	return structTpl
}

func genNewFunc() string {
	funcTpl := fmt.Sprintf(`
// NewClient 只是获取一个client对象，此时并没有连接服务发现中心
// 调用时才会去连接服务发现中心
// 不需传入注册中心addr，默认读取ini配置
//[Registration]
//addrs=127.0.0.1:2379 127.0.0.1:2379
//group=online/gray/release/dev
func NewClient() (*Client,error) {
	if client != nil {
		return client, nil
	}

	once.Do(func(){
		client = &Client{}
	})

	return client, nil
}

func GetClient() *Client {
	return client
}

func (c *Client) getRpcxClient() (wrapc *rpcxutil.WrapClient) {
	if c.wrapclient != nil {
		return c.wrapclient
	}

	//logger.D("rpcClient", "获取rpcxclient:%v", rpcxutil.REGMonitor.IsInFault())
	if c.localWrapclient != nil && rpcxutil.REGMonitor.IsInFault() {
		//logger.I("rpcClient", "启用本地服务发现")
		return c.localWrapclient
	}

	c.lock.Lock()
	defer c.lock.Unlock()

	//防止重复实例化, 获得锁之后再次检查
	if c.wrapclient != nil {
		return c.wrapclient
	}

	if c.localWrapclient != nil && rpcxutil.REGMonitor.IsInFault() {
		return c.localWrapclient
	}

	opt := rpcxutil.GetClientOption()
	defer func() {
		if r := recover(); r != nil {
			if rpcxutil.REGMonitor.IsInFault() {
				c.localWrapclient = rpcxutil.NewLocalWrapClient(BasePath, ServiceName, rpcxutil.GetFailMode(), rpcxutil.GetSelectMode(), opt)
				wrapc = c.localWrapclient
				return
			}

			//未启用服务发现兜底，则保持原先控制逻辑，抛出异常阻止服务启动
			panic("注册中心故障！")
		}
	}()

	c.wrapclient = rpcxutil.NewWrapClient(BasePath, ServiceName, rpcxutil.GetFailMode(), rpcxutil.GetSelectMode(), opt)

	return c.wrapclient
}`)

	return funcTpl
}

func genFunc(ctx Arg, fn string, arg Arg, reply Arg) string {
	argCtx := fmt.Sprintf("%s %s%s.%s", ctx.Name, ctx.Star, ctx.X, ctx.Sel)
	argArg := fmt.Sprintf("%s %s%s.%s", arg.Name, arg.Star, arg.X, arg.Sel)
	argReply := fmt.Sprintf("%s %s%s.%s", reply.Name, reply.Star, arg.X, reply.Sel)

	as := fmt.Sprintf("%s, %s, %s", argCtx, argArg, argReply)

	body := fmt.Sprintf(`
func (c *Client) %s(%s) error{
	wrapclient := c.getRpcxClient()
	return wrapclient.WrapCall(%s, "%s", %s, %s)
}`, fn, as, ctx.Name, fn, arg.Name, reply.Name)

	return body
}

func getServiceName(fs *ast.File) string {
	serviceName := ""
	for _, v := range fs.Scope.Objects {
		if v.Kind != ast.Typ {
			continue
		}
		serviceName = v.Name
		break
	}

	return serviceName
}

func parseFunc(decl ast.Decl) *Fn {
	fn := &Fn{}

	fd, ok := decl.(*ast.FuncDecl)
	if ok {
		fn.Name = fd.Name.Name
		for _, field := range fd.Type.Params.List {
			switch getFieldType(field.Type) {
			case "*ast.SelectorExpr":
				ft, _ := field.Type.(*ast.SelectorExpr)
				arg := Arg{
					Star: "",
					X:    fmt.Sprintf("%s", ft.X),
					Name: field.Names[0].Name,
					Sel:  fmt.Sprintf("%s", ft.Sel),
				}
				fn.Args = append(fn.Args, arg)
			case "*ast.StarExpr":
				ft, _ := field.Type.(*ast.StarExpr)
				ftx, _ := ft.X.(*ast.SelectorExpr)
				arg := Arg{
					Star: "*",
					Name: field.Names[0].Name,
					X:    fmt.Sprintf("%s", ftx.X),
					Sel:  fmt.Sprintf("%s", ftx.Sel),
				}
				fn.Args = append(fn.Args, arg)
			}
		}
		return fn
	} else {
		return nil
	}
}
