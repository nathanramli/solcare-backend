package models

import "time"

type Campaign struct {
	Address      string `gorm:"primaryKey;size:44"`
	CreatedAt    time.Time
	OwnerAddress string `gorm:"size:44"`
	Owner        Users  `gorm:"foreignKey:OwnerAddress"`
	Title        string `gorm:"size:255"`
	Description  string `gorm:"size:5000"`
	CategoryId   uint
	Categories   Categories `gorm:"foreignKey:CategoryId"`
	AmountTarget uint64
	DateTarget   time.Time
	Banner       string `gorm:"size:255"`
	DonatedFunds uint64
	Status       uint8
	Delisted     bool
}
