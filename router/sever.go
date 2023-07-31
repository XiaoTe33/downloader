package router

import (
	"downloader/pkg/errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"os"
	"strconv"
)

// upload 上传文件到网盘
func upload(c *gin.Context) {
	username := c.GetString("username")
	fileHeader, err := c.FormFile("file")

	if err != nil {
		jsonError(c, errors.ErrUploadFiles)
		return
	}
	dir := "./fs/user/" + username
	filename := fileHeader.Filename
	if handleError(c, c.SaveUploadedFile(fileHeader, dir+"/"+filename)) {
		return
	}
	jsonSuccess(c)
}

func merge(c *gin.Context) {
	filename := c.PostForm("filename")
	part := c.PostForm("part")
	partNum, err := strconv.Atoi(part)
	if handleError(c, err) {
		return
	}
	file, err := os.Create(filename)
	if handleError(c, err) {
		return
	}
	for i := 0; i < partNum; i++ {
		partialFileName := fmt.Sprintf("%s.part%d", filename, i)
		partialFile, err := os.Open(partialFileName)
		if handleError(c, err) {
			return
		}
		buf := make([]byte, 1024*1024)
		_, err = io.CopyBuffer(file, partialFile, buf)
		if handleError(c, err) {
			return
		}
	}
}

func downloadSlice(c *gin.Context) {
	f, e := os.Open(getFilepath(c))
	if handleError(c, e) {
		return
	}

	offsetStr := c.PostForm("offset")
	bufStr := c.PostForm("buf")
	offset, e := strconv.ParseInt(offsetStr, 10, 64)
	if handleError(c, e) {
		return
	}
	bufSize, e := strconv.ParseInt(bufStr, 10, 64)
	if handleError(c, e) {
		return
	}

	buf := make([]byte, bufSize)
	_, e = f.ReadAt(buf, offset)

	if e == io.EOF {
		c.Header("Done", "OK")
	}
	_, _ = c.Writer.Write(buf)

}

func fileSize(c *gin.Context) {
	f, e := os.Open(getFilepath(c))
	if handleError(c, e) {
		return
	}
	stat, e := f.Stat()
	if handleError(c, e) {
		return
	}
	_, e = c.Writer.Write([]byte(strconv.FormatInt(stat.Size(), 10)))
	if handleError(c, e) {
		return
	}

}
