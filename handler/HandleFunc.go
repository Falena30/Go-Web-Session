package handler

import (
	"database/sql"
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

	db, err1 := Connect()
	if err1 != nil {
		fmt.Println(err1.Error())
	}
	defer db.Close()
	//parse dan decode request body menjadi `credentials` intence yang baru
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

func HandlerSignIn(w http.ResponseWriter, r *http.Request) {
	//hubungkan db
	db, err := Connect()
	if err != nil {
		fmt.Println(err.Error())
	}
	defer db.Close()

	//parse dan decode request body menjadi `credentials` intence yang baru
	cerds := &Credentials{}
	err = json.NewDecoder(r.Body).Decode(cerds)
	if err != nil {
		//jika terjadi kesalahan di body akan mengembalikan 400/ bad request
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)

		return
	}
	//panggil data dari db
	selectData := db.QueryRow("SELECT Password FROM User_DB WHERE Username = ?", cerds.Username)

	if err != nil {
		//jika ada error mengembalikan nilai 500
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//buat instence baru untuk credentials untuk menampung data dari database
	storageCerds := &Credentials{}
	err = selectData.Scan(&storageCerds.Password)
	if err != nil {
		//jika tidak ditemukan usernamya balikkan nilai 401
		if err == sql.ErrNoRows {
			fmt.Println(err.Error())
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		//jika ada error lain akan mengembalikan nilai 500
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//lakukan perbandingan inputan user hash dengan db
	if err = bcrypt.CompareHashAndPassword([]byte(storageCerds.Password), []byte(cerds.Password)); err != nil {
		//jika tidak cocok kemablikan nilai 401
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusUnauthorized)
	}

}
