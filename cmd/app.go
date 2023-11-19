package cmd

import (
	"github.com/Attsun1031/sqlc-query-gen/api"
	"github.com/urfave/cli/v2"
)

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
			appConfig := api.Config{
				TemplateConfigs: []api.TemplateConfig{
					{
						TemplatePath: "work/templates/sqlc-query.gtpl",
						TargetSchema: "public",
					},
				},
			}
			generator := api.NewGenerator(api.WithFuncMap(
				map[string]interface{}{
					"PsqlTypeToGqlType": psqlTypeToGqlType,
				},
			))
			return generator.Generate(c.Context, appConfig, *dbConfig)
		},
	}
	return app
}

func psqlTypeToGqlType(typeName string) string {
	return ""
}
