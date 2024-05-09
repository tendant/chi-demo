package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/ggicci/httpin"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/tendant/chi-demo/app"
	"github.com/tendant/chi-demo/dbconn"
	"github.com/tendant/chi-demo/handler"
	"github.com/tendant/chi-demo/tutorial"
	"go.uber.org/zap"
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
	slog.Warn("demo error", "stack", string(debug.Stack()))
	driver := "postgres"
	dsn := "host=localhost port=5432 user=demo password=pwd dbname=demo_db sslmode=disable"
	settings := dbconn.DBConnSettings{}
	db, _ := dbconn.OpenDBConn(driver, dsn, settings)
	queries := tutorial.New(nil)
	handle := handler.Handle{
		DB:      db,
		Queries: queries,
	}
	handler.Routes(newApp.R, handle)
	newApp.Run()
}

type Handle struct {
	DB      *sql.DB
	Queries *tutorial.Queries
}

func (handle *Handle) Demo(w http.ResponseWriter, r *http.Request) {
	slog.Debug("Demo Log")
	fmt.Fprintf(w, "Hello, Demo")
}

type QueryInput struct {
	Q string `in:"query=q"`
}

func (handle *Handle) Query(w http.ResponseWriter, r *http.Request) {
	slog.Debug("Query...")
	q := r.Context().Value(httpin.Input).(*QueryInput)
	slog.Debug("Q:", zap.Any("q", q))
	render.PlainText(w, r, "Query")
}

type DemoPostInput struct {
	Name string `in:"form=name"`
}

func (handle *Handle) DemoPost(w http.ResponseWriter, r *http.Request) {
	slog.Debug("Form...")
	q := r.Context().Value(httpin.Input).(*DemoPostInput)
	slog.Debug("Q:", zap.Any("q", q))
	render.PlainText(w, r, "DemoPost")
}

type DemoJsonInput struct {
	httpin.JSONBody
	Login    string   `json:"login"`
	Password string   `json:"password"`
	Amount   *float64 `json:"amount,omitempty"`
}

func (handle *Handle) DemoJson(w http.ResponseWriter, r *http.Request) {
	slog.Debug("Json...")
	q := r.Context().Value(httpin.Input).(*DemoJsonInput)
	slog.Info("Q:", zap.Any("q", q))
	render.JSON(w, r, "OK")
}

type DemoListInput struct {
	httpin.JSONBody
	Emails []string `json:"emails"`
}

func (handle *Handle) DemoList(w http.ResponseWriter, r *http.Request) {
	slog.Debug("Json list...")
	q := r.Context().Value(httpin.Input).(*DemoListInput)
	slog.Debug("Q:", zap.Any("q", q))
	render.JSON(w, r, "OK")
}

func (handle *Handle) Ready(w http.ResponseWriter, r *http.Request) {
	ready := false
	var err error
	if handle.DB != nil {
		ready, err = dbconn.CheckDBConn(handle.DB)
		slog.Warn("Failed checking database connection", "err", err)
	}
	if ready {
		render.PlainText(w, r, "Ready")
	} else {
		render.PlainText(w, r, "Not Ready")
	}
}

func Routes(r *chi.Mux, handle Handle) {
	r.Get("/demo", handle.Demo)
	r.Get("/healthz/ready", handle.Ready)
	r.With(httpin.NewInput(QueryInput{})).Get("/query", handle.Query)
	r.With(httpin.NewInput(DemoPostInput{})).Post("/post", handle.DemoPost)
	r.With(httpin.NewInput(DemoJsonInput{})).Post("/json", handle.DemoJson)
	r.With(httpin.NewInput(DemoListInput{})).Post("/json/list", handle.DemoList)
}
