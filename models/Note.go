package models

import "github.com/gofrs/uuid"

type Note struct{
	User_id		uuid.UUID		`json:"user_id"`
	Word		string			`json:"word"`
	Note		string			`json:"note"`
}

