/**
  @author:panliang
  @data:2021/9/3
  @note
**/
package tcp

import (
	"go_im/pkg/config"
	"go_im/pkg/pool"
	"io"
	"log"
	"net"
	"os"
)

func StartTcpClient()  {
	tcpAddr, err := net.ResolveTCPAddr("tcp", ":"+config.GetString("app.tcp_port"))
	if err != nil {
		log.Fatal(err)
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)

	if err != nil {
		log.Fatal(err)
	}
	done := make(chan struct{})
	pool.AntsPool.Submit(func() {
		io.Copy(os.Stdout, conn) // 注意：忽略错误
		log.Println("done")
		done <- struct{}{} // 通知主 goroutine 的信号
	})
	mustCopy(conn, os.Stdin)
	conn.CloseWrite()
	<-done // 等待后台 goroutine 完成
}

func mustCopy(dst io.Writer,src io.Reader)  {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}
