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

type Imports struct {
	Text []string
}

type ServiceFn struct {
	Args    []Arg
	Name    string
	Comment string
}

type ParseResult struct {
	Imports []string
	Fns     []ServiceFn
}

func GenParseResult(input string) (*ParseResult, error) {
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

	result := new(ParseResult)

	imports := parseRpcImport(fs.Imports)

	result.Imports = imports

	for _, decl := range fs.Decls {
		fns := parseRpcFunc(decl)
		if len(fns) > 0 {
			result.Fns = fns
		}
	}
	if len(result.Fns) <= 0 {
		return nil, errors.New("Parse Error")
	}
	return result, nil
}

func GenService(parse *ParseResult, serviceName string) *bytes.Buffer {
	//imports
	buffer := bytes.NewBufferString(genServiceImport(parse.Imports, []string{`rpcxplugin "github.com/tal-tech/odinPlugin"`, `"github.com/tal-tech/odinPlugin/wrap"`, fmt.Sprintf(`"%s/app/serviceInterface"`, serviceName)}))

	//struct
	structTpl := `//rpcx服务注册类型
type %s struct {
	*rpcxplugin.RpcxPlugin
	service serviceInterface.Service
}`
	buffer.WriteString(fmt.Sprintf(structTpl+"\n\n", strings.Title(serviceName)))

	//Newstruct
	newStructTpl := `//传入实现Service接口的类型
func New%s(service serviceInterface.Service) *%s {
	this := new(%s)
	this.service = service
	return this
}`
	buffer.WriteString(fmt.Sprintf(newStructTpl+"\n\n", strings.Title(serviceName), strings.Title(serviceName), strings.Title(serviceName)))

	//func
	for _, fn := range parse.Fns {
		buffer.WriteString(genServiceFunc(serviceName, fn.Name, fn.Args[0], fn.Args[1], fn.Args[2]))
	}

	return buffer
}

func GenServiceBridge(parse *ParseResult) *bytes.Buffer {
	//imports
	buffer := bytes.NewBufferString(genServiceImport(parse.Imports, []string{}))

	//interface
	for _, fn := range parse.Fns {
		buffer.WriteString(genServiceBridgeInterface(fn.Name, fn.Args[0], fn.Args[1], fn.Args[2]))
	}

	//struct
	buffer.WriteString("type serviceBridge struct {\n")
	for _, fn := range parse.Fns {
		buffer.WriteString(fmt.Sprintf("	%sImpl    %s\n", fn.Name, fn.Name))
	}
	buffer.WriteString("}\n\n")

	//newStruct
	buffer.WriteString(`func NewServiceBridge() *serviceBridge {
	return new(serviceBridge)
}`)
	buffer.WriteString("\n\n")

	//func
	for _, fn := range parse.Fns {
		buffer.WriteString(genServiceBridgeFunc(fn.Name, fn.Args[0], fn.Args[1], fn.Args[2]))
	}

	return buffer
}

func GenServiceInit(parse *ParseResult, serviceName string) *bytes.Buffer {
	//import
	buffer := bytes.NewBufferString("package app\n\n")
	buffer.WriteString("import (\n")
	buffer.WriteString(fmt.Sprintf("	\"%s/app/service\"\n", serviceName))
	importFilter := make(map[string]bool, 0)
	for _, fn := range parse.Fns {
		impl := fmt.Sprintf("	\"%s/app/serviceImpl", serviceName)
		paths := parseComment(fn.Comment)
		for i := 1; i < len(paths); i++ {
			impl = impl + "/" + paths[i-1]
		}
		if _, ok := importFilter[impl]; !ok {
			importFilter[impl] = true
			buffer.WriteString(impl + "\"\n")
		}
	}
	buffer.WriteString(")\n\n")

	//newService
	buffer.WriteString(fmt.Sprintf("func NewService() *service.%s {\n", strings.Title(serviceName)))
	buffer.WriteString("	s := service.NewServiceBridge()\n")
	newFilter := make(map[string]bool, 0)
	for _, fn := range parse.Fns {
		paths := parseComment(fn.Comment)
		serviceName := paths[len(paths)-1]
		packageName := "serviceImpl"
		if len(paths) > 1 {
			packageName = paths[len(paths)-2]
		}
		if _, ok := newFilter[fn.Comment]; !ok {
			newFilter[fn.Comment] = true
			buffer.WriteString(fmt.Sprintf("	%s := %s.New%s()\n", serviceName, packageName, strings.Title(serviceName)))
		}
		buffer.WriteString(fmt.Sprintf("	s.%sImpl = %s\n", fn.Name, serviceName))
	}
	buffer.WriteString(fmt.Sprintf("	return service.New%s(s)\n", strings.Title(serviceName)))
	buffer.WriteString("}")

	return buffer
}

func GenImplFile(imports []string, fn ServiceFn) *bytes.Buffer {
	paths := parseComment(fn.Comment)
	serviceName := paths[len(paths)-1]
	packageName := "serviceImpl"
	if len(paths) > 1 {
		packageName = paths[len(paths)-2]
	}
	//imports
	buffer := bytes.NewBufferString("package " + packageName + "\n\n")
	buffer.WriteString("import (\n")
	for _, path := range imports {
		buffer.WriteString("	" + path + "\n")
	}
	buffer.WriteString(")\n\n")

	//struct
	buffer.WriteString(fmt.Sprintf("type %s struct {\n", strings.Title(serviceName)))
	buffer.WriteString("}\n\n")

	//newstruct
	buffer.WriteString(fmt.Sprintf("func New%s() *%s {\n", strings.Title(serviceName), strings.Title(serviceName)))
	buffer.WriteString(fmt.Sprintf("	return new(%s)\n", strings.Title(serviceName)))
	buffer.WriteString("}\n\n")
	return buffer
}

