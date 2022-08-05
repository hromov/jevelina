package models

import (
	"time"

	"github.com/hromov/jevelina/domain/users"
	"gorm.io/gorm"
)

type User struct {
	ID           uint64 `gorm:"primaryKey"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	Name         string         `gorm:"size:32"`
	Email        string         `gorm:"size:128; unique"`
	Hash         string         `gorm:"size:128; unique"`
	Distribution float32        `gorm:"type:decimal(2,2);"`
	// Events    []Event
	RoleID *uint8
	Role   Role `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

func (u *User) ToDomain() users.User {
	return users.User{
		ID:           u.ID,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
		DeletedAt:    u.DeletedAt.Time,
		Name:         u.Name,
		Email:        u.Email,
		Hash:         u.Hash,
		Distribution: u.Distribution,
		Role:         u.Role.Role,
	}
}

type Role struct {
	ID        uint8 `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Priority  uint8
	Role      string `gorm:"unique;size:32"`
}

func (r *Role) ToDomain() users.Role {
	return users.Role{
		ID:        r.ID,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
		DeletedAt: r.DeletedAt.Time,
		Priority:  r.Priority,
		Role:      r.Role,
	}
}
