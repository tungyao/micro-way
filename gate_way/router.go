package gate_way

import (
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"

	"./util"
)

// The role of routing here is to forward requests to the backend
// Still go to the registration center to get the registered address
// register center default connect to url:port localhost:80
// Router in addition to receiving and sending
// so it's pressure so bigger
type WayConfig struct {
	RegisterLocation string `register center location ,default localhost:9000`
	TimeOut          int    `set timeout ,default => 60 ms`
	IsCache          bool   `is cache service`
	CacheTime        int    `set cache time ,default 10s`
	MaxCap           int    ``
}
type RT struct {
	http.Handler
}

// GLOBAL PARAMETERS
//      GLOBAL_ALL_CONNECT  Count the number of service
//
var (
	GLOBAL_ALL_CONNECT int
	FPOOL              *FPool
	registerLocation   string
	timeOut            int
	isCache            bool
	cacheTime          int
)

// We need to get routing information from the registry
// use TCP protocol to connect registry
// in 0.1 version
func StartRouter(config *WayConfig) {
	util.CheckConfig(config, WayConfig{
		RegisterLocation: "localhost:6000",
		TimeOut:          60,
		IsCache:          true,
		CacheTime:        10,
	})
	util.CheckConfig(&registerLocation, config.RegisterLocation)
	util.CheckConfig(&timeOut, config.TimeOut)
	util.CheckConfig(&isCache, config.IsCache)
	util.CheckConfig(&cacheTime, config.CacheTime)
}
func (rt *RT) Router() http.Handler {
	return rt
}
func (rt *RT) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	op:=r.URL.Query().Get("wechat")
	if op != "" {
		sendQuery(w,op,r)
	}else {
		sendRouter(w, r)
	}
}
func sendQuery(w http.ResponseWriter,op string,r *http.Request)  {

}
func sendRouter(w http.ResponseWriter, r *http.Request) {
	n, err := net.Dial("tcp", registerLocation)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(501)
		_, _ = w.Write(template(501))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	data, err := ioutil.ReadAll(r.Body)
	err = r.Body.Close()
	if err != nil && err != io.EOF {
		log.Println("router.go -> 73", err)
	}
	_, _ = n.Write(data)
	_, _ = w.Write(getData(n))
	_ = n.Close()
	return
}
func template(n int) []byte {
	switch n {
	case 501:
		return []byte(`{"error":"501"}`)
	}
	return []byte(`{"error":"not found"}`)
}
func getData(a net.Conn) []byte {
	out := make([][]byte, 0)
	o := make([]byte, 0)
	for {
		data := make([]byte, 4096)
		n, err := a.Read(data)
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
	//fmt.Println(string(o))
	if string(o) == "false" {
		return template(501)
	}
	return o
}

// Initialize parameters
func init() {
	util.CheckConfig(&GLOBAL_ALL_CONNECT, 1)
}
