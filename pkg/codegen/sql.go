package codegen

import (
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/Attsun1031/sqlc-query-gen/pkg/db"
	"github.com/iancoleman/strcase"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type param struct {
	TableName    string
	TableNameFCU string
}

func GenerateSql(tableName string, columnDefs []db.GetColumnDefinitionsRow, templateString string) {
	funcMap := template.FuncMap{
		"ToUpper":    strings.ToUpper,
		"FirstUpper": cases.Title(language.Und, cases.NoLower).String,
	}

	tmpl, err := template.New("sql").Funcs(funcMap).Parse(templateString)
	if err != nil {
		panic(err)
	}
	p := param{
		TableName:    tableName,
		TableNameFCU: cases.Title(language.Und, cases.NoLower).String(strcase.ToCamel(tableName)),
	}
	err = tmpl.Execute(os.Stdout, p)
	if err != nil {
		fmt.Printf("exeFCUting template: %v", err)
		panic(err)
	}
}
