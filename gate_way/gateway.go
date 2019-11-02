package gate_way

import (
	"log"
	"net"
	"net/http"
)

func StartGateWay() {
	l, err := net.Listen("tcp", ":80")
	if err != nil {
		log.Panicln(err)
	}

	l = Limiter(&Config{MaxConn: 2000}, l)
	l = StartBlackShield(&BlackShieldConfig{}, l).Next()
	go StartRouter(&WayConfig{})
	r := new(RT)
	_ = http.Serve(l, r.Router())
}
