package gopsql_test

import (
	"example.com/m/v2"
	"fmt"
	"testing"
)

func TestConfig_ConnectionUrl(t *testing.T) {
	cnf := gopsql.Config{}
	_, e := cnf.ConnectionURL()

	if e == nil {
		t.Fatal("must be an error")
	}

	fmt.Println(e)
}
