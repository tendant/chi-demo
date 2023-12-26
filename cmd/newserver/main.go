package main

import (
	"github.com/ilyakaznacheev/cleanenv"
	server "github.com/tendant/chi-demo/sever"
)

func main() {
	var config server.Config
	cleanenv.ReadEnv(&config)
	s := server.Default(config)
	s.Run()
}
