package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	var err error

	host := os.Getenv("HOST_POSTGRES")
	port := os.Getenv("PORT_POSTGRES")
	user := os.Getenv("USER_POSTGRES")
	password := os.Getenv("PASSWORD_POSTGRES")
	dbname := os.Getenv("DATABASE_POSTGRES")

	// Buat connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s", host, port, user, password, dbname)

	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Connect Postgre Error: ", err)
	}
	err = DB.Ping()
	if err != nil {
		log.Fatal("Ping Postgre Error: ", err)
	}

	fmt.Println("Connected to Postgres!")
}
