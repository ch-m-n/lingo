package models

import "time"
type User struct {
    Id			string		`json:"id"`
    Username	string		`json:"username"`
    Email 		string		`json:"email"`
    Pwd 	    string		`json:"pwd"`
    Created_At	time.Time	`json:"created_at" gorm:"column:created_at"`
}
