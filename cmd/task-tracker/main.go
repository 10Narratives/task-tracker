package main

import (
	"fmt"
	"log"

	"github.com/10Narratives/task-tracker/internal/config"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Can not load environment variables.")
	}

	cfg := config.MustLoad()
	fmt.Println(cfg)

	// TODO: Initialize logger

	// TODO: Initialize storage

	// TODO: Initialize http server
}
