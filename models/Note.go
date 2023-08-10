package models

import "github.com/gofrs/uuid"

type Note struct{
	User_id		uuid.UUID		`json:"user_id"`
	Word		string			`json:"word"`
	Note		string			`json:"note"`
}

type InputGetNote struct {
	User_id		string			`json:"user_id"`
	Words		[]string		`json:"words"`
}

type OutputGetNote struct {
	Notes		[]Note			`json:"notes"`
}

type InputUpdateNote struct {
	Notes		[]Note			`json:"input_notes"`	
}