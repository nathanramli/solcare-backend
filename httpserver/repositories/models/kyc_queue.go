package models

import "time"

type KycQueues struct {
	Id                      uint   `gorm:"primaryKey;autoIncrement"`
	UsersWalletAddress      string `gorm:"size:44"`
	Users                   Users  `gorm:"foreignKey:UsersWalletAddress"`
	RequestedAt             time.Time
	IdCardPicture           string
	FacePicture             string
	SelfieWithIdCardPicture string
	Status                  uint8
}
