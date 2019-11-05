package test

import (
	"fmt"
	"testing"

	g "../gate_way"
	u "../gate_way/util"
)

func TestUtil(t *testing.T) {
	data := u.SplitString([]byte("127.0.0.1**4545"), []byte("**"))
	fmt.Println(string(data[0]))
	fmt.Println(string(data[1]))
}
func TestRouter(t *testing.T) {
	g.StartGateWay(&g.Config{}, &g.BlackShieldConfig{}, &g.WayConfig{RegisterLocation: "localhost:6000"})
}
