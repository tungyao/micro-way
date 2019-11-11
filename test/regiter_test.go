package test

import (
	r "../register"
	"testing"
)

func TestRegister(t *testing.T) {
	r.StartServer(r.Config{
		PollingType: r.FILE,
		Address:     ":6000",
		File:        "C:\\github\\micro-way\\test",
	})
}
func TestClient(t *testing.T) {
	// r.StartClient()
}
