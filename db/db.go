package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	HOST = "localhost"
	PORT = 5432
)

var ErrNoMatch = fmt.Errorf("no matching record")

type Database struct {
	Conn *sql.Tx
}
type Database2 struct {
	Conn *sql.DB
}
var Conn *sql.DB
var Err error
func Initialize2 () (Database, error) {
	db := Database{}

	Conn, Err = sql.Open("postgres", "user=postgres password=postgres dbname=social_network sslmode=disable")
	if Err != nil {
		return db, Err
	}
	x,_ := Conn.Begin()
	db.Conn = x
	if Err!= nil {
		return db, Err
	}
	log.Println("Database connection established")
	return db, nil
}
func Initialize () () {
	Conn, Err = sql.Open("postgres", "user=postgres password=postgres dbname=social_network sslmode=disable")
}