package bot

import (
	"fmt"

	"gopkg.in/telebot.v3"
)

func commandCompound(c telebot.Context) error {
	u, err := getUserOrCreate(c)
	if err != nil {
		loge(err)
	}

	u.compound()

	r := u.Referrer

	if r != nil && r.ID != 0 && u.TMU >= 10100000 && !u.ReferralActive {
		r.TMU += 2500000
		if err := db.Save(r).Error; err != nil {
			loge(err)
		}
		msg := fmt.Sprintf(lNewRefFCS, float64(2500000)/float64(Mul9))
		notify(msg, r.TelegramId)

		u.ReferralActive = true
		if err := db.Save(u).Error; err != nil {
			loge(err)
		}
	}

	ab := getAppButton()

	b.Send(c.Sender(), lCompounded, ab)

	return nil
}
