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

	rb := getRestartButtons(c)

	frenAmount := formatNumber(float64(u.TMU) / float64(Mul9))
	msg := fmt.Sprintf(lCompounded+"\n\n<b>Staked FREN:</b> <code>%s</code>", frenAmount)

	b.Send(c.Sender(), msg, rb)

	return nil
}
