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

func TestConfig_Connect(t *testing.T) {
	cnf := getConfig()
	db, e := cnf.Connect()
	if e != nil {
		t.Fatal(e)
	}
	_ = db
}

func getConfig() Config {
	return Config{
		Host:           "eledam.clkgw5sqzopc.eu-central-1.rds.amazonaws.com",
		Port:           5432,
		UserName:       "files",
		DbName:         "files",
		Password:       "Filedb1",
		SslMode:        "disabled",
		MaxConnections: 50,
		ConnectionIdle: 10,
	}
}
