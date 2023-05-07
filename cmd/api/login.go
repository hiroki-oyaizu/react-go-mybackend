package main

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"react-go-mybackend/database"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type LoginInfo struct {
	Mail     string `json:"mail"`
	Password string `json:"password"`
}

func (a *application) LoginUser(w http.ResponseWriter, r *http.Request) {
	var loginInfo LoginInfo
	err := json.NewDecoder(io.LimitReader(r.Body, 1048576)).Decode(&loginInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db := database.Connect()
	defer db.Close()

	var userID int
	var firstName, lastName string
	err = db.QueryRow("SELECT id, firstName, lastName FROM users WHERE mail = ? AND password = ?", loginInfo.Mail, loginInfo.Password).Scan(&userID, &firstName, &lastName)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":        userID,
		"email":     loginInfo.Mail,
		"firstName": firstName,
		"lastName":  lastName,
		"exp":       time.Now().Add(time.Hour * 24).Unix(),
	})

	secretKey := "your_secret_key"
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token":     signedToken,
		"userId":    strconv.Itoa(userID), // 追加
		"firstName": firstName,
		"lastName":  lastName,
	})

}
