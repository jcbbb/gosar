package common

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	ErrInternal         = Error{Message: "Internal server error", Code: http.StatusInternalServerError}
	ErrMethodNotAllowed = Error{Message: "Method not allowed", Code: http.StatusMethodNotAllowed}
)

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func MakeHandlerFunc(f ApiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			if e, ok := err.(Error); ok {
				WriteJSON(w, e.Code, e)
				return
			}
			WriteJSON(w, http.StatusInternalServerError, ErrInternal)
		}
	}
}

func MakeHandlerFuncMap(funcs map[string]ApiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if f, ok := funcs[r.Method]; ok {
			MakeHandlerFunc(f)(w, r)
			return
		}
		WriteJSON(w, http.StatusMethodNotAllowed, ErrMethodNotAllowed)
	}
}

func ErrBadRequest(message string) Error {
	return Error{
		Message: message,
		Code:    http.StatusBadRequest,
	}
}

func ErrValidation(message string, errors map[string]string) Error {
	return Error{
		Message: message,
		Code:    http.StatusUnprocessableEntity,
		Errors:  errors,
	}
}

func ErrConflict(message string) Error {
	return Error{
		Message: message,
		Code:    http.StatusConflict,
	}
}

func GetEnvInt(key string, def int) int {
	s := os.Getenv(key)
	if len(s) == 0 {
		return def
	}
	v, err := strconv.Atoi(s)

	if err != nil {
		return def
	}

	return v
}

func GetEnvStr(key string, def string) string {
	s := os.Getenv(key)
	if len(s) == 0 {
		return def
	}

	return s
}

func Perftimer(name string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s took %v\n", name, time.Since(start))
	}
}
