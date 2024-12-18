package bot

import (
	"fmt"

	"gopkg.in/telebot.v3"
)

func commandRanks(c telebot.Context) error {
	user := getUser(c.Sender().ID)
	var users []*User
	msg := lRanks
	showCaller := true

	db.Order("tmu desc").Find(&users)

	for i, u := range users {
		if i < 10 {
			if user.TelegramId == u.TelegramId {
				showCaller = false
				tmu := float64(u.TMU) / float64(Mul9)
				msg += fmt.Sprintf("\n<b><i>%d - %s</i></b> - <code>%.9f TMU (%d)</code>", i+1, u.Name, tmu, u.CompoundCount)
			} else {
				tmu := float64(u.TMU) / float64(Mul9)
				msg += fmt.Sprintf("\n<b>%d - %s</b> - <code>%.9f TMU (%d)</code>", i+1, u.Name, tmu, u.CompoundCount)
			}
		}
	}

	if showCaller {
		msg += "\n..."

		for i, u := range users {
			if user.TelegramId == u.TelegramId {
				tmu := float64(0)

				if i > 0 {
					u = users[i-1]

					tmu = float64(u.TMU) / float64(Mul9)
					msg += fmt.Sprintf("\n<b>%d - %s</b> - <code>%.9f TMU (%d)</code>", i, u.Name, tmu, u.CompoundCount)
				}

				u = users[i]

				tmu = float64(u.TMU) / float64(Mul9)
				msg += fmt.Sprintf("\n<b><i>%d - %s</i></b> - <code>%.9f TMU (%d)</code>", i+1, u.Name, tmu, u.CompoundCount)

				if i < len(users) {
					u = users[i+1]

					tmu = float64(u.TMU) / float64(Mul9)
					msg += fmt.Sprintf("\n<b>%d - %s</b> - <code>%.9f TMU (%d)</code>", i+2, u.Name, tmu, u.CompoundCount)
				}
			}
		}
	}

	_, err := b.Send(c.Chat(), msg, telebot.NoPreview)
	return err
}
