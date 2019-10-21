package test

import (
	"testing"

	"../register"
)

func TestDiff(t *testing.T) {
	a := "zbcdefashimn\r\nopqrstu\r\nvwxyz\r\nqwefdsdsa\r\nsdadasd\r\nasdas"
	b := "zbcdefashimn\r\nopqrstu\r\nvwxyz\r\nqwefdsdsa\r\nsdadasd\r\nsdadasd"
	boo := register.Diff([]byte(a), []byte(b))
	t.Log(boo)
}
