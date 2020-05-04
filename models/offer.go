package models

import (
	"github.com/google/uuid"
	"time"
)

type Trade struct {
	ID         uint       `gorm:"primary_key" json:"-"`
	PID        uuid.UUID  `gorm:"column:pid;type:varchar(36);index" json:"pid"`
	CreatedAt  time.Time  `gorm:"not null" json:"created_at"`
	UpdatedAt  *time.Time `json:"-"`
	DeletedAt  *time.Time `sql:"index" json:"-"`
	User       uint       `gorm:"column:user_id" json:"user_id"`
	FromWallet uint       `gorm:"column:from_wallet_id" json:"from_wallet_id"`
	ToWallet   uint       `gorm:"column:to_wallet_id" json:"to_wallet_id"`
	From       string     `gorm:"column:from" json:"from"`
	To         string     `gorm:"column:to" json:"to"`
	Accepted   *time.Time `gorm:"column:accepted" json:"accepted"`
	Executed   *time.Time `gorm:"column:executed" json:"executed"`
}
