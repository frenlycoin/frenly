package bot

import "gorm.io/gorm"

type Channel struct {
	gorm.Model
	TelegramId int64 `gorm:"size:255;uniqueIndex"`
	OwnerId    *uint
	Owner      *User
}
