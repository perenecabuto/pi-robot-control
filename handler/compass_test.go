package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCompassHandler(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/compass", nil)

	CompassHandler(rec, req)

	var resp []byte
	if _, err := rec.Result().Body.Read(resp); err != nil {
		t.Fatal("Error reading response")
	}
	if rec.Code != http.StatusOK {
		t.Fatal("Status code should be:", http.StatusOK, "returned:", rec.Code, "resp:", resp)
	}
	payload := &CompassPayload{}
	if err := json.Unmarshal(resp, payload); err != nil {
		t.Fatal("Compass payload error:", err, "response:", string(resp))
	}
	if payload.Direction == "" {
		t.Fatal("Direction is empty")
	}
	t.Log("Direction", payload.Direction)
}
