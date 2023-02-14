package main

import (
	"database/sql"
	"fmt"

	"github.com/tendant/chi-demo/app"
	"github.com/tendant/chi-demo/handler"
	"github.com/tendant/chi-demo/tutorial"
)

func main() {
	newApp := app.Default()
	db, err := sql.Open("postgres", "host=localhost port=5432 user=demo password=pwd dbname=demo_db sslmode=disable")
	if err != nil {
		panic("connect to database failed")
	} else {
		fmt.Println("connect to database successed")
	}
	queries := tutorial.New(db)
	handle := handler.Handle{
		Log:     newApp.Log,
		DB:      db,
		Queries: queries,
	}
	handler.Routes(newApp.R, handle)
	newApp.Run()
}
