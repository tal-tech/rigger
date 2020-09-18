# rigger
Very convenient scaffold components, support Gaea, Odin framework code generation, support process management, support code management

### 安装

```shell
go get github.com/tal-tech/rigger
```

### 命令介绍

#### build
* 用途：

    编译当前项目

* 命令：
    ```shell
    rigger build
    ```
* 说明：

    命令为在当前路径下执行make，需在项目根路径下执行，且要求Makefile一定存在


#### clean
* 用途：

    清理当前项目编译产生文件

* 命令：
    ```shell
    rigger clean
    ```
* 说明：

    命令为在当前路径下执行make clean，Makefile应设置clean target

#### example
* 用途：

    执行项目路径下examples的实例代码

* 命令：
    ```shell
    rigger example [tags]
    ```
* 说明：

    命令为在当前路径下执行go run examles/main.go，需要存在main.go文件，且执行时自动使用项目下conf部分项目，如xesMicro，编译时需要加入tags，此时指定第三个参数

#### genc
* 用途：

    生成rpc服务sdk代码

* 命令：
    ```shell
    rigger genc [go|php] yourservice [flags]
    ```
* 说明：
    
    * 命令通过Micro服务项目下指定路径app/rpc内的Service定义，生成对应调用sdk

    * 可生成go和php两种语言sdk，yourservice为rpc服务的项目名

    * 生成sdk路径为当前路径下rpc文件夹内，部分代码如下图：

    ![pic](https://wiki.zhiyinlou.com/download/attachments/34029786/image2019-11-6_16-37-14.png?version=1&modificationDate=1573029435000&api=v2)
         
    **flag说明：**

    * -b, --basepath string service BasePath (default "xes_xueyan_hudong")

    * -i, --importpath string service proto import path (default "git.100tal.com/wangxiao_xueyan_hudong/common")

    * -p, --projectpath string your project path
    ```
    BasePath为在注册中心中注册服务的前缀，可区分服务组，-b 参数修改默认值

    Importpath是sdk内引用服务公共代码，参数结构体proto路径，如图，为importpath/[servicename]/proto，-i 参数修改默认值

    projectpath为指定Micro服务端项目路径，默认值为$GOPATH/[servicename]，如不在此路径下，需通过-p 参数指定
    ```

#### gens
* 用途：

    生成rpc服务层代码

* 命令：
    ```shell
    rigger gens yourservice [flags]
    ```
* 说明：

    * rpc项目目录下执行，根据app/serviceInterfece/interface.go定义对外方法，一键生成项目代码，只填入逻辑代码即可
    

#### help
* 用途：
    
    同 -h

* 命令：

    ```shell
    rigger help
    ```

#### new
* 用途：
    
    根据工程模板创建项目

* 命令：
    ```shell
    rigger new [micro|api|async|custom] servicename [flags]
    ```
* 说明：
    * 命令可通过clone模板生成新的项目
    
    * servicename为新生成项目名，生成路径为$GOPATH/servicename

    * 目前模板可支持rpc服务odin(micro)、api服务gaea(api)、异步消费服务asyncworker(async)

    * 新生成项目内都附有简单demo，可通过rigger start、rigger example测试

    * custom生成自定义模板项目，-g gitlib地址，-d 默认项目待替换名称 -t 默认项目待替换为大写名称

#### start
* 用途：
    
    启动项目

* 命令：
    ```shell
    rigger start
    ```
* 说明：

    * 命令需在项目根路径下执行，启动服务为后台启动

    * 启动后会在当前路径下创建run/servername.pid文件，stop、status命令均通过pid进行操作
    
    * -f foreground 指定在前台运行



#### status
* 用途：

    当前服务的运行状态

* 命令：
    ```shell
    rigger status
    ```
* 说明：

    命令执行路径下需存在run/servername.pid文件

#### stop
* 用途：
        
    停止项目

* 命令：
    ```shell
    rigger stop
    ```
* 说明：

    命令执行路径下需存在run/servername.pid文件


#### restart
* 用途：

    重启项目

* 命令：
    ```shell
    rigger restart [flags]
    ```

* 说明：

    命令执行路径下需存在run/servername.pid文件

    -r flag执行顺序为stop、build、start，方便调试



#### tag
* 用途：

    用git tag对项目进行标签操作

* 命令：
    ```shell
    rigger tag [subcommand] [flags]
    ```
* 说明：
    
    tag命令，需指定子命令

    * init 初始化一个tag

    * now 展示当前tag

    * push 推送到远端

    * up 升级tag,使用up x或up y或up z

    `up子命令，项目推荐使用go标准tag格式，为vX.Y.Z，如v1.1.1，up命令后跟x、y、z，升级指定位置`


#### fswatch
* 用途：
    
    开发过程，watch文件变化，自动重新编译、重启

* 命令：
    ```
    rigger fswatch
    ```

* 说明：

    需先安装fswatch，go get github.com/codeskyblue/fswatch

    目前支持gaea和odin框架，请确认版本，可用版本为目录下有.fsw.yml，或从最新版内copy
    
    停止命令使用rigger stop

#### frame
* 用途：

    gaea&odin 启动项插件化管理，命令添加可选插件

* 命令：
    ```
    rigger frame [Plugin|Middleware] (pprof/perf/expvar/maxfd|perf/trace)
    ```

* 说明：
     
    支持gaea&odin的启动插件，以及gaea框架中的http中间件

    插件包括pprof性能分析插件、perf耗时打点插件、expvr内存分析插件、maxfd最大文件打开数管理


#### reverse
* 用途：

    一键生成MySQL表对象实体文件

* 命令：
    ```
    rigger reverse [-s] [-t tmplPath] driverName datasourceName [generatedPath] [tableFilterReg]
    ```

* 说明：
   
    -s 指定是否生成单一文件

    -t 指定生成模板，不指定使用默认模板

    generatedPath 生成文件目录

    tableFilterReg 表名匹配过滤

#### tree
* 用途：

   查看golang生态组件 

* 命令：
    ```
    rigger tree
    rigger trem [name]
    ```

* 说明：
   
   rigger tree命令查看所有golang组件

   rigger tree 组件名，可查看组件详情并下载组件
