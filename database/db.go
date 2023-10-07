package database

import (
	"github.com/jmoiron/sqlx"
)

func ConnDB() *sqlx.DB {
	db, err := sqlx.Connect("postgres",
		"user=postgres password=Hades330! dbname=postgres host=lingo.co05gj6uni3r.us-east-1.rds.amazonaws.com port=5432 sslmode=require")
	if err != nil {
		panic(err)
	}
	return db
}
