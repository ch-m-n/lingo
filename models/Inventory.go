package models

import "github.com/gofrs/uuid"

type Inventory struct {
	User_id		uuid.UUID		`json:"user_id" db:"user_id" `
	Head_id		uuid.UUID		`json:"head_id" db:"head_id" `
	Lang_iso	string			`json:"lang_iso" db:"lang_iso" `
}

type GetInventory struct {
	User_id		string			`json:"user_id"`
	Lang_iso	string			`json:"lang_iso"`
}	