package http

import (
	"strings"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/GeertJohan/go.rice"
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

func (e *Env) getStaticHandler() http.Handler {
	box := rice.MustFindBox("../../react/build")
	handler := http.FileServer(box.HTTPBox())

	// TODO: cleanup this code. generate data previously.
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("x-frame-options", "SAMEORIGIN")
		w.Header().Set("x-xss-protection", "1; mode=block")

		// TODO: prefix URL

		if _, err := box.Open(r.URL.Path); err != nil {
			r.URL.Path = "/"
		}

		if r.URL.Path == "/" || r.URL.Path == "/index.html" {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			index := template.Must(template.New("index").Delims("$", "$").Parse(box.MustString("/index.html")))
			data := map[string]interface{}{"HOMEPAGE": strings.TrimSuffix(e.Settings.BaseURL, "/")}
			err := index.Execute(w, data)
			if err != nil {
				httpErr(w, http.StatusInternalServerError, err)
			}
			return
		}

		handler.ServeHTTP(w, r)
	})
}

// Handler ...
func Handler(e *Env) http.Handler {
	r := mux.NewRouter()

	r.NotFoundHandler = e.getStaticHandler()

	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/login", e.loginHandler)
	api.HandleFunc("/signup", e.signupHandler)

	users := api.PathPrefix("/users").Subrouter()
	users.HandleFunc("", e.auth(e.usersGetHandler)).Methods("GET")
	users.HandleFunc("", e.auth(e.userPostHandler)).Methods("POST")
	users.HandleFunc("/{id:[0-9]+}", e.auth(e.userPutHandler)).Methods("PUT")
	users.HandleFunc("/{id:[0-9]+}", e.auth(e.userGetHandler)).Methods("GET")
	users.HandleFunc("/{id:[0-9]+}", e.auth(e.userDeleteHandler)).Methods("DELETE")

	api.PathPrefix("/resources").HandlerFunc(e.auth(e.resourceGetHandler)).Methods("GET")
	api.PathPrefix("/resources").HandlerFunc(e.auth(e.resourceDeleteHandler)).Methods("DELETE")
	api.PathPrefix("/resources").HandlerFunc(e.auth(e.resourcePostPutHandler)).Methods("POST")
	api.PathPrefix("/resources").HandlerFunc(e.auth(e.resourcePostPutHandler)).Methods("PUT")
	api.PathPrefix("/resources").HandlerFunc(e.auth(e.resourcePatchHandler)).Methods("PATCH")

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
