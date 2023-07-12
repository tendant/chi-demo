package main

import (
	"runtime/debug"

	"github.com/tendant/chi-demo/app"
	"github.com/tendant/chi-demo/dbconn"
	"github.com/tendant/chi-demo/handler"
	"github.com/tendant/chi-demo/tutorial"
	"golang.org/x/exp/slog"
)

func main() {
	newApp := app.Default()
	// db, err := sql.Open("postgres", "host=localhost port=5432 user=demo password=pwd dbname=demo_db sslmode=disable")
	// if err != nil {
	// 	panic("connect to database failed")
	// } else {
	// 	fmt.Println("connect to database successed")
	// }
	slog.Error("demo error", "stack", debug.Stack())
	driver := "postgres"
	dsn := "host=localhost port=5432 user=demo password=pwd dbname=demo_db sslmode=disable"
	settings := dbconn.DBConnSettings{}
	db, _ := dbconn.OpenDBConn(driver, dsn, settings)
	queries := tutorial.New(nil)
	handle := handler.Handle{
		Log:     newApp.Log,
		DB:      db,
		Queries: queries,
	}
	handler.Routes(newApp.R, handle)
	newApp.Run()
}
