package znet

import (
	"fmt"
	"net"
	"testing"
	"time"
)

func ClientTest() {
	time.Sleep(time.Second * 3)

	conn, _ := net.Dial("tcp", "127.0.0.1:7777")

	for {
		_, _ = conn.Write([]byte("hello zinx"))

		buf := make([]byte, 512)
		_, _ = conn.Read(buf)

		fmt.Printf("server call back:%s", buf)

		time.Sleep(time.Second)
	}

}

func TestServer(t *testing.T) {
	s := NewServer()
	go ClientTest()
	s.Serve()
}
