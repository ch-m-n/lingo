package models

import (
	"time"

	"github.com/gofrs/uuid"
)


type Content struct{
	Id			uuid.UUID		`json:"id"`
	User_id		uuid.UUID		`json:"user_id"`
	Title		string			`json:"title"`
	Lang_iso 	string			`json:"lang_iso"`
	Body		string			`json:"body"`
	Created_at	time.Time		`json:"created_at"`
	Edited_at	time.Time		`json:"edited_at"`
	Img			string			`json:"img"`
}

type CreateContent struct{
	User_id		uuid.UUID		`json:"user_id"`
	Title		string			`json:"title"`
	Lang_iso 	string			`json:"lang_iso"`
	Body		string			`json:"body"`
	Img			string			`json:"img"`
}
type EditContent struct{
	Content_id	uuid.UUID		`json:"content_id"`
	User_id		uuid.UUID		`json:"user_id"`
	Title		string			`json:"title"`
	Lang_iso 	string			`json:"lang_iso"`
	Body		string			`json:"body"`
	Img			string			`json:"img"`
}

type RequestContent struct {
	Title		string			`json:"title"`
	Lang_iso 	string			`json:"lang_iso"`
}