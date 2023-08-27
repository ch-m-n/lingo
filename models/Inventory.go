package models

import "github.com/gofrs/uuid"

type Inventory struct {
	User_id		uuid.UUID		`json:"user_id"`
	Head_id		uuid.UUID		`json:"head_id"`
	Lang_iso	string			`json:"lang_iso"`
}

