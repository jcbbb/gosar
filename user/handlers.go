package user

import (
	"encoding/json"
	"net/http"

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

	phone, err := addPhone(AddPhoneOpts{
		phone:       req.phone,
		description: req.description,
		isFax:       req.isFax,
	})

	if err != nil {
		return err
	}

	return common.WriteJSON(w, http.StatusCreated, phone)
}
