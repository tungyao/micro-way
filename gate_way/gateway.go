package gate_way

import (
	"log"
	"net"
	"net/http"
)

func StartGateWay() {
	l, err := net.Listen("tcp", ":7000")
	if err != nil {
		log.Panicln(err)
	}
	l = Limiter(&Config{MaxConn: 3000}, l)
	l = StartBlackShield(&BlackShieldConfig{}, l).Next()
	go StartRouter(&WayConfig{RegisterLocation: ":6000"})
	r := new(RT)
	err = http.Serve(l, r.Router())
	if err != nil {
		log.Println(err)
	}
}
