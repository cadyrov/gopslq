package gopsql

import (
	"database/sql"
	"fmt"
	"github.com/cadyrov/goerr"
	"log"
	"strconv"
	"strings"
)

type Queryer interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Close() error
}

type Builder struct {
	Queryer
}

func NewBuilder(q Queryer) *Builder {
	return &Builder{q}
}

func (b *Builder) Query(query string, args ...interface{}) (rows *sql.Rows, e goerr.IError) {
	if b.Queryer == nil {
		e = goerr.New(" Queryer is nil ")
		return
	}
	logQuery(query, args...)
	rows, err := b.Queryer.Query(b.Prepare(query), args...)
	if err != nil {
		e = goerr.New(err.Error())
	}
	return
}

func (b *Builder) QueryRow(query string, args ...interface{}) (row *sql.Row, e goerr.IError) {
	if b.Queryer == nil {
		e = goerr.New(" Queryer is nil ")
		return
	}
	logQuery(query, args...)
	row = b.Queryer.QueryRow(b.Prepare(query), args...)
	if row == nil {
		e = goerr.New("no data")
	}
	return
}

func (b *Builder) Exec(query string, args ...interface{}) (res sql.Result, e goerr.IError) {
	if b.Queryer == nil {
		e = goerr.New(" Queryer is nil ")
		return
	}
	logQuery(query, args...)
	res, err := b.Queryer.Exec(b.Prepare(query), args...)
	if err != nil {
		e = goerr.New(err.Error())
	}
	return
}

func (b Builder) Prepare(statement string) (prepared string) {
	pieces := strings.Split(statement, "?")
	for i := range pieces {
		if i < (len(pieces) - 1) {
			pieces[i] += "$" + strconv.Itoa(i+1)
		}
	}
	prepared = strings.Join(pieces, "")
	return
}

func logQuery(statement string, args ...interface{}) {
	pieces := strings.Split(statement, "?")
	for i := range pieces {
		if i < (len(pieces) - 1) {
			pieces[i] += fmt.Sprintf("\"%v\"", args[i])
		}
	}
	log.Println(strings.Join(pieces, ""))
	return
}
