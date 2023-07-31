package auth

import (
	"encoding/json"
	"net/http"

	"github.com/jcbbb/gosar/common"
)

var (
	ENV                 = common.GetEnvStr("ENV", "development")
	JWT_EXPIRATION      = common.GetEnvInt("JWT_EXPIRATION", 86400)
	SESSION_COOKIE_NAME = common.GetEnvStr("SESSION_COOKIE_NAME", "SESSTOKEN")
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

	session, err := signup(req)

	if err != nil {
		return err
	}

	cookie := http.Cookie{
		Name:     SESSION_COOKIE_NAME,
		Value:    session.Token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		MaxAge:   JWT_EXPIRATION,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, &cookie)

	return common.WriteJSON(w, http.StatusCreated, session)
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

	session, err := login(req)

	if err != nil {
		return err
	}

	cookie := http.Cookie{
		Name:     SESSION_COOKIE_NAME,
		Value:    session.Token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		MaxAge:   JWT_EXPIRATION,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, &cookie)
	return common.WriteJSON(w, http.StatusOK, session)
}
