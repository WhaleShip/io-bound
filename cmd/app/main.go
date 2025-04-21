package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/whaleship/io-bound/internal/config"
	"github.com/whaleship/io-bound/internal/database"
	"github.com/whaleship/io-bound/internal/handlers"
	"github.com/whaleship/io-bound/internal/service"
)

func main() {
	app := fiber.New()
	cfg := config.Load()
	rConn := database.GetInitRedis(cfg)
	tskSvc := service.NewTaskService(rConn)
	tskHandl := handlers.NewTaskHandler(tskSvc)

	app.Use(logger.New())

	app.Post("/tasks", tskHandl.HandleCreateTask)
	app.Get("/tasks/:id", tskHandl.HandleGetTask)

	if err := app.Listen(":8080"); err != nil {
		log.Fatalln(err)
	}
}
