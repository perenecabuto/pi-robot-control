package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/perenecabuto/robot-control/device"
)

// CompassPayload represents the compass message
type CompassPayload struct {
	Direction string `json:"direction"`
	Degress   int    `json:"degress"`
	X, Y, Z   int
}

// CompassHandler handle compass actions
type CompassHandler struct {
	compass device.Compass
}

// NewCompassHandler creates a new http handler for compass device
func NewCompassHandler(c device.Compass) *CompassHandler {
	return &CompassHandler{c}
}

func (h CompassHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data, err := h.compass.Read()
	if err != nil {
		log.Println("[CompassHandler] Error to read compass:", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	payload := &CompassPayload{string(data.Direction), data.Degress, data.X, data.Y, data.Z}
	resp, _ := json.Marshal(payload)
	if _, err := w.Write(resp); err != nil {
		log.Println("[CompassHandler] Error to write compass payload:", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
