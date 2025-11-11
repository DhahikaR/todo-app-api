package test

import (
	"context"
	"fmt"
	"testing"

	"todo-app-api/models/domain"
	"todo-app-api/repository"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// helper to create an in-memory gorm DB and migrate the Todo model
func setupTestDB(t *testing.T) *gorm.DB {
	// create a unique in-memory database per test to avoid cross-test pollution
	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	err = db.AutoMigrate(&domain.Todo{})
	if err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	return db
}



func TestTodoRepository_SaveAndFindById(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewTodoRepository(db)
	ctx := context.Background()

	tx := db.Begin()
	todo := domain.Todo{Title: "Test Repo", Description: "Repository test", Status: "pending"}
	saved := repo.Save(ctx, tx, todo)
	if err := tx.Commit().Error; err != nil {
		t.Fatalf("commit failed: %v", err)
	}

	assert.NotZero(t, saved.Id)

	found, err := repo.FindById(ctx, db, saved.Id)
	assert.NoError(t, err)
	assert.Equal(t, "Test Repo", found.Title)
}

func TestTodoRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewTodoRepository(db)
	ctx := context.Background()

	// create record
	tx := db.Begin()
	todo := domain.Todo{Title: "ToUpdate", Description: "desc", Status: "pending"}
	saved := repo.Save(ctx, tx, todo)
	if err := tx.Commit().Error; err != nil {
		t.Fatalf("commit failed: %v", err)
	}

	// update
	saved.Title = "Updated"
	tx2 := db.Begin()
	updated := repo.Update(ctx, tx2, saved)
	if err := tx2.Commit().Error; err != nil {
		t.Fatalf("commit failed: %v", err)
	}

	assert.Equal(t, "Updated", updated.Title)
}

func TestTodoRepository_FindAll(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewTodoRepository(db)
	ctx := context.Background()

	tx := db.Begin()
	repo.Save(ctx, tx, domain.Todo{Title: "One", Description: "d1", Status: "pending"})
	repo.Save(ctx, tx, domain.Todo{Title: "Two", Description: "d2", Status: "done"})
	if err := tx.Commit().Error; err != nil {
		t.Fatalf("commit failed: %v", err)
	}

	all := repo.FindAll(ctx, db)
	assert.Len(t, all, 2)
}

func TestTodoRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewTodoRepository(db)
	ctx := context.Background()

	tx := db.Begin()
	todo := domain.Todo{Title: "ToDelete", Description: "d", Status: "pending"}
	saved := repo.Save(ctx, tx, todo)
	if err := tx.Commit().Error; err != nil {
		t.Fatalf("commit failed: %v", err)
	}

	tx2 := db.Begin()
	repo.Delete(ctx, tx2, saved)
	if err := tx2.Commit().Error; err != nil {
		t.Fatalf("commit failed: %v", err)
	}

	_, err := repo.FindById(ctx, db, saved.Id)
	assert.Error(t, err)
}
