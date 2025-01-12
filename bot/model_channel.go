package bot

import (
	"gorm.io/gorm"
)

type Channel struct {
	gorm.Model
	TelegramId int64 `gorm:"size:255;uniqueIndex"`
	Type       uint8 `gorm:"default:1"`
	OwnerId    *uint
	Owner      *User
}

func getChannelOrCreate(tgid int64, ownerId *uint) (*Channel, error) {
	c := &Channel{}

	res := db.Preload("Owner").Where(&User{TelegramId: tgid}).Attrs(
		&Channel{
			OwnerId: ownerId,
		}).FirstOrCreate(c)

	if res.Error != nil {
		loge(res.Error)
		return c, res.Error
	}

	return c, nil
}
