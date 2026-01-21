package helper

import (
	"net/http/httptest"
	"testing"
)

func TestHelper_RespondJSON(t *testing.T) {
	rr := httptest.NewRecorder()
	data := map[string]string{"message": "success"}
	RespondJSON(rr, 200, data)
}

func TestHelper_RespondError(t *testing.T) {
	rr := httptest.NewRecorder()
	RespondError(rr, 400, "bad request")
}