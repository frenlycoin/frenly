package bot

import (
	"fmt"
	"strconv"
	"strings"

	"gopkg.in/telebot.v3"
)

func commandBoost(c telebot.Context, p string) error {
	u := getUser(c.Sender().ID)

	// if u.ID == 0 {
	// 	_, err := getUserOrCreate(c)
	// 	if err != nil {
	// 		loge(err)
	// 	}

	// 	ab := getAppButton()
	// 	b.Send(c.Sender(), lStart, ab)
	// 	return nil
	// }

	msg := lBoosted
	var btn *telebot.ReplyMarkup

	pids := strings.Split(p, "-")[1]
	pid, err := strconv.Atoi(pids)
	if err != nil {
		loge(err)
	}

	po := getPost(pid)

	if po.ID != 0 {
		for _, p := range u.Boosts {
			if p.ID == po.ID {
				msg = lAreadyBoosted
			}
		}

		if msg == lBoosted {
			err := db.Model(u).Association("Boosts").Append(po)
			if err != nil {
				loge(err)
			}
		}
	} else {
		msg = lAreadyBoosted
	}

	unb := u.getUnboosted()

	if len(unb) > 0 {
		btn = getButtonsBoost(unb[0].Name, unb[0].Link)
		msg += fmt.Sprintf("\n\nBoosts Left: %d\n\nClick the button bellow for the next boost:", len(unb))
	} else {
		msg += "\n\nYou have no more boosts available. 👍"
	}

	notifyWithButton(msg, c.Sender().ID, btn)

	return nil
}
