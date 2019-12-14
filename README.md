## 文件管理接口

## go 安装包时需要配置的代理

```shell script
git config --global http.proxy socks5://127.0.0.1:1080
export http_proxy=socks5://127.0.0.1:1080
```

## 安装依赖并启动

```shell script
go mod tidy
go mod download
go run main.go
```

## 参考文章

- [Example to handle GET and POST request in Golang](https://www.golangprograms.com/example-to-handle-get-and-post-request-in-golang.html)
- [golang 设置 http response 响应头与坑](https://www.jianshu.com/p/4a26a4681464)
- [Tip for "easy" logging errors](https://stackoverflow.com/questions/43976140/check-errors-when-calling-http-responsewriter-write)
- [golang JWT 包生成 Token, 验证 Token](https://hacpai.com/article/1540349739379)