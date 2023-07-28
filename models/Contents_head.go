package models

import (
	"time"

	"github.com/gofrs/uuid"
)


type Content_head struct{
	Id		uuid.UUID		`json:"id"`
	User_id	uuid.UUID		`json:"user_id"`
	Title	string			`json:"title"`
	Lang_iso string			`json:"lang_iso"`
	Created_at	time.Time	`json:"created_at"`
	Edited_at	time.Time	`json:"edited_at"`
}