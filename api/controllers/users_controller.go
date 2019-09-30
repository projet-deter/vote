package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	//"strings"
	"strconv"
	"os"

	"github.com/gorilla/mux"
	"github.com/m2fof/vote/api/auth"
	"github.com/m2fof/vote/api/models"
	"github.com/m2fof/vote/api/responses"
	"github.com/m2fof/vote/api/utils/formaterror"
)

func init() {
	//globalSessions, _ = session.NewManager("memory", `{"cookieName":"gosessionid", "enableSetCookie,omitempty": true, "gclifetime":3600, "maxLifetime": 3600, "secure": false, "sessionIDHashFunc": "sha1", "sessionIDHashKey": "", "cookieLifeTime": 3600, "providerConfig": ""}`)
	//go globalSessions.GC()
}

func (server *Server) CreateUser(w http.ResponseWriter, r *http.Request) {
	log.Println(os.Getenv("currentUserFirst_name"))
	log.Println(os.Getenv("currentUserAccessLevel"))

	if os.Getenv("currentUserAccessLevel") == strconv.FormatInt(0, 10) {

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			responses.ERROR(w, http.StatusUnprocessableEntity, err)
		}
		user := models.User{}
		err = json.Unmarshal(body, &user)
		if err != nil {
			responses.ERROR(w, http.StatusUnprocessableEntity, err)
			return
		}
		user.Prepare()

		err = user.Validate("")
		if err != nil {
			responses.ERROR(w, http.StatusUnprocessableEntity, err)
			return
		}
		userCreated, err := user.SaveUser(server.DB)

		if err != nil {

			formattedError := formaterror.FormatError(err.Error())

			responses.ERROR(w, http.StatusInternalServerError, formattedError)
			return
		}
		w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, userCreated.ID))
		responses.JSON(w, http.StatusCreated, userCreated)
	} else {
		log.Println("You don't have access right to perform this action:")
	}
}

func (server *Server) GetUsers(w http.ResponseWriter, r *http.Request) {

	user := models.User{}

	users, err := user.FindAllUsers(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, users)
}

func (server *Server) GetUser(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	user := models.User{}
	userGotten, err := user.FindUserByID(server.DB, uint32(uid))
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	responses.JSON(w, http.StatusOK, userGotten)
}

func (server *Server) UpdateUser(w http.ResponseWriter, r *http.Request) {
	if os.Getenv("currentUserAccessLevel") == strconv.FormatInt(0, 10) {

		vars := mux.Vars(r)
		uid, err := strconv.ParseUint(vars["id"], 10, 32)
		if err != nil {
			responses.ERROR(w, http.StatusBadRequest, err)
			return
		}
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
		tokenID, err := auth.ExtractTokenID(r)
		if err != nil {
			responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
			return
		}
		if tokenID != uint32(uid) {
			responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
			return
		}
		user.Prepare()
		err = user.Validate("update")
		if err != nil {
			responses.ERROR(w, http.StatusUnprocessableEntity, err)
			return
		}
		updatedUser, err := user.UpdateAUser(server.DB, uint32(uid))
		if err != nil {
			formattedError := formaterror.FormatError(err.Error())
			responses.ERROR(w, http.StatusInternalServerError, formattedError)
			return
		}
		responses.JSON(w, http.StatusOK, updatedUser)

	} else {
		log.Println("You don't have access right to perform this action:")
	}
}

func (server *Server) DeleteUser(w http.ResponseWriter, r *http.Request) {

	if os.Getenv("currentUserAccessLevel") == strconv.FormatInt(0, 10) {

		vars := mux.Vars(r)

		user := models.User{}

		uid, err := strconv.ParseUint(vars["id"], 10, 32)
		if err != nil {
			responses.ERROR(w, http.StatusBadRequest, err)
			return
		}
		tokenID, err := auth.ExtractTokenID(r)
		if err != nil {
			responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
			return
		}
		if tokenID != 0 && tokenID != uint32(uid) {
			responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
			return
		}
		_, err = user.DeleteAUser(server.DB, uint32(uid))
		if err != nil {
			responses.ERROR(w, http.StatusInternalServerError, err)
			return
		}
		w.Header().Set("Entity", fmt.Sprintf("%d", uid))
		responses.JSON(w, http.StatusNoContent, "")

	} else {
		log.Println("You don't have access right to perform this action:")
	}
}
