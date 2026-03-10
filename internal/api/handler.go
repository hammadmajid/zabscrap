package api

import (
	"fmt"
	"log"
	"net/http"
)

type Handler struct {
	logger *log.Logger
}

func NewHandler(logger *log.Logger) Handler {
	return Handler{
		logger: logger,
	}
}

//goland:noinspection GoUnusedParameter
func (h Handler) Health(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "Status is available")
	if err != nil {
		h.logger.Println(err)
		return
	}
}
