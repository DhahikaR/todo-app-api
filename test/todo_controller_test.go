package test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"todo-app-api/controller"
	"todo-app-api/exception"
	"todo-app-api/models/web"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTodoService struct {
	mock.Mock
}

func (m *MockTodoService) Create(context context.Context, request web.TodoCreateRequest) web.TodoResponse {
	args := m.Called(context, request)
	return args.Get(0).(web.TodoResponse)
}

func (m *MockTodoService) Update(context context.Context, request web.TodoUpdateRequest) web.TodoResponse {
	args := m.Called(context, request)
	return args.Get(0).(web.TodoResponse)
}

func (m *MockTodoService) Delete(context context.Context, todoId int) {
	m.Called(context, todoId)
}

func (m *MockTodoService) FindById(context context.Context, todoId int) web.TodoResponse {
	args := m.Called(context, todoId)
	return args.Get(0).(web.TodoResponse)
}

func (m *MockTodoService) FindAll(context context.Context) []web.TodoResponse {
	args := m.Called(context)
	return args.Get(0).([]web.TodoResponse)
}

func setupFiberApp(todoController controller.TodoController) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: exception.NewErrorHandler,
	})
	app.Post("/todos", todoController.Create)
	app.Put("/todos/:todoId", todoController.Update)
	app.Delete("/todos/:todoId", todoController.Delete)
	app.Get("/todos", todoController.FindAll)
	app.Get("/todos/:todoId", todoController.FindById)

	return app
}

func TestControllerCreateSuccess(t *testing.T) {
	mockService := new(MockTodoService)
	todoController := controller.NewTodoController(mockService)
	app := setupFiberApp(todoController)

	requestBody := web.TodoCreateRequest{
		Title:       "Test Controller",
		Description: "Description Test",
	}

	requestJSON, _ := json.Marshal(requestBody)

	expected := web.TodoResponse{
		Id:          1,
		Title:       "Test Controller",
		Description: "Description Test",
		Status:      "pending",
	}

	mockService.On("Create", mock.Anything, requestBody).Return(expected)

	request := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewReader(requestJSON))
	request.Header.Set("Content-Type", "application/json")

	response, _ := app.Test(request, -1)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	mockService.AssertExpectations(t)
}

func TestControllerCreateFailed(t *testing.T) {
	mockService := new(MockTodoService)
	todoController := controller.NewTodoController(mockService)
	app := setupFiberApp(todoController)

	request := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewBufferString(`{invalid_json}`))
	request.Header.Set("Content-Type", "application/json")

	response, _ := app.Test(request, -1)
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
}

func TestControllerUpdateSuccess(t *testing.T) {
	mockService := new(MockTodoService)
	todoController := controller.NewTodoController(mockService)
	app := setupFiberApp(todoController)

	requestBody := web.TodoUpdateRequest{
		Id:          1,
		Title:       "Test Controller New",
		Description: "Description Test New",
		Status:      "done",
	}
	requestJSON, _ := json.Marshal(requestBody)

	expected := web.TodoResponse{
		Id:          1,
		Title:       "Test Controller New",
		Description: "Description Test New",
		Status:      "done",
	}
	mockService.On("Update", mock.Anything, requestBody).Return(expected)

	request := httptest.NewRequest(http.MethodPut, "/todos/1", bytes.NewReader(requestJSON))
	request.Header.Set("Content-Type", "application/json")

	response, _ := app.Test(request, -1)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	mockService.AssertExpectations(t)
}

func TestControllerUpdateFailed(t *testing.T) {
	mockService := new(MockTodoService)
	todoController := controller.NewTodoController(mockService)
	app := setupFiberApp(todoController)

	request := httptest.NewRequest(http.MethodPut, "/todos/abc", nil)
	response, _ := app.Test(request, -1)
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
}

func TestControllerFindByIdSuccess(t *testing.T) {
	mockService := new(MockTodoService)
	todoController := controller.NewTodoController(mockService)
	app := setupFiberApp(todoController)

	expected := web.TodoResponse{
		Id:          1,
		Title:       "Tes Controller",
		Description: "Description Test",
		Status:      "pending",
	}
	mockService.On("FindById", mock.Anything, 1).Return(expected)

	request := httptest.NewRequest(http.MethodGet, "/todos/1", nil)
	response, _ := app.Test(request, -1)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	mockService.AssertExpectations(t)
}

func TestControllerFindByIdFailed(t *testing.T) {
	mockService := new(MockTodoService)
	todoController := controller.NewTodoController(mockService)
	app := setupFiberApp(todoController)

	req := httptest.NewRequest(http.MethodGet, "/todos/abc", nil)
	resp, _ := app.Test(req, -1)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestControllerFindAllSuccess(t *testing.T) {
	mockService := new(MockTodoService)
	todoController := controller.NewTodoController(mockService)
	app := setupFiberApp(todoController)

	expected := []web.TodoResponse{
		{Id: 1,
			Title:       "Tes Controller",
			Description: "Description Test",
			Status:      "pending",
		},
		{
			Id:          1,
			Title:       "Tes Controller 2",
			Description: "Description Test 2",
			Status:      "done",
		},
	}
	mockService.On("FindAll", mock.Anything).Return(expected)

	request := httptest.NewRequest(http.MethodGet, "/todos", nil)
	response, _ := app.Test(request, -1)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	mockService.AssertExpectations(t)
}

func TestControllerDeletedSuccess(t *testing.T) {
	mockService := new(MockTodoService)
	todoController := controller.NewTodoController(mockService)
	app := setupFiberApp(todoController)

	mockService.On("Delete", mock.Anything, 1).Return()

	request := httptest.NewRequest(http.MethodDelete, "/todos/1", nil)
	response, _ := app.Test(request, -1)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	mockService.AssertExpectations(t)
}

func TestControllerDeletedFailed(t *testing.T) {
	mockService := new(MockTodoService)
	todoController := controller.NewTodoController(mockService)
	app := setupFiberApp(todoController)

	mockService.On("Delete", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			panic(fiber.NewError(fiber.StatusNotFound, "todo not found"))
		})

	request := httptest.NewRequest(http.MethodDelete, "/todos/1", nil)
	response, _ := app.Test(request, -1)

	assert.Equal(t, http.StatusNotFound, response.StatusCode)
}
