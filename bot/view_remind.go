package bot

import (
	"fmt"
	"time"

	macaron "gopkg.in/macaron.v1"
)

func viewRemind(ctx *macaron.Context) {
	rr := &RemindResponse{}
	tgid := getTgId(ctx)
	ctx.Header().Add("Access-Control-Allow-Origin", "*")

	if tgid == 0 {
		rr.Success = false
		rr.ErrorMessage = "Invalid telegram id"
		ctx.JSON(200, rr)
		return
	}

	u := getUser(tgid)
	if u.ID == 0 {
		rr.Success = false
		rr.ErrorMessage = "User not found"
		ctx.JSON(200, rr)
		return
	}

	if time.Since(u.LastReminder).Hours() < 24 {
		rr.Success = false
		rr.ErrorMessage = fmt.Sprintf("Reminder already sent. Try again in %.0f hours.", 24-time.Since(u.LastReminder).Hours())
		ctx.JSON(200, rr)
		return
	}

	var refUsers []*User
	db.Where("referrer_id = ?", u.ID).Find(&refUsers)

	sent := 0
	for _, ru := range refUsers {
		if ru.isActive() || ru.BotBlocked {
			continue
		}

		weeklyFREN := uint64(float64(ru.TMU) * 604800 / (2400 * 3600))
		weeklyFRENFormatted := formatNumber(float64(weeklyFREN) / float64(Mul9))

		senderName := u.Name
		if u.Code != "" {
			senderName = "@" + u.Code
		}

		priceKv := &KeyValue{Key: "dexLastPrice"}
		if err := db.Where("key = ?", priceKv.Key).FirstOrCreate(priceKv).Error; err == nil && priceKv.ValueInt > 0 {
			price := float64(priceKv.ValueInt) / float64(Mul9)
			tonValue := (float64(weeklyFREN) / float64(Mul9)) * price
			msg := fmt.Sprintf(
				"🔵 <b><u>Your miner is inactive!</u></b>\n\n"+
					"You are losing <b>%s FREN (%s GRAM)</b> every week by not mining!\n\n"+
					"<i>This notification was sent by your referrer %s.</i>\n\n"+
					"Restart your miner and compound your rewards to earn more.",
				weeklyFRENFormatted,
				formatNumber(tonValue),
				senderName,
			)
			notifyRestartInactive(msg, ru.TelegramId)
		} else {
			msg := fmt.Sprintf(
				"🔵 <b><u>Your miner is inactive!</u></b>\n\n"+
					"You are losing <b>%s FREN</b> every week by not mining!\n\n"+
					"<i>This notification was sent by your referrer %s.</i>\n\n"+
					"Restart your miner and compound your rewards to earn more.",
				weeklyFRENFormatted,
				senderName,
			)
			notifyRestartInactive(msg, ru.TelegramId)
		}

		sent++
	}

	u.LastReminder = time.Now()
	if err := db.Save(u).Error; err != nil {
		loge(err)
	}

	if sent == 0 {
		rr.Success = false
		rr.ErrorMessage = "No inactive referred users to remind."
		ctx.JSON(200, rr)
		return
	}

	rr.Success = true
	ctx.JSON(200, rr)
}

type RemindResponse struct {
	Success      bool   `json:"success"`
	ErrorMessage string `json:"errorMessage,omitempty"`
}
