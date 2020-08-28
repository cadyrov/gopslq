package gopsql

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

const (
	typeString = "string"
)

type Column struct {
	Name            string
	ModelName       string
	Default         *string
	IsNullable      bool
	DataType        string
	ModelType       string
	Schema          string
	Table           string
	Sequence        *string
	Description     *string
	IsPrimaryKey    bool
	Tags            string
	Import          []string
	IsArray         bool
	HasUniqueIndex  bool
	UniqueIndexName *string
}

// Array of columns.
type Columns []Column
type Imports []string

func (im *Imports) Add(value string) {
	for i := range *im {
		if (*im)[i] == value {
			return
		}
	}

	*im = append(*im, value)
}

func parse(rows *sql.Rows) (column *Column, e error) {
	column = &Column{}
	e = rows.Scan(
		&column.Name,
		&column.DataType,
		&column.IsNullable,
		&column.Schema,
		&column.Table,
		&column.IsPrimaryKey,
		&column.Default,
		&column.Sequence,
		&column.Description,
		&column.HasUniqueIndex,
		&column.UniqueIndexName,
	)

	return
}

// Get table columns from db.
func GetTableColumns(q Queryer, schema string, table string) (columns *Columns, imports *Imports, e error) {
	var hasPrimary bool

	query := fmt.Sprintf(`
SELECT a.attname                                                                       AS column_name,
       format_type(a.atttypid, a.atttypmod)                                            AS data_type,
       CASE WHEN a.attnotnull THEN FALSE ELSE TRUE END                                 AS is_nullable,
       s.nspname                                                                       AS schema,
       t.relname                                                                       AS table,
       (SELECT EXISTS(SELECT i.indisprimary
                      FROM pg_index i
                      WHERE i.indrelid = a.attrelid
                        AND a.attnum = ANY (i.indkey)
                        AND i.indisprimary IS TRUE))                                   AS is_primary,
       ic.column_default,
       pg_get_serial_sequence(ic.table_schema || '.' || ic.table_name, ic.column_name) AS sequence,
       col_description(t.oid, ic.ordinal_position)                                     AS description,
       (SELECT EXISTS(SELECT i.indisunique
                      FROM pg_index i
                      WHERE i.indrelid = a.attrelid
                        AND i.indisunique IS TRUE
                        AND a.attnum = ANY (i.indkey)))                                AS has_unique_index,
       (SELECT ins.indexname
        FROM pg_indexes ins
                 JOIN pg_index i ON ins.indexdef = pg_get_indexdef(i.indexrelid)
        WHERE i.indisunique IS TRUE
          AND i.indrelid = a.attrelid
          AND a.attnum = ANY (i.indkey))                                               AS unique_index_name
FROM pg_attribute a
         JOIN pg_class t ON a.attrelid = t.oid
         JOIN pg_namespace s ON t.relnamespace = s.oid
         LEFT JOIN information_schema.columns AS ic
                   ON ic.column_name = a.attname AND ic.table_name = t.relname AND ic.table_schema = s.nspname
         LEFT JOIN information_schema.key_column_usage AS kcu
                   ON kcu.table_name = t.relname AND kcu.column_name = a.attname
         LEFT JOIN information_schema.table_constraints AS tc
                   ON tc.constraint_name = kcu.constraint_name AND tc.constraint_type = 'FOREIGN KEY'
         LEFT JOIN information_schema.constraint_column_usage AS ccu ON ccu.constraint_name = tc.constraint_name
WHERE a.attnum > 0
  AND NOT a.attisdropped
  AND s.nspname = '%s'
  AND t.relname = '%s'
GROUP BY a.attname, a.atttypid, a.attrelid, a.atttypmod, a.attnotnull, s.nspname, t.relname, ic.column_default,
         ic.table_schema, ic.table_name, ic.column_name, a.attnum, t.oid, ic.ordinal_position
ORDER BY a.attnum;`, schema, table)
	rows, e := q.Query(query)

	if e != nil {
		return
	}

	columns = &Columns{}
	imports = &Imports{}

	for rows.Next() {
		column, err := parse(rows)
		if err != nil {
			e = err

			return
		}

		json := SnakeToCamel(column.Name, false)
		column.ModelName = SnakeToCamelWithGOData(column.Name, true)
		column.Tags = fmt.Sprintf(`%ccolumn:"%s" json:"%s"%c`, '`', column.Name, json, '`')

		uncnownColumnErr := fmt.Errorf("unknown column type: %s", column.DataType)

		switch {
		case column.DataType == "bigint":
			column.ModelType = "int64"
		case column.DataType == "integer":
			column.ModelType = "int"
		case column.DataType == "text":
			column.ModelType = typeString
		case column.DataType == "double precision":
			column.ModelType = "float64"
		case column.DataType == "boolean":
			column.ModelType = "bool"
		case column.DataType == "ARRAY":
			column.ModelType = "[]interface{}"
		case column.DataType == "json":
			column.ModelType = "json.RawMessage"

			imports.Add(`"encoding/json"`)
		case column.DataType == "smallint":
			column.ModelType = "int16"
		case column.DataType == "date":
			column.ModelType = "time.Time"

			imports.Add(`"time"`)
		case strings.Contains(column.DataType, "character varying"):
			column.ModelType = "string"
		case strings.Contains(column.DataType, "numeric"):
			column.ModelType = "float32"
		case column.DataType == "uuid":
			column.ModelType = "string"
		case column.DataType == "jsonb":
			column.ModelType = "json.RawMessage"

			imports.Add(`"encoding/json"`)
		case column.DataType == "uuid[]":
			column.ModelType = "[]string"
			column.IsArray = true

			imports.Add(`"github.com/lib/pq"`)
		case column.DataType == "integer[]":
			column.ModelType = "[]int64"
			column.IsArray = true

			imports.Add(`"github.com/lib/pq"`)
		case column.DataType == "bigint[]":
			column.ModelType = "[]int64"
			column.IsArray = true

			imports.Add(`"github.com/lib/pq"`)
		case column.DataType == "text[]":
			column.ModelType = "[]string"
			column.IsArray = true

			imports.Add(`"github.com/lib/pq"`)
		case strings.Contains(column.DataType, "timestamp"):
			column.ModelType = "time.Time"

			imports.Add(`"time"`)
		default:
			e = uncnownColumnErr

			return
		}

		if column.IsNullable && !column.IsArray {
			column.ModelType = "*" + column.ModelType
		}

		if column.IsPrimaryKey {
			hasPrimary = true
		}

		*columns = append(*columns, *column)
	}

	if !hasPrimary {
		for key, column := range *columns {
			if column.Name == "id" {
				(*columns)[key].IsPrimaryKey = true
				if (*columns)[key].ModelType[0] == '*' {
					(*columns)[key].ModelType = (*columns)[key].ModelType[1:]
					hasPrimary = true
				}

				break
			}
		}

		if !hasPrimary {
			var uniqueIndexName *string

			for key, column := range *columns {
				if column.HasUniqueIndex {
					if uniqueIndexName == nil {
						uniqueIndexName = column.UniqueIndexName
					}

					if *uniqueIndexName == *column.UniqueIndexName {
						(*columns)[key].IsPrimaryKey = true
						if (*columns)[key].ModelType[0] == '*' {
							(*columns)[key].ModelType = (*columns)[key].ModelType[1:]
						}
					}
				}
			}
		}
	}
	return
}

