package MiscFunc

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

//GetPwd adalah fungsi yang digunakan untuk mendapatkan inputan password user
func GetPwd() []byte {
	//buat user memasukkan password apapun itu
	fmt.Println("Masukkan password")

	//variabel untuk menyimpan inputan user
	var storage string

	//membaca inputan dari user
	_, err := fmt.Scan(&storage)
	if err != nil {
		fmt.Println(err)
	}
	//mengembalikan nilai yang diinputkan oleh user beruba byte slice
	//nanti
	return []byte(storage)
}

//HashAndSalt digunakan untuk merubah string password menjadi Hash + salt
func HashAndSalt(pwd []byte) string {
	//gunakan GenerateFromPassword untuk hash dan salt password.
	//MinCost hanyalah integer yang bernilai constant yang diberikan oleh bcrypt paket
	//bersama dengan defaultcost dan maxcost
	//costnya bisa berapapun sesuai keingan tapi tidak boleh lebih kecil dari mincost
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		fmt.Println(err)
	}

	//kemablikan nilai hash
	//karena hash bertipe data []byte maka rubah jadi string
	return string(hash)
}

func ComparePassword(hashPwd string, plainText []byte) bool {
	//karena biasanya password di DB berupa string kita harus merubahnya terlebih dahulu
	//merubah menjadi byte
	byteHash := []byte(hashPwd)

	//lakukan perbandingan antara bytehash dan plaintext
	err := bcrypt.CompareHashAndPassword(byteHash, plainText)
	//jika ada err tampilkan
	if err != nil {
		fmt.Println(err)
		//kembalikan nilai false jika ada error
		return false
	}

	//jika error tidak nil maka kembalikan nilai true
	return true
}
