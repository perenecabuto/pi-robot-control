package handler

import (
	"errors"
	"log"
	"net/http"

	"../device"

	"golang.org/x/net/websocket"
)

type RobotHandler struct {
	robot *device.Robot
}

func NewRobotHandler(r *device.Robot) *RobotHandler {
	return &RobotHandler{r}
}

func (h RobotHandler) ListenWS() websocket.Handler {
	return websocket.Handler(func(ws *websocket.Conn) {
		log.Println("New WS connection from:" + ws.Request().Host)
		defer ws.Close()
		for {
			if _, err := ws.Write([]byte("")); err != nil {
				log.Println("Error reading ws")
				break
			}
			var msg string
			if err := websocket.Message.Receive(ws, &msg); err != nil {
				log.Println("Error reading ws")
				break
			}
			h.parseAction(msg)
		}
		log.Println("Close WS connection from:" + ws.Request().Host)
	})
}

func (h RobotHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	action := r.URL.Query().Get("direction")
	err := h.parseAction(action)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func (h RobotHandler) parseAction(action string) (err error) {
	log.Println("Got move action:", action)
	defer func() {
		if r := recover(); r != nil {
			switch r.(type) {
			case string:
				log.Println(err)
				err = errors.New(r.(string))
			default:
				err = errors.New("Unknown error parsing" + action)
			}
		}
	}()
	switch action {
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
		return errors.New("action " + action + " is not allowed")
	}
	return nil
}
