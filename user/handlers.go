package user

import (
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
