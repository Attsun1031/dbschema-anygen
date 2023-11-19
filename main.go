package main

import (
	"log"
	"os"

	"github.com/Attsun1031/sqlc-query-gen/cmd"
)

func main() {
	app := cmd.NewApp()
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
