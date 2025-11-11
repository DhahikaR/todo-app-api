package service

import (
	"context"
	"todo-app-api/models/web"
)

type TodoService interface {
	Create(context context.Context, request web.TodoCreateRequest) web.TodoResponse
	Update(context context.Context, request web.TodoUpdateRequest) web.TodoResponse
	Delete(context context.Context, todoId int)
	FindById(context context.Context, todoId int) web.TodoResponse
	FindAll(context context.Context) []web.TodoResponse
}
