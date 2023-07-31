package user

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jcbbb/gosar/common"
	"github.com/jcbbb/gosar/db"
	"golang.org/x/crypto/bcrypt"
)

func New(opts UserOpts) *User {
	return &User{
		Name:     opts.Name,
		Password: opts.Password,
		Age:      opts.Age,
		Login:    opts.Login,
	}
}

func Save(user *User) (*User, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)

	if err != nil {
		return nil, err
	}

	user.Password = string(bytes)

	row := db.Pool.QueryRow(
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
