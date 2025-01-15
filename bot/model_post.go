package bot

import (
	"time"

	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	TelegramId int `gorm:"size:255;uniqueIndex"`
	ChannelId  uint
	Channel    Channel
	Boosted    []*User `gorm:"many2many:boosts;"`
}

func getPostOrCreate(msgId int, c *Channel) (*Post, error) {
	p := &Post{}

	res := db.Preload("Channel").Where(&Post{TelegramId: msgId, ChannelId: c.ID}).FirstOrCreate(p)

	if res.Error != nil {
		loge(res.Error)
		return p, res.Error
	}

	return p, nil
}

func getPost(id int) *Post {
	p := &Post{}

	db.First(p, id)

	return p
}

func getBoostTasks(t time.Time) []*Post {
	var posts []*Post

	db.Where("created_at > ?", t.Add(-48*time.Hour)).Find(&posts)

	return posts
}
