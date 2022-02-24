package main

import (
	"backend/cmd/api/dtos"
	"backend/models"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/pascaldekloe/jwt"
	"golang.org/x/crypto/bcrypt"
)

var validUser = models.User{
	Id:       1,
	Email:    "email@there.com",
	Password: "$2a$14$ajq8Q7fbtFRQvXpdCq7Jcuy.Rx1h/L4J60Otx.gyNLbAYctGMJ9tK",
}

func (app *application) Signin(w http.ResponseWriter, r *http.Request) {
	var creds dtos.CredentialsRequest

	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		app.errorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}

	hashedPassword := validUser.Password

	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(creds.Password))
	if err != nil {
		app.errorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}

	var claims jwt.Claims
	claims.Subject = fmt.Sprint(validUser.Id)
	claims.Issued = jwt.NewNumericTime(time.Now())
	claims.NotBefore = jwt.NewNumericTime(time.Now())
	claims.Expires = jwt.NewNumericTime(time.Now().Add(24 * time.Hour))
	claims.Issuer = "localhost"
	claims.Audiences = []string{"localhost"}
	claims.Set = map[string]interface{}{"profile": "admin"}

	jwtBytes, err := claims.HMACSign(jwt.HS512, []byte(app.config.jwt.secret))
	if err != nil {
		app.errorJSON(w, errors.New("error signin"))
		return
	}

	app.writeJSON(w, http.StatusOK, string(jwtBytes), "token")
}
