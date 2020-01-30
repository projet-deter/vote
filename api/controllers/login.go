package controllers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/m2fof/vote/api/auth"
	"github.com/m2fof/vote/api/models"
	"github.com/m2fof/vote/api/responses"
	"github.com/m2fof/vote/api/utils/formaterror"
	"golang.org/x/crypto/bcrypt"
)

func (server *Server) Login(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	user.Prepare()
	err = user.Validate("login")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	var token string
	user, token, err = server.SignIn(user.Email, user.Password)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusUnprocessableEntity, formattedError)
		return
	}

	type data struct {
		User_actif   models.User
		Access_token string
	}
	dataActif := data{
		User_actif:   user,
		Access_token: token,
	}

	responses.JSON(w, http.StatusOK, dataActif)
}

func (server *Server) SignIn(email, password string) (models.User, string, error) {
	var err error
	user := models.User{}

	err = server.DB.Debug().Model(models.User{}).Where("email = ?", email).Take(&user).Error
	if err != nil {
		return user, "", err
	}

	err = models.VerifyPassword(user.Password, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return user, "", err
	}

	log.Println("Signin")
	log.Println(user.ID)
	log.Println(user.First_name)
	log.Println(user.Last_name)
	log.Println(user.Email)
	log.Println("End Signin")

	token, err := auth.CreateToken(user.ID)

	return user, token, err
}
