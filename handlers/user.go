package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/spectrum-team/tachartas/commons"
	"github.com/spectrum-team/tachartas/interfaces"
	"github.com/spectrum-team/tachartas/models"
	"github.com/spectrum-team/tachartas/repos"
)

type UserHandler struct {
	dbConfig *models.DatabaseConfig
	userRepo interfaces.UserRepository
}

func NewUserHandler(dbConfig *models.DatabaseConfig) *UserHandler {
	return &UserHandler{
		dbConfig: dbConfig,
		userRepo: repos.NewUserRepository(dbConfig),
	}
}

func (u *UserHandler) SignIn(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user := &models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		log.Println("klok wawawa => ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = u.userRepo.UpsertUser(user)
	if err != nil {
		log.Println("klok wawawa => ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	token := commons.SignIn(user)

	w.Header().Set("Authorization", token)
	w.WriteHeader(http.StatusOK)
}
