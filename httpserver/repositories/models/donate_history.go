package models

import "time"

type DonateHistories struct {
	Id                   uint     `gorm:"primaryKey;autoIncrement"`
	UsersWalletAddress   string   `gorm:"size:44"`
	Users                Users    `gorm:"foreignKey:UsersWalletAddress"`
	CampaignAddress      string   `gorm:"size:44"`
	Campaign             Campaign `gorm:"foreignKey:CampaignAddress"`
	Date                 time.Time
	Amount               uint64
	TransactionSignature string `gorm:"size:256"`
}
