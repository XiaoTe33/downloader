package router

import (
	"downloader/pkg/model"
	"downloader/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"sync"
	"time"

	"net/http"
	"strings"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken := c.Request.Header.Get("Token")
		id, err := utils.IsAccessToken(accessToken)
		if handleError(c, err) {
			c.Abort()
			return
		}
		c.Set("id", id)
		u := model.User{}
		if handleError(c, db.Where("id = ?", id).Find(&u).Error) {
			return
		}
		c.Set("username", u.Username)
		c.Next()
	}
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method               //请求方法
		origin := c.Request.Header.Get("Origin") //请求头部
		var headerKeys []string                  // 声明请求头keys
		for k := range c.Request.Header {
			headerKeys = append(headerKeys, k)
		}
		headerStr := strings.Join(headerKeys, ", ")
		if headerStr != "" {
			headerStr = fmt.Sprintf("access-control-allow-origin, access-control-allow-headers, %s", headerStr)
		} else {
			headerStr = "access-control-allow-origin, access-control-allow-headers"
		}
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Origin", "*")                                       // 允许访问所有域
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE") //服务器支持的所有跨域请求的方法,为了避免浏览次请求的多次'预检'请求
			//  header的类型
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")
			//              允许跨域设置                                                                                                      可以返回其他子段
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar") // 跨域关键设置 让浏览器可以解析
			c.Header("Access-Control-Max-Age", "172800")                                                                                                                                                           // 缓存请求信息 单位为秒
			c.Header("Access-Control-Allow-Credentials", "false")                                                                                                                                                  //  跨域请求是否需要带cookie信息 默认设置为true
			c.Set("content-type", "application/json")                                                                                                                                                              // 设置返回格式是json
		}

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
		}
		// 处理请求
		c.Next() //  处理请求
	}
}

type RateLimiter struct {
	mu       sync.Mutex           // 互斥锁，保证并发安全
	limit    int                  // 限制次数
	interval time.Duration        // 限制时间间隔
	buckets  map[string]time.Time // 存储请求时间戳的哈希表
}

// 创建一个新的速率限制器
func NewRateLimiter(limit int, interval time.Duration) *RateLimiter {
	return &RateLimiter{
		limit:    limit,
		interval: interval,
		buckets:  make(map[string]time.Time),
	}
}

// 实现速率限制器的逻辑
func (r *RateLimiter) Limit(c *gin.Context) {
	// 获取请求者的IP地址
	ip := c.ClientIP()
	// 上锁
	r.mu.Lock()
	defer r.mu.Unlock()
	// 获取当前时间
	now := time.Now()
	// 如果IP地址不存在于哈希表中，或者上次请求时间已经超过限制时间间隔，则更新哈希表
	if _, ok := r.buckets[ip]; !ok || now.Sub(r.buckets[ip]) > r.interval {
		r.buckets[ip] = now
		c.Next() // 继续处理请求
		return
	}
	// 如果IP地址存在于哈希表中，且上次请求时间没有超过限制时间间隔，则判断请求次数是否超过限制
	count := 0 // 记录请求次数
	for _, t := range r.buckets {
		// 如果请求时间在限制时间间隔内，则计数加一
		if now.Sub(t) <= r.interval {
			count++
		}
	}
	// 如果请求次数超过限制，则返回429状态码，表示太多请求
	if count >= r.limit {
		c.AbortWithStatus(429) // 中止处理请求
		return
	}
	c.Next() // 继续处理请求
}
