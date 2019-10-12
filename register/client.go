package register

import (
	"log"
	"net"
)

func StartClient() {
	listen, _ := net.Dial("tcp", ":6000")
	listen.Write([]byte("set service ababsdbasd"))
	data := make([]byte, 128)
	n, _ := listen.Read(data)
	log.Println(string(data[:n]))
	listen.Close()
}
