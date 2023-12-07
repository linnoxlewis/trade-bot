package domain

import "github.com/google/uuid"

type User struct {
	Id       uuid.UUID
	TgId     int64
	Username string
}
