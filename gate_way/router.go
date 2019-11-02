package gate_way

import (
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
	registerLocation   *string
	timeOut            *int
	isCache            *bool
	cacheTime          *int
)

// We need to get routing information from the registry
// use TCP protocol to connect registry
// in 0.1 version
func StartRouter(config *WayConfig) {
	util.CheckConfig(config, WayConfig{
		RegisterLocation: "localhost:81",
		TimeOut:          60,
		IsCache:          true,
		CacheTime:        10,
	})
	util.CheckConfig(registerLocation, config.RegisterLocation)
	util.CheckConfig(timeOut, config.TimeOut)
	util.CheckConfig(isCache, config.IsCache)
	util.CheckConfig(cacheTime, config.CacheTime)
}
func (rt *RT) Router() http.Handler {

	return rt
}
func (rt *RT) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sendRouter(w, r)
}
func sendRouter(w http.ResponseWriter, r *http.Request) {
	n, err := net.Dial("tcp", *registerLocation)
	if err != nil {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(501)
		_, _ = w.Write(template(501))
	}
}
func template(n int) []byte {
	switch n {
	case 501:
		return []byte(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Title</title>
</head>
<body>
<h1>Connect Service Error</h1>
</body>
</html>`)
	}
	return []byte("not found")
}

// Initialize parameters
func init() {
	util.CheckConfig(&GLOBAL_ALL_CONNECT, 1)
}
