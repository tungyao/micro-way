package test

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	g "../gate_way"
	u "../gate_way/util"
)

func TestGateWay(t *testing.T) {
	t.Log("Starting")
	l, err := net.Listen("tcp", ":80")
	if err != nil {
		t.Fatal(err)
	}
	l = g.Limiter(&g.Config{MaxConn: 2000}, l)
	l = g.StartBlackShield(&g.BlackShieldConfig{}, l).Next()
	var open int32
	http.Serve(l, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if n := atomic.AddInt32(&open, 1); n > 5 {
			t.Errorf("%d open connections, want <= %d", n, 5)
		}
		time.Sleep(time.Millisecond)
		defer atomic.AddInt32(&open, -1)
		_, _ = fmt.Fprint(w, "some body")
	}))
}
func TestSplit(t *testing.T) {
	for i := 0; i < 10000000; i++ {
		strings.Split("127.0.0.1:61142", ":")
	}
}
func TestSplit2(t *testing.T) {
	for i := 0; i < 10000000; i++ {
		u.SplitString([]byte("ashdahs**djhajksd"), []byte("**"))
	}
}
