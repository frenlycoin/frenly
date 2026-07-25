package bot

import (
	"gopkg.in/telebot.v3"
)

func commandCompound(c telebot.Context) error {
	u, err := getUserOrCreate(c)
	if err != nil {
		loge(err)
	}

	u.compound()

	rb := getRestartButtons(c)

	b.Send(c.Sender(), lCompounded, rb)

	return nil
}
