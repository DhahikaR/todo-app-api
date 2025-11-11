package controller

import (
	"fmt"
	"strconv"
	"todo-app-api/helper"
	"todo-app-api/models/web"
	"todo-app-api/service"

	"github.com/gofiber/fiber/v2"
)

type TodoControllerImpl struct {
	todoService service.TodoService
}

func NewTodoController(todoService service.TodoService) TodoController {
	return &TodoControllerImpl{
		todoService: todoService,
	}
}

func (controller *TodoControllerImpl) Create(c *fiber.Ctx) error {
	todoCreateRequest := web.TodoCreateRequest{}
	if err := helper.ReadFromRequestBody(c, &todoCreateRequest); err != nil {
		return helper.BadRequest(c, err.Error())
	}

	todoResponse := controller.todoService.Create(c.Context(), todoCreateRequest)
	return helper.ResponseSuccess(c, todoResponse)
}

func (controller *TodoControllerImpl) Update(c *fiber.Ctx) (err error) {
	todoUpdateRequest := web.TodoUpdateRequest{}
	if err := helper.ReadFromRequestBody(c, &todoUpdateRequest); err != nil {
		return helper.BadRequest(c, err.Error())
	}

	todoId := c.Params("todoId")
	id, errConv := strconv.Atoi(todoId)
	if errConv != nil {
		return helper.BadRequest(c, "todoId must be a number")
	}

	todoUpdateRequest.Id = id

	defer func() {
		if r := recover(); r != nil {
			switch e := r.(type) {
			case error:
				err = e
			default:
				err = fmt.Errorf("%v", e)
			}
		}
	}()

	todoResponse := controller.todoService.Update(c.Context(), todoUpdateRequest)
	return helper.ResponseSuccess(c, todoResponse)
}

func (controller *TodoControllerImpl) Delete(c *fiber.Ctx) (err error) {
	todoId := c.Params("todoId")
	id, errConv := strconv.Atoi(todoId)
	if errConv != nil {
		return helper.BadRequest(c, "todoId must be a number")
	}

	defer func() {
		if r := recover(); r != nil {
			switch e := r.(type) {
			case error:
				err = e
			default:
				err = fmt.Errorf("%v", e)
			}
		}
	}()

	controller.todoService.Delete(c.Context(), id)
	return c.Status(fiber.StatusOK).JSON(web.WebResponse{
		Code:   200,
		Status: "Success",
	})
}

func (controller *TodoControllerImpl) FindById(c *fiber.Ctx) (err error) {
	todoId := c.Params("todoId")
	id, errConv := strconv.Atoi(todoId)
	if errConv != nil {
		return helper.BadRequest(c, "todoId must a be number")
	}

	defer func() {
		if r := recover(); r != nil {
			switch e := r.(type) {
			case error:
				err = e
			default:
				err = fmt.Errorf("%v", e)
			}
		}
	}()

	todoResponse := controller.todoService.FindById(c.Context(), id)
	return helper.ResponseSuccess(c, todoResponse)
}

func (controller *TodoControllerImpl) FindAll(c *fiber.Ctx) error {
	todoResponse := controller.todoService.FindAll(c.Context())
	return helper.ResponseSuccess(c, todoResponse)
}
