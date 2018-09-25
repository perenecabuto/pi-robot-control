package handler

import (
	"io/ioutil"
	"net/http"
)

// IndexHandler renders robot control web page
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	file, err := ioutil.ReadFile("webapp/index.html")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(file)
}
