package models

import "github.com/gofrs/uuid"

type Overview struct{
	User_id		uuid.UUID		`json:"user_id"`
	Lang_iso	string			`json:"lang_iso"`
	Word		string			`json:"word"`
	Literacy	int				`json:"literacy"`
}


