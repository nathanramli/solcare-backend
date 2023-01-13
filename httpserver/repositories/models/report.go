package models

import "time"

type Reports struct {
	Id                 uint `gorm:"primaryKey;autoIncrement"`
	CreatedAt          time.Time
	UsersWalletAddress string   `gorm:"size:44"`
	Users              Users    `gorm:"foreignKey:UsersWalletAddress"`
	CampaignAddress    string   `gorm:"size:44"`
	Campaign           Campaign `gorm:"foreignKey:CampaignAddress"`
	Description        string
}
