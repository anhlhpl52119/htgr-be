package api

import (
	"log"
	"net/http"
)

type TodoHandler struct {
	logger *log.Logger
}

func NewTodoHandler(logger *log.Logger) *TodoHandler {
	return &TodoHandler{
		logger: logger,
	}
}

func (h TodoHandler) HandleCreateTodo(w http.ResponseWriter, r *http.Request) {

}
