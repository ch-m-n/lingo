package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// func ConnDB() *gorm.DB {

// 	dsn := "host=" + "lingo.co05gj6uni3r.us-east-1.rds.amazonaws.com" +
// 		" user=" + "postgres" +
// 		" password=" + "Hades330!" +
// 		" dbname=" + "postgres" +
// 		" port=" + "5432" + " sslmode=require TimeZone=Asia/Shanghai"
// 	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
// 	if err != nil {
// 		panic("Failed to connect to database")
// 	}

// 	return db
// }

func ConnDB() *gorm.DB {

	dsn := "host=" + "localhost" +
		" user=" + "admin" +
		" password=" + "6457" +
		" dbname=" + "lingo" +
		" port=" + "5432" + " sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}

	return db
}