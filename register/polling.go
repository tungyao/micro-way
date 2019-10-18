package register

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// Loop Time 30s
// 轮询器 自动开启轮询 和 手动开启轮询
// 同时还将具备收集 服务状态  自动报告 和 收集 时延

// TODO MD5
func StatusConfigFile(config Config) {
	err := filepath.Walk(config.File, func(path string, info os.FileInfo, err error) error {
		if path[len(path)-3:] == ".wm" {
			d, err := ioutil.ReadFile(path)
			if err != nil {
				log.Println(err)
			}
			dd, err := ioutil.ReadFile(path + ".run")
			if err != nil {
				log.Println(err)
			}
			// for k, v := range dd {
			// 	dd[k] = v << 1
			// }
			fmt.Println(string(dd))
			if Diff(d, dd) {
				// for k, v := range d {
				// 	d[k] = v << 1
				// }
				_ = ioutil.WriteFile(path+".run", d, 766)
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
			} else {
				out = append(out, append(v, '='))
			}
		}
		OutDisplayOrFile(out)
		if len(out) != 0 {
			return true
		}
	} else if len(og) == len(dt) {
		out := make([][]byte, 0)
		for k, v := range og {
			if string(v) != string(dt[k]) {
				out = append(out, append(v, '*'))
				continue
			}
			out = append(out, append(v, '='))
		}
		OutDisplayOrFile(out)
		if len(out) != 0 {
			return true
		}
	} else if len(og) < len(dt) {
		out := make([][]byte, 0)
		for k, v := range dt {
			if len(og) > k && string(v) != string(og[k]) {
				out = append(out, append(v, '*'))
			} else if len(og) < k+1 {
				out = append(out, append(v, '+'))
			} else {
				out = append(out, append(v, '='))
			}
		}
		OutDisplayOrFile(out)
		if len(out) != 0 {
			return true
		}
	}
	return false
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
