package emit_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thisisthemurph/emit"
)

func TestResponseBuilder_Status(t *testing.T) {
	rr := httptest.NewRecorder()
	rb := emit.New(rr).Status(http.StatusCreated)
	rb.Text("Created")

	assert.Equal(t, http.StatusCreated, rr.Code)
}

func TestResponseBuilder_Header(t *testing.T) {
	rr := httptest.NewRecorder()
	rb := emit.New(rr).Header("X-Custom-Header", "Value")
	rb.Text("OK")

	assert.Equal(t, "Value", rr.Header().Get("X-Custom-Header"))
}

func TestResponseBuilder_WithCookie(t *testing.T) {
	rr := httptest.NewRecorder()
	cookie := &http.Cookie{Name: "session", Value: "abc123"}

	emit.New(rr).Cookie(cookie).Text("OK")
	cookies := rr.Result().Cookies()

	assert.Len(t, cookies, 1)
	assert.Equal(t, "session", cookies[0].Name)
	assert.Equal(t, "abc123", cookies[0].Value)
}

func TestResponseBuilder_Text(t *testing.T) {
	rr := httptest.NewRecorder()
	emit.New(rr).Text("Hello, world!")

	assert.Equal(t, "text/plain", rr.Header().Get("Content-Type"))
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "Hello, world!", rr.Body.String())
}

func TestResponseBuilder_Text_WithSetStatus(t *testing.T) {
	rr := httptest.NewRecorder()
	emit.New(rr).Status(http.StatusAccepted).Text("Hello, world!")

	assert.Equal(t, "text/plain", rr.Header().Get("Content-Type"))
	assert.Equal(t, http.StatusAccepted, rr.Code)
	assert.Equal(t, "Hello, world!", rr.Body.String())
}

func TestResponseBuilder_JSON(t *testing.T) {
	rr := httptest.NewRecorder()
	data := map[string]string{"message": "Hello"}

	emit.New(rr).JSON(data)

	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
	assert.Equal(t, http.StatusOK, rr.Code)

	var responseData map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &responseData)
	assert.NoError(t, err)
	assert.Equal(t, "Hello", responseData["message"])
}

func TestResponseBuilder_JSON_WithSetStatus(t *testing.T) {
	rr := httptest.NewRecorder()
	data := map[string]string{"message": "Hello"}

	emit.New(rr).Status(http.StatusCreated).JSON(data)

	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
	assert.Equal(t, http.StatusCreated, rr.Code)

	var responseData map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &responseData)
	assert.NoError(t, err)
	assert.Equal(t, "Hello", responseData["message"])
}

func TestResponseBuilder_ErrorJSON(t *testing.T) {
	rr := httptest.NewRecorder()
	emit.New(rr).ErrorJSON("Something went wrong")

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	var responseData map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &responseData)
	assert.NoError(t, err)
	assert.Equal(t, "Something went wrong", responseData["error"])
}

func TestResponseBuilder_ErrorJSON_WithSetStatus(t *testing.T) {
	rr := httptest.NewRecorder()
	emit.New(rr).Status(http.StatusBadRequest).ErrorJSON("Something went wrong")

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	var responseData map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &responseData)
	assert.NoError(t, err)
	assert.Equal(t, "Something went wrong", responseData["error"])
}

func TestResponseBuilder_NoContent(t *testing.T) {
	rr := httptest.NewRecorder()
	emit.New(rr).NoContent()

	assert.Equal(t, http.StatusNoContent, rr.Code)
	assert.Empty(t, rr.Body.String())
}
