package models

import "time"

type Asset struct {
	Name      string    `gorm:"primary_key"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt *time.Time
	DeletedAt *time.Time `sql:"index" json:"-"`
	Bid       uint64     `gorm:"not null"`
	Ask       uint64     `gorm:"not null"`
	Precision uint64     `gorm:"not null"`
}
