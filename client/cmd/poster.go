package main

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

type Poster struct {
	Url    string
	Fields map[string]string
	Files  map[string]string
	Header map[string]string
}

func (r Poster) Post() (*http.Response, error) {
	var buff bytes.Buffer

	writer := multipart.NewWriter(&buff)

	for k, v := range r.Fields {
		_ = writer.WriteField(k, v)
	}
	for k, v := range r.Files {
		w, err := writer.CreateFormFile(k, v)
		file, err := os.Open(v)
		if err != nil {
			return nil, errors.New("打开文件失败")
		}
		_, _ = io.Copy(w, file)

	}
	_ = writer.Close()
	req, err := http.NewRequest(http.MethodPost, r.Url, &buff)
	if err != nil {
		return nil, errors.New("创建请求失败")
	}
	for k, v := range r.Header {
		req.Header.Set(k, v)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	var client = &http.Client{}
	resp, err := client.Do(req)
	return resp, err
}
