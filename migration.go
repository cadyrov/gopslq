package gopsql

import (
	"github.com/cadyrov/goerr"
	"time"
)

type Migration struct {
	Name    string   `json:"-" yaml:"-"`
	UpSql   []string `json:"upSql" yaml:"upSql"`
	DownSql []string `json:"downSql" yaml:"downSql"`
}

func (m *Migration) Up(q Queryer) (e goerr.IError) {
	for i := range m.UpSql {
		if m.UpSql[i] == "" {
			continue
		}
		if _, e = q.Exec(m.UpSql[i]); e != nil {
			return
		}
	}
	_, e = q.Exec(sqlAddMigration(), m.Name, time.Now().UnixNano()/int64(time.Second))
	return
}

func (m *Migration) Down(q Queryer) (e goerr.IError) {
	for i := range m.DownSql {
		if m.DownSql[i] == "" {
			continue
		}
		if _, e = q.Exec(m.DownSql[i]); e != nil {
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
