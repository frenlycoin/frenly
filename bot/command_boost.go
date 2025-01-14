package bot

import (
	"strconv"
	"strings"

	"gopkg.in/telebot.v3"
)

func commandBoost(c telebot.Context, p string) error {
	u := getUser(c.Sender().ID)
	msg := lBoosted
	var btn *telebot.ReplyMarkup

	pids := strings.Split(p, "-")[1]
	pid, err := strconv.Atoi(pids)
	if err != nil {
		loge(err)
	}

	po := getPost(pid)

	for _, p := range u.Boosts {
		if p.ID == po.ID {
			msg = lAreadyBoosted
		}
	}

	if msg == lBoosted {
		u.Boosts = append(u.Boosts, po)
		err = db.Save(u).Error
		if err != nil {
			loge(err)
		}
	}

	unb := u.getUnboosted()

	if len(unb) > 0 {
		btn = getButtonLink(unb[0].Name, unb[0].Link)
		msg += "\n\nClick the button bellow for the next boost:"
	} else {
		msg += "\n\nYou have no more boosts available. 👍"
	}

	notifyWithButton(msg, c.Sender().ID, btn)

	return nil
}
