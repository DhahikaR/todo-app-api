package repository

import (
	"context"
	"todo-app-api/models/domain"

	"gorm.io/gorm"
)

type TodoRepositoryImpl struct {
	DB *gorm.DB
}

func NewTodoRepository(db *gorm.DB) TodoRepository {
	return &TodoRepositoryImpl{
		DB: db,
	}
}

func (repository *TodoRepositoryImpl) Save(ctx context.Context, tx *gorm.DB, todo domain.Todo) domain.Todo {
	tx.WithContext(ctx).Create(&todo)
	return todo
}

func (repository *TodoRepositoryImpl) Update(ctx context.Context, tx *gorm.DB, todo domain.Todo) domain.Todo {
	tx.WithContext(ctx).Save(&todo)
	return todo
}

func (repository *TodoRepositoryImpl) Delete(ctx context.Context, tx *gorm.DB, todo domain.Todo) {
	tx.WithContext(ctx).Delete(&todo)
}

func (repository *TodoRepositoryImpl) FindById(ctx context.Context, tx *gorm.DB, todoId int) (domain.Todo, error) {
	var todo domain.Todo
	result := tx.WithContext(ctx).First(&todo, todoId)

	if result.Error != nil {
		return todo, result.Error
	}

	return todo, result.Error
}

func (repository *TodoRepositoryImpl) FindAll(ctx context.Context, tx *gorm.DB) []domain.Todo {
	var todos []domain.Todo
	tx.WithContext(ctx).Order("id ASC").Find(&todos)
	return todos
}
