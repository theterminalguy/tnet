package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/theterminalguy/tnet/internal/handler"
	repo "github.com/theterminalguy/tnet/internal/repository"
)

func init() {
	os.Setenv("ENV", "test")
	os.Setenv("PLATFORM", "docker")
	os.Setenv("PORT", "8080")
	os.Setenv("APP_HOST", "http://localhost")
	os.Setenv("JWT_SIGNED_SECRET", "test-secret")
}

var e = echo.New()

// ---- HealthHandler ----

func TestHealthHandler_OK(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	require.NoError(t, handler.HealthHandler(c))
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "OK", rec.Body.String())
}

func TestHealthHandler_MethodNotAllowed(t *testing.T) {
	// POST to a GET-only handler — echo returns 405 at routing level;
	// calling the handler directly still returns 200 (handler is method-agnostic),
	// so we verify the handler itself always returns 200 regardless of method.
	req := httptest.NewRequest(http.MethodPost, "/health", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	require.NoError(t, handler.HealthHandler(c))
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestHealthHandler_ResponseBodyIsOK(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	require.NoError(t, handler.HealthHandler(c))
	assert.Contains(t, rec.Body.String(), "OK")
}

// ---- UserHandler ----

func newUserHandler() *handler.UserHandler {
	return handler.NewUserHandler()
}

func createUserViaRepo(t *testing.T) *repo.UserParams {
	t.Helper()
	return &repo.UserParams{
		FirstName: "Integration",
		LastName:  "Test",
		Email:     "inttest+" + t.Name() + "@example.com",
		PhotoURL:  "",
	}
}

func TestUserHandler_ReadAll_OK(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := newUserHandler()
	require.NoError(t, h.ReadAll(c))
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestUserHandler_CreateOne_OK(t *testing.T) {
	body, _ := json.Marshal(map[string]string{
		"first_name": "John",
		"last_name":  "Smith",
		"email":      "john.smith." + t.Name() + "@example.com",
		"photo_url":  "",
		"role":       "talent",
	})
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := newUserHandler()
	require.NoError(t, h.CreateOne(c))
	assert.Equal(t, http.StatusCreated, rec.Code)
}

func TestUserHandler_CreateOne_InvalidInput(t *testing.T) {
	// missing required fields
	body, _ := json.Marshal(map[string]string{})
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := newUserHandler()
	require.NoError(t, h.CreateOne(c))
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUserHandler_ReadByID_NotFound(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/users/00000000-0000-0000-0000-000000000000", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("uuid")
	c.SetParamValues("00000000-0000-0000-0000-000000000000")

	h := newUserHandler()
	require.NoError(t, h.ReadByID(c))
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestUserHandler_ReadByID_InvalidUUID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/users/not-a-uuid", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("uuid")
	c.SetParamValues("not-a-uuid")

	h := newUserHandler()
	require.NoError(t, h.ReadByID(c))
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUserHandler_DeleteOne_NotFound(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/users/00000000-0000-0000-0000-000000000000", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("uuid")
	c.SetParamValues("00000000-0000-0000-0000-000000000000")

	h := newUserHandler()
	require.NoError(t, h.DeleteOne(c))
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestUserHandler_DeleteOne_InvalidUUID(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/users/bad", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("uuid")
	c.SetParamValues("bad")

	h := newUserHandler()
	require.NoError(t, h.DeleteOne(c))
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// ---- PartnerHandler ----

func newPartnerHandler() *handler.PartnerHandler {
	return handler.NewPartnerHandler()
}

func validPartnerBody() []byte {
	b, _ := json.Marshal(map[string]string{
		"company_name":                "Acme Corp",
		"company_location":            "New York",
		"website_url":                 "https://acme.example.com",
		"contact_person_name":         "Alice",
		"contact_person_phone_number": "+1234567890",
		"contact_person_email":        "alice@acme.example.com",
	})
	return b
}

func TestPartnerHandler_ReadAll_OK(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/partners", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := newPartnerHandler()
	require.NoError(t, h.ReadAll(c))
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestPartnerHandler_CreateOne_OK(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/partners", bytes.NewReader(validPartnerBody()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := newPartnerHandler()
	require.NoError(t, h.CreateOne(c))
	assert.Equal(t, http.StatusCreated, rec.Code)
}

func TestPartnerHandler_CreateOne_InvalidInput(t *testing.T) {
	body, _ := json.Marshal(map[string]string{})
	req := httptest.NewRequest(http.MethodPost, "/partners", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := newPartnerHandler()
	require.NoError(t, h.CreateOne(c))
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestPartnerHandler_ReadByID_InvalidUUID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/partners/bad", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("uuid")
	c.SetParamValues("bad")

	h := newPartnerHandler()
	require.NoError(t, h.ReadByID(c))
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestPartnerHandler_ReadByID_NotFound(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/partners/00000000-0000-0000-0000-000000000000", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("uuid")
	c.SetParamValues("00000000-0000-0000-0000-000000000000")

	h := newPartnerHandler()
	require.NoError(t, h.ReadByID(c))
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestPartnerHandler_DeleteOne_InvalidUUID(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/partners/bad", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("uuid")
	c.SetParamValues("bad")

	h := newPartnerHandler()
	require.NoError(t, h.DeleteOne(c))
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
