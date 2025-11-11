package main

import (
	"log"
	"os"
	"todo-app-api/config"
	"todo-app-api/controller"
	"todo-app-api/exception"
	"todo-app-api/repository"
	"todo-app-api/routes"
	"todo-app-api/service"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	app := fiber.New(fiber.Config{
		ErrorHandler: exception.NewErrorHandler,
	})

	app.Use(recover.New())

	db := config.NewDB()
	validate := validator.New()

	todoRepository := repository.NewTodoRepository(db)
	todoService := service.NewTodoService(todoRepository, db, validate)
	todoController := controller.NewTodoController(todoService)

	routes.NewRouter(app, todoController)

	app.Listen(":" + os.Getenv("APP_PORT"))

}
