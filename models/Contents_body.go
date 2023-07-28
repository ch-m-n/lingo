package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type Content_body struct{
	Id		uuid.UUID		`json:"id"`
	Head_id	uuid.UUID		`json:"head_id"`
	Body	string			`json:"body"`
	Created_at	time.Time	`json:"created_at"`
	Edited_at	time.Time	`json:"edited_at"`
}