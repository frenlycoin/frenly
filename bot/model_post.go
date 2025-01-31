package bot

import (
	"time"

	"gopkg.in/telebot.v3"
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	TelegramId int `gorm:"size:255"`
	ChannelId  uint
	Channel    Channel
	Boosted    []*User `gorm:"many2many:boosts;"`
	AlbumId    string  `gorm:"size:255;uniqueIndex"`
}

func getPostOrCreate(msg *telebot.Message, c *Channel) (*Post, error) {
	p := &Post{}

	aid := msg.AlbumID
	if aid == "" {
		aid = generateCode()
	}

	res := db.Preload("Channel").Where(&Post{TelegramId: msg.ID, ChannelId: c.ID, AlbumId: aid}).FirstOrCreate(p)

	if res.Error != nil {
		// loge(res.Error)
		return p, res.Error
	}

	return p, nil
}

func getPost(id int) *Post {
	p := &Post{}

	db.Preload("Channel").First(p, id)

	return p
}

func getPostByAlbumId(albumId string) *Post {
	p := &Post{}

	db.Preload("Channel").Where("album_id = ?", albumId).First(p)

	return p
}

func getBoostTasks(t time.Time) []*Post {
	var posts []*Post

	db.Where("created_at > ?", t).Find(&posts)

	return posts
}
