package gopsql

import (
	"database/sql"
	"fmt"
	"github.com/cadyrov/goerr"
	_ "github.com/lib/pq"
	"net/http"
)

func (c *Config) Connect() (db *sql.DB, e goerr.IError) {
	cu, e := c.ConnectionUrl()
	if e != nil {
		return
	}
	db, err := sql.Open("postgres", cu)
	if err != nil {
		e = goerr.New(err.Error())
	}
	return
}

func (c *Config) ConnectionUrl() (url string, e goerr.IError) {
	url = "host=%s port=%d user=%s password=%s dbname=%s"
	if c.Host == "" || c.Port == 0 || c.UserName == "" || c.DbName == "" || c.Password == "" {
		e = goerr.New(fmt.Sprintf("config isn't full "+url, c.Host, c.Port, c.UserName, c.Password, c.DbName)).Http(http.StatusBadRequest)
		return
	}
	if c.SslMode != "" {
		url += " sslmode=" + c.SslMode
	}
	if c.Binary {
		url += " binary_parameters=yes"
	}
	url = fmt.Sprintf(url, c.Host, c.Port, c.UserName, c.Password, c.DbName)
	return
}
