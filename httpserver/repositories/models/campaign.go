package models

import (
	"github.com/gagliardetto/solana-go"
	"time"
)

type Campaign struct {
	Address      string `gorm:"primaryKey;size:44"`
	CreatedAt    time.Time
	OwnerAddress string `gorm:"size:44"`
	Owner        Users  `gorm:"foreignKey:OwnerAddress"`
	Title        string `gorm:"size:255"`
	Description  string `gorm:"size:5000"`
	CategoryId   uint
	Categories   Categories `gorm:"foreignKey:CategoryId"`
	Banner       string     `gorm:"size:255"`
	Evidence     string     `gorm:"size:255"`
	Status       uint8
	Delisted     *bool `gorm:"default:false"`
}

const (
	CAMPAIGN_STATUS_FUNDED = 4

	EVIDENCE_STATUS_WAITING   = 0
	EVIDENCE_STATUS_REQUESTED = 1
	EVIDENCE_STATUS_SUCCESS   = 2
	EVIDENCE_STATUS_FAILED    = 3
)

type CampaignBlockchainData struct {
	Discriminator [8]byte
	Owner         solana.PublicKey
	CreatedAt     int64
	HeldDuration  int64
	TargetAmount  uint64
	FundedAmount  uint64
	CampaignVault solana.PublicKey
	Status        uint8
}
