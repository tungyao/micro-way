package test

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

func TestHttp(t *testing.T) {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
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
