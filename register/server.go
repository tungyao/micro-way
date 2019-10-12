package register

import (
	"log"
	"net"
)

// server default port 6000
const (
	REDIS  = iota //通过redis注册服务
	FILE          //通过本地文件注册服务
	MEMORY        //通过服务端内存注册服务
)

type Ruler struct {
	IsDie   bool    `是否死亡`
	TimeOut int64   `服务超时`
	Status  int     `服务状态 -1 宕机 , 1-100 负载状态`
	Service Service `运行的服务`
}
type Config struct {
	Address     string    `listen port : default 6000`
	MaxCap      int       `max connect`
	PollingType int       `polling type`
	Redis       *net.Conn `redis connect`
	File        string    `file path`
}

func StartServer(config Config) {
	listen, err := net.Listen("tcp", config.Address)
	if err != nil {
		log.Println(err)
	}
	pool := NewPool(config.MaxCap)
	go func() {
		for {
			con, err := listen.Accept()
			if err != nil {
				log.Println(err)
			}
			pool.EntryChannel <- NewTask(func() error {
				data := make([]byte, 128)
				n, err := con.Read(data)
				log.Println(string(data[:n]))
				_, err = con.Write(data[:n])
				if err != nil {

				}
				return nil
			})
		}
	}()
	pool.Run()
}

func runFile() {

}
func Memory() {

}
