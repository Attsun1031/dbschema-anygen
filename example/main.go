package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Attsun1031/dbschema-anygen/api"
	"github.com/urfave/cli/v2"
)

func main() {
	app := NewApp()
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func Flags(d *api.DbConfig) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "db-host",
			Required:    false,
			Usage:       "db host",
			EnvVars:     []string{"DB_HOST"},
			Destination: &d.Host,
		},
		&cli.IntFlag{
			Name:        "db-port",
			Required:    false,
			Usage:       "db port",
			EnvVars:     []string{"DB_PORT"},
			Destination: &d.Port,
		},
		&cli.StringFlag{
			Name:        "db-user",
			Required:    false,
			Usage:       "db user",
			EnvVars:     []string{"DB_USER"},
			Destination: &d.User,
		},
		&cli.StringFlag{
			Name:        "db-password",
			Required:    false,
			Usage:       "db password",
			EnvVars:     []string{"DB_PASSWORD"},
			Destination: &d.Password,
		},
		&cli.StringFlag{
			Name:        "db-name",
			Required:    false,
			Usage:       "db name",
			EnvVars:     []string{"DB_NAME"},
			Destination: &d.DbName,
		},
	}
}

func NewApp() *cli.App {
	dbConfig := &api.DbConfig{}
	app := &cli.App{
		Name:  "sqlc-query-gen",
		Usage: "Generate sqlc queries",
		Flags: Flags(dbConfig),
		Action: func(c *cli.Context) error {
			// TODO: Read from config file
			cfg := api.Config{
				TargetSchema: "public",
				TemplateConfigs: []api.TemplateConfig{
					{
						TemplatePath: "templates/sqlc-query.gtpl",
						OutputPath:   "out/generated.sql",
					},
					{
						TemplatePath: "templates/graphql.gtpl",
						OutputPath:   "out/generated.graphql",
					},
				},
				DbConfig: *dbConfig,
			}
			typeBodyGen := &typeBodyGenerator{
				cfg: TypeConfig{
					AdditionalFieldDef: []GraphQLFieldDef{
						{
							Name:        "example",
							TypeDef:     "Example",
							Description: "Example reference",
						},
						{
							Name:        "example2",
							TypeDef:     "Example2",
							Description: "Example2 reference",
						},
					},
				},
			}

			generator := api.NewGenerator(api.WithFuncMap(
				map[string]interface{}{
					"GenGraphQLTypeBody": typeBodyGen.genTypeBody,
				},
			))
			return generator.Generate(c.Context, cfg)
		},
	}
	return app
}

func psqlTypeToGqlType(typeName string) string {
	t := strings.ToLower(typeName)
	switch t {
	case "integer", "bigint", "smallint":
		return "Int"
	case "character", "character varying", "text", "json", "jsonb":
		return "String"
	case "boolean":
		return "Boolean"
	case "numeric", "real", "double precision":
		return "Float"
	case "date", "timestamp without time zone", "timestamp with time zone":
		return "Time"
	case "uuid":
		return "UUID"
	}
	panic("unknown type: " + typeName)
}

type GraphQLConfig struct {
}

type TypeConfig struct {
	FieldConfigs       map[string]FieldConfig
	AdditionalFieldDef []GraphQLFieldDef
}

type FieldConfig struct {
	Ignore bool
}

type GraphQLFieldDef struct {
	Name        string
	TypeDef     string
	Description string
}

type typeBodyGenerator struct {
	cfg TypeConfig
}

func (x *typeBodyGenerator) genTypeBody(tableParam api.TableParam) []string {
	lines := []string{}
	for _, colParam := range tableParam.Columns {
		fieldCfg := x.cfg.FieldConfigs[colParam.ColumnName]
		if fieldCfg.Ignore {
			continue
		}
		lines = append(lines, toGqlFieldDef(colParam))
	}
	for _, additionalFieldDef := range x.cfg.AdditionalFieldDef {
		if additionalFieldDef.Description != "" {
			lines = append(lines, fmt.Sprintf(`"""%s"""`, additionalFieldDef.Description))
		}
		lines = append(lines, fmt.Sprintf("%s: %s", additionalFieldDef.Name, additionalFieldDef.TypeDef))
	}
	return lines
}

func toGqlFieldDef(colParam api.ColumnParam) string {
	if colParam.ColumnName == "id" {
		return "id: ID!"
	}

	fieldName := colParam.ColumnNameFCU
	typeDef := psqlTypeToGqlType(colParam.ColumnType)
	requiredSign := ""
	if !colParam.IsNullable {
		requiredSign = "!"
	}
	return fmt.Sprintf("%s: %s%s", fieldName, typeDef, requiredSign)
}
