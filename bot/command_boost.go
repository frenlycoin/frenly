package bot

import (
	"fmt"
	"strconv"
	"strings"

	"gopkg.in/telebot.v3"
)

func commandBoost(c telebot.Context, p string, bg bool) error {
	u := getUser(c.Sender().ID)

	pids := strings.Split(p, "-")[1]
	pid, err := strconv.Atoi(pids)
	if err != nil {
		loge(err)
	}

	po := getPost(pid)

	msg := lBoosted
	var btn *telebot.ReplyMarkup

	if po.ID != 0 {
		if u.MiningTime.Before(po.CreatedAt) {
			for _, p := range u.Boosts {
				if p.ID == po.ID {
					msg = lAlreadyBoosted
				}
			}

			if msg == lBoosted {
				err := db.Model(u).Association("Boosts").Append(po)
				if err != nil {
					loge(err)
				}
			}
		} else {
			msg = lBoostTooOld
		}
	} else {
		msg = lAlreadyBoosted
	}

	unb := u.getUnboosted()

	if len(unb) > 0 {
		btn = getButtonsBoost(unb[0].Name, unb[0].Link)
		msg += fmt.Sprintf("\n\nBoosts Left: %d\n\nClick the button bellow for the next boost:", len(unb))
	} else {
		msg += "\n\nYou have no more boosts available. 👍"
	}

	if !bg || len(unb) == 0 {
		notifyWithButton(msg, c.Sender().ID, btn)
	}

	return nil
}
