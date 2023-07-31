# 下载器

> - 基于gin框架
> - 并发分片下载


- quickstart

> 配置./etc/conf/yaml   
> go run ./main.go         
> apipost 注册登录获取token     
- 下载
```go

func main() {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJZCI6MjI2OTcxNjk3MTUyLCJUb2tlblR5cGUiOiJhY2Nlc3NfdG9rZW4iLCJleHAiOjE2OTA3ODE3MDEsImlhdCI6MTY5MDY5NTMwMX0.Pa27PbYh7O10yKE_mNkMu3s19CQ5ySUAwPlmatutir8"
	client := NewClient(token)
	client.Download(DownloadOption{
		SrcPath: "test.exe",
		DstPath: "./fs/sever/test.exe",
		Buffer:  1024,
	})
}
```

- 上传
> POST localhost:8080/user/upload

- 鉴权
> 采用`jwt`实现

- p2p
> `./client/p2p/main.go`

- 心跳检测
> `./client/p2p/main.go` `./router/p2p.go`

- 进度条
> `./client/cmd/load.go`

