package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gomodule/redigo/redis"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

//Credentials untuk menampung json dan db
type Credentials struct {
	Password string `json:"password", db:"Password"`
	Username string `json:"username", db:"Username"`
}

//simpan redis package menjadi sebuah variabel
var Chace redis.Conn

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

func SignInSession(w http.ResponseWriter, r *http.Request) {
	//hubungkan db
	db, err := Connect()
	if err != nil {
		fmt.Println(err.Error())
	}
	defer db.Close()

	//buat variabel untuk menampung struct cres
	var cres Credentials
	//dapatkan nilai JSON body lalu decode kedalam cres
	err = json.NewDecoder(r.Body).Decode(&cres)
	if err != nil {
		//jika terjadi error kembalikan nilai 500 atau http bad request
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//Read data berdasarakan JSON jika ada
	selectData := db.QueryRow("SELECT * FROM User_DB WHERE Username = ?", cres.Username)

	//jika password ada dan sama dengan password yang ada sama dengan password
	//yang diberikan maka bisa ke langkah selanjunya
	//tetapi jika tidak kembalikan status "Unauthorized"
	if selectData == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	//buat random session token
	sessionToken := uuid.NewV4().String()
	//masukkan token kedalam cache dan juha passwordnya
	//token memiliki waktu kadaluarsa 120 detik
	_, err = Chace.Do("SETEX", sessionToken, "120", cres.Username)
	if err != nil {
		//jika terjadi error pada saat setting cache kemabalikan nilai interval server error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//set client cookie untuk "session_Token" sebagai session token yang sebelumya kita genereate
	//dan set juga waktu exp menjadi 120 sama dengan cache
	http.SetCookie(w, &http.Cookie{
		Name:    "session_Token",
		Value:   sessionToken,
		Expires: time.Now().Add(120 * time.Second),
	})
}

func InitCache() {
	//init redis connection dengan localhost
	conn, err := redis.DialURL("redis://localhost")
	if err != nil {
		panic(err)
	}
	//masukkan nilai conn ke variabel global cache
	Chace = conn
}

func Welcome(w http.ResponseWriter, r *http.Request) {
	//ambil cookies dari request
	c, err := r.Cookie("session_Token")
	if err != nil {
		if err == http.ErrNoCookie {
			//jika cookie tidak ada kemablikan nilai status Unautorized
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		//untuk tipe lainnya kemablikan nilai bad request
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//simpan nilai sessionToken
	sessionToken := c.Value
	//sekarang kita bisa mendapatkan nama dari cache
	respones, err := Chace.Do("GET", sessionToken)
	if err != nil {
		//jika pada saat fetched chace terjadi error kemablikan nilai internas server error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if respones == nil {
		//jika session Token tidak ada di cache maka kemabalikan status Unautoried
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	//kembalikan pesan welcome kepada user
	w.Write([]byte(fmt.Sprintf("Welcome %s", respones)))
}

func Refresh(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_Token")
	if err != nil {
		if err != http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	storageSession := c.Value
	response, err := Chace.Do("GET", storageSession)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if response == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	//buat session baru untuk user sekarng
	newSessionToken := uuid.NewV4().String()
	_, err = Chace.Do("SETEX", newSessionToken, "120", response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//hapus session yang lama
	_, err = Chace.Do("DEl", storageSession)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//jadikan newsessiontoken menjadi session baru
	http.SetCookie(w, &http.Cookie{
		Name:    "session_Token",
		Value:   newSessionToken,
		Expires: time.Now().Add(120 * time.Second),
	})
}
