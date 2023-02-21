package models

import "time"

type Users struct {
	WalletAddress  string `gorm:"primaryKey"`
	CreatedAt      time.Time
	Email          string
	FirstName      string
	LastName       string
	Gender         *bool
	IsVerified     *bool
	IsWarned       *bool
	IdCardNumber   string
	ProfilePicture string
}
