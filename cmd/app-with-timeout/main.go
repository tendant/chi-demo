package main

import (
	"github.com/tendant/chi-demo/app"
)

func main() {
	server := app.DefaultAppWithTimeout()
	app.RoutesDefault(server.R)
	app.RoutesVersion(server.R)
	app.RoutesHealthz(server.R)
	app.RoutesHealthzReady(server.R)
	server.RunWithTimeout()
}