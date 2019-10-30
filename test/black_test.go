package test

import (
	"testing"

	g "../gate_way"
)

func TestBlack(t *testing.T) {
	g.StartBlackShield(&g.BlackShieldConfig{})
}
