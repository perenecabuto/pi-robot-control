package handler

import (
	"encoding/json"
	"net/http"

	"../device"
)

type CompassPayload struct {
	Direction device.Direction `json:"direction"`
	X, Y, Z   int
}

var compass *device.Compass

func CompassHandler(w http.ResponseWriter, r *http.Request) {
	if compass == nil {
		compass = device.NewCompass(1, 0x1e, 1.3)
	}
	data, err := compass.Read()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	payload := &CompassPayload{device.North, data.X, data.Y, data.Z}
	if resp, err := json.Marshal(payload); err == nil {
		w.Write(resp)
	} else {
		http.Error(w, err.Error(), 500)
	}
}
