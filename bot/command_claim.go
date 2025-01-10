package bot

import (
	"fmt"

	"gopkg.in/telebot.v3"
)

func commandClaim(c telebot.Context) error {
	msg := ""
	u, err := getUserOrCreate(c)
	if err != nil {
		loge(err)
	}

	kv := &KeyValue{Key: "prizeWinner"}
	db.FirstOrCreate(kv, kv)

	if u.ID == uint(kv.ValueInt) {
		if u.AddressWithdraw != u.Code {
			msg = fmt.Sprintf(lClaimSuccess, u.Name)

			kv.ValueInt = 0
			db.Save(kv)

			send(100000000, u.AddressWithdraw, conf.Seed)
		} else {
			msg = fmt.Sprintf(lClaimError, u.Name)
		}
	} else {
		msg = fmt.Sprintf(lClaimFail, u.Name)
	}

	_, err = b.Send(c.Chat(), msg, telebot.NoPreview)
	if err != nil {
		loge(err)
	}

	return nil
}
