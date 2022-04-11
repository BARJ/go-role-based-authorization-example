package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"strings"
)

type Route struct {
	path   string
	method string
	handle httprouter.Handle
	roles  []Role
}

func newRouter() http.Handler {
	router := httprouter.New()

	routes := []Route{
		{
			path:   "/ping",
			method: http.MethodGet,
			handle: func(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
				_, _ = fmt.Fprintf(w, "Pong!\n")
			},
		},
		{
			path:   "/users.sign-in",
			method: http.MethodPost,
			handle: signInHandler,
		},
		{
			path:   "/users.whoami",
			method: http.MethodGet,
			handle: whoAmIHandler,
		},
		{
			path:   "/jokes",
			method: http.MethodGet,
			handle: searchJokesHandler,
		},
		{
			path:   "/jokes",
			method: http.MethodPost,
			handle: createJokeHandler,
			roles:  []Role{User},
		},
		{
			path:   "/jokes/:id",
			method: http.MethodDelete,
			handle: deleteJokeHandler,
			roles:  []Role{Admin},
		},
	}

	for _, route := range routes {
		h := authorization(route.handle, route.roles)
		h = authentication(h)
		router.Handle(route.method, route.path, h)
	}

	return router
}

func authentication(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		var name string
		var roles []Role

		tokenString, err := getTokenString(r)
		if err != nil {
			log.Printf("authentication: %v\n", err)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		} else if tokenString == "" {
			name = "Anonymous"
			roles = []Role{}
		} else {
			_, claims, err := decodeJwt(tokenString)
			if err != nil {
				log.Printf("authentication: %v\n", err)
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			name = claims.Name
			roles = claims.Roles
		}

		r = r.WithContext(context.WithValue(r.Context(), "username", name))
		r = r.WithContext(context.WithValue(r.Context(), "roles", roles))

		h(w, r, ps)
	}
}

func authorization(h httprouter.Handle, roles []Role) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		grantedRoles := r.Context().Value("roles").([]Role)

		if len(roles) == 0 {
			h(w, r, ps)
			return
		}

		for _, wantRole := range roles {
			for _, gotRole := range grantedRoles {
				if gotRole == wantRole {
					h(w, r, ps)
					return
				}
			}
		}
		log.Printf("authorization: has any authority %v but got %v\n", roles, grantedRoles)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)

	}
}

func getTokenString(r *http.Request) (string, error) {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		return "", nil
	}
	auths := strings.Split(auth, " ")
	if len(auths) != 2 || auths[0] != "Bearer" {
		return "", errors.New("getTokenString: missing or invalid bearer token")
	}
	return auths[1], nil
}
