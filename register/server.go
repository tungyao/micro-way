package register

import (
	"io/ioutil"
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

func StartServer(config Config) {
	listen, err := net.Listen("tcp", config.Address)
	if err != nil {
		log.Println(err)
	}
	// 检查配置config的错误
	checkParameter(&config)
	containerMap := make(map[string]*Ruler, 0)
	// 根据config类型加载配置文件
	switch config.PollingType {
	case FILE:
		containerMap = createContainer(getAllConfigFile(config.File), containerMap) // 获取配置文件
	case REDIS:
		log.Panic("next time")
	case MEMORY:
		log.Panic("next time")
	}
	// 到这里加载配置文件
	// 注册中心提供服务注册和服务发现功能
	// 注册中心解决单点故障问题
	// 注册中心需要保存服务注册信息以及服务发现时的筛选和简单计算能力

	// 加载配置文件
	LoadGlobalService(containerMap)
	// 开启线程池
	pool := NewPool(config.MaxCap)
	go func() {
		for {
			con, err := listen.Accept()
			if err != nil {
				log.Println(err)
			}
			pool.EntryChannel <- NewTask(func() error {
				data := make([]byte, 2048) // 默认读取大小为2kb
				nc := con
				n, _ := nc.Read(data)
				if n == 0 {
					_ = nc.Close()
					return nil
				}
				_, _ = nc.Write(data)
				// LoadSingleService("hello", Ruler{
				// 	IsDie:   true,
				// 	TimeOut: 100,
				// 	Status:  -1,
				// 	Name:    "hello",
				// 	Service: &Service{
				// 		Name: "hello",
				// 		DNS:  "www.yaop.ink/hello",
				// 		URL:  "127.0.0.1:4000",
				// 		Note: "this is hello",
				// 	},
				// })
				// a, v := GetStatusSingleService("wechat")
				// fmt.Println(a, v)
				// _, err = con.Write([]byte(containerMap[string(data[:n])].Service.DNS))
				// if err != nil {
				// 	log.Println(err)
				// }
				// _ = nc.Close()
				return nil
			})
		}
	}()
	// 开启文件轮询模式
	CheckConfigFile(config)
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
			log.Panic("Server polling type is FILE , but FILE path is empty")
		}
	}
}

// 获取配置文件夹下的配置文件 并且创建 运行文件
func getAllConfigFile(path string) *ConfigFiles {
	files := new(ConfigFiles)
	err := filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if info == nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if info.Name()[len(info.Name())-3:] == ".wm" {
			f, err := os.Open(p)
			if err != nil {
				log.Panicln(err)
			}
			data, _ := ioutil.ReadAll(f)
			for k, v := range data {
				data[k] = v << 1
			}
			_ = ioutil.WriteFile(p+".run", data, 766)
			_ = f.Close()
			files.Count = files.Count + 1
			pcf := ParseConfigFile(p)
			files.ConfigFile = append(files.ConfigFile, &ConfigFile{
				Count:    len(pcf),
				Size:     info.Size(),
				Path:     p,
				Services: pcf,
			})
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	if files != nil {
		return files
	}
	return nil
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
			url := FindString(column[i], []byte("URL="))
			if url != nil {
				ser.URL = string(url.([]byte))
			}
			dns := FindString(column[i], []byte("DNS="))
			if dns != nil {
				ser.DNS = string(dns.([]byte))
			}
			note := FindString(column[i], []byte("Note="))
			if note != nil {
				ser.Note = string(note.([]byte))
			}
			method := FindString(column[i], []byte("Method="))
			if method != nil {
				ser.Method = string(note.([]byte))
			} else {
				ser.Method = "GET"
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
	// switch v.(type) {
	// case []byte:
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
	// case string:
	// 	sr := v.(string)
	// }
	// return nil
}

// 建立一个大型容器来运行，增加，删除容器
func createContainer(configs *ConfigFiles, mp map[string]*Ruler) map[string]*Ruler {
	for i := 0; i < len(configs.ConfigFile); i++ {
		for _, v := range configs.ConfigFile[i].Services {
			mp[v.Name] = &Ruler{
				IsDie:   true,
				TimeOut: 100,
				Name:    v.Name,
				Status:  -1,
				Service: v,
			}
		}
	}
	return mp
}
