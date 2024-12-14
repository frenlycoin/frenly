package bot

import (
	"time"

	"gopkg.in/telebot.v3"
)

func commandStart(c telebot.Context) error {
	p := c.Message().Payload
	u, err := getUserOrCreate(c)
	if err != nil {
		loge(err)
	}

	if p == "restart" {
		rb := getRestartButtons(c)
		if time.Since(u.MiningTime).Minutes() > 1410 {
			u.MiningTime = time.Now()
			u.CycleCount++
			if err := db.Save(u).Error; err != nil {
				loge(err)
			}
			b.Send(c.Sender(), lCycleRestarted, rb)
		} else {
			b.Send(c.Sender(), lCycleRunning, rb)
		}
	} else {
		ab := getAppButton()
		b.Send(c.Sender(), lStart, ab)
	}

	return nil
}

func getAppButton() *telebot.ReplyMarkup {
	rm := &telebot.ReplyMarkup{}
	btn := rm.URL("⚪️ Launch Frenly App", "https://t.me/FrenlyRobot/miner")

	rm.Inline(
		rm.Row(btn),
	)

	return rm
}

func getRestartButtons(c telebot.Context) *telebot.ReplyMarkup {
	rm := &telebot.ReplyMarkup{}

	btn1 := rm.Data("Compound", "compound")
	btn2 := rm.URL("Launch App", "https://t.me/FrenlyRobot/miner")

	rm.Inline(
		rm.Row(btn1, btn2),
	)

	b.Handle(&btn1, commandCompound)

	return rm
}