func CreateFile(schema string, table string, path string) (*os.File, string, error) {
	name := table

	if schema != "public" {
		name = fmt.Sprintf("%s_%s", schema, table)
	}

	filePath := fmt.Sprintf("%s.go", name)

	if path != "" {
		folderPath := path
		err := os.MkdirAll(folderPath, os.ModePerm)

		if err != nil {
			return nil, "", err
		}

		filePath = fmt.Sprintf("%s/%s.go", folderPath, name)
	}

	f, err := os.Create(filePath)

	if err != nil {
		return nil, "", err
	}

	return f, filePath, nil
}

func MakeModel(db Queryer, path string, schema string, table string, templatePath string) error {
	// Imports in model file
	var imports = []string{
		`"strings"`,
		`"database/sql"`,
		`"net/http"`,
		`"github.com/cadyrov/goerr"`,
		`"github.com/cadyrov/govalidation"`,
		`"github.com/cadyrov/gopsql"`,
	}

	var (
		tableFoundErr  = errors.New("No table found or no columns in table ")
		tableNameEmpty = errors.New("table name is empty")
	)

	var name = table

	if table == "" {
		return tableFoundErr
	}

	tmpl := template.New("model")

	templateFile, err := os.Open(templatePath)
	if err != nil {
		return err
	}

	data, err := ioutil.ReadAll(templateFile)
	if err != nil {
		return err
	}

	tmpl = template.Must(tmpl.Parse(string(data)))

	columns, tableImport, err := GetTableColumns(db, schema, table)
	if err != nil {
		return err
	}

	if columns == nil || len(*columns) == 0 {
		return tableNameEmpty
	}

	for _, imp := range *tableImport {
		imports = append(imports, imp)
	}

	if schema != "public" && schema != "" {
		name = schema + "_" + table
	}

	// To camel case
	modelName := SnakeToCamel(name, true)

	var hasSequence bool
	// Check for sequence and primary key
	for _, column := range *columns {
		if column.IsPrimaryKey && column.Sequence != nil {
			hasSequence = true

			break
		}
	}

	// Create file
	file, path, err := CreateFile(schema, table, path)
	if err != nil {
		return err
	}

	primaryColumns := Columns{}
	nonPrimaryColumns := Columns{}

	for i := range *columns {
		if (*columns)[i].IsPrimaryKey {
			primaryColumns = append(primaryColumns, (*columns)[i])

			continue
		}

		nonPrimaryColumns = append(nonPrimaryColumns, (*columns)[i])
	}

	// Parse template to file
	err = tmpl.Execute(file, struct {
		Model             string
		Table             string
		Name              string
		PrimaryColumns    Columns
		NonPrimaryColumns Columns
		Columns           Columns
		HasSequence       bool
		Imports           []string
	}{
		Model:             modelName,
		Table:             schema + "." + table,
		Name:              table,
		PrimaryColumns:    primaryColumns,
		NonPrimaryColumns: nonPrimaryColumns,
		Columns:           *columns,
		HasSequence:       hasSequence,
		Imports:           imports,
	})

	if err != nil {
		return err
	}

	err = file.Close()

	if err != nil {
		return err
	}

	cmd := exec.Command("go", "fmt", path)
	err = cmd.Run()

	if err != nil {
		return err
	}

	fmt.Printf("file created: %s", path)

	return nil
}

func SnakeToCamel(value string, firstTitle bool) (res string) {
	splitted := strings.Split(strings.Trim(value, ""), "_")
	for i := range splitted {
		if !firstTitle && i == 0 {
			res += splitted[i]
		} else {
			res += strings.Title(splitted[i])
		}
	}

	return
}

func SnakeToCamelWithGOData(value string, firstTitle bool) (res string) {
	splitted := strings.Split(strings.Trim(value, ""), "_")

	for i := range splitted {
		if strings.EqualFold(splitted[i], "id") ||
			strings.EqualFold(splitted[i], "db") ||
			strings.EqualFold(splitted[i], "sql") ||
			strings.EqualFold(splitted[i], "url") {
			res += strings.ToUpper(splitted[i])

			continue
		}

		if !firstTitle && i == 0 {
			res += splitted[i]

			continue
		}

		res += strings.Title(splitted[i])
	}

	return
}
