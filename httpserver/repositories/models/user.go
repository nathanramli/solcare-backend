package models

import "time"

type Users struct {
	WalletAddress  string `gorm:"primaryKey"`
	CreatedAt      time.Time
	Email          string
	FirstName      string
	LastName       string
	Gender         *bool `gorm:"default:true"`
	IsVerified     *bool `gorm:"default:false"`
	IsWarned       *bool `gorm:"default:false"`
	IdCardNumber   string
	ProfilePicture string
}
