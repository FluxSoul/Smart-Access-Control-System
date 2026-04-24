package models

type Node struct {
	ID     int `json:"id" db:"id"`
	UserId int `json:"user_id" db:"user_id"`
}
