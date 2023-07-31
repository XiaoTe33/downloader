package main

import (
	"context"
	"downloader/pkg/myLog"
	"io"
	"net"
	"time"
)

const heartbeat = "ping"
const interval = 5 * time.Second
const timeout = 10 * time.Second

type Client struct {
	conn net.Conn
	ch   chan []byte
	ctx  context.Context
}

func NewClient(conn net.Conn) *Client {
	return &Client{
		conn: conn,
		ch:   make(chan []byte),
		ctx:  context.Background(),
	}
}

func (c *Client) Start() {
	// 创建一个子上下文，用于取消心跳检测
	ctx, cancel := context.WithCancel(c.ctx)
	defer cancel()
	// 发送心跳包
	go c.sendHeartbeat(ctx)
	// 接收心跳响应
	go c.receiveHeartbeat(ctx)
	// 处理业务逻辑
	go c.handle(ctx)

	<-ctx.Done()
}

func (c *Client) sendHeartbeat(ctx context.Context) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			//myLog.Log.Info("send heartbeat")
			if _, err := c.conn.Write([]byte(heartbeat)); err != nil {
				myLog.Log.Error("send heartbeat error:", err)
				return
			}
		}
	}
}

func (c *Client) receiveHeartbeat(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			buf := make([]byte, len(heartbeat))
			if _, err := io.ReadFull(c.conn, buf); err != nil {
				myLog.Log.Error("receive heartbeat error:", err)
				return
			}
			c.ch <- buf
		}
	}
}

// 处理业务逻辑
func (c *Client) handle(ctx context.Context) {
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-c.ch:
			if string(msg) == heartbeat {
				timer.Reset(timeout)
			} else {

				// 处理其他业务逻辑
			}
		case <-timer.C:
			myLog.Log.Info("heartbeat timeout")
			return // 心跳超时
		}
	}
}

func main() {
	//注册到8888
	register, e := net.Dial("tcp", "127.0.0.1:8888")
	if e != nil {
		myLog.Log.Error(e)
		return
	}
	//维持心跳
	NewClient(register).Start()

	listener, err := net.Listen("tcp", ":9090")
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	for {
		// 接受客户端连接
		conn, err := listener.Accept()
		if err != nil {
			myLog.Log.Error("accept error:", err)
			continue
		}
		myLog.Log.Info("client connected:", conn.RemoteAddr())
		// 创建一个客户端对象
		client := NewClient(conn)
		// 启动客户端
		client.Start()
	}
}
