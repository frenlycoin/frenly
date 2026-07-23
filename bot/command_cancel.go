package bot

import (
	"gopkg.in/telebot.v3"
)

func commandCancel(c telebot.Context) error {
	if c.Message() == nil {
		return nil
	}

	if c.Sender() != nil && !isAdmin(c.Sender().ID) {
		_, err := b.Send(c.Message().Chat, "Forbidden. Only admins are allowed to click this button.", telebot.NoPreview)
		if err != nil {
			return err
		}
		return nil
	}

	keys := []string{"dexFren", "dexGram", "dexLastPrice"}
	for _, key := range keys {
		mainKv := &KeyValue{Key: key}
		oldKv := &KeyValue{Key: key + "Old"}

		if err := db.Where("key = ?", mainKv.Key).FirstOrCreate(mainKv).Error; err != nil {
			return err
		}
		if err := db.Where("key = ?", oldKv.Key).FirstOrCreate(oldKv).Error; err != nil {
			return err
		}

		mainKv.ValueInt = oldKv.ValueInt
		if err := db.Save(mainKv).Error; err != nil {
			return err
		}
	}

	if c.Sender() != nil {
		u := getUser(c.Sender().ID)
		if u.ID != 0 {
			u.TMU += u.PayoutAmount
			u.PayoutAmount = 0
			if err := db.Save(u).Error; err != nil {
				return err
			}
		}
	}

	emptyMarkup := &telebot.ReplyMarkup{}
	_, err := b.Edit(c.Message(), c.Message().Text, emptyMarkup, telebot.NoPreview)
	if err != nil {
		return err
	}

	notify(lTradeCanceled, GroupHall)

	return nil
}
