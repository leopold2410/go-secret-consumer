package handler

import (
	"encoding/json"
	"net/http"
)

type EchoHandler struct{}

func (srv *EchoHandler) Get(w http.ResponseWriter, req *http.Request) {
	http.Error(w, "not implemented", 500)
}

func (srv *EchoHandler) Put(w http.ResponseWriter, req *http.Request) {
	http.Error(w, "not implemented", 500)
}

func (srv *EchoHandler) Post(w http.ResponseWriter, req *http.Request) {
	payload, err := ParseRequestPayload(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payload)
}
