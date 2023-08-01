package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jcbbb/gosar/common"
	"github.com/jcbbb/gosar/db"
	"golang.org/x/crypto/bcrypt"
)

func (apr AddPhoneReq) validate() error {
	errors := make(map[string]string, 0)

	if len(apr.Phone) > 12 {
		errors["phone"] = "Phone number must be maximum 12 characters"
	}

	if len(apr.Phone) == 0 {
		errors["phone"] = "Phone is required"
	}

	if len(apr.Description) == 0 {
		errors["description"] = "Description is required"
	}

	if len(errors) > 0 {
		return common.ErrValidation("Validation failed", errors)
	}

	return nil
}

func (upr UpdatePhoneReq) validate() error {
	errors := make(map[string]string, 0)

	if len(upr.Phone) > 12 {
		errors["phone"] = "Phone number must be maximum 12 characters"
	}

	if len(upr.Phone) == 0 {
		errors["phone"] = "Phone is required"
	}

	if len(upr.Description) == 0 {
		errors["description"] = "Description is required"
	}

	if upr.PhoneID == nil {
		errors["phone_id"] = "Phone id is required"
	}

	if len(errors) > 0 {
		return common.ErrValidation("Validation failed", errors)
	}

	return nil
}

func New(opts Opts) *User {
	return &User{
		Name:     opts.Name,
		Password: opts.Password,
		Age:      opts.Age,
		Login:    opts.Login,
	}
}

func (user *User) Save(tx pgx.Tx) (*User, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)

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

func addPhone(opts AddPhoneOpts) (*Phone, error) {
	phone := Phone{
		Phone:       opts.phone,
		UserID:      opts.userId,
		Description: opts.description,
		IsFax:       opts.isFax,
	}

	row := db.Pool.QueryRow(
		context.Background(),
		"insert into user_phones (user_id, phone, description, is_fax) values ($1, $2, $3, $4) returning id",
		opts.userId, opts.phone, opts.description, opts.isFax,
	)

	// TODO: properly handle errors
	if err := row.Scan(&phone.ID); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return nil, common.ErrConflict("User already exists")
		}

		return nil, common.ErrInternal
	}

	return &phone, nil
}

func getPhones(query string) ([]Phone, error) {
	phones := make([]Phone, 0)

	rows, err := db.Pool.Query(
		context.Background(),
		"select id, user_id, phone, description, is_fax from user_phones where phone ilike $1",
		"%"+query+"%",
	)

	defer rows.Close()

	// TODO: handle errors properly
	if err != nil {
		return nil, common.ErrInternal
	}

	for rows.Next() {
		var phone Phone
		if err := rows.Scan(&phone.ID, &phone.UserID, &phone.Phone, &phone.Description, &phone.IsFax); err != nil {
			return nil, common.ErrInternal
		}

		phones = append(phones, phone)
	}

	if err := rows.Err(); err != nil {
		return nil, common.ErrInternal
	}

	return phones, nil
}

func updatePhone(opts UpdatePhoneOpts) (*Phone, error) {
	phone := Phone{
		ID:          opts.id,
		UserID:      opts.userId,
		Description: opts.description,
		Phone:       opts.phone,
		IsFax:       opts.isFax,
	}

	res, err := db.Pool.Exec(
		context.Background(),
		"update user_phones set description = $1, phone = $2, is_fax = $3 where id = $4 returning id",
		opts.description, opts.phone, opts.isFax, opts.id,
	)

	if err != nil {
		return nil, common.ErrInternal
	}

	if res.RowsAffected() == 0 {
		return nil, common.ErrNotFound(fmt.Sprintf("Phone with id %v not found", opts.id))
	}

	return &phone, nil
}

func deletePhone(opts DeletePhoneOpts) (int, error) {
	res, err := db.Pool.Exec(context.Background(), "delete from user_phones where id = $1 and user_id = $2", opts.id, opts.userId)

	if err != nil {
		return 0, common.ErrInternal
	}

	if res.RowsAffected() == 0 {
		return 0, common.ErrNotFound(fmt.Sprintf("Phone with id %v not found", opts.id))
	}

	return opts.id, nil
}
