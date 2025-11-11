package test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"todo-app-api/exception"
	"todo-app-api/models/web"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func setupApp() *fiber.App {
	return fiber.New(fiber.Config{ErrorHandler: exception.NewErrorHandler})
}

func decodeResponse(t *testing.T, response *http.Response) web.WebResponse {
	var webResponse web.WebResponse
	defer response.Body.Close()
	if err := json.NewDecoder(response.Body).Decode(&webResponse); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	return webResponse
}

func TestErrorHandler_ValidationError(t *testing.T) {
	app := setupApp()

	app.Get("/v", func(c *fiber.Ctx) error {
		v := validator.New()
		type R struct {
			Name string `validate:"required"`
		}
		err := v.Struct(R{})
		return err
	})

	req := httptest.NewRequest(http.MethodGet, "/v", nil)
	resp, _ := app.Test(req, -1)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	wr := decodeResponse(t, resp)
	assert.Equal(t, http.StatusBadRequest, wr.Code)
	assert.Equal(t, "BAD REQUEST", wr.Status)
}

func TestErrorHandler_NotFoundError(t *testing.T) {
	app := setupApp()

	app.Get("/nf", func(c *fiber.Ctx) error {
		return exception.NotFoundError{Message: "todo not found"}
	})

	req := httptest.NewRequest(http.MethodGet, "/nf", nil)
	resp, _ := app.Test(req, -1)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	wr := decodeResponse(t, resp)
	assert.Equal(t, http.StatusNotFound, wr.Code)
	assert.Equal(t, "NOT FOUND", wr.Status)
	// ensure the custom message is forwarded in the response data
	if msg, ok := wr.Data.(string); ok {
		assert.Equal(t, "todo not found", msg)
	} else {
		t.Fatalf("expected Data to be string, got %T", wr.Data)
	}
}

func TestErrorHandler_ValidationMultipleFields(t *testing.T) {
	app := setupApp()

	app.Get("/vm", func(c *fiber.Ctx) error {
		v := validator.New()
		type R struct {
			Name     string `validate:"required"`
			Password string `validate:"required,min=6"`
		}
		err := v.Struct(R{})
		return err
	})

	req := httptest.NewRequest(http.MethodGet, "/vm", nil)
	resp, _ := app.Test(req, -1)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	wr := decodeResponse(t, resp)
	assert.Equal(t, http.StatusBadRequest, wr.Code)
	assert.Equal(t, "BAD REQUEST", wr.Status)
	// Data should contain validation messages for both fields
	if s, ok := wr.Data.(string); ok {
		assert.Contains(t, s, "Name")
		assert.Contains(t, s, "Password")
	} else {
		t.Fatalf("expected Data to be string, got %T", wr.Data)
	}
}

func TestErrorHandler_FiberError(t *testing.T) {
	app := setupApp()

	app.Get("/fe", func(c *fiber.Ctx) error {
		return fiber.NewError(fiber.StatusNotFound, "not found")
	})

	req := httptest.NewRequest(http.MethodGet, "/fe", nil)
	resp, _ := app.Test(req, -1)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	wr := decodeResponse(t, resp)
	assert.Equal(t, http.StatusNotFound, wr.Code)
	assert.Equal(t, "NOT FOUND", wr.Status)
}

func TestErrorHandler_InternalServerError(t *testing.T) {
	app := setupApp()

	app.Get("/ie", func(c *fiber.Ctx) error {
		return errors.New("boom")
	})

	req := httptest.NewRequest(http.MethodGet, "/ie", nil)
	resp, _ := app.Test(req, -1)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	wr := decodeResponse(t, resp)
	assert.Equal(t, http.StatusInternalServerError, wr.Code)
	assert.Equal(t, "INTERNAL SERVICE ERROR", wr.Status)
}
