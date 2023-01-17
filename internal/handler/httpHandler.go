package handler

import (
	"fmt"
	"log"
	"net/http"
)

type HttpMethodHandler interface {
	Get(w http.ResponseWriter, req *http.Request)
	Post(w http.ResponseWriter, req *http.Request)
	Put(w http.ResponseWriter, req *http.Request)
}

type HttpHandler struct {
	methodHandlerByPath map[string]*HttpMethodHandler
	server              *http.ServeMux
}

func (handler *HttpHandler) Init() {
	handler.server = http.NewServeMux()
	handler.methodHandlerByPath = make(map[string]*HttpMethodHandler)
}

func (handler *HttpHandler) StartServer(port int) {
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), handler.server); err != nil {
		log.Fatal(err)
	}
}

func (handler *HttpHandler) RegisterHandler(path string, httpMethodHandler HttpMethodHandler) {
	_, exists := handler.methodHandlerByPath[path]
	if !exists {
		handler.methodHandlerByPath[path] = &httpMethodHandler
		handler.server.HandleFunc(path, handler.Handle)
	}
}

func (handler *HttpHandler) Handle(w http.ResponseWriter, req *http.Request) {
	methodHandler, exists := handler.methodHandlerByPath[req.URL.Path]
	if !exists {
		http.Error(w, "resource not found", 404)
		return
	}
	switch req.Method {
	case "PUT":
		(*methodHandler).Put(w, req)
	case "POST":
		(*methodHandler).Post(w, req)
	case "GET":
		(*methodHandler).Get(w, req)
	default:
		http.Error(w, fmt.Sprintf("invalid method: %s", req.Method), 500)
	}
}
