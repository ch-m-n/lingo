package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type Head struct {
	Id			uuid.UUID		`json:"id" db:"id"`
	User_id		uuid.UUID		`json:"user_id" db:"user_id" `
	Title		string			`json:"title" db:"title" `
	Lang_iso	string			`json:"lang_iso" db:"lang_iso" `
	Img			string			`json:"img" db:"img" `
}

type Content struct{
	Id			uuid.UUID		`json:"id" db:"id" `
	User_id		uuid.UUID		`json:"user_id" db:"user_id" `
	Head_id		uuid.UUID		`json:"head_id" db:"head_id" `
	Lang_iso 	string			`json:"lang_iso" db:"lang_iso" `
	Body		string			`json:"body" db:"body" `
	Created_at	time.Time		`json:"created_at" db:"created_at" `
	Edited_at	time.Time		`json:"edited_at" db:"edited_at" `
}

type CreateContent struct {
	User_id		uuid.UUID		`json:"user_id" db:"user_id" `
	Title		string			`json:"title" db:"title" `
	Img			string			`json:"img" db:"img" `
	Lang_iso 	string			`json:"lang_iso" db:"lang_iso" `
	Body		string			`json:"body" db:"body" `
}

type RequestContent struct {
	My_id		string			`json:"my_id"`
	Head_id 	string			`json:"head_id"`
	Lang_iso	string			`json:"lang_iso"`
}

type RequestInventoryContent struct {
	User_id		string			`json:"user_id"`
	Lang_iso	string			`json:"lang_iso"`
}

type RequestAllContent struct {
	My_id 		string			`json:"my_id"`			
	Lang_iso 	string			`json:"lang_iso"`
}

