package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func signInHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	request := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	token, err := func(username, password string) (string, error) {
		var name string
		var roles []Role
		if username == "John Doe" && password == "12345678" {
			name = "John Doe"
			roles = []Role{User}
		} else if username == "Jane Doe" && password == "87654321" {
			name = "Jane Doe"
			roles = []Role{Admin, User}
		} else {
			return "", errors.New("invalid user credentials")
		}
		return encodeJwt(JwtClaims{
			Name:  name,
			Roles: roles,
		})
	}(request.Username, request.Password)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	response := struct {
		Token string `json:"token"`
	}{token}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}

func whoAmIHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	username := r.Context().Value("username").(string)
	roles := r.Context().Value("roles").([]Role)
	var roleNames = make([]string, len(roles))
	for _, role := range roles {
		roleNames = append(roleNames, fmt.Sprint(role))
	}
	response := struct {
		Username string   `json:"username"`
		Roles    []string `json:"roles"`
	}{
		username,
		roleNames,
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}
