package api

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/Attsun1031/sqlc-query-gen/pkg/db"
	"github.com/iancoleman/strcase"
	"github.com/jackc/pgx/v5"
	"github.com/samber/lo"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Generator struct {
	FuncMaps template.FuncMap
}

type GeneratorOption func(*Generator)

func WithFuncMap(funcMap template.FuncMap) GeneratorOption {
	return func(g *Generator) {
		g.FuncMaps = lo.Assign(g.FuncMaps, funcMap)
	}
}

func NewGenerator(opts ...GeneratorOption) *Generator {
	g := &Generator{
		FuncMaps: defaultFuncMap,
	}
	for _, opt := range opts {
		opt(g)
	}
	fmt.Printf("FuncMaps: %+v\n", g.FuncMaps)
	return g
}

func addNum(num int, a int) int {
	return num + a
}

var defaultFuncMap = template.FuncMap{
	"ToUpper":    strings.ToUpper,
	"FirstUpper": cases.Title(language.Und, cases.NoLower).String,
	"AddNum":     addNum,
}

func (x *Generator) Generate(ctx context.Context, appCfg Config, dbCfg DbConfig) error {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
		dbCfg.Host, dbCfg.Port, dbCfg.User, dbCfg.Password, dbCfg.DbName)
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	queries := db.New(conn)
	columnDefs, err := queries.GetColumnDefinitions(ctx, appCfg.TargetSchema)
	if err != nil {
		return err
	}
	tmplParam := columnDefsToTemplateParam(columnDefs)
	for _, templateCfg := range appCfg.TemplateConfigs {
		fmt.Printf("Generate %s\n", templateCfg.TemplatePath)

		// Parse template
		dat, err := os.ReadFile(templateCfg.TemplatePath)
		if err != nil {
			return err
		}
		templateString := string(dat)
		tmpl, err := template.New(templateCfg.TemplatePath).Funcs(x.FuncMaps).Parse(templateString)
		if err != nil {
			panic(err)
		}

		// Execute template
		buf := new(bytes.Buffer)
		err = tmpl.Execute(buf, tmplParam)
		if err != nil {
			return err
		}
		ret := strings.Trim(buf.String(), "\n ")

		// Output genereated contents
		fmt.Println(ret)
	}
	return nil
}

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
	IsNullable    bool
}

func columnDefsToTemplateParam(columnDefs []db.GetColumnDefinitionsRow) Param {
	tableToColumnDefs := lo.GroupBy(columnDefs, func(c db.GetColumnDefinitionsRow) string {
		return c.TableName
	})
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
					IsNullable:    columnDef.IsNullable,
				}
			}),
		})
	}
	return param
}
