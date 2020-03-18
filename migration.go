package gopsql

import "github.com/cadyrov/goerr"

type Migration struct {
	UpSql   string `json:"upSql"`
	DownSql string `json:"upSql"`
}

func (m *Migration) Up(q Queryer) (e goerr.IError) {
	if m.UpSql == "" {
		return
	}
	_, e = q.Exec(m.UpSql)
	return
}

func (m *Migration) Down(q Queryer) (e goerr.IError) {
	if m.UpSql == "" {
		return
	}
	_, e = q.Exec(m.DownSql)
	return
}
