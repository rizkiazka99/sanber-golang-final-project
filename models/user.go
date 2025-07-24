package models

import (
	"database/sql"
	"time"
)

type User struct {
	Id         int            `json:"id"`
	Username   string         `json:"username"`
	Password   string         `json:"password"`
	Token      sql.NullString `json:"token"`
	ExpireTime sql.NullTime   `json:"expire_time"`
	Role       string         `json:"role"`
}

type UserResponse struct {
	Id         int64      `json:"id"`
	Username   string     `json:"username"`
	Password   string     `json:"password"`
	Role       string     `json:"role"`
	Token      *string    `json:"token"`
	ExpireTime *time.Time `json:"expire_time"`
}

type LoggedIn struct {
	Id                  int       `json:"id"`
	Username            string    `json:"username"`
	Role                string    `json:"role"`
	AccessToken         string    `json:"access_token"`
	TokenExpirationTime time.Time `json:"token_expiration_time"`
}

func BuildUserResponse(u User) UserResponse {
	var token *string
	if u.Token.Valid {
		token = &u.Token.String
	}

	var expireTime *time.Time
	if u.ExpireTime.Valid {
		expireTime = &u.ExpireTime.Time
	}

	return UserResponse{
		Id:         int64(u.Id),
		Username:   u.Username,
		Password:   u.Password,
		Role:       u.Role,
		Token:      token,
		ExpireTime: expireTime,
	}
}
