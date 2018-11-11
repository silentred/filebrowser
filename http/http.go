package http

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/filebrowser/filebrowser/types"
	"github.com/gorilla/mux"
)

type key int

const (
	keyUserID key = iota
)

// Env ...
type Env struct {
	Auther   types.Auther
	Runner   *types.Runner
	Settings *types.Settings
	Store    *types.Store
}

// Handler ...
func Handler(e *Env) http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/api/login", e.loginHandler)

	users := r.PathPrefix("/api/users").Subrouter()
	users.HandleFunc("", e.auth(e.usersGetHandler)).Methods("GET")
	users.HandleFunc("", e.auth(e.userPostHandler)).Methods("POST")
	users.HandleFunc("/{id:[0-9]+}", e.auth(e.userPutHandler)).Methods("PUT")
	users.HandleFunc("/{id:[0-9]+}", e.auth(e.userGetHandler)).Methods("GET")
	users.HandleFunc("/{id:[0-9]+}", e.auth(e.userDeleteHandler)).Methods("DELETE")

	r.PathPrefix("/api/resources").HandlerFunc(e.auth(e.resourceGetHandler)).Methods("GET")
	r.PathPrefix("/api/resources").HandlerFunc(e.auth(e.resourceDeleteHandler)).Methods("DELETE")
	r.PathPrefix("/api/resources").HandlerFunc(e.auth(e.resourcePostPutHandler)).Methods("POST")
	r.PathPrefix("/api/resources").HandlerFunc(e.auth(e.resourcePostPutHandler)).Methods("PUT")
	r.PathPrefix("/api/resources").HandlerFunc(e.auth(e.resourcePatchHandler)).Methods("PATCH")
	return r
}

func httpErr(w http.ResponseWriter, status int, err error) {
	txt := http.StatusText(status)
	if err != nil {
		log.Printf("%v", err)
	}
	http.Error(w, strconv.Itoa(status)+" "+txt, status)
}

func renderJSON(w http.ResponseWriter, data interface{}) {
	marsh, err := json.Marshal(data)
	if err != nil {
		httpErr(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if _, err := w.Write(marsh); err != nil {
		httpErr(w, http.StatusInternalServerError, err)
	}
}

func (e *Env) getUser(w http.ResponseWriter, r *http.Request) (*types.User, bool) {
	id := r.Context().Value(keyUserID).(uint)
	user, err := e.Store.Users.Get(id)
	if err == types.ErrNotExist {
		httpErr(w, http.StatusForbidden, nil)
		return nil, false
	}

	if err != nil {
		httpErr(w, http.StatusInternalServerError, err)
		return nil, false
	}

	return user, true
}
