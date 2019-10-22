package register

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// Loop Time 30s
// 轮询器 自动开启轮询 和 手动开启轮询
// 同时还将具备收集 服务状态  自动报告 和 收集 时延

// TODO timer from time package
func StartPolling(t time.Duration, config Config) {
	ti := time.NewTicker(time.Second * t)
	for v := range ti.C {
		go func() {
			log.Println("waiting time =>", t, v)
			CheckConfigFile(config)
		}()
	}
}

// TODO Check File Different
func CheckConfigFile(config Config) {
	err := filepath.Walk(config.File, func(path string, info os.FileInfo, err error) error {
		if len(path) > 3 && path[len(path)-3:] == ".wm" {
			d, err := ioutil.ReadFile(path)
			if err != nil {
				log.Println(err)
				return err
			}
			dd, err := ioutil.ReadFile(path + ".run")
			if err != nil {
				log.Println(err)
				return err
			}
			for k, v := range dd {
				dd[k] = v >> 1
			}
			if !Diff(d, dd) {
				log.Println("文件不同")
				pcf := ParseConfigFile(path)
				for _, v := range pcf {
					LoadSingleService(v.Name, Ruler{
						IsDie:   true,
						TimeOut: 60,
						Status:  -1,
						Name:    v.Name,
						Service: v,
					})
				}
				for k, v := range d {
					d[k] = v << 1
				}
				_ = ioutil.WriteFile(path+".run", d, 766)
			}
		}
		return nil
	})
	if err != nil {
		log.Println(err)
	}

}

// TODO DIFF alg
func Diff(a []byte, b []byte) bool { // 直接对比 每一行
	var og = SplitString(a, []uint8{13, 10}) // 源文件
	var dt = SplitString(b, []uint8{13, 10}) // 新文件
	if len(og) > len(dt) {
		out := make([]int, 0)
		for k, v := range og {
			if len(dt) > k && string(v) != string(dt[k]) {
				out = append(out, 0)
			} else if len(dt) < k+1 {
				out = append(out, 0)
			}
		}
		if len(out) != 0 {
			println(string(a))
			return false
		}
	} else if len(og) == len(dt) {
		out := make([]int, 0)
		for k, v := range og {
			if string(v) != string(dt[k]) {
				out = append(out, 0)
				continue
			}
		}
		if len(out) != 0 {
			return false
		}
	} else if len(og) < len(dt) {
		out := make([]int, 0)
		for k, v := range dt {
			if len(og) > k && string(v) != string(og[k]) {
				out = append(out, 0)
			} else if len(og) < k+1 {
				out = append(out, 0)
			}
		}
		if len(out) != 0 {
			return false
		}
	}
	return true
}
func OutDisplayOrFile(n [][]byte) { // 将这个玩意打印出来
	for k, v := range n {
		if v[len(v)-1:][0] == '*' {
			fmt.Println(k, "*", string(v[:len(v)-1]))
		} else if v[len(v)-1:][0] == '+' {
			fmt.Println(k, "+", string(v[:len(v)-1]))
		} else if v[len(v)-1:][0] == '=' {
			fmt.Println(k, "=", string(v[:len(v)-1]))
		} else if v[len(v)-1:][0] == '-' {
			fmt.Println(k, "-", string(v[:len(v)-1]))
		}
	}
}

// 监控服务 状态 应该是个自动化服务
func StartMonitorService() {
	stats := &runtime.MemStats{}
	runtime.ReadMemStats(stats)
	// log.Println(stats.Sys)
}

//  return number more small ,service is more health
func URLMonitor() int {
	for k, v := range GlobalContainer.Rulers {

	}
	return 0
}
func sendReq(method string, url string, body io.Reader) {
	req, err := http.NewRequest(strings.ToUpper(method), url, body)
	if err != nil {
		return nil, err
	}
}
