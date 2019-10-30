package gate_way

import (
	"fmt"
	"log"
	"net"

	"./util"
)

// The role of routing here is to forward requests to the backend
// Still go to the registration center to get the registered address
// register center default connect to url:port localhost:80
type WayConfig struct {
	RegisterLocation string `register center location ,default localhost:9000`
	TimeOut          int    `set timeout ,default => 1s`
	IsCache          bool   `is cache service`
	CacheTime        int    `set cache time ,default 10s`
}

// GLOBAL PARAMETERS
//      GLOBAL_ALL_CONNECT  Count the number of service
//
var (
	GLOBAL_ALL_CONNECT int
)

// We need to get routing information from the registry
// use TCP protocol to connect registry
// in 0.1 version
func NewGateWayRouter(config *WayConfig) {
	util.CheckConfig(config, WayConfig{
		RegisterLocation: "localhost:80",
		TimeOut:          1,
		IsCache:          true,
		CacheTime:        10,
	})
	con, err := net.Dial("tcp", config.RegisterLocation)
	if err != nil {
		log.Panicln(err)
	}
	data := make([]byte, 4096)
	con.Read(data)
	fmt.Println(data)
}

// Initialize parameters
func init() {
	util.CheckConfig(&GLOBAL_ALL_CONNECT, 1)
}
