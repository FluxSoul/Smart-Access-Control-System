package models

import (
	"time"
)

type EmpxAdmin struct {
	ID             int       `json:"id" db:"id"`
	Username       string    `json:"username" db:"username"`
	Password       string    `json:"password" db:"password"`
	Status         int8      `json:"status" db:"status"`
	CreatedTime    time.Time `json:"createdTime" db:"created_time"`
	Token          string    `json:"token" db:"token"`
	TokenExpiresAt time.Time `json:"tokenExpiresAt" db:"token_expires_at"`
}

type GetEmpxAdmin struct {
	ID       int    `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
	Status   int8   `json:"status" db:"status"`
}
