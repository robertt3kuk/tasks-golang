package main

import (
	"log"

	"github.com/robertt3kuk/tasks-golang/server/config"
	"github.com/robertt3kuk/tasks-golang/server/internal/app"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalln(err)
	}
	app.Run(cfg)
}
