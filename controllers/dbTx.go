package controllers

import "lingo/database"

var tx = database.ConnDB().MustBegin()