package router

import (
	"downloader/pkg/dao"
	"github.com/gin-gonic/gin"
)

var (
	db = dao.DB
)

func InitRouters() {
	r := gin.Default()
	r.Use(Cors())

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"hello": "world",
		})
	})
	r.POST("/user/download/filesize", JWT(), fileSize)
	r.POST("/user/download/slice", JWT(), downloadSlice)
	r.POST("/user/upload", JWT(), upload)
	r.POST("/user/register", register)
	r.POST("/user/login", login)
	r.GET("/user/token/refresh", refreshToken)
	_ = r.Run()
}
