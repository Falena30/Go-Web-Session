package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

//Credentials untuk menampung json dan db
type Credentials struct {
	Password string `json:"password", db:"Password"`
	Username string `json:"username", db:"Username"`
}

func HandleSignUp(w http.ResponseWriter, r *http.Request) {
	//parse dan decode request body menjadi `credentials` intence yang baru
	db, err1 := Connect()
	if err1 != nil {
		fmt.Println(err1.Error())
	}
	defer db.Close()
	creds := &Credentials{}
	err := json.NewDecoder(r.Body).Decode(creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//salt dan has denga menggunakan algoritma bcrypt
	//argumen kedua adalah cost untuk hash, untuk sekarang gunakan saja 8 tapi tergantung
	//nanti bisa dikurangin
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), 8)

	if _, err := db.Query("INSERT INTO User_DB (`ID_User`, `Username`, `Password`) VALUES (?,?,?)", nil, creds.Username, string(hashedPassword)); err != nil {
		//jika terdapat error pada db kemabalikan nilai 500
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err.Error())
		return
	}
}
