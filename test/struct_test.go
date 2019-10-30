package test

import (
	"testing"

	g "../gate_way"
)

//
// func TestStruct(t *testing.T) {
// 	n:=&g.WayConfig{}
// 	g.CheckConfig(n,g.WayConfig{
// 		RegisterLocation: "lcao",
// 		TimeOut:          123,
// 		IsCache:          true,
// 		CacheTime:        10,
// 	})
// 	fmt.Println(n)
// }

func TestS(t *testing.T) {
	g.NewGateWayRouter(&g.WayConfig{})
}
