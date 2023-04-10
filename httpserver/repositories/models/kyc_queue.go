package models

import "time"

type KycQueues struct {
	Id                      uint   `gorm:"primaryKey;autoIncrement"`
	UsersWalletAddress      string `gorm:"size:44"`
	Users                   Users  `gorm:"foreignKey:UsersWalletAddress"`
	RequestedAt             time.Time
	Nin                     string
	IdCardPicture           string
	FacePicture             string
	SelfieWithIdCardPicture string
	Status                  uint8
}

const (
	KYC_STATUS_REQUESTED = 0
	KYC_STATUS_APPROVED  = 1
	KYC_STATUS_DECLINED  = 2
	KYC_STATUS_REMOVED   = 3
)
