package models

import "time"

type Transaction struct {
	Signature          string   `gorm:"size:256;primaryKey"`
	UsersWalletAddress string   `gorm:"size:44"`
	Users              Users    `gorm:"foreignKey:UsersWalletAddress"`
	CampaignAddress    string   `gorm:"size:44"`
	Campaign           Campaign `gorm:"foreignKey:CampaignAddress"`
	CreatedAt          time.Time
	Amount             uint64
	Type               uint8
}
