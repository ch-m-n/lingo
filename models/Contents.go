package models

import (
	"time"

	"github.com/gofrs/uuid"
)


type Content struct{
	Id			uuid.UUID		`json:"id" db:"id" `
	User_id		uuid.UUID		`json:"user_id" db:"user_id" `
	Title		string			`json:"title" db:"title" `
	Lang_iso 	string			`json:"lang_iso" db:"lang_iso" `
	Body		string			`json:"body" db:"body" `
	Created_at	time.Time		`json:"created_at" db:"created_at" `
	Edited_at	time.Time		`json:"edited_at" db:"edited_at" `
	Img			string			`json:"img" db:"img" `
}

type RequestContent struct {
	User_id		string			`json:"user_id"`
	Title		string			`json:"title"`
	Id 			string			`json:"id"`
	Lang_iso	string			`json:"lang_iso"`
}

type RequestAllContent struct {
	User_id		string			`json:"user_id"`
	Lang_iso 	string			`json:"lang_iso"`
}
