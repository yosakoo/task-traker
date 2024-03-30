package main

import (
	"log"

	"github.com/yosakoo/task-traker/internal/app"
	"github.com/yosakoo/task-traker/internal/config"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	app.Run(cfg)
}
