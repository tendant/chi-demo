package main

import (
	"github.com/tendant/chi-demo/app"
	"github.com/tendant/chi-demo/handler"
)

func main() {
	newApp := app.Default()
	handle := handler.Handle{
		Log: newApp.Log,
	}
	handler.Routes(newApp.R, handle)
	newApp.Run()
}
