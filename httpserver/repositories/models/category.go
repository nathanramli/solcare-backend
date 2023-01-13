package models

type Categories struct {
	Id          uint   `gorm:"primaryKey;autoIncrement"`
	Name        string `gorm:"size:24"`
	Description string `gorm:"size:255"`
}
