package user

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jcbbb/gosar/common"
)

func HandleGetUser(w http.ResponseWriter, r *http.Request) error {
	name := r.URL.Path[len("/user/"):]
	user, err := getByName(name)

	if err != nil {
		return err
	}

	return common.WriteJSON(w, http.StatusOK, user)
}

func HandleAddPhone(w http.ResponseWriter, r *http.Request) error {
	var req AddPhoneReq

	err := json.NewDecoder(r.Body).Decode(&req)
	defer r.Body.Close()

	if err != nil {
		return common.ErrInternal
	}

	if err := req.validate(); err != nil {
		return err
	}

	claims := r.Context().Value("claims").(*jwt.RegisteredClaims)
	sub, _ := strconv.Atoi(claims.Subject)

	phone, err := addPhone(AddPhoneOpts{
		phone:       req.Phone,
		description: req.Description,
		isFax:       req.IsFax,
		userId:      sub,
	})

	if err != nil {
		return err
	}

	return common.WriteJSON(w, http.StatusCreated, phone)
}

func HandleGetPhones(w http.ResponseWriter, r *http.Request) error {
	q := r.URL.Query().Get("q")

	phones, err := getPhones(q)
	if err != nil {
		return err
	}

	return common.WriteJSON(w, http.StatusOK, phones)
}

func HandleUpdatePhone(w http.ResponseWriter, r *http.Request) error {
	var req UpdatePhoneReq

	err := json.NewDecoder(r.Body).Decode(&req)
	defer r.Body.Close()

	if err != nil {
		return common.ErrInternal
	}

	if err := req.validate(); err != nil {
		return err
	}

	claims := r.Context().Value("claims").(*jwt.RegisteredClaims)
	sub, _ := strconv.Atoi(claims.Subject)

	phone, err := updatePhone(UpdatePhoneOpts{
		id:          *req.PhoneID,
		description: req.Description,
		userId:      sub,
		isFax:       req.IsFax,
		phone:       req.Phone,
	})

	if err != nil {
		return err
	}

	return common.WriteJSON(w, http.StatusOK, phone)
}

func HandleDeletePhone(w http.ResponseWriter, r *http.Request) error {
	//phoneId := r.URL.Path[len("/user/phone/"):]

	return nil
}
