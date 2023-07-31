package main

import (
	"fmt"
)

var setFilename = "1.mp4"
var dstFilename = "./fs/sever/test.mp4"

func main04() {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJZCI6MjI2OTcxNjk3MTUyLCJUb2tlblR5cGUiOiJhY2Nlc3NfdG9rZW4iLCJleHAiOjE2OTA3ODE3MDEsImlhdCI6MTY5MDY5NTMwMX0.Pa27PbYh7O10yKE_mNkMu3s19CQ5ySUAwPlmatutir8"
	client := NewClient(token)
	fmt.Println(client.GetSize("test.exe"))
}
func main() {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJZCI6MjI2OTcxNjk3MTUyLCJUb2tlblR5cGUiOiJhY2Nlc3NfdG9rZW4iLCJleHAiOjE2OTA3ODE3MDEsImlhdCI6MTY5MDY5NTMwMX0.Pa27PbYh7O10yKE_mNkMu3s19CQ5ySUAwPlmatutir8"
	client := NewClient(token)
	client.Download(DownloadOption{
		SrcPath: "test.exe",
		DstPath: "./fs/sever/test.exe",
		Buffer:  1024,
	})
}
func main02() {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJZCI6MjI2OTcxNjk3MTUyLCJUb2tlblR5cGUiOiJhY2Nlc3NfdG9rZW4iLCJleHAiOjE2OTA3ODE3MDEsImlhdCI6MTY5MDY5NTMwMX0.Pa27PbYh7O10yKE_mNkMu3s19CQ5ySUAwPlmatutir8"
	client := NewClient(token)
	client.Download(DownloadOption{
		SrcPath: "1.mp4",
		DstPath: "./fs/sever/test.mp4",
		Buffer:  1024 * 1024 * 5,
	})
}

type DownloadOption struct {
	SrcPath string
	DstPath string
	Buffer  int64
}
