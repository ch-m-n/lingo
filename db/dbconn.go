package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Dbconn()(*gorm.DB){
	dsn := "host=localhost user=postgres password=password dbname=lingo port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	switch err != nil{
	case true:
		print("Unable to connect")
	}
	return db
}

