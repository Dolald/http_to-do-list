package main

import (
	"todolist/internal/app"

	_ "github.com/lib/pq"
)

// @title Todo App API
// @version 1.0
// @description API Server for TodoList Application

// @host localhost:8090
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	app.Run()
}
