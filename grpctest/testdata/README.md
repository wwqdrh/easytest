## protoc编译命令

安装: `go install google.golang.org/protobuf/cmd/protoc-gen-go@latest`

- `--go_out` => `protoc-gen-go`插件
- `--go-grpc_out` => `protoc-gen-go-grpc`插件
- `*_out` => `protoc-gen-*`

当使用参数 --go_out=xxx --go-grpc_out=xxx 生成时，会生成两个文件 *.pb.go 和 *._grpc.pb.go ，它们分别是消息序列化代码和 gRPC 代码。

如果在proto文件中定义了import，那么就需要引入相关依赖

`protoc -I ./ -I $GOPATH/src -I $GOPATH/src/google/api --go_out=plugins=grpc:.`

protoc3中要么在proto文件中加入`option go_package = "[包的全路径]"`要么在命令行中添加m参数

```shell
Please specify either:
        • a "go_package" option in the .proto source file, or
        • a "M" argument on the command line.
```

```bash
protoc --proto_path=./ \
  --go_opt=Mgreeter.proto=./helloworld \
  greeter.proto
```

~~`--go_out=plugins=grpc:.`~~, 已经不能这么写了，需要拆成`--go-grpc_out`

--go_out主要的两个参数为plugins 和 paths，分别表示生成Go代码所使用的插件，以及生成的Go代码的位置。--go_out的写法是，参数之间用逗号隔开，最后加上冒号来指定代码的生成位置，比如--go_out=plugins=grpc,paths=import:.

# 在greeter中添加go_package

`option go_package = "./helloworld";`

protoc --go_out=. --go-grpc_out=. greeter.proto