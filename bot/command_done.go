package bot

import (
	"gopkg.in/telebot.v3"
)

func commandDone(c telebot.Context) error {
	if c.Message() == nil {
		return nil
	}

	emptyMarkup := &telebot.ReplyMarkup{}
	_, err := b.Edit(c.Message(), c.Message().Text, emptyMarkup, telebot.NoPreview)
	if err != nil {
		return err
	}

	return nil
}
