package main

import (
	"bytes"
	"downloader/pkg/myLog"
	"fmt"
	"io"
	"os"
	"strconv"
)

var setFilename = "1.mp4"
var dstFilename = "./fs/sever/test.mp4"

func main() {
	size, err := getSize(setFilename)
	if err != nil {
		myLog.Log.Error(err)
		return
	}
	i, err := strconv.ParseInt(size, 10, 64)
	if err != nil {
		myLog.Log.Error(err)
		return
	}
	downloadSlice(setFilename, 0, i-i+1024*1024*10)

}

func getSize(filename string) (string, error) {
	p := Poster{
		Url: "http://localhost:8080/user/download/filesize",
		Fields: map[string]string{
			"filename": filename,
		},
		Header: map[string]string{
			"Token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJZCI6MjI2OTcxNjk3MTUyLCJUb2tlblR5cGUiOiJhY2Nlc3NfdG9rZW4iLCJleHAiOjE2OTA3ODE3MDEsImlhdCI6MTY5MDY5NTMwMX0.Pa27PbYh7O10yKE_mNkMu3s19CQ5ySUAwPlmatutir8",
		},
	}
	resp, e := p.Post()
	var b = &bytes.Buffer{}
	_, e = io.Copy(b, resp.Body)
	return fmt.Sprint(string(b.Bytes())), e
}

func downloadSlice(filename string, offset, buf int64) {
	file, e := os.Create(dstFilename)
	if e != nil {
		myLog.Log.Error(e)
		return
	}
	for {
		p := Poster{
			Url: "http://localhost:8080/user/download/slice",
			Fields: map[string]string{
				"filename": filename,
				"offset":   strconv.FormatInt(offset, 10),
				"buf":      strconv.FormatInt(buf, 10),
			},
			Header: map[string]string{
				"Token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJZCI6MjI2OTcxNjk3MTUyLCJUb2tlblR5cGUiOiJhY2Nlc3NfdG9rZW4iLCJleHAiOjE2OTA3ODE3MDEsImlhdCI6MTY5MDY5NTMwMX0.Pa27PbYh7O10yKE_mNkMu3s19CQ5ySUAwPlmatutir8",
			},
		}
		resp, e := p.Post()
		if e != nil {
			myLog.Log.Error(e)
			return
		}
		buf := &bytes.Buffer{}
		_, e = io.Copy(buf, resp.Body)
		if buf.Len() == 0 {
			break
		}

		offset += int64(buf.Len())
		_, e = io.Copy(file, buf)
		if e != nil {
			myLog.Log.Error(e)
			return
		}
		if resp.StatusCode == 201 {
			break
		}
	}

}
