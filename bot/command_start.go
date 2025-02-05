package bot

import (
	"strconv"
	"strings"
	"time"

	"gopkg.in/telebot.v3"
)

func commandStart(c telebot.Context) error {
	var err error
	p := c.Message().Payload
	u := getUser(c.Sender().ID)
	if u.ID == 0 || p == "" {
		ab := getAppButton()
		b.Send(c.Sender(), lStart, ab)

		u, err = getUserOrCreate(c)
		if err != nil {
			loge(err)
		}

		if strings.HasPrefix(p, "b-") {
			n := uint(0)

			pids := strings.Split(p, "-")[1]
			pid, err := strconv.Atoi(pids)
			if err != nil {
				loge(err)
			}

			po := getPost(pid)

			if po.Channel.OwnerId != nil && po.Channel.OwnerId != &n {
				u.ReferrerID = po.Channel.OwnerId
				err := db.Save(u).Error
				if err != nil {
					loge(err)
				}
			}
		}
	} else {
		if p == "restart" {
			rb := getRestartButtons(c)
			if time.Since(u.MiningTime).Minutes() > 1410 {
				u.MiningTime = time.Now()
				u.LastNotification = time.Now()
				u.CycleCount++
				if err := db.Save(u).Error; err != nil {
					loge(err)
				}
				b.Send(c.Sender(), lCycleRestarted, rb)
			} else {
				b.Send(c.Sender(), lCycleRunning, rb)
			}
		} else if p == "claim" {
			commandClaim(c)
		} else if strings.HasPrefix(p, "b-") {
			commandBoost(c, p, false)
		}
	}

	return nil
}

func getAppButton() *telebot.ReplyMarkup {
	rm := &telebot.ReplyMarkup{}
	btn1 := rm.URL("Subscribe", "https://t.me/FrenlyNews")
	btn2 := rm.URL("Start Mining", "https://t.me/FrenlyRobot/miner")

	rm.Inline(
		rm.Row(btn1, btn2),
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

	// b.Handle(&btn1, commandCompound)

	return rm
}
