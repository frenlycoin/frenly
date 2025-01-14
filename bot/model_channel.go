package bot

import (
	"gopkg.in/telebot.v3"
	"gorm.io/gorm"
)

type Channel struct {
	gorm.Model
	TelegramId int64 `gorm:"uniqueIndex"`
	Type       uint8 `gorm:"default:1"`
	OwnerId    *uint
	Owner      *User
	Name       string `gorm:"size:255"`
	Link       string `gorm:"size:255"`
}

func getChannelOrCreate(c telebot.Context, ownerId *uint) (*Channel, error) {
	ch := &Channel{}

	res := db.Preload("Owner").Where(&User{TelegramId: c.Chat().ID}).Attrs(
		&Channel{
			OwnerId: ownerId,
		}).FirstOrCreate(ch)

	if res.Error != nil {
		loge(res.Error)
		return ch, res.Error
	}

	ch.Name = c.Chat().Title
	ch.Link = c.Chat().Username
	db.Save(ch)

	return ch, nil
}

func getChannel(id int) *Channel {
	c := &Channel{}

	db.First(c, id)

	return c
}
