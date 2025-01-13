package bot

import (
	"gopkg.in/telebot.v3"
)

func commandChannelPost(c telebot.Context) error {
	var err error

	if c.Chat().ID != News {
		ch, err := getChannelOrCreate(c.Chat().ID, nil)
		if err != nil {
			loge(err)
		}

		if ch.Type == TypePost {
			fb := getFrenlyButton()
			_, err = b.Send(c.Chat(), lBoost, fb, telebot.NoPreview)
			if err != nil {
				loge(err)
			}
		} else if ch.Type == TypeButton {
			fb := getFrenlyButton()
			msg := c.Message()
			_, err = b.Edit(msg, fb)
			if err != nil {
				loge(err)
			}
		} else if ch.Type == TypeLink {
			msg := c.Message()
			text := msg.Text
			text += "\n\n<b><u>Boost Frenly Miner</u></b> 🚀\nt.me/FrenlyRobot?start=boost"
			_, err = b.Edit(msg, text, telebot.NoPreview)
			if err != nil {
				loge(err)
			}
		}
	}

	return err
}
