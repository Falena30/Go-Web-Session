package handler

import (
	"database/sql"

	//digunakan untuk mengubungkan ke mysql
	_ "github.com/go-sql-driver/mysql"
)

//Connect untuk menghubungkan ke mysql
func Connect() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/Go_Simple_Session")
	if err != nil {
		return nil, err
	}
	return db, nil
}
