package handler

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/perenecabuto/robot-control/device"

	"golang.org/x/net/websocket"
)

// RobotHandler handles robot actions
type RobotHandler struct {
	robot *device.Robot
}

// NewRobotHandler create a new handler for the robot
func NewRobotHandler(r *device.Robot) *RobotHandler {
	return &RobotHandler{r}
}

// LookTo handles camera axis messages
func (h RobotHandler) LookTo() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		qs := r.URL.Query()
		x, _ := strconv.Atoi(qs.Get("x"))
		y, _ := strconv.Atoi(qs.Get("y"))
		if err := h.robot.Look.To(uint8(x), uint8(y)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Println("Look.To", x, y)
		w.Write([]byte("OK"))
	})
}

// ListenWS stabilish a websocket connection to robot controls
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
			if err := h.parseAction(msg); err != nil {
				log.Println("Error:", err)
			}
		}
		log.Println("Close WS connection from:" + ws.Request().Host)
	})
}

func (h RobotHandler) parseAction(msg string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			switch r.(type) {
			case string:
				log.Println(err)
				err = errors.New(r.(string))
			default:
				err = errors.New("Unknown error parsing" + msg)
			}
		}
	}()

	params := strings.Split(msg, ":")
	action := params[0]
	params = params[1:]

	switch action {
	case "move-right":
		h.robot.Move.Right()
	case "move-left":
		h.robot.Move.Left()
	case "move-forward":
		h.robot.Move.Forward()
	case "move-backward":
		h.robot.Move.Backward()
	case "move-stop":
		h.robot.Move.Stop()
	case "look-to":
		if len(params) < 2 {
			return errors.New("look-to must receive x and y")
		}
		angles := make([]uint8, 2)
		for i, p := range params[:2] {
			a, err := strconv.Atoi(p)
			if err != nil {
				return err
			}
			angles[i] = uint8(a)
		}
		return h.robot.Look.To(angles[0], angles[1])
	default:
		return errors.New("action " + action + " is not allowed")
	}

	return nil
}
