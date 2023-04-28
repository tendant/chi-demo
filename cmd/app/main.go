package main

import (
	"github.com/tendant/chi-demo/app"
	"github.com/tendant/chi-demo/user"
)

func main() {
	newApp := app.Default()
	var userImpl user.UserImpl
	newApp.R.Mount("/", user.Handler(&userImpl))
	newApp.Run()
}
