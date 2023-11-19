package api

import (
	"context"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/Attsun1031/sqlc-query-gen/pkg/codegen"
	"github.com/Attsun1031/sqlc-query-gen/pkg/db"
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
		g.FuncMaps = lo.Assign(defaultFuncMap, funcMap)
	}
}

func NewGenerator(opts ...GeneratorOption) *Generator {
	g := &Generator{
		FuncMaps: defaultFuncMap,
	}
	for _, opt := range opts {
		opt(g)
	}
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
	for _, templateCfg := range appCfg.TemplateConfigs {
		dat, err := os.ReadFile(templateCfg.TemplatePath)
		if err != nil {
			return err
		}
		templateString := string(dat)

		columnDefs, err := queries.GetColumnDefinitions(ctx, templateCfg.TargetSchema)
		if err != nil {
			return err
		}
		tableToColumns := lo.GroupBy(columnDefs, func(c db.GetColumnDefinitionsRow) string {
			return c.TableName
		})
		ret, err := codegen.Generate(tableToColumns, templateString, templateCfg.TemplatePath)
		if err != nil {
			return err
		}
		fmt.Println(ret)
	}
	return nil
}
