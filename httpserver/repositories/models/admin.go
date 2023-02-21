package models

type Admin struct {
	WalletAddress string `gorm:"size:44"`
	User          Users  `gorm:"foreignKey:WalletAddress"`
	IsSuperAdmin  *bool
}
