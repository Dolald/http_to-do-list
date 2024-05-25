package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v3"
)

func main() {
	app := fiber.New()
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Fatalln(app.Listen(fmt.Sprintf(":%v", port)))
}
