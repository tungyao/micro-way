package gate_way

import (
	"log"
	"net"
	"net/http"
)

func StartGateWay(config *Config, shieldConfig *BlackShieldConfig, wayConfig *WayConfig) {
	l, err := net.Listen("tcp", ":7000")
	if err != nil {
		log.Panicln(err)
	}
	l = Limiter(config, l)
	l = StartBlackShield(shieldConfig, l).Next()
	go StartRouter(wayConfig)
	r := new(RT)
	err = http.Serve(l, r.Router())
	if err != nil {
		log.Println(err)
	}
}
