package user

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jcbbb/gosar/common"
	"github.com/jcbbb/gosar/db"
	"golang.org/x/crypto/bcrypt"
)

func New(opts Opts) *User {
	return &User{
		Name:     opts.Name,
		Password: opts.Password,
		Age:      opts.Age,
		Login:    opts.Login,
	}
}

func (user *User) Save(tx pgx.Tx) (*User, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)

	if err != nil {
		return nil, err
	}

	user.Password = string(bytes)

	row := tx.QueryRow(
		context.Background(),
		"insert into users (login, password, name, age) values ($1, $2, $3, $4) returning id",
		user.Login, user.Password, user.Name, user.Age,
	)

	if err := row.Scan(&user.ID); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return nil, common.ErrConflict("User already exists")
		}

		return nil, common.ErrInternal
	}

	return user, nil
}

func (user *User) VerifyPassword(password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return common.ErrBadRequest("Login or password is incorrect")
	}
	return nil
}

func GetByLogin(login string) (*User, error) {
	var user User
	row := db.Pool.QueryRow(
		context.Background(),
		"select id, password from users where login = $1",
		login,
	)

	if err := row.Scan(&user.ID, &user.Password); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return nil, common.ErrConflict("User already exists")
		}

		return nil, common.ErrInternal
	}

	return &user, nil
}

func getByName(name string) (*User, error) {
	var user User
	row := db.Pool.QueryRow(
		context.Background(),
		"select id, name, login, age, password from users where name = $1",
		name,
	)

	err := row.Scan(&user.ID, &user.Name, &user.Login, &user.Age, &user.Password)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return nil, common.ErrConflict("User already exists")
		}

		return nil, common.ErrInternal
	}

	return &user, nil
}
