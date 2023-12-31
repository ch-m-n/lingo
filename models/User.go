package models

import (
	"time"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct{
	Id			uuid.UUID  	`json:"id" db:"id"`
	Username	string		`json:"username" db:"username"`
	Email		string		`json:"email" db:"email"`
	Pwd			string		`json:"pwd" db:"pwd"`
	Created_at	time.Time	`json:"created_at" db:"created_at"`
	Edited_at	time.Time	`json:"edited_at" db:"edited_at"`
	Verified	bool		`json:"verified" db:"verified"`
}

type CreateUser struct{
	// Position	string		`json:"position"`
	Username	string		`json:"username"`
	Email		string		`json:"email"`
	Pwd			string		`json:"pwd"`
}

type GetUser struct {
	Email		string		`json:"email"`
	Pwd 		string		`json:"pwd"`
}

type EditUser struct {
	Id 			string		`json:"user_id"`
	Username	string 		`json:"username"`
	Pwd			string		`json:"pwd"`
}

func PassHash(pwd string) []byte {
	hashPwd, err := bcrypt.GenerateFromPassword([]byte(pwd), 12)
	if err != nil {
		panic(err)
	}
	return hashPwd
}

func VerifyHash(pwd string, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(pwd))
	if err != nil {
		return false
	}
	return true
}