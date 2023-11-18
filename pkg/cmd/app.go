package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/Attsun1031/sqlc-query-gen/pkg/codegen"
	"github.com/Attsun1031/sqlc-query-gen/pkg/db"
	"github.com/jackc/pgx/v5"
	"github.com/samber/lo"
	"github.com/urfave/cli/v2"
)

type DbConfig struct {
	// Connection configuration
	Host     string
	Port     int
	User     string
	Password string
	DbName   string
}

func (d *DbConfig) Flags() []cli.Flag {
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
	dbConfig := &DbConfig{}
	app := &cli.App{
		Name:  "sqlc-query-gen",
		Usage: "Generate sqlc queries",
		Flags: dbConfig.Flags(),
		Action: func(c *cli.Context) error {
			// TODO: Read from config file
			appConfig := Config{
				TemplatePath: "work/templates/generated.gtpl",
				TargetSchema: "public",
			}
			return appAction(c.Context, appConfig, *dbConfig)
		},
	}
	return app
}

func appAction(ctx context.Context, appCfg Config, dbCfg DbConfig) error {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
		dbCfg.Host, dbCfg.Port, dbCfg.User, dbCfg.Password, dbCfg.DbName)
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	queries := db.New(conn)
	dat, err := os.ReadFile(appCfg.TemplatePath)
	if err != nil {
		return err
	}
	templateString := string(dat)

	columnDefs, err := queries.GetColumnDefinitions(ctx, appCfg.TargetSchema)
	if err != nil {
		return err
	}
	tableToColumns := lo.GroupBy(columnDefs, func(c db.GetColumnDefinitionsRow) string {
		return c.TableName
	})
	ret, err := codegen.Generate(tableToColumns, templateString, appCfg.TemplatePath)
	fmt.Println(ret)
	return err
}
