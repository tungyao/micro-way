package test

import (
	"testing"

	"../register"
)

func TestDiff(t *testing.T) {
	a := "zbcdefghijklmn\r\nopqrt\r\nvwxyz\r\nsdadasd\r\nqwefdsdsa\r\nsdadasd\r\nsdadasd"
	b := "zbcdefashimn\r\nopqrstu\r\nvwxyz\r\nqwefdsdsa\r\nsdadasd\r\nsdadasd"
	register.Diff([]byte(a), []byte(b))
}
