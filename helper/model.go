package helper

import (
	"todo-app-api/models/domain"
	"todo-app-api/models/web"
)

func ToTodoResponse(todo domain.Todo) web.TodoResponse {
	return web.TodoResponse{
		Id:          int(todo.Id),
		Title:       todo.Title,
		Description: todo.Description,
		Status:      todo.Status,
	}
}

func ToTodoResponses(todos []domain.Todo) []web.TodoResponse {
	var todoResponses []web.TodoResponse
	for _, todo := range todos {
		todoResponses = append(todoResponses, ToTodoResponse(todo))
	}

	return todoResponses
}