func GenImplFunc(fn ServiceFn) string {
	paths := parseComment(fn.Comment)
	serviceName := paths[len(paths)-1]
	return genServiceImplFunc(serviceName, fn.Name, fn.Args[0], fn.Args[1], fn.Args[2])
}

func parseComment(in string) []string {
	in = strings.TrimLeft(in, "//")
	return strings.Split(in, ".")
}

func genServiceImport(parseImports, otherImports []string) string {
	ret := `package service

import (
`
	for _, path := range parseImports {
		ret = ret + "	" + path + "\n"
	}
	for _, path := range otherImports {
		ret = ret + "	" + path + "\n"
	}
	return ret + ")\n\n"
}

func parseRpcImport(specs []*ast.ImportSpec) []string {
	ret := make([]string, 0)
	for _, spec := range specs {
		ret = append(ret, spec.Path.Value)
	}
	return ret
}

func parseRpcFunc(decl ast.Decl) []ServiceFn {
	fns := make([]ServiceFn, 0)
	var comment string
	fd, ok := decl.(*ast.GenDecl)
	if ok {
		for _, spec := range fd.Specs {
			s, ok := spec.(*ast.TypeSpec)
			if ok {
				inter, ok := s.Type.(*ast.InterfaceType)
				if ok {
					for _, method := range inter.Methods.List {
						var fn ServiceFn
						fn.Name = method.Names[0].Name
						if method.Doc != nil {
							for _, item := range method.Doc.List {
								comment = item.Text
							}
						}
						fn.Comment = comment
						funcType, ok := method.Type.(*ast.FuncType)
						if ok {
							for _, field := range funcType.Params.List {
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
						}
						if len(fn.Args) != 3 {
							continue
						}
						fns = append(fns, fn)
					}
				} else {
					return nil
				}
			} else {
				return nil
			}
		}
	} else {
		return nil
	}
	return fns
}

func genServiceFunc(serviceName, fnName string, ctx Arg, arg Arg, reply Arg) string {
	argCtx := fmt.Sprintf("%s %s%s.%s", ctx.Name, ctx.Star, ctx.X, ctx.Sel)
	argArg := fmt.Sprintf("%s %s%s.%s", arg.Name, arg.Star, arg.X, arg.Sel)
	argReply := fmt.Sprintf("%s %s%s.%s", reply.Name, reply.Star, arg.X, reply.Sel)

	args := fmt.Sprintf("%s, %s, %s", argCtx, argArg, argReply)
	callArgs := fmt.Sprintf("%s, %s, %s", ctx.Name, arg.Name, reply.Name)

	fnTpl := `//服务方法调用入口，先通过Wrapcall执行plugin方法，后调用service实现业务方法
func (this *%s) %s(%s) error {
	fn := func(w *wrap.Wrap) error {
		ctx := w.GetCtx()
		err := this.service.%s(%s)
		return err
	}
	return this.Wrapcall(ctx, "%s", fn)
}`
	return fmt.Sprintf(fnTpl+"\n\n", strings.Title(serviceName), fnName, args, fnName, callArgs, fnName)
}

func genServiceBridgeInterface(fnName string, ctx Arg, arg Arg, reply Arg) string {
	argCtx := fmt.Sprintf("%s%s.%s", ctx.Star, ctx.X, ctx.Sel)
	argArg := fmt.Sprintf("%s%s.%s", arg.Star, arg.X, arg.Sel)
	argReply := fmt.Sprintf("%s%s.%s", reply.Star, arg.X, reply.Sel)

	args := fmt.Sprintf("%s, %s, %s", argCtx, argArg, argReply)

	interfaceTpl := `type %s interface {
	%s(%s) error
}`
	return fmt.Sprintf(interfaceTpl+"\n\n", fnName, fnName, args)
}

func genServiceBridgeFunc(fnName string, ctx Arg, arg Arg, reply Arg) string {
	argCtx := fmt.Sprintf("%s %s%s.%s", ctx.Name, ctx.Star, ctx.X, ctx.Sel)
	argArg := fmt.Sprintf("%s %s%s.%s", arg.Name, arg.Star, arg.X, arg.Sel)
	argReply := fmt.Sprintf("%s %s%s.%s", reply.Name, reply.Star, arg.X, reply.Sel)

	args := fmt.Sprintf("%s, %s, %s", argCtx, argArg, argReply)
	callArgs := fmt.Sprintf("%s, %s, %s", ctx.Name, arg.Name, reply.Name)

	fnTpl := `func (s *serviceBridge) %s(%s) error {
	return s.%sImpl.%s(%s)
}`
	return fmt.Sprintf(fnTpl+"\n\n", fnName, args, fnName, fnName, callArgs)
}

func genServiceImplFunc(serviceName, fnName string, ctx Arg, arg Arg, reply Arg) string {
	argCtx := fmt.Sprintf("%s %s%s.%s", ctx.Name, ctx.Star, ctx.X, ctx.Sel)
	argArg := fmt.Sprintf("%s %s%s.%s", arg.Name, arg.Star, arg.X, arg.Sel)
	argReply := fmt.Sprintf("%s %s%s.%s", reply.Name, reply.Star, arg.X, reply.Sel)

	args := fmt.Sprintf("%s, %s, %s", argCtx, argArg, argReply)
	funcTpl := `func (this *%s) %s(%s) error {
	return nil
}`
	return fmt.Sprintf(funcTpl+"\n\n", strings.Title(serviceName), fnName, args)
}
