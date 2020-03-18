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

func CreateMigrationTable(q Queryer) (e goerr.IError) {
	_, e = q.Exec(sqlCreateTableMigration())
	return
}

func sqlCreateTableMigration() string {
	return `CREATE TABLE IF NOT EXISTS migration
	(
		version  text NOT NULL PRIMARY KEY,
		apply_time  int       NOT NULL,
	);
	create index on migration (apply_time);
	`
}
