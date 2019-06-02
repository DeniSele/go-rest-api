// Package handler presents Simple API for operations with client data
// CRUD operations
package handler

import (
	"GoTask/database"
	"GoTask/model"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// GetUsers returns all users
func GetUsers(w http.ResponseWriter, r *http.Request) {
	users := []model.UserInfo{}

	users, err := database.SelectAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	bytes, err := json.Marshal(users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	writeJSONResponse(w, bytes)
}

// GetUserBySName returns a user by last name
func GetUserBySName(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	curRecord, err := database.SelectBySecondName(params["secondname"])
	if err != nil {
		http.NotFound(w, r)
		return
	}

	bytes, err := json.Marshal(curRecord)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	writeJSONResponse(w, bytes)
}

// CreateUser creates a new user entry
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var newUser model.UserInfo

	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = database.AddNewUser(newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// UpdateUser updates user record
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var record model.UserInfo
	err := json.NewDecoder(r.Body).Decode(&record)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = database.UpdateUserDB(id, record)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

// DeleteUser deletes the user record
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = database.DeleteUserByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
