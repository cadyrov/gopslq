package gopsql_test

import (
	"fmt"
	"testing"

	"github.com/cadyrov/gopsql"
)

func TestConfig_ConnectionUrl(t *testing.T) {
	cnf := gopsql.Config{}
	_, e := cnf.ConnectionURL()

	if e == nil {
		t.Fatal("must be an error")
	}

	fmt.Println(e)
}
