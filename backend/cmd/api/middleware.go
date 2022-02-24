package main

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/pascaldekloe/jwt"
)

func (app *application) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")

		next.ServeHTTP(w, r)
	})
}

func (app *application) checkToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Add("Vary", "Authorization")

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			app.errorJSON(w, errors.New("unauthorized"), http.StatusForbidden)
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 {
			app.errorJSON(w, errors.New("unauthorized"), http.StatusForbidden)
			return
		}

		if strings.ToLower(headerParts[0]) != "bearer" {
			app.errorJSON(w, errors.New("unauthorized"), http.StatusForbidden)
			return
		}

		token := headerParts[1]

		claims, err := jwt.HMACCheck([]byte(token), []byte(app.config.jwt.secret))
		if err != nil {
			app.errorJSON(w, errors.New("unauthorized"), http.StatusForbidden)
			return
		}

		if !claims.Valid(time.Now()) {
			app.errorJSON(w, errors.New("unauthorized"), http.StatusForbidden)
			return
		}

		if !claims.AcceptAudience("localhost") {
			app.errorJSON(w, errors.New("unauthorized"), http.StatusForbidden)
			return
		}

		if claims.Issuer != "localhost" {
			app.errorJSON(w, errors.New("unauthorized"), http.StatusForbidden)
			return
		}

		// userId, err := strconv.ParseInt(claims.Subject, 10, 64)
		// if err != nil {
		//	app.errorJSON(w, errors.New("unauthorized"), http.StatusForbidden)
		// 	return
		// }

		next.ServeHTTP(w, r)
	})

}
