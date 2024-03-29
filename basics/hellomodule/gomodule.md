# Go module
Go module 构建模式是在 Go 1.11 版本正式引入的，为的是彻底解决 Go 项目复杂版本依赖的问题，
在 Go 1.16 版本中，Go module 已经成为了 Go 默认的包依赖管理机制和 Go 源码构建机制。

GO111MODULE 有三个值：off, on和auto（默认值）:
  1. GO111MODULE=off，go命令行将不会支持module功能，寻找依赖包的方式将会沿用旧版本那种通过vendor目录或者GOPATH模式来查找。
  2. GO111MODULE=on，go命令行会使用modules，而一点也不会去GOPATH目录下查找。
  3. GO111MODULE=auto，默认值，go命令行将会根据当前目录来决定是否启用module功能。

这种情况下可以分为两种情形：
1. 当前目录在GOPATH/src之外且该目录包含go.mod文件
2. 当前文件在包含go.mod文件的目录下面。

当modules功能启用时，依赖包的存放位置变更为$GOPATH/pkg，允许同一个package多个版本并存，且多个项目可以共享缓存的 module

修改GO111MODULE的值：
```
go env -w GO111MODULE=on
```

## 基本使用

Go Module 的核心是一个名为 go.mod 的文件，在这个文件中存储了这个 module 对第三方依赖的全部信息。

创建hellomodule.go, 直接go build会报错，先执行go mod：
```
$go mod init github.com/luoji_demo/hellomodule
go: creating new go.mod: module github.com/bigwhite/hellomodule
go: to add module requirements and sums:
go mod tidy
```

这时候生成了一个，go.mod文件
一个 module 就是一个包的集合，这些包和 module 一起打版本、发布和分发。
go.mod 所在的目录被我们称为它声明的 module 的根目录。

由于在hellomodule.go中引入了两个第三方包，我们可以手动引入，也可以使用如下命令引入：
```
$go mod tidy       
go: downloading go.uber.org/zap v1.18.1
go: downloading github.com/valyala/fasthttp v1.28.0
go: downloading github.com/andybalholm/brotli v1.0.2
... ...
```
执行完后发现，go.mod中也加入了三方包的引用, 并且 require 段中依赖了具体的版本号, 比如v1.32.0， 其中1是主版本，32是次版本，0是补丁版本
还生成了一个名为 go.sum 的文件，这个文件记录了 hellomodule 的直接依赖和间接依赖包的相关版本的 hash 值，
用来校验本地包的正确性。

对于一个处于初始状态的 module 而言，go mod tidy 分析了当前 main module 的所有源文件，
找出了当前 main module 的所有第三方依赖，确定第三方依赖的版本，还下载了当前 main module 的直接依赖包（比如 fasthttp），
以及相关间接依赖包。

go mod tidy 下载的依赖 module 会被放置在本地的 module 缓存路径下，默认值为 $GOPATH[0]/pkg/mod，
Go 1.15 及以后版本可以通过 GOMODCACHE 环境变量，自定义本地 module 的缓存路径。

这时候就可以 go build hellomodule.go,
在windows下生成的是hellomoudle.exe, Linux、Mac等是hellomoudle,
直接即可执行。

后续如果添加了新的依赖，有两种方式处理，
1还是使用go mod tidy， 2是可以直接使用go get, 这时候go get执行后会更新go.mod
这种场景建议使用方法1

另一种场景是依赖包需要回退到指定版本，这时候就使用go get指定版本号，比如：
```
$go get github.com/sirupsen/logrus@v1.7.0
或者
$go mod edit -require=github.com/sirupsen/logrus@v1.7.0
$go mod tidy
```

通过 go mod edit 命令可以修改 go.mod 文件中 require、replace 块的内容。

删除依赖包，除了要在go.mod中删除对应的依赖外，还要重新执行go mod tidy命令

