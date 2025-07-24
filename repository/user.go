package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"golang-final-project/config"
	"golang-final-project/middleware"
	"golang-final-project/models"
	"time"
)

func CreateUser(user models.User) string {
	var exists bool
	var userCredentials models.User
	fmt.Println("role input:", user.Role)

	existQuery := `SELECT EXISTS (SELECT 1 FROM users WHERE username = $1)`
	e := config.Db.QueryRow(existQuery, user.Username).Scan(&exists)
	if e != nil {
		return "Something went wrong"
	} else if exists {
		return "Username has been taken"
	} else {
		sqlStatement := `
		INSERT INTO users (id, username, password, token, expire_time, role)
		VALUES ($1, $2, $3, $4, $5, $6)
		Returning id, username, password, token, expire_time, role
		`
		config.Err = config.Db.QueryRow(
			sqlStatement,
			user.Id,
			user.Username,
			user.Password,
			nil,
			nil,
			user.Role,
		).Scan(
			&userCredentials.Id,
			&userCredentials.Username,
			&userCredentials.Password,
			&userCredentials.Token,
			&userCredentials.ExpireTime,
			&userCredentials.Role,
		)

		if config.Err != nil {
			panic(config.Err)
		} else {
			fmt.Printf("User: %+v\n", userCredentials)
		}

		return ""
	}
}

func Login(username string, password string) (*models.User, error) {
	var user models.User

	sqlStatement := `
	SELECT id, username, password, role
	FROM users 
	WHERE username = $1 LIMIT 1`

	config.Err = config.Db.QueryRow(sqlStatement, username).Scan(
		&user.Id,
		&user.Username,
		&user.Password,
		&user.Role,
	)

	if config.Err != nil {
		if errors.Is(config.Err, sql.ErrNoRows) {
			return nil, errors.New("account doesn't exist")
		} else {
			return nil, config.Err
		}
	} else {
		if err := middleware.CheckPasswordHash(password, user.Password); !err {
			return nil, errors.New("incorrect username or password")
		} else {
			return &user, nil
		}
	}
}

func AssignAccessToken(id int, token string, expire_time time.Time) (int64, error) {
	sqlStatement := `
	UPDATE users
	set token = $2, expire_time = $3
	WHERE id = $1
	`

	res, err := config.Db.Exec(
		sqlStatement,
		id,
		token,
		expire_time,
	)

	if err != nil {
		return 0, err
	} else {
		count, e := res.RowsAffected()
		if e != nil {
			return 0, err
		} else {
			return count, nil
		}
	}
}
