package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type Database struct{
	db *sql.DB
}

func NewDatabase() (*Database, error){
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=disable", dbUsername, dbPassword, dbPort, dbName)
	db, err := sql.Open("postgres", dsn)
	if err != nil{
		return nil, err
	}
	return &Database{db: db}, nil
}

func (d *Database) Close(){
	d.db.Close()
}

func (d *Database) GetDB() *sql.DB{
	return d.db
}