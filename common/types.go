package common

import (
	"net/http"
)

type Error struct {
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors"`
}

func (e Error) Error() string {
	return e.Message
}

type ApiFunc func(w http.ResponseWriter, r *http.Request) error
