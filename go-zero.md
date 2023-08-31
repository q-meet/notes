# go-zero 使用

## 个人demo地址

[https://github.com/q-meet/go-zero-mall](https://github.com/q-meet/go-zero-mall)

## 安装

```shell
go install github.com/zeromicro/go-zero/tools/go-ctl@v1.5.4
```

检查是否安装 protoc, protoc-gen-go, protoc--gen-go-grpc 没安装利用 goctl 安装

```shell
goctl env check -i -f --verbose
```

## 常用命令

- api demo 代码生成

```shell
goctl api new demo
```

- gRPC demo 代码生成
  
```shell
goctl rpc new demo
```

[什么是微服务](./go-zero-md/01/01.md)

[go-zero 入门](./go-zero-md/02/02.md)

[go-zero 存储](./go-zero-md/03/03.md)

[go-zero 功能:jwt,中间件,自定义错误,http自定义返回,修改模板](./go-zero-md/04/04.md)

[go-zero goctl使用](./go-zero-md/05/05.md)

[go-zero 集成:日志,Prometheus,jaeger,分布式事务dtm](./go-zero-md/04/04.md)
