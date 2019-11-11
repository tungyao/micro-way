package register

import (
	"fmt"
	"net"
)

// 手动注册器 通过命令行 以及监控控制
func main() {
	var operation string
	var key string
	for {
		listen, _ := net.Dial("tcp", ":6000")
		fmt.Print("127.0.0.1:6000>> ")
		_, _ = fmt.Scanln(&operation, &key)
		if operation == "exit" {
			break
		}
		_, _ = listen.Write([]byte("====" + operation + "*" + key))
		data := make([]byte, 128)
		n, _ := listen.Read(data)
		fmt.Println(string(data[:n]))
		_ = listen.Close()
	}

}
