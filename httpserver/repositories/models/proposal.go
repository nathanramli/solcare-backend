package models

import "time"

type Proposal struct {
	Address            string `gorm:"size:44,primaryKey"`
	CreatedAt          time.Time
	UsersWalletAddress string   `gorm:"size:44"`
	Users              Users    `gorm:"foreignKey:UsersWalletAddress"`
	CampaignAddress    string   `gorm:"size:44"`
	Campaign           Campaign `gorm:"foreignKey:CampaignAddress"`
	Title              string   `gorm:"size:255"`
	Url                string   `gorm:"size:255"`
}
