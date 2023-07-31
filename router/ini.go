package router

import (
	"downloader/pkg/dao"
	"github.com/gin-gonic/gin"
	"time"
)

var (
	db = dao.DB
)

func InitRouters() {
	r := gin.Default()
	r.Use(Cors())
	limiter := NewRateLimiter(10, time.Second)
	r.POST("/user/download/filesize", JWT(), fileSize)
	r.POST("/user/download/slice", JWT(), limiter.Limit, downloadSlice)
	r.POST("/user/upload", JWT(), limiter.Limit, upload)
	r.POST("/user/register", register)
	r.POST("/user/login", login)
	r.GET("/user/token/refresh", refreshToken)
	_ = r.Run()
}
