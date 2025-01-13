package bot

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	TelegramId int `gorm:"size:255;uniqueIndex"`
	ChannelId  uint
	Channel    Channel
}

func getPostOrCreate(msgId int, c *Channel) (*Post, error) {
	p := &Post{}

	res := db.Preload("Channel").Where(&Post{TelegramId: msgId}).Attrs(
		&Post{
			ChannelId: c.ID,
		}).FirstOrCreate(p)

	if res.Error != nil {
		loge(res.Error)
		return p, res.Error
	}

	return p, nil
}
