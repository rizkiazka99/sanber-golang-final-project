package config

import "database/sql"

var (
	Db  *sql.DB
	Err error
)

var BaseUrl = "http://localhost:8080/"
