package routes

import (
	"todo-app-api/controller"

	"github.com/gofiber/fiber/v2"
)

func NewRouter(app *fiber.App, todoController controller.TodoController) {
	todo := app.Group("/todos")

	todo.Get("/", todoController.FindAll)
	todo.Get("/:todoId", todoController.FindById)
	todo.Post("/", todoController.Create)
	todo.Put("/:todoId", todoController.Update)
	todo.Delete("/:todoId", todoController.Delete)
}
