package gopsql

import (
	"time"

	"github.com/cadyrov/goerr/v2"
)

type Migration struct {
	Name    string   `json:"-" yaml:"-"`
	UpSQL   []string `json:"upSql" yaml:"upSql"`
	DownSQL []string `json:"downSql" yaml:"downSql"`
}

func (m *Migration) Up(q Queryer) (e goerr.IError) {
	for i := range m.UpSQL {
		if m.UpSQL[i] == "" {
			continue
		}

		if _, e = q.Exec(m.UpSQL[i]); e != nil {
			return
		}
	}

	_, e = q.Exec(sqlAddMigration(), m.Name, time.Now().UnixNano()/int64(time.Second))

	return
}

func (m *Migration) Down(q Queryer) (e goerr.IError) {
	for i := range m.DownSQL {
		if m.DownSQL[i] == "" {
			continue
		}

		if _, e = q.Exec(m.DownSQL[i]); e != nil {
			return
		}
	}

	_, e = q.Exec(sqlDropMigration(), m.Name)

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
		apply_time  int       NOT NULL
	);
	create index on migration (apply_time);
	`
}

func sqlAddMigration() string {
	return `
		INSERT INTO migration (version, apply_time) values (?, ?);
	`
}

func sqlDropMigration() string {
	return `
		DELETE FROM migration WHERE version = ?;
	`
}
