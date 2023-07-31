package auth

import (
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/jcbbb/gosar/common"
)

var (
	JWT_SECRET     = os.Getenv("JWT_SECRET")
	JWT_EXPIRATION = os.Getenv("JWT_EXPIRATION")
)

func (sr SignupReq) validate() error {
	errors := make(map[string]string)

	if len(sr.Password) == 0 {
		errors["password"] = "Password is required"
	}

	if len(sr.Name) == 0 {
		errors["name"] = "Name is required"
	}

	if len(sr.Login) == 0 {
		errors["login"] = "Login is required"
	}

	if len(sr.Age) == 0 {
		errors["age"] = "Age is required"
	}

	if _, err := strconv.Atoi(sr.Age); err != nil {
		errors["age"] = "Age must be an integer"
	}

	if len(errors) > 0 {
		return common.ErrValidation("Validation failed", errors)
	}

	return nil
}

func (lr LoginReq) validate() error {
	errors := make(map[string]string)

	if len(lr.Login) == 0 {
		errors["login"] = "Login is required"
	}

	if len(lr.Password) == 0 {
		errors["password"] = "Password is required"
	}

	if len(errors) > 0 {
		return common.ErrValidation("Validation failed", errors)
	}

	return nil
}

func NewSession(userID int) *Session {
	return &Session{
		UserID: userID,
	}
}

func Save(session *Session) (*Session, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": session.UserID,
		"exp": time.Now().Add(time.Duration(common.GetEnvInt("JWT_EXPIRATION", 1440)) * time.Minute),
	})

	tokenString, err := token.SignedString(os.Getenv("JWT_SECRET"))

	return nil, nil
}
