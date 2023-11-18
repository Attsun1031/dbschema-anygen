package codegen

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/Attsun1031/sqlc-query-gen/pkg/db"
	"github.com/iancoleman/strcase"
	"github.com/samber/lo"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Param struct {
	TableParams []TableParam
}

type TableParam struct {
	TableName    string
	TableNameFCU string
	Columns      []ColumnParam
}

type ColumnParam struct {
	ColumnName    string
	ColumnNameFCU string
	ColumnType    string
}

func addNum(num int, a int) int {
	return num + a
}

func Generate(tableToColumnDefs map[string][]db.GetColumnDefinitionsRow, templateString string, templateName string) (string, error) {
	funcMap := template.FuncMap{
		"ToUpper":    strings.ToUpper,
		"FirstUpper": cases.Title(language.Und, cases.NoLower).String,
		"AddNum":     addNum,
	}

	tmpl, err := template.New(templateName).Funcs(funcMap).Parse(templateString)
	if err != nil {
		panic(err)
	}
	var param Param
	for tableName, columnDefs := range tableToColumnDefs {
		param.TableParams = append(param.TableParams, TableParam{
			TableName:    tableName,
			TableNameFCU: cases.Title(language.Und, cases.NoLower).String(strcase.ToCamel(tableName)),
			Columns: lo.Map(columnDefs, func(columnDef db.GetColumnDefinitionsRow, idx int) ColumnParam {
				return ColumnParam{
					ColumnName:    columnDef.ColumnName,
					ColumnNameFCU: cases.Title(language.Und, cases.NoLower).String(strcase.ToCamel(columnDef.ColumnName)),
					ColumnType:    columnDef.DataType,
				}
			}),
		})
	}

	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, param)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
