package bot

import (
	"fmt"
	"time"

	"gopkg.in/macaron.v1"
)

func viewWithdraw(ctx *macaron.Context) {
	wr := &GeneralResponse{Success: true}
	tgid := getTgId(ctx)

	if tgid != 0 {
		u := getUser(tgid)
		amount := int64((u.rewards(true) / 10) - 5000000)

		if amount > 0 {
			u.LastUpdated = time.Now()
			u.CycleCountTotal += u.CycleCount
			u.CycleCount = 1
			// u.delayedUpdateBalance()
			if err := db.Save(u).Error; err != nil {
				loge(err)
			} else {
				send(amount, u.AddressWithdraw, conf.Seed)

				notify(fmt.Sprintf("Withdraw: %s (%.9f TON)", u.Name, float64(amount)/float64(Mul9)), Frenly)
			}
		}
	}

	ctx.Header().Add("Access-Control-Allow-Origin", "*")

	ctx.JSON(200, wr)
}
