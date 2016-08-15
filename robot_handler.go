package main

import (
	"log"
	"net/http"
)

type RobotHandler struct {
	robot *Robot
}

func NewRobotHandler(r *Robot) *RobotHandler {
	return &RobotHandler{r}
}

func (h RobotHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
			http.Error(w, "Unexpcted error", 500)
		}
	}()

	direction := r.URL.Query().Get("direction")
	switch direction {
	case "right":
		h.robot.Right()
	case "left":
		h.robot.Left()
	case "forward":
		h.robot.Forward()
	case "backward":
		h.robot.Backward()
	case "stop":
		h.robot.Stop()
	default:
		http.Error(w, "action "+direction+" is not allowed", 500)
		return
	}

	log.Println("Got move action:", direction)
}
