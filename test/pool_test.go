package test

import (
	"../register"
	"testing"
)

func TestPool2(t *testing.T) {
	register.StartServer(register.Config{
		Address: ":6000",
		MaxCap:  1000,
		File:    `D:\Tung\Github\micro-way\test`,
	})
}
func TestPool3(t *testing.T) {
	// register.StartClient("normal")
	register.StartPolling(3, register.Config{
		Address: ":6000",
		MaxCap:  1000,
		File:    `D:\Tung\Github\micro-way\test`,
	})
}
