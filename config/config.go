package config

import "database/sql"

var (
	Db  *sql.DB
	Err error
)
