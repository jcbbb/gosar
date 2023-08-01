package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
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

func EnsureAuth(next common.ApiFunc) common.ApiFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		tokenCookie, err := r.Cookie(SESSION_COOKIE_NAME)

		if err != nil {
			return common.ErrUnauthenticated("Session cookie missing")
		}

		token, err := jwt.ParseWithClaims(tokenCookie.Value, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, common.ErrBadRequest(fmt.Sprintf("Unexpected signing method: %v", token.Header["alg"]))
			}

			return []byte(JWT_SECRET), nil
		})

		// TODO: validate claims
		claims, ok := token.Claims.(*jwt.RegisteredClaims)

		if !ok && !token.Valid {
			return common.ErrUnauthenticated("Invalid token")
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, "claims", claims)

		return next(w, r.WithContext(ctx))
	}
}
