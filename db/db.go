package db

import (
	"database/sql" // Provides an interface for database operations
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB // database connection

func ConnectDb() {
	dsn := "root:root@123@tcp(127.0.0.1:3306)/ecommerce"
	var err error
	DB, err = sql.Open("mysql", dsn) //  initializes a database connection (only validates arguments)

	if err != nil {
		log.Fatal("Error connting to databse: ", err) // Logs error message  and exits program
	}

	err = DB.Ping() // verifies the connection
	if err != nil {
		log.Fatal("Database connection failed", err)
	}
	fmt.Println("Database connected successfully!")
}

func GetDb() *sql.DB {
	return DB
}
