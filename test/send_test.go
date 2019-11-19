package test

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"testing"
)

func TestHttp(t *testing.T) {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Println(request.URL.String())
		fmt.Println(request.Method)
		data := make([]byte, 0)
		for i := 0; i < 10; i++ {
			data = append(data, 'a')
		}
		fmt.Fprint(writer, string(data))

	})
	http.ListenAndServe(":80", nil)
}

func TestSend(t *testing.T) {
	n, err := http.Get("http://127.0.0.1:80")
	if err != nil {
		log.Println("server -> 116", err)
	}
	data, _ := ioutil.ReadAll(n.Body)
	fmt.Println(data)
}
func TestNet(t *testing.T) {
	a, _ := net.Listen("tcp", ":82")
	for {
		c, _ := a.Accept()
		data := make([]byte, 0)
		for i := 0; i < 10; i++ {
			data = append(data, 'a')
		}
		c.Write(data)
		c.Close()
	}
}
func TestSendNet(t *testing.T) {
	a, _ := net.Dial("tcp", "localhost:80")
	ns := 0
	out := make([][]byte, 0)
	o := make([]byte, 0)
	for {
		data := make([]byte, 1024)
		n, err := a.Read(data)
		ns += n
		out = append(out, data)
		if n == 0 || err == io.EOF {
			break
		}
	}
	for _, v := range out {
		for _, j := range v {
			if j == 0 {
				continue
			}
			o = append(o, j)
		}
	}
	fmt.Println(len(o))
}
