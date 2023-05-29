package models

import (
	"time"
)

type User struct {
	Id        int64     `json:"id"`
	Username  string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"passwordHash"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type Session struct {
	Id           int64     `json:"id"`
	UserId       int64     `json:"user_id"`
	SessionToken string    `json:"session_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}
type CreateUserReq struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}
type CreateUserRes struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}
type LoginUserReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUserRes struct {
	accessToken string
	Id          int64  `json:"id"`
	Username    string `json:"username"`
	Email       string `json:"email"`
}
