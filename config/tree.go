package config

type ModuleInfo struct {
	Name     string
	Describe string
	Summary  string
	MainPage string
	GitPath  string
	GoGet    string
}

var Modules []ModuleInfo

func init() {
	Modules = make([]ModuleInfo, 0)
	Modules = append(Modules, gaea)
	Modules = append(Modules, odin)
	Modules = append(Modules, pan)
	Modules = append(Modules, triton)
	Modules = append(Modules, redis)
	Modules = append(Modules, mysql)
	Modules = append(Modules, logger)
	Modules = append(Modules, hera)
	Modules = append(Modules, tools)
	Modules = append(Modules, microPlugin)
}

var gaea ModuleInfo = ModuleInfo{Name: "gaea", Describe: "简单好用,高性能Web框架", MainPage: "", GitPath: "https://github.com/tal-tech/gaea", Summary: "请使用：rigger new api yourservicename 创建新的API项目"}
var odin ModuleInfo = ModuleInfo{Name: "odin", Describe: "高性能RPC服务框架,使用简单,功能强大,界面化管理", MainPage: "", GitPath: "https://github.com/tal-tech/odin", Summary: "请使用：rigger new micro yourservicename 创建新的RPC项目"}
var pan ModuleInfo = ModuleInfo{Name: "pan", Describe: "[MQ生产管家]高性能高稳定MQ代理服务", MainPage: "", GitPath: "https://github.com/tal-tech/pan"}
var triton ModuleInfo = ModuleInfo{Name: "triton", Describe: "[MQ消费管家]高性能,高稳定,模板化配置,插件化接入,多样化处理机制,分布式支持", MainPage: "", GitPath: "https://github.com/tal-tech/triton", Summary: "请使用：rigger new async yourservicename 创建新的队列消费项目"}
var redis ModuleInfo = ModuleInfo{Name: "xredis", Describe: "[redis管家,redis客户端]高性能,高稳定,简单接入,灵活管理", MainPage: "", GitPath: "https://github.com/tal-tech/xredis", GoGet: "github.com/tal-tech/xredis"}
var mysql ModuleInfo = ModuleInfo{Name: "torm", Describe: "[mysql管家,mysql客户端]高性能,高稳定,配置灵活,上手简单", MainPage: "", GitPath: "https://github.com/tal-tech/torm", GoGet: "github.com/tal-tech/torm"}
var logger ModuleInfo = ModuleInfo{Name: "loggerX", Describe: "强大的日志组件,高性能磁盘写入,配置多样化,插件化支持,多种落地方案", MainPage: "", GitPath: "https://github.com/tal-tech/loggerX", GoGet: "github.com/tal-tech/loggerX"}
var hera ModuleInfo = ModuleInfo{Name: "hera", Describe: "[服务孵化组件]快速搭建http,rpc,kafka服务框架", MainPage: "", GitPath: "https://github.com/tal-tech/hera", GoGet: "github.com/tal-tech/hera"}
var tools ModuleInfo = ModuleInfo{Name: "xtools", Describe: "超齐全的golang工具库,限流库,打点工具......", MainPage: "", GitPath: "https://github.com/tal-tech/xtools", GoGet: "github.com/tal-tech/xtools"}
var microPlugin ModuleInfo = ModuleInfo{Name: "odinPugin", Describe: "odin插件管理,例如监控,限流......", MainPage: "", GitPath: "https://github.com/tal-tech/odinPlugin", GoGet: "github.com/tal-tech/odinPlugin"}
