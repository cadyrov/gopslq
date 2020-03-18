package gopsql

import (
	"fmt"
	"testing"
)

func TestConfig_ConnectionUrl(t *testing.T) {
	cnf := Config{}
	_, e := cnf.ConnectionUrl()
	if e == nil {
		t.Fatal("must be an error")
	}
	fmt.Println(e)

}
