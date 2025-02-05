package bot

import (
	"strings"

	"gopkg.in/telebot.v3"
)

func commandCallback(c telebot.Context) error {
	d := strings.Replace(c.Data(), "\f", "", -1)

	if d == "compound" {
		return commandCompound(c)
	} else if strings.HasPrefix(d, "t.me/FrenlyRobot") {
		p := strings.Replace(d, "t.me/FrenlyRobot?start=", "", -1)
		return commandBoost(c, p, false)
	} else if strings.HasPrefix(d, "t.me/") {
		return commandChannelDelete(c)
	} else if strings.HasPrefix(d, "b-") {
		return commandBoost(c, d, true)
	}

	return nil
}
