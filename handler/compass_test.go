package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"../device"
)

type FakeCompass struct {
	X, Y, Z int
}

func (c FakeCompass) Read() (*device.CompassData, error) {
	return &device.CompassData{c.X, c.Y, c.Z}, nil
}

func TestCompassHandler(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/compass", nil)

	handler := NewCompassHandler(&FakeCompass{1, 2, 3})
	handler.ServeHTTP(rec, req)

	var resp []byte
	var err error
	if resp, err = ioutil.ReadAll(rec.Body); err != nil {
		t.Fatal("Error reading response", string(resp))
	}
	if rec.Code != http.StatusOK {
		t.Fatal("Status code should be:", http.StatusOK, "returned:", rec.Code, "resp:", resp)
	}
	payload := &CompassPayload{}
	if err = json.Unmarshal(resp, payload); err != nil {
		t.Fatal("Compass payload error:", err, "response:", string(resp))
	}
	if payload.Direction == "" {
		t.Fatal("Direction is empty")
	}
	t.Log("Direction", payload.Direction)
}
