package main

import "github.com/tendant/chi-demo/app"

func main() {
	newApp := app.Default()
	newApp.Run()
}
