package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var (
	use    string
	path   string
	create bool
	start  bool
	stop   bool
)

func init() {
	fmt.Println("We strongly recommend using the configuration file is configured ,more tips -h")
	flag.StringVar(&use, "use", "file", "help message for flag name")
	flag.StringVar(&path, "path", "./config.wm", "config file path ,default current folder")
	flag.BoolVar(&create, "create", false, "create new config file")
	flag.BoolVar(&start, "start", false, "start gate-way")
	flag.BoolVar(&stop, "stop", false, "stop gate-way")
}
func main() {
	flag.Parse()
	if create {
		// if err != nil {
		// 	fmt.Println(26,err)
		// }
		f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 777)
		if err != nil {
			fmt.Println(30, err)
			f.Close()
			os.Exit(0)
		}
		_, err = f.Write([]byte("{\r\nname=wayconfig\r\nmaxcon=10\r\nmaxbufflow=10\r\n}\r\n{\r\nname=blackconfig\r\nLoadBlackPath=\"\"\r\nMonitorInterval=0\r\nMonitoringPeriod=0\r\nMonitorPipebuf=0\r\nMonitorPeriosBuf=0\r\n}\r\n{\r\nRegisterLocation=\"\"\r\nTimeOut=0\r\nIsCache=false\r\nCacheTime=0\r\nMaxCap=0\r\n}"))
		err = f.Close()
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}

	}
	if start {
		fmt.Println("start")
		ParseConfigFile(path)
	}
	// gate_way.StartGateWay(&gate_way.Config{
	// 	MaxConn:     0,
	// 	MaxBuffFlow: 0,
	// },&gate_way.BlackShieldConfig{
	// 	LoadBlackPath:    "",
	// 	MonitorInterval:  0,
	// 	MonitoringPeriod: 0,
	// 	MonitorPipebuf:   0,
	// 	MonitorPeriosBuf: 0,
	// },&gate_way.WayConfig{
	// 	RegisterLocation: "",
	// 	TimeOut:          0,
	// 	IsCache:          false,
	// 	CacheTime:        0,
	// 	MaxCap:           0,
	// })
	// gate_way.StartRouter(&gate_way.WayConfig{
	// 	RegisterLocation: "",
	// 	TimeOut:          0,
	// 	IsCache:          false,
	// 	CacheTime:        0,
	// 	MaxCap:           0,
	// })
}
func ParseConfigFile(path string) {
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
	isgroup := false
	str := make([]byte, 0)
	for i := 0; i < len(get); i++ {
		if get[i] == 32 {
			continue
		}
		if get[i] == 123 {
			isgroup = true
		}
		if get[i] == 125 {
			isgroup = false
		}
		if isgroup && get[i] != 123 {
			str = append(str, get[i])
		}
	}
	group := SplitString(str, []uint8{13, 10})
	fmt.Println(len(group))
	// 到这一部可以开始解析数据到出来
	for _, v := range group {
		column := SplitString(v, []uint8{13, 10})
		fmt.Println(column)
	}
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
func SplitString(str []byte, p []byte) [][]byte {
	group := make([][]byte, 0)
	for i := 0; i < len(str); i++ {
		if str[i] == p[0] && i < len(str)-len(p) {
			if len(p) == 1 {
				return [][]byte{str[:i+1], str[i:]}
			} else {
				for j := 1; j < len(p); i++ {
					if str[i+j] != p[j] {
						continue
					}
					return [][]byte{str[:i], str[i+len(p):]}
				}
			}
		} else {
			continue
		}
	}
	return group
}
