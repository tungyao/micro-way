package test

import (
	"fmt"
	"net"
	"net/http"
	"sync/atomic"
	"testing"
	"time"

	g "../gate_way"
)

func TestGateWay(t *testing.T) {
	t.Log("Starting")
	l, err := net.Listen("tcp", ":80")
	if err != nil {
		t.Fatal(err)
	}
	l = g.Limiter(g.Config{MaxConn: 2000}, l)
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
