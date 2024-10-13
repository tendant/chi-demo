package main

import (
	"github.com/tendant/chi-demo/app"
	"github.com/tendant/chi-demo/user"
)

func main() {
	// newApp := app.Default()
	server := app.DefaultApp()
	var userImpl user.UserImpl
	server.R.Mount("/", user.Handler(&userImpl))
	server.Run()
}
