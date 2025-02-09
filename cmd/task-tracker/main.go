package main

import "github.com/10Narratives/task-tracker/internal/app"

// @title Task Tracker App
// @version 1.0
// @description REST API server for TODO Application
// @host localhost:8080
// @BasePath /
func main() {
	app := app.New()
	app.Run()
}