另外说明的是，
如果我们要为 Go 项目添加主版本号大于 1 的依赖，我们就需要使用“语义导入版本”机制，
在声明它的导入路径的基础上，加上版本号信息，比如下面redis的v7就表示redis的7.x.y版本：
```
$go get github.com/go-redis/redis/v7
go: downloading github.com/go-redis/redis/v7 v7.4.1
go: downloading github.com/go-redis/redis v6.15.9+incompatible
go get: added github.com/go-redis/redis/v7 v7.4.1
```

Go 程序的入口函数只能是 main，这是 Go 语言的要求.
有了 Go Module 构建模式后，Go 项目可以放在任意目录，不一定要放在GOPATH下。

go mod tidy 不区分私有 Go Module 还是公共 Go Module，只要满足Go Module proxy 的下载协议，
就可以从中拉取依赖包。

Go 是一种编译型语言，这意味着只有你编译完 Go 程序之后，才可以将生成的可执行文件交付于其他人，
并运行在没有安装 Go 的环境中。

<br>

## 演进过程：
在1.11之前的Go版本中，
go get 下载的包只是那个时刻各个依赖包的最新主线版本，这样会给后续 Go 程序的构建带来一些问题。
比如，依赖包持续演进，可能会导致不同开发者在不同时间获取和编译同一个 Go 包时，得到不同的结果，也就是不能保证可重现的构建（Reproduceable Build）。
又比如，如果依赖包引入了不兼容代码，程序将无法通过编译。

可重现构建，顾名思义，就是针对同一份go module的源码进行构建，不同人，在不同机器(同一架构，比如都是x86-64)，
相同os上，在不同时间点都能得到相同的二进制文件。

Go 在 1.5 版本中引入 vendor 机制。vendor 机制本质上就是在 Go 项目的某个特定目录下，
将项目的所有依赖包缓存起来，这个特定目录名就是 vendor。
Go 编译器会优先感知和使用 vendor 目录下缓存的第三方包版本，而不是 GOPATH 环境变量所配置的路径下的第三方包版本。
但是，你还需要手工管理 vendor 下面的 Go 依赖包，包括项目依赖包的分析、版本的记录、依赖包获取和存放，等等，最让开发者头疼的就是这一点。

为了解决这个问题，Go 核心团队与社区将 Go 构建的重点转移到如何解决包依赖管理上。
Go 社区先后开发了诸如 gb、glide、dep 等工具，来帮助 Go 开发者对 vendor 下的第三方包进行自动依赖分析和管理，
但这些工具也都有自身的问题。


为什么现在还保留了vendor模式？
通常我们直接使用go module(非vendor)模式即可满足大部分需求。
如果是那种开发环境受限，因无法访问外部代理而无法通过go命令自动解决依赖和下载依赖的环境下，
我们通过vendor来辅助解决。
```
$go mod vendor
$tree -LF 2 vendor
vendor
├── github.com/
│   ├── google/
│   ├── magefile/
│   └── sirupsen/
├── golang.org/
│   └── x/
└── modules.txt
```

## 补充：replace

在项目的go.mod中, 使用replace命令的话, 可以对import路径进行替换. 也就是可以达到, 我们import的是a, 但在构建的时候, 实际使用的是b.
可以replace是成VCS(github或者其他地方), 或者文件系统路径(可以是绝对路径, 也可以是相对于项目根目录的相对路径). 使用文件系统路径这种, 在开发调试的使用, 非常有用.
最顶层go.mod的replace命令, 将影响到自身, 以及自身的所有依赖. 也就是可以间接改变依赖的项目的import. 
这样我们在fork的项目的import, 不用在fork项目里面的go.mod进行replace, 直接在使用的项目里replace即可.

例子：

编辑go.mod:
```
module nt-pdf-generator

go 1.12

require (
    github.com/xxx/abc v0.2.0 # replace的依赖, 必须要require
)

replace github.com/xxx/abc => ../github.com/hanFengSan/abc # 文件路径不需要附带版本信息
```
然后运行一下:
```
go mod tidy
```
大功告成.