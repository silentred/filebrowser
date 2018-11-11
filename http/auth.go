package http

import (
	"context"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/filebrowser/filebrowser/types"
)

func (e *Env) loginHandler(w http.ResponseWriter, r *http.Request) {
	user, err := e.Auther.Auth(r)
	if err == types.ErrNoPermission {
		httpErr(w, http.StatusForbidden, nil)
	} else if err != nil {
		httpErr(w, http.StatusInternalServerError, err)
	} else {
		e.printToken(w, user.ID)
	}
}

type authToken struct {
	UserID uint
	jwt.StandardClaims
}

func (e *Env) auth(next http.HandlerFunc) http.HandlerFunc {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		return e.Settings.Key, nil
	}

	tkExtractor := request.AuthorizationHeaderExtractor

	nextWithUser := func(w http.ResponseWriter, r *http.Request, id uint) {
		ctx := context.WithValue(r.Context(), keyUserID, id)
		next(w, r.WithContext(ctx))
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var tk authToken
		token, err := request.ParseFromRequestWithClaims(r, tkExtractor, &tk, keyFunc)

		if err != nil || !token.Valid {
			httpErr(w, http.StatusForbidden, nil)
			return
		}

		nextWithUser(w, r, tk.UserID)
	}
}

func (e *Env) printToken(w http.ResponseWriter, id uint) {
	claims := &authToken{
		UserID: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			Issuer:    "File Browser",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(e.Settings.Key)

	if err != nil {
		httpErr(w, http.StatusInternalServerError, err)
	} else {
		w.Header().Set("Content-Type", "cty")
		w.Write([]byte(signed))
	}
}
