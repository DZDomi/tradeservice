package models

import (
	"github.com/google/uuid"
	"time"
)

type Trade struct {
	ID        uint       `gorm:"primary_key" json:"-"`
	PID       uuid.UUID  `gorm:"column:pid;type:varchar(36);index" json:"pid"`
	CreatedAt time.Time  `gorm:"not null" json:"created_at"`
	UpdatedAt *time.Time `json:"-"`
	DeletedAt *time.Time `sql:"index" json:"-"`
	User      uint       `gorm:"column:user_id" json:"user_id"`
	Wallet    uint       `gorm:"column:wallet_id" json:"wallet_id"`
	From      string     `gorm:"column:from" json:"from"`
	To        string     `gorm:"column:to" json:"to"`
	Executed  *time.Time `gorm:"column:executed" json:"executed"`
}
