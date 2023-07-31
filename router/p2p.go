package router

import (
	"bytes"
	"downloader/pkg/myLog"
	"io"
	"net"
	"sync"
	"time"
)

const heartbeat = "ping"

var (
	cli map[string]time.Time
	mu  sync.Mutex
)

func Accept() {
	listen, err := net.Listen("tcp", ":8888")
	if err != nil {
		myLog.Log.Error(err)
		return
	}
	for {
		conn, err := listen.Accept()
		if err != nil {
			myLog.Log.Error(err)
			return
		}
		go KeepAlive(conn)
	}
}

func KeepAlive(conn net.Conn) {
	for {
		buf := &bytes.Buffer{}
		_, err := io.Copy(buf, conn)
		if err != nil {
			myLog.Log.Error(err)
			return
		}
		if buf.String() == heartbeat {
			mu.Lock()
			cli[conn.RemoteAddr().String()] = time.Now()
			mu.Unlock()
		} else {
			//传输数据
		}
	}
}
