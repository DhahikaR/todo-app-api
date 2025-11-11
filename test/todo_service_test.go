package test

import (
	"context"
	"testing"
	"todo-app-api/exception"
	"todo-app-api/models/domain"
	"todo-app-api/models/web"
	"todo-app-api/service"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type TodoRepositoryMock struct {
	mock.Mock
}

func (m *TodoRepositoryMock) Save(ctx context.Context, tx *gorm.DB, todo domain.Todo) domain.Todo {
	args := m.Called(ctx, tx, todo)
	return args.Get(0).(domain.Todo)
}

func (m *TodoRepositoryMock) Update(ctx context.Context, tx *gorm.DB, todo domain.Todo) domain.Todo {
	args := m.Called(ctx, tx, todo)
	return args.Get(0).(domain.Todo)
}

func (m *TodoRepositoryMock) Delete(ctx context.Context, tx *gorm.DB, todo domain.Todo) {
	m.Called(ctx, tx, todo)
}

func (m *TodoRepositoryMock) FindById(ctx context.Context, tx *gorm.DB, id int) (domain.Todo, error) {
	args := m.Called(ctx, tx, id)
	return args.Get(0).(domain.Todo), args.Error(1)
}

func (m *TodoRepositoryMock) FindAll(ctx context.Context, tx *gorm.DB) []domain.Todo {
	args := m.Called(ctx, tx)
	return args.Get(0).([]domain.Todo)
}

func TestServiceCreateSuccess(t *testing.T) {
	mockRepo := new(TodoRepositoryMock)

	db, err := gorm.Open(sqlite.Open(":memory"), &gorm.Config{})
	assert.Nil(t, err)

	validate := validator.New()

	request := web.TodoCreateRequest{
		Title:       "Test",
		Description: "Description Test",
	}
	expected := domain.Todo{
		Id:          1,
		Title:       "Test",
		Description: "Description Test",
		Status:      "pending",
	}
	mockRepo.On("Save", mock.Anything, mock.Anything, mock.AnythingOfType("domain.Todo")).Return(expected)

	todoService := service.NewTodoService(mockRepo, db, validate)
	result := todoService.Create(context.Background(), request)

	assert.Equal(t, "Test", result.Title)
	assert.Equal(t, "pending", result.Status)
	mockRepo.AssertExpectations(t)
}

func TestServiceCreateFailed(t *testing.T) {
	mockRepo := new(TodoRepositoryMock)
	validate := validator.New()
	db, _ := gorm.Open(sqlite.Open(":memory"), &gorm.Config{})
	todoService := service.NewTodoService(mockRepo, db, validate)

	request := web.TodoCreateRequest{
		Title: "",
	}

	assert.Panics(t, func() {
		todoService.Create(context.Background(), request)
	})
}

func TestServiceUpdateSuccess(t *testing.T) {
	mockRepo := new(TodoRepositoryMock)
	db, _ := gorm.Open(sqlite.Open(":memory"), &gorm.Config{})
	validate := validator.New()

	todoService := service.NewTodoService(mockRepo, db, validate)

	existing := domain.Todo{
		Id:          1,
		Title:       "Test",
		Description: "Description Test",
		Status:      "pending",
	}

	updated := domain.Todo{
		Id:          1,
		Title:       "New Test",
		Description: "Description Test New",
		Status:      "done",
	}

	request := web.TodoUpdateRequest{
		Id:          1,
		Title:       "New Test",
		Description: "Description Test New",
		Status:      "done",
	}

	mockRepo.On("FindById", mock.Anything, mock.Anything, 1).Return(existing, nil)
	mockRepo.On("Update", mock.Anything, mock.Anything, mock.AnythingOfType("domain.Todo")).Return(updated)

	result := todoService.Update(context.Background(), request)

	assert.Equal(t, updated.Id, result.Id)
	assert.Equal(t, "New Test", result.Title)
	assert.Equal(t, "Description Test New", result.Description)
	assert.Equal(t, "done", result.Status)

	mockRepo.AssertExpectations(t)
}

func TestServiceUpdateFailed(t *testing.T) {
	mockRepo := new(TodoRepositoryMock)
	db, _ := gorm.Open(sqlite.Open(":memory"), &gorm.Config{})
	validate := validator.New()

	request := web.TodoUpdateRequest{
		Id:          99,
		Title:       "Failed Title",
		Description: "Test Description",
		Status:      "done",
	}
	mockRepo.On("FindById", mock.Anything, mock.Anything, 99).Return(domain.Todo{}, gorm.ErrRecordNotFound)

	todoService := service.NewTodoService(mockRepo, db, validate)

	assert.PanicsWithValue(t, exception.NotFoundError{Message: "todo not found"}, func() {
		todoService.Update(context.Background(), request)
	})
}

func TestServiceDeleteSuccess(t *testing.T) {
	mockRepo := new(TodoRepositoryMock)
	db, _ := gorm.Open(sqlite.Open(":memory"), &gorm.Config{})
	validate := validator.New()

	todoService := service.NewTodoService(mockRepo, db, validate)

	existing := domain.Todo{
		Id:          1,
		Title:       "Test",
		Description: "Description Test",
		Status:      "pending",
	}

	mockRepo.On("FindById", mock.Anything, mock.Anything, 1).Return(existing, nil)
	mockRepo.On("Delete", mock.Anything, mock.Anything, mock.AnythingOfType("domain.Todo")).Return()

	todoService.Delete(context.Background(), 1)

	mockRepo.AssertExpectations(t)
}

func TestServiceDeleteFailed(t *testing.T) {
	mockRepo := new(TodoRepositoryMock)
	db, _ := gorm.Open(sqlite.Open(":memory"), &gorm.Config{})
	validate := validator.New()

	todoService := service.NewTodoService(mockRepo, db, validate)

	mockRepo.On("FindById", mock.Anything, mock.Anything, 99).Return(domain.Todo{}, gorm.ErrRecordNotFound)

	assert.PanicsWithValue(t, exception.NotFoundError{Message: "todo not found"}, func() {
		todoService.Delete(context.Background(), 99)
	})
}

func TestServiceFindByIdSuccess(t *testing.T) {
	mockRepo := new(TodoRepositoryMock)
	db, _ := gorm.Open(sqlite.Open(":memory"), &gorm.Config{})
	validate := validator.New()

	todoService := service.NewTodoService(mockRepo, db, validate)

	existing := domain.Todo{
		Id:          1,
		Title:       "Test",
		Description: "Description Test",
		Status:      "pending",
	}

	mockRepo.On("FindById", mock.Anything, mock.Anything, 1).Return(existing, nil)

	result := todoService.FindById(context.Background(), 1)
	assert.Equal(t, existing.Id, result.Id)
	assert.Equal(t, "Test", result.Title)
	assert.Equal(t, "Description Test", result.Description)
	assert.Equal(t, "pending", result.Status)

	mockRepo.AssertExpectations(t)
}

func TestServiceFindByIdFailed(t *testing.T) {
	mockRepo := new(TodoRepositoryMock)
	db, _ := gorm.Open(sqlite.Open(":memory"), &gorm.Config{})
	validate := validator.New()

	todoService := service.NewTodoService(mockRepo, db, validate)

	mockRepo.On("FindById", mock.Anything, mock.Anything, 99).Return(domain.Todo{}, gorm.ErrRecordNotFound)

	assert.PanicsWithValue(t, exception.NotFoundError{Message: "todo not found"}, func() {
		todoService.FindById(context.Background(), 99)
	})
}

func TestServiceFindAllSuccess(t *testing.T) {
	mockRepo := new(TodoRepositoryMock)
	db, _ := gorm.Open(sqlite.Open(":memory"), &gorm.Config{})
	validate := validator.New()

	todoService := service.NewTodoService(mockRepo, db, validate)

	existing := []domain.Todo{}

	mockRepo.On("FindAll", mock.Anything, mock.Anything).Return(existing)

	todoService.FindAll(context.Background())

	assert.Len(t, existing, 0)

	mockRepo.AssertExpectations(t)
}
