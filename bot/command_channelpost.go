package bot

import (
	"fmt"

	"gopkg.in/telebot.v3"
)

func commandChannelPost(c telebot.Context) error {
	var err error

	if c.Chat().ID != News {
		ch, err := getChannelOrCreate(c, nil)
		if err != nil {
			loge(err)
		}

		p, err := getPostOrCreate(c.Message().ID, ch)
		if err != nil {
			loge(err)
		}

		link := fmt.Sprintf("t.me/FrenlyRobot?start=b-%d", p.ID)

		if ch.Type == TypePost {
			fb := getFrenlyButton(link)
			_, err = b.Send(c.Chat(), lBoost, fb, telebot.NoPreview)
			if err != nil {
				loge(err)
			}
		} else if ch.Type == TypeButton {
			fb := getFrenlyButton(link)
			msg := c.Message()
			_, err = b.Edit(msg, fb)
			if err != nil {
				loge(err)
				db.Delete(p)
			}
		} else if ch.Type == TypeLink {
			msg := c.Message()
			text := msg.Text
			text += "\n\n<b><u>Boost Frenly Miner</u></b> 🚀\n" + link
			_, err = b.Edit(msg, text, telebot.NoPreview)
			if err != nil {
				loge(err)
				db.Delete(p)
			}
		}
	}

	return err
}
