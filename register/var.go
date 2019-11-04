package register

import (
	"crypto/sha1"
	"fmt"
	"net"
	"sync"
)

type Service struct {
	Name     string `service name`
	DNS      string `localhost:localhost`
	URL      string `127.0.0.1:80`
	Method   string `method`
	Note     string `desc`
	PassWord string
}

type ConfigFile struct {
	Count    int    `配置文件中有多少 配置文件`
	Size     int64  `配置文件有多大`
	Path     string `储存的具体位置`
	Services []*Service
}
type ConfigFiles struct {
	Count      int `文件下有都个配置文件 .wm`
	ConfigFile []*ConfigFile
}
type Ruler struct {
	IsDie   bool  `是否死亡`
	TimeOut int64 `服务超时`
	Status  int   `服务状态 -1 宕机 , 1-100 负载状态`
	Name    string
	Service *Service `运行的服务`
}
type Config struct {
	Address     string    `listen port : default 6000`
	MaxCap      int       `max connect`
	PollingType int       `polling type`
	Redis       *net.Conn `redis connect`
	File        string    `file path`
}

// 存放 service 容器
type Container struct {
	mux    sync.Mutex
	Rulers []*Ruler
	Number int
}

type Monitor struct {
	Mux  sync.Mutex
	Sort []int `直接排序`
}

func MD(s string) string {
	Sha1Inst := sha1.New()
	Sha1Inst.Write([]byte(s))
	Result := Sha1Inst.Sum([]byte(""))
	return fmt.Sprintf("%x", Result)
}
