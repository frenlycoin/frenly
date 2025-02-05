package bot

import (
	"fmt"
	"log"
	"strings"

	"gopkg.in/telebot.v3"
)

func commandChannelPost(c telebot.Context) error {
	var err error
	// var msg *telebot.Message

	if c.Chat().ID != News {
		ch, err := getChannelOrCreate(c, nil)
		if err != nil {
			loge(err)
		}

		p := getPostByAlbumId(c.Message().AlbumID)

		if p.ID == 0 || len(c.Message().AlbumID) == 0 {
			p, err = getPostOrCreate(c.Message(), ch)
			if err != nil {
				if strings.Contains(err.Error(), "duplicate key value violates") {
					log.Println(err)
				} else {
					loge(err)
				}
				return err
			}

			link := fmt.Sprintf("t.me/FrenlyRobot?start=b-%d", p.ID)
			bb := fmt.Sprintf("b-%d", p.ID)

			if len(c.Message().AlbumID) > 0 && ch.Type != TypePost {
				ch.Type = TypePost
			}

			if ch.Type == TypePost {
				fb := getFrenlyButtons(bb)
				m, err := b.Send(c.Chat(), lBoost, fb, telebot.NoPreview)
				if err != nil {
					loge(err)
				}
				p.TelegramId = m.ID
				db.Save(p)
			} else if ch.Type == TypeButton {
				fb := getFrenlyButtons(bb)
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
	}

	return err
}
