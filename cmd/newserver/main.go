package main

import (
	"github.com/ilyakaznacheev/cleanenv"
	server "github.com/tendant/chi-demo/server"
)

func main() {
	var config server.Config
	cleanenv.ReadEnv(&config)
	s := server.Default(config)
	server.Routes(s.R)
	s.Run()
}
