package test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"todo-app-api/helper"
	"todo-app-api/models/domain"
	"todo-app-api/models/web"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestPanicIfError(t *testing.T) {

	helper.PanicIfError(nil)

	defer func() {
		r := recover()
		if r == nil {
			t.Fatalf("expected panic")
		}

		if ferr, ok := r.(*fiber.Error); ok {
			assert.Equal(t, fiber.StatusBadRequest, ferr.Code)
			assert.Equal(t, "boom", ferr.Message)
			return
		}
		t.Fatalf("unexpected panic type: %T", r)
	}()

	helper.PanicIfError(errors.New("boom"))
}

func TestToTodoResponseHelpers(t *testing.T) {
	todo := domain.Todo{Id: 10, Title: "T", Description: "D", Status: "pending"}
	resp := helper.ToTodoResponse(todo)
	assert.Equal(t, 10, resp.Id)
	assert.Equal(t, "T", resp.Title)

	todos := []domain.Todo{todo, domain.Todo{Id: 11, Title: "T2", Description: "D2", Status: "done"}}
	respList := helper.ToTodoResponses(todos)
	assert.Len(t, respList, 2)
	assert.Equal(t, "T2", respList[1].Title)
}

func TestResponseHelpersAndReadFromRequestBody(t *testing.T) {
	app := fiber.New()

	app.Post("/parse", func(c *fiber.Ctx) error {
		var req web.TodoCreateRequest
		if err := helper.ReadFromRequestBody(c, &req); err != nil {
			return helper.BadRequest(c, err.Error())
		}
		return helper.ResponseSuccess(c, req)
	})

	// valid JSON
	body, _ := json.Marshal(web.TodoCreateRequest{Title: "X", Description: "D"})
	req := httptest.NewRequest(http.MethodPost, "/parse", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, -1)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// invalid JSON
	req2 := httptest.NewRequest(http.MethodPost, "/parse", bytes.NewReader([]byte(`{invalid`)))
	req2.Header.Set("Content-Type", "application/json")
	r2, _ := app.Test(req2, -1)
	assert.Equal(t, http.StatusBadRequest, r2.StatusCode)
}

func TestCommitOrRollback(t *testing.T) {
	// setup DB
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if err := db.AutoMigrate(&domain.Todo{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	// commit case
	tx := db.Begin()
	tx.Create(&domain.Todo{Title: "C", Description: "c"})
	helper.CommitOrRollback(tx)
	var count int64
	db.Model(&domain.Todo{}).Count(&count)
	assert.EqualValues(t, 1, count)

	// rollback case: simulate panic inside function with deferred CommitOrRollback
	func() {
		tx2 := db.Begin()
		defer func() {
			if r := recover(); r != nil {
			}
		}()
		defer helper.CommitOrRollback(tx2)
		tx2.Create(&domain.Todo{Title: "R", Description: "r"})
		panic("boom")
	}()

	db.Model(&domain.Todo{}).Count(&count)

	assert.EqualValues(t, 1, count)
}
