package models

import "time"

type Organization struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"unique;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Users     []User
}

type User struct {
	ID           uint      `gorm:"primaryKey"`
	Email        string    `gorm:"unique;not null"`
	PasswordHash string    `gorm:"not null"`
	OrgID        uint      `gorm:"index"`
	Role         string    `gorm:"default:user"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Room struct {
	ID        uint      `gorm:"primaryKey"`
	OrgID     uint      `gorm:"index"`
	Name      string    `gorm:"not null"`
	Active    bool      `gorm:"default:true"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Queue     []QueueEntry
}

type QueueEntry struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"index"`
	RoomID    uint      `gorm:"index"`
	JoinedAt  time.Time
	Resolved  bool      `gorm:"default:false"`
}