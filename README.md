# rigger
Very convenient scaffold components, support Gaea, Odin framework code generation, support process management, support code management

安装

```shell
go get github.com/tal-tech/rigger
```

命令介绍

```shell
rigger help //查看帮助命令
rigger build  //使用makefile编译项目
rigger clean  //使用makefile清理编译产生的文件 bin目录
rigger eg //运行example/rpc目录中的main.go调用示例
rigger genc //生成sdk代码
rigger genfn //生成xesSDK Func代码
rigger gens //生成service代码
rigger help //帮助
rigger new yourservicename//根据xes-micro工程模板创建项目
rigger shell //执行shell命令
rigger start //启动项目
rigger status //查看当前服务的运行状态
rigger stop //停止项目
rigger tag //使用git tag给项目打标签
```
