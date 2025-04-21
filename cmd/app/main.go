package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/whaleship/io-bound/internal/handlers"
	"github.com/whaleship/io-bound/internal/service"
)

func main() {
	app := fiber.New()

	tskSvc := service.NewTaskService()
	tskHandl := handlers.NewTaskHandler(tskSvc)

	app.Post("/tasks", tskHandl.HandleCreateTask)

	app.Get("/tasks/:id", tskHandl.HandleGetTask)

	if err := app.Listen(":3000"); err != nil {
		log.Fatalln(err)
	}
}
