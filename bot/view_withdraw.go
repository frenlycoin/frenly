package bot

import (
	"fmt"
	"math"

	"gopkg.in/macaron.v1"
)

func viewWithdraw(ctx *macaron.Context) {
	wr := &GeneralResponse{Success: true}
	tgid := getTgId(ctx)

	if tgid != 0 {
		u := getUser(tgid)
		amount := int64((u.rewards(true) / 1000))

		if amount > 0 {
			amountOut, err := exchange(u)
			if err != nil {
				loge(err)
			} else {
				tonAmount := int64(math.Round(amountOut * float64(Mul9)))

				notifyCashout(u, tonAmount, Frenly)

				notify(fmt.Sprintf(lCashOutPublic, u.Name, float64(tonAmount)/float64(Mul9)), Group)
			}
		}
	}

	ctx.Header().Add("Access-Control-Allow-Origin", "*")

	ctx.JSON(200, wr)
}
