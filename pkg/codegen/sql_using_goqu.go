package codegen

import (
	"fmt"

	"github.com/doug-martin/goqu/v9"
	// import the dialect
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
)

type Sqler interface {
	ToSQL() (string, []interface{}, error)
}

func GenerateSqlByGoqu() {
	dialect := goqu.Dialect("postgres")
	db := dialect.From("test").Prepared(true)

	// select
	singleValue := "ps"
	// multiValue := []string{"ps", "ps2"}
	ss := db.Where(goqu.Ex{"id": singleValue})
	outSql("select by id", ss)
	mss := db.Where(goqu.L(`"id" = ANY($1::UUID [])`))
	outSql("select by ids", mss)

	// insert
	is := db.Insert().Rows(goqu.Record{"id": singleValue, "name": singleValue})
	outSql("insert", is)

	// delete
	ds := db.Delete().Where(goqu.Ex{"id": singleValue})
	outSql("delete by id", ds)
}

func outSql(title string, sqler Sqler) {
	fmt.Printf("---- %v ----\n", title)
	sql, args, err := sqler.ToSQL()
	if err != nil {
		fmt.Println("An error occurred while generating the SQL", err.Error())
	} else {
		fmt.Printf("SQL: %v\n", sql)
		fmt.Printf("ARGS: %v\n", args)
	}
}
