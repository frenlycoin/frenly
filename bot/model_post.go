package bot

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	TelegramId int64 `gorm:"size:255;uniqueIndex"`
	ChannelId  uint
	Channel    Channel
}
