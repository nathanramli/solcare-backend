package models

import "time"

type Proposal struct {
	Address         string `gorm:"size:44;primaryKey"`
	CreatedAt       time.Time
	CampaignAddress string   `gorm:"size:44"`
	Campaign        Campaign `gorm:"foreignKey:CampaignAddress"`
	Url             string   `gorm:"size:255"`
}
