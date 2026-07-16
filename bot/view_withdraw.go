package bot

import (
	"fmt"

	"gopkg.in/macaron.v1"
)

func viewWithdraw(ctx *macaron.Context) {
	wr := &GeneralResponse{Success: true}
	tgid := getTgId(ctx)

	if tgid != 0 {
		u := getUser(tgid)
		amount := int64((u.rewards(true) / 1000))

		if amount > 0 {
			exchange(u)

			notify(fmt.Sprintf("Withdraw: %s (%.9f TON)", u.Name, float64(amount)/float64(Mul9)), Frenly)

			notify(fmt.Sprintf(lCashOut, u.Name, float64(amount)/float64(Mul9)), Group)
		}
	}

	ctx.Header().Add("Access-Control-Allow-Origin", "*")

	ctx.JSON(200, wr)
}
