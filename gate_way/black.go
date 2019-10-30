package gate_way

import (
	"bufio"
	"io"
	"log"
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
	LoadBlackPath   string `load black file path`
	MonitorInterval int
}

var (
	AccessDefineList    map[string]int // use it to monitor black list
	AccessDefineListNew map[string]int
	MonitoringPeriod    int         // second , default 5s
	MonitorInterval     int         // second , default 5s
	MonitorPipe         chan func() // yes ,use chan to monitor data
)

// init func ,find of local file ,name with *.black
func StartBlackShield(config *BlackShieldConfig) {
	util.CheckConfig(config, BlackShieldConfig{LoadBlackPath: "./black.lst"})
	loadList(config)
	go func() {
		for {

			time.Sleep(time.Second * time.Duration(config.MonitorInterval))
		}
	}()
}
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
func BlackShieldAlg(ip string) {
	if ips := AccessDefineList[ip]; ips > 8 {

	}
}
