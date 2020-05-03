package models

import "time"

type Asset struct {
	Name      string    `gorm:"primary_key"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt *time.Time
	DeletedAt *time.Time `sql:"index" json:"-"`
	Bid       uint       `gorm:"not null"`
	Ask       uint       `gorm:"not null"`
	Precision uint       `gorm:"not null"`
}
