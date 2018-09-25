package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/perenecabuto/robot-control/device"
)

// CompassPayload represents the compass message
type CompassPayload struct {
	Direction device.Direction `json:"direction"`
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
		log.Println(err)
		http.Error(w, err.Error(), 500)
		return
	}
	payload := &CompassPayload{device.North, data.X, data.Y, data.Z}
	if resp, err := json.Marshal(payload); err == nil {
		w.Header().Set("Content-Type", "application/json")
		if i, err := w.Write(resp); err != nil {
			log.Println("Error - writen(", i, ")", err)
		}
	} else {
		http.Error(w, err.Error(), 500)
	}
}
