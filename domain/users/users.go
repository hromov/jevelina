package users

import "time"

type User struct {
	ID           uint64
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    time.Time
	Name         string
	Email        string
	Hash         string
	Distribution float32
	Role         string
}

type Role struct {
	ID        uint8
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
	Priority  uint8
	Role      string
}
