package bot

import (
	"strconv"
	"strings"

	"gopkg.in/telebot.v3"
)

func commandBoost(c telebot.Context, p string) error {
	u := getUser(c.Sender().ID)

	pids := strings.Split(p, "-")[1]
	pid, err := strconv.Atoi(pids)
	if err != nil {
		loge(err)
	}

	po := getPost(pid)

	u.Boosts = append(u.Boosts, po)
	db.Save(u)

	notify(lBoosted, c.Sender().ID)

	return nil
}
