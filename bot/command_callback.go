package bot

import (
	"fmt"
	"strings"

	"gopkg.in/telebot.v3"
)

func commandCallback(c telebot.Context) error {
	d := strings.Replace(c.Data(), "\f", "", -1)

	if d == "compound" {
		return commandCompound(c)
	} else if d == "done" {
		return commandDone(c)
	} else if d == "cancel" {
		return commandCancel(c)
	} else if strings.HasPrefix(d, fmt.Sprintf("t.me/%s", getRobotName())) {
		p := strings.Replace(d, fmt.Sprintf("t.me/%s?start=", getRobotName()), "", -1)
		return commandBoost(c, p, false)
	} else if strings.HasPrefix(d, "t.me/") {
		return commandChannelDelete(c)
	} else if strings.HasPrefix(d, "b-") {
		return commandBoost(c, d, false)
	}

	return nil
}
