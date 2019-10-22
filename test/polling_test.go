package test

import (
	"net/http"
	"testing"
	"time"

	"../gate_way"
)

func TestMonitor(t *testing.T) {
	gate_way.StartGate()
	go func() {
		startTime := time.Now().UnixNano() / 0xf4240
		http.Get("https://www.yaop.ink")
		t.Log(time.Now().UnixNano()/0xf4240 - startTime)
	}()
	time.Sleep(time.Second * 2)
}
