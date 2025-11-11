package service

import (
	"context"
	"errors"
	"todo-app-api/exception"
	"todo-app-api/helper"
	"todo-app-api/models/domain"
	"todo-app-api/models/web"
	"todo-app-api/repository"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type TodoServiceImpl struct {
	TodoRepository repository.TodoRepository
	DB             *gorm.DB
	Validate       *validator.Validate
}

func NewTodoService(todoRepository repository.TodoRepository, DB *gorm.DB, validate *validator.Validate) TodoService {
	return &TodoServiceImpl{
		TodoRepository: todoRepository,
		DB:             DB,
		Validate:       validate,
	}
}

func (service *TodoServiceImpl) Create(ctx context.Context, request web.TodoCreateRequest) web.TodoResponse {
	err := service.Validate.Struct(request)
	helper.PanicIfError(err)

	tx := service.DB.Begin()
	defer helper.CommitOrRollback(tx)

	todo := domain.Todo{
		Title:       request.Title,
		Description: request.Description,
		Status:      request.Status,
	}

	if todo.Status == "" {
		todo.Status = "pending"
	}

	service.TodoRepository.Save(ctx, tx, todo)

	return helper.ToTodoResponse(todo)
}

func (service *TodoServiceImpl) Update(ctx context.Context, request web.TodoUpdateRequest) web.TodoResponse {
	err := service.Validate.Struct(request)
	helper.PanicIfError(err)

	tx := service.DB.Begin()
	defer helper.CommitOrRollback(tx)

	todo, err := service.TodoRepository.FindById(ctx, tx, request.Id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			panic(exception.NotFoundError{Message: "todo not found"})
		}
		panic(err)
	}

	todo.Title = request.Title
	todo.Description = request.Description
	todo.Status = request.Status

	todo = service.TodoRepository.Update(ctx, tx, todo)

	if todo.Status == "" {
		todo.Status = "pending" // default
	}

	return helper.ToTodoResponse(todo)
}

func (service *TodoServiceImpl) Delete(ctx context.Context, todoId int) {
	tx := service.DB.Begin()
	defer helper.CommitOrRollback(tx)

	todo, err := service.TodoRepository.FindById(ctx, tx, todoId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			panic(exception.NotFoundError{Message: "todo not found"})
		}
		panic(err)
	}

	service.TodoRepository.Delete(ctx, tx, todo)
}

func (service *TodoServiceImpl) FindById(ctx context.Context, todoId int) web.TodoResponse {
	tx := service.DB.Begin()
	defer helper.CommitOrRollback(tx)

	todo, err := service.TodoRepository.FindById(ctx, tx, todoId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			panic(exception.NotFoundError{Message: "todo not found"})
		}
		panic(err)
	}

	return helper.ToTodoResponse(todo)
}

func (service *TodoServiceImpl) FindAll(ctx context.Context) []web.TodoResponse {
	tx := service.DB.Begin()
	defer helper.CommitOrRollback(tx)

	todos := service.TodoRepository.FindAll(ctx, tx)

	return helper.ToTodoResponses(todos)
}
