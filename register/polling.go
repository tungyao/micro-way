package register

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
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
			log.Println(v)
			StatusConfigFile(config)
		}()
	}
}

// TODO Check File Different
func StatusConfigFile(config Config) {
	err := filepath.Walk(config.File, func(path string, info os.FileInfo, err error) error {
		if path[len(path)-3:] == ".wm" {
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
				_ = ioutil.WriteFile(path+".run", dd, 766)
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
	og := SplitString(a, []uint8{13, 10}) // 源文件
	dt := SplitString(b, []uint8{13, 10}) // 新文件
	if len(og) > len(dt) {
		out := make([][]byte, 0)
		for k, v := range og {
			if len(dt) > k && string(v) != string(dt[k]) {
				out = append(out, append(v, '*'))
			} else if len(dt) < k+1 {
				out = append(out, append(v, '-'))
			}
		}
		if len(out) != 0 {
			return false
		}
	} else if len(og) == len(dt) {
		out := make([][]byte, 0)
		for k, v := range og {
			if string(v) != string(dt[k]) {
				out = append(out, append(v, '*'))
				continue
			}
		}
		if len(out) != 0 {
			return false
		}
	} else if len(og) < len(dt) {
		out := make([][]byte, 0)
		for k, v := range dt {
			if len(og) > k && string(v) != string(og[k]) {
				out = append(out, append(v, '*'))
			} else if len(og) < k+1 {
				out = append(out, append(v, '+'))
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
