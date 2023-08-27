package models

type Literacy struct {
	User_id		string 		`json:"user_id" db:"user_id" `
	Word		string		`json:"word" db:"word" `
	Lang_iso	string		`json:"lang_iso" db:"lang_iso" `
	Known_level	int			`json:"known_level" db:"known_level" `
}

type InputGetLiteracy struct {
	User_id		string 		`json:"user_id"`
	Words		[]string	`json:"words"`
	Lang_iso	string		`json:"lang_iso"`
}

type InputParagraph struct {
	Words		[]Literacy	`json:"words"`
}
