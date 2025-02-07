package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func ConnectDb() {
	dsn := "root:root@123@tcp(127.0.0.1:3306)/ecommerce"
	var err error
	DB, err = sql.Open("mysql", dsn)

	if err != nil {
		log.Fatal("Error connting to databse: ", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("Database connection failed", err)
	}
	fmt.Println("Database connected successfully!")
}

func GetDb() *sql.DB {
	return DB
}
