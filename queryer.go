package gopsql

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/cadyrov/goerr"
)

type Queryer interface {
	Query(query string, args ...interface{}) (*sql.Rows, goerr.IError)
	Exec(query string, args ...interface{}) (sql.Result, goerr.IError)
	QueryRow(query string, args ...interface{}) (*sql.Row, goerr.IError)
}

type DB struct {
	Debug bool
	*sql.DB
}

type Tx struct {
	Debug bool
	*sql.Tx
}

var (
	ErrQuerierIsNil = errors.New("Queryer is nil ")
	ErrNoData = errors.New("No data")
)

func (b *Tx) Query(query string, args ...interface{}) (rows *sql.Rows, e goerr.IError) {
	if b.Tx == nil {
		e = goerr.Internal(ErrQuerierIsNil)

		return
	}

	if b.Debug {
		logQuery(query, args...)
	}

	rows, err := b.Tx.Query(prepare(query), args...)
	if err != nil {
		e = goerr.Internal(err)
	}

	return
}

func (b *Tx) QueryRow(query string, args ...interface{}) (row *sql.Row, e goerr.IError) {
	if b.Tx == nil {
		e = goerr.Internal(ErrQuerierIsNil)

		return
	}

	if b.Debug {
		logQuery(query, args...)
	}

	row = b.Tx.QueryRow(prepare(query), args...)

	if row == nil {
		e = goerr.Internal(ErrNoData)
	}

	return
}

func (b *Tx) Exec(query string, args ...interface{}) (res sql.Result, e goerr.IError) {
	if b.Tx == nil {
		e = goerr.Internal(ErrQuerierIsNil)

		return
	}

	if b.Debug {
		logQuery(query, args...)
	}

	res, err := b.Tx.Exec(prepare(query), args...)
	if err != nil {
		e = goerr.Internal(err)
	}

	return
}

func (b *DB) Begin() (tx *Tx, e goerr.IError) {
	if b.DB == nil {
		e = goerr.Internal(ErrQuerierIsNil)

		return
	}

	transaction, err := b.DB.Begin()
	if err != nil {
		e = goerr.Internal(err)

		return
	}

	tx = &Tx{b.Debug, transaction}

	return
}

func (b *DB) Query(query string, args ...interface{}) (rows *sql.Rows, e goerr.IError) {
	if b.DB == nil {
		e = goerr.Internal(ErrQuerierIsNil)

		return
	}

	if b.Debug {
		logQuery(query, args...)
	}

	rows, err := b.DB.Query(prepare(query), args...)
	if err != nil {
		e = goerr.Internal(err)
	}

	return
}

func (b *DB) QueryRow(query string, args ...interface{}) (row *sql.Row, e goerr.IError) {
	if b.DB == nil {
		e = goerr.Internal(ErrQuerierIsNil)

		return
	}

	if b.Debug {
		logQuery(query, args...)
	}

	row = b.DB.QueryRow(prepare(query), args...)

	if row == nil {
		e = goerr.Internal(ErrNoData)
	}

	return
}

func (b *DB) Exec(query string, args ...interface{}) (res sql.Result, e goerr.IError) {
	if b.DB == nil {
		e = goerr.Internal(ErrQuerierIsNil)

		return
	}

	if b.Debug {
		logQuery(query, args...)
	}

	res, err := b.DB.Exec(prepare(query), args...)
	if err != nil {
		e = goerr.Internal(err)
	}

	return
}

func prepare(statement string) (prepared string) {
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
			pieces[i] += fmt.Sprintf("'%v'", args[i])
		}
	}

	log.Println(strings.Join(pieces, ""))
}
