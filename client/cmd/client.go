package main

import (
	"bytes"
	"downloader/pkg/myLog"
	"fmt"
	"io"
	"os"
	"strconv"
	"sync"
	"time"
)

type Client struct {
	Token string
	Wg    *sync.WaitGroup
}

func NewClient(token string) *Client {
	return &Client{
		Token: token,
		Wg:    new(sync.WaitGroup),
	}
}

func (c *Client) Download(opt DownloadOption) {
	var (
		offset int64
		part   int
		done   bool
		size   = c.GetSize(opt.SrcPath)
		line   = NewLine(opt.SrcPath, size, int(opt.Buffer), time.Now())
	)

	for {
		if done {
			break
		}
		file, err := os.Create(opt.DstPath + fmt.Sprintf(".part%d", part))
		if err != nil {
			myLog.Log.Error(err)
			return
		}
		p := Poster{
			Url: "http://localhost:8080/user/download/slice",
			Fields: map[string]string{
				"filename": opt.SrcPath,
				"offset":   strconv.FormatInt(offset, 10),
				"buf":      strconv.FormatInt(opt.Buffer, 10),
			},
			Header: map[string]string{
				"Token": c.Token,
			},
		}
		resp, e := p.Post()
		if e != nil {
			myLog.Log.Error(e)
			return
		}
		if resp.Header.Get("Done") == "OK" {
			done = true
		}
		offset += opt.Buffer
		part++
		c.Wg.Add(1)
		go func(writer io.Writer, reader io.Reader) {
			_, e := io.Copy(writer, reader)
			if e != nil {
				myLog.Log.Error(e)
				return
			}
			line <- struct{}{}
			c.Wg.Done()
		}(file, resp.Body)

	}
	c.Wg.Wait()
	c.Merge(opt.DstPath, part)
}

func (c *Client) Merge(filename string, parts int) {
	file, err := os.Create(filename)
	if err != nil {
		myLog.Log.Error(err)
		return
	}
	for i := 0; i < parts; i++ {
		partFilename := fmt.Sprintf("%s.part%d", filename, i)
		partFile, err := os.Open(partFilename)
		if err != nil {
			myLog.Log.Error(err)
			return
		}
		_, err = io.CopyBuffer(file, partFile, make([]byte, 1024*1024))
		if err != nil {
			myLog.Log.Error(err)
			return
		}
		err = partFile.Close()
		if err != nil {
			myLog.Log.Error(err)
			return
		}
		c.Wg.Add(1)
		go func(name string) {
			err := os.Remove(name)
			if err != nil {
				myLog.Log.Error(err)
				c.Wg.Done()
				return
			}
			c.Wg.Done()
		}(partFilename)
		c.Wg.Wait()
	}

}

func (c *Client) GetSize(filename string) int64 {
	p := Poster{
		Url: "http://localhost:8080/user/download/filesize",
		Fields: map[string]string{
			"filename": filename,
		},
		Header: map[string]string{
			"Token": c.Token,
		},
	}
	resp, _ := p.Post()
	var b = &bytes.Buffer{}
	_, _ = io.Copy(b, resp.Body)
	size := fmt.Sprint(string(b.Bytes()))
	parseInt, _ := strconv.ParseInt(size, 10, 64)
	return parseInt
}
