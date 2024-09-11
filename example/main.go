package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Attsun1031/dbschema-anygen/api"
	"github.com/jackc/pgx/v5"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/urfave/cli/v2"
)

func main() {
	app := NewApp()
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func NewApp() *cli.App {
	app := &cli.App{
		Name:  "sqlc-query-gen",
		Usage: "Generate sqlc queries",
		Action: func(c *cli.Context) error {
			ctx := c.Context
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
			}
			typeBodyGen := &typeBodyGenerator{
				cfg: TypeConfig{
					AdditionalFieldDef: []GraphQLFieldDef{
						{
							Name:        "example",
							TypeDef:     "Example",
							Description: "Example reference",
						},
					},
				},
			}

			fmt.Println("Starting testcontainer...")
			pgContainer, err := postgres.Run(ctx,
				"docker.io/postgres:15-alpine",
				postgres.WithInitScripts(filepath.Join("database", "init-db.sql")),
				postgres.WithDatabase("anygen"),
				postgres.WithUsername("postgres"),
				postgres.WithPassword("postgres"),
				testcontainers.WithWaitStrategy(
					wait.ForListeningPort("5432/tcp"),
					wait.ForLog("database system is ready to accept connections").
						WithOccurrence(2).
						WithStartupTimeout(5*time.Second),
				),
			)
			if err != nil {
				log.Fatalf("failed to start container: %s", err)
			}

			// Clean up the container
			defer func() {
				fmt.Println("Terminating container...")
				if err := pgContainer.Terminate(ctx); err != nil {
					log.Fatalf("failed to terminate container: %s", err)
				}
			}()
			fmt.Println("Container started")

			connStr := pgContainer.MustConnectionString(ctx, "sslmode=disable")
			conn, err := pgx.Connect(ctx, connStr)
			if err != nil {
				log.Fatalf("failed to connect to container: %s", err)
			}
			defer func() {
				_ = conn.Close(ctx)
			}()

			generator := api.NewGenerator(api.WithFuncMap(
				map[string]interface{}{
					"GenGraphQLTypeBody": typeBodyGen.genTypeBody,
				},
			))
			return generator.Generate(ctx, cfg, conn)
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
	case "user-defined":
		return "String"
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

	fieldName := colParam.ColumnNameCamel
	typeDef := psqlTypeToGqlType(colParam.ColumnType)
	requiredSign := ""
	if !colParam.IsNullable {
		requiredSign = "!"
	}
	return fmt.Sprintf("%s: %s%s", fieldName, typeDef, requiredSign)
}
