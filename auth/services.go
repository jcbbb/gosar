package auth

import (
	"context"
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jcbbb/gosar/common"
	"github.com/jcbbb/gosar/db"
	"github.com/jcbbb/gosar/user"
)

var (
	JWT_SECRET = os.Getenv("JWT_SECRET")
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

func (session *Session) save(tx pgx.Tx) (*Session, error) {
	claims := &jwt.RegisteredClaims{
		Subject:   strconv.Itoa(session.UserID),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(JWT_EXPIRATION) * time.Second)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Audience:  []string{"https://gosar.homeleess.dev"},
		Issuer:    "https://gosar.homeless.dev",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(JWT_SECRET))
	if err != nil {
		return nil, common.ErrInternal
	}

	session.Token = tokenString

	row := tx.QueryRow(
		context.Background(),
		"insert into sessions (user_id, token) values ($1, $2) returning id",
		session.UserID, session.Token,
	)

	if err := row.Scan(&session.ID); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return nil, common.ErrConflict("User already exists")
		}

		return nil, common.ErrInternal
	}

	return session, nil
}

func login(req LoginReq) (*Session, error) {
	user, err := user.GetByLogin(req.Login)

	if err != nil {
		return nil, err
	}

	if err := user.VerifyPassword(req.Password); err != nil {
		return nil, err
	}

	tx, err := db.Pool.BeginTx(context.TODO(), pgx.TxOptions{})

	if err != nil {
		return nil, common.ErrInternal
	}

	defer func() {
		if err != nil {
			tx.Rollback(context.TODO())
		} else {
			tx.Commit(context.TODO())
		}
	}()

	session, err := NewSession(user.ID).save(tx)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func signup(req SignupReq) (*Session, error) {
	// ignore error as it's already validated
	age, _ := strconv.Atoi(req.Age)

	tx, err := db.Pool.Begin(context.TODO())

	if err != nil {
		return nil, common.ErrInternal
	}

	defer func() {
		if err != nil {
			tx.Rollback(context.TODO())
		} else {
			tx.Commit(context.TODO())
		}
	}()

	user, err := user.New(user.Opts{
		Name:     req.Name,
		Password: req.Password,
		Age:      age,
		Login:    req.Login,
	}).Save(tx)

	if err != nil {
		return nil, err
	}

	session, err := NewSession(user.ID).save(tx)

	if err != nil {
		return nil, err
	}

	return session, nil
}
