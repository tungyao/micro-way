package test

import (
	"../register"
	"testing"
)

func TestPool2(t *testing.T) {
	register.StartServer(register.Config{
		Address:     ":6000",
		MaxCap:      1000,
		PollingType: register.MEMORY,
	})
}
func TestPool3(t *testing.T) {
	register.StartClient()
}
