package auth

import (
	"encoding/json"
	"net/http"

	"github.com/jcbbb/gosar/common"
)

func HandleSignup(w http.ResponseWriter, r *http.Request) error {
	req := SignupReq{
		Login:    r.FormValue("login"),
		Password: r.FormValue("password"),
		Age:      r.FormValue("age"),
		Name:     r.FormValue("name"),
	}

	if err := req.validate(); err != nil {
		return err
	}

	return nil
}

func HandleLogin(w http.ResponseWriter, r *http.Request) error {
	var req LoginReq
	err := json.NewDecoder(r.Body).Decode(&req)

	defer r.Body.Close()

	if err != nil {
		return common.ErrInternal
	}

	if err := req.validate(); err != nil {
		return err
	}

	return nil
}
