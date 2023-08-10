package models

type Literacy struct {
	User_id		string 		`json:"user_id"`
	Word		string		`json:"word"`
	Lang_iso	string		`json:"lang_iso"`
	Known_level	int			`json:"known_level"`
}

type InputGetLiteracy struct {
	User_id		string 		`json:"user_id"`
	Words		[]string	`json:"words"`
	Lang_iso	string		`json:"lang_iso"`
}

type InputParagraph struct {
	User_id		string 		`json:"user_id"`
	Words		[]Literacy	`json:"words"`
}

type OutputParagraph struct {
	Level		[]Literacy	`json:"level"`
}