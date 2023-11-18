package codegen

import (
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/Attsun1031/sqlc-query-gen/pkg/db"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const templateString = `
{{ $tableName := .TableName }}
-- name: Get{{ $tableName | FirstUpper }} :one
SELECT * FROM {{ $tableName }} WHERE id = $1;
`

type param struct {
	TableName string
}

func GenerateSql(tableName string, columnDefs []db.GetColumnDefinitionsRow) {
	funcMap := template.FuncMap{
		"ToUpper":    strings.ToUpper,
		"FirstUpper": cases.Title(language.Und, cases.NoLower).String,
	}

	tmpl, err := template.New("sql").Funcs(funcMap).Parse(templateString)
	if err != nil {
		panic(err)
	}
	p := param{TableName: tableName}
	err = tmpl.Execute(os.Stdout, p)
	if err != nil {
		fmt.Printf("executing template: %v", err)
		panic(err)
	}
}
