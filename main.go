package main

import (
	"fmt"
	"net/http"

	"Go-Web-Session-Vendor/handler"

	"github.com/gorilla/mux"
)

func main() {
	/*
		for {
			//masukkan password dan lanjutkan dengan generate hash dan salt
			//masukkan password yang pertama
			pwd := MiscFunc.GetPwd()
			hash := MiscFunc.HashAndSalt(pwd)

			//masukkan password yang kedua
			//dan bandingkan apakah sama atau tidak

			pwd2 := MiscFunc.GetPwd()
			CompareHash := MiscFunc.ComparePassword(hash, pwd2)
			fmt.Println("apakah password cocok? ", CompareHash)

		}
	*/

	r := mux.NewRouter()

	r.HandleFunc("/signup", handler.HandleSignUp)
	r.HandleFunc("/coba", handler.HandlerSignIn)
	fmt.Println("server start at localhost:8080")
	http.ListenAndServe(":8080", r)
}
