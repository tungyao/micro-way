package register

import (
	"log"
	"net"
)

// 手动注册器 通过命令行 以及监控控制
func StartClient(str string) {
	listen, _ := net.Dial("tcp", ":6000")
	listen.Write([]byte(str))
	data := make([]byte, 128)
	n, _ := listen.Read(data)
	log.Println(string(data[:n]))
	listen.Close()
}
