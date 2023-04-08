package main

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"react-go-mybackend/database"
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
	err = db.QueryRow("SELECT id FROM users WHERE mail = ? AND password = ?", loginInfo.Mail, loginInfo.Password).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}
