package bot

import (
	"fmt"
	"time"

	"gopkg.in/telebot.v3"
)

func commandJoin(c telebot.Context) error {
	var err error

	msg := fmt.Sprintf(`<b><u>Welcome, %s!</u></b>

You can read more about Frenly App here:

https://t.me/FrenlyNews/8`, c.Message().Sender.FirstName)

	m, err := b.Send(c.Chat(), msg, telebot.NoPreview)

	go func(m *telebot.Message) {
		time.Sleep(time.Second * 120)
		err := b.Delete(m)
		if err != nil {
			loge(err)
		}
	}(m)

	return err
}
