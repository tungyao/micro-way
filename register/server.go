package register

import (
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
)

// server default port 6000
const (
	MEMORY = iota // 通过内存注册服务 如果使用内存注册,每次开机则需要手动重新注册,所以默认不采用该方法,采用FILE方式
	FILE          // 通过本地文件注册服务 , file 文件 可以通过bin文件的registerFile.go 快速得到
	REDIS         // 通过redis注册服务
)

type ConfigFile struct {
	Count    int    `配置文件中有多少 配置文件`
	Size     int64  `配置文件有多大`
	Path     string `储存的具体位置`
	Services []*Service
}
type ConfigFiles struct {
	Count      int `文件下有都个配置文件 .wm`
	ConfigFile []ConfigFile
}
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
	checkParameter(&config)
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

// TODO tool is next
func checkParameter(config *Config) {
	if config.MaxCap == 0 {
		config.MaxCap = 4000
	}
	if config.Address == "" {
		config.Address = ":6000"
	}
	if config.PollingType == 0 {
		config.PollingType = FILE
	}
	if config.PollingType == REDIS {
		if config.Redis == nil {
			log.Panic("Server polling type is redis , but redis config is nil")
		}
	}
	if config.PollingType == FILE {
		if config.File == "" {
			log.Panic("Server polling type is FILE , but FILE path is nil")
		}
		getAllConfigFile(config.File)
	}
}

// 获取配置文件夹下的配置文件
func getAllConfigFile(path string) *ConfigFiles {
	files := new(ConfigFiles)
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if info == nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if info.Name()[len(info.Name())-3:] == ".wm" {
			files.Count = files.Count + 1
			pcf := ParseConfigFile(path)
			files.ConfigFile = append(files.ConfigFile, ConfigFile{
				Count:    len(pcf),
				Size:     info.Size(),
				Path:     path,
				Services: pcf,
			})
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(files)
	return files
}
func ParseConfigFile(path string) []*Service {
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		log.Println("open config file error", err)
	}
	stat, _ := f.Stat()
	get := make([]byte, stat.Size())
	_, err = f.Read(get)
	if err != nil {
		log.Panic(err)
	}
	str := make([]byte, 0)
	isGroup := false
	for i := 0; i < len(get); i++ {
		if get[i] == 0x7b {
			isGroup = true
		}
		if get[i] == 0x7d {
			isGroup = false
		}
		if isGroup && get[i] != 0x7b && get[i] != 0x20 {
			str = append(str, get[i])
		}
	}
	group := SplitString(str, []uint8{13, 10, 13, 10})
	// 到这一部可以开始解析数据到出来
	service := make([]*Service, 0)
	for _, v := range group {
		column := SplitString(v, []uint8{13, 10})
		ser := new(Service)
		for i := 1; i < len(column); i++ {
			name := FindString(column[i], []byte("Name="))
			if name != nil {
				ser.Name = string(name.([]byte))
			}
			port := FindString(column[i], []byte("Port="))
			if port != nil {
				ser.Port = string(port.([]byte))
			}
			dns := FindString(column[i], []byte("DNS="))
			if dns != nil {
				ser.DNS = string(dns.([]byte))
			}
			note := FindString(column[i], []byte("Note="))
			if note != nil {
				ser.Note = string(note.([]byte))
			}
		}
		service = append(service, ser)
	}
	return service
}
func SplitString(str []byte, p []byte) [][]byte {
	group := make([][]byte, 0)
	ps := make([]int, 0)
	for i := 0; i < len(str); i++ {
		ist := make([]int, len(p))
		for k, v := range p {
			if i < len(str)-len(p) && str[i+k] == v {
				ist[k] = 1
			}
		}
		st := true
		for _, v := range ist {
			if v != 1 {
				st = false
			}
		}
		if st {
			ps = append(ps, i)
		}
	}
	ps = append(ps, len(str))
	sto := 0
	for i := 0; i < len(ps); i++ {
		group = append(group, str[sto:ps[i]])
		sto = ps[i] + 2
	}
	return group
}
func FindString(v interface{}, p []byte) interface{} {
	switch v.(type) {
	case []byte:
		bt := v.([]byte)
		for i := 0; i < len(bt); i++ {
			ist := make([]int, len(p))
			for k, v := range p {
				if i < len(bt)-len(p) && bt[i+k] == v {
					ist[k] = 1
				}
			}
			st := true
			for _, v := range ist {
				if v != 1 {
					st = false
				}
			}
			if st {
				return bt[i+len(p):]
			}
		}
		return nil
	case string:
		// sr := v.(string)
	}
	return nil
}

// 更安全的分配内存
func malloc(v interface{}) {

}
