package gate_way

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"./util"
)

// Used to monitor malicious requests and block
// Reference paper => https://wenku.baidu.com/view/3076f3d7b14e852458fb57ad.html

// Every ip access must go through BLACK SHIELD, preventing the attack from being triggered immediately,
// but delaying the trigger.
// you can put black shield to http handler
// access define list struct
type ADL struct {
	IP        string
	Frequency int
}
type BlackShieldConfig struct {
	LoadBlackPath    string `load black file path`
	MonitorInterval  int
	MonitoringPeriod int
	MonitorPipebuf   int `default 2048`
	MonitorPeriosBuf int `default 10240 ,`
}
type Black struct {
	net.Listener
	l net.Listener
}

var (
	AccessDefineList     map[string]int // use it to monitor black list
	AccessDefineListNew  map[string]int
	MonitoringPeriod     int         // second , default 5s
	MonitorInterval      int         // second , default 3s
	MonitorPipe          chan func() // yes ,use chan to monitor data
	MonitorLastTime      int
	MonitoringPeriodTime int
	Break                bool
	emptySlice           []string
	MonitorPeriosBuf     int
)

// init func ,find of local file ,name with *.black
func StartBlackShield(config *BlackShieldConfig, l net.Listener) *Black {
	util.CheckConfig(config, BlackShieldConfig{MonitorPeriosBuf: 10240, LoadBlackPath: "./black.lst", MonitorPipebuf: 2048, MonitoringPeriod: 5, MonitorInterval: 5})
	loadList(config)
	MonitorPipe = make(chan func(), config.MonitorPipebuf)
	MonitoringPeriod = config.MonitoringPeriod
	MonitorInterval = config.MonitorInterval
	MonitorPeriosBuf = config.MonitorPeriosBuf

	go func() {
		for k := range MonitorPipe {
			k()
		}
	}()
	// During the independent control of threads to keep reading
	go func() {
		for {
			Break = false
			BlackShieldAlg(emptySlice)
			time.Sleep(time.Second * time.Duration(MonitorInterval))
			Break = true
			emptySlice = make([]string, 0)
			time.Sleep(time.Second * time.Duration(MonitoringPeriod))
		}
	}()
	return &Black{
		l: l,
	}
}
func (b *Black) Next() net.Listener {
	return b
}
func (b *Black) Accept() (net.Conn, error) {
	a, err := b.l.Accept()
	if Break {
		MonitorPipe <- func() {
			SendFlow(a.RemoteAddr().String())
		}
	}
	if len(emptySlice) > 200 {
		if AccessDefineList[string(util.SplitString([]byte(a.RemoteAddr().String()), []byte(":"))[0])] != 0 {
			err = a.Close()
			return nil, err
		}
	}
	return a, err
}

// read black list from local file
func loadList(config *BlackShieldConfig) {
	f, err := os.OpenFile(config.LoadBlackPath, os.O_CREATE|os.O_RDONLY, 777)
	if err != nil {
		log.Panicln(err)
	}
	buf := bufio.NewReader(f)
	AccessDefineList = make(map[string]int, 0)
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		g := strings.Split(line, " ")
		in, _ := strconv.Atoi(string(g[1]))
		AccessDefineList[string(g[0])] = in
		if err != nil {
			if err == io.EOF {
				break
			}
		}
	}
}

// core algorithm
func BlackShieldAlg(s []string) {
	fmt.Println("\n\n\n\n\n")
	fmt.Print("\t\tBLACK LIST IS RUNNING\n------------------------------------\n")
	for k, _ := range AccessDefineList {
		fmt.Printf("|\t\t%s\t\t\t\t\t\t\n", k)
	}
	fmt.Println("------------------------------------")
	buf := map[string]int{}
	for _, v := range emptySlice {
		buf[v] = buf[v] + 1
		if buf[v] > 100 {
			AccessDefineList[v] = AccessDefineList[v] + 1
		}
	}

	// for k,v:=range s{
	// 	fmt.Println(k,v)
	// }
}

// send ip address to algorithm and get threat index
// we should create a new slice every time ,until the free
func SendFlow(ip string) {
	ip = string(util.SplitString([]byte(ip), []byte(":"))[0])
	if len(emptySlice) <= MonitorPeriosBuf {
		emptySlice = append(emptySlice, ip)
	}
}
