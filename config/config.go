package config

import "database/sql"

var (
	Db  *sql.DB
	Err error
)

// var BaseUrl = "http://localhost:8080/"
var BaseUrl = "https://sanber-golang-final-project-production.up.railway.app/"
