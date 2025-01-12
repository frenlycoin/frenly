package bot

import (
	"gopkg.in/telebot.v3"
)

func commandChannelPost(c telebot.Context) error {
	var err error

	getChannelOrCreate(c.Chat().ID, nil)

	fb := getFrenlyButton()

	msg := c.Message()
	_, err = b.Edit(msg, fb)
	if err != nil {
		loge(err)
	}

	return err
}
