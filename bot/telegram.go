package bot

import (
	"fmt"
	"log"
	"strings"
	"time"

	"gopkg.in/telebot.v3"
)

func initTelegram(key string) *telebot.Bot {
	b, err := telebot.NewBot(telebot.Settings{
		Token:     key,
		Poller:    &telebot.LongPoller{Timeout: 30 * time.Second},
		Verbose:   false,
		ParseMode: "html",
	})

	if err != nil {
		log.Fatal(err)
	}

	b.Handle("/start", commandStart)
	b.Handle("/stats", commandStats)
	b.Handle("/ranks", commandRanks)
	b.Handle("/check", commandCheck)

	b.Handle(telebot.OnText, commandFilter)
	b.Handle(telebot.OnChannelPost, commandChannelPost)
	b.Handle(telebot.OnCallback, commandCallback)
	b.Handle(telebot.OnUserJoined, commandJoin)

	return b
}

func notify(msg string, tgid int64) {
	rec := &telebot.Chat{
		ID: tgid,
	}
	b.Send(rec, msg, telebot.NoPreview)
}

func notifyWithButton(msg string, tgid int64, btn *telebot.ReplyMarkup) {
	rec := &telebot.Chat{
		ID: tgid,
	}
	b.Send(rec, msg, btn, telebot.NoPreview)
}

func notifyRestart() {
	// mb := getRestartButton()
	rec := &telebot.Chat{
		ID: News,
	}

	msg := lRestartMining

	m, err := b.Send(rec, msg, telebot.Silent)
	if err != nil {
		loge(err)
	}

	kv := &KeyValue{Key: "restartPostId"}
	db.FirstOrCreate(kv, kv)

	kv.ValueInt = int64(m.ID)
	db.Save(kv)
}

func notifystart(msg string, tgid int64) {
	sb := getStartButton()
	rec := &telebot.Chat{
		ID: tgid,
	}
	_, err := b.Send(rec, msg, sb, telebot.NoPreview)
	if err != nil {
		loge(err)
	}
}

func getRestartButton() *telebot.ReplyMarkup {
	rm := &telebot.ReplyMarkup{}
	btn := rm.URL("⚪️ Restart Mining", "https://t.me/FrenlyRobot?start=restart")

	rm.Inline(
		rm.Row(btn),
	)

	return rm
}

func getRestartButtonChannel() *telebot.ReplyMarkup {
	kv := &KeyValue{Key: "restartPostId"}
	db.FirstOrCreate(kv, kv)
	link := fmt.Sprintf("https://t.me/FrenlyNews/%d", kv.ValueInt)

	rm := &telebot.ReplyMarkup{}
	btn := rm.URL("⚪️ Restart Mining", link)

	rm.Inline(
		rm.Row(btn),
	)

	return rm
}

func getRestartButtonBoost(boostLink string) *telebot.ReplyMarkup {
	link := "https://t.me/FrenlyRobot?start=restart"

	rm := &telebot.ReplyMarkup{}
	btn1 := rm.URL("Boost Miner", boostLink)
	btn2 := rm.URL("Restart Mining", link)
	// btn2 := rm.URL("Test", "https://bot.frenlycoin.com/data/7422140567/unknown/unknown/Frenly")

	rm.Inline(
		rm.Row(btn1, btn2),
	)

	return rm
}

func getStartButton() *telebot.ReplyMarkup {
	rm := &telebot.ReplyMarkup{}
	btn := rm.URL("Start Mining Now 🚀", "t.me/FrenlyRobot")

	rm.Inline(
		rm.Row(btn),
	)

	return rm
}

func getFrenlyButtons(boostId string) *telebot.ReplyMarkup {
	rm := &telebot.ReplyMarkup{}
	// app := &telebot.WebApp{
	// 	URL: fmt.Sprintf("https://t.me/DevFrenlyRobot/miner?startapp=%s", boostId),
	// }
	// btn1 := rm.WebApp("Boost Frenly Miner 🚀", app)
	btn1 := rm.URL("Boost Frenly Miner 🚀", fmt.Sprintf("https://t.me/DevFrenlyRobot/miner?startapp=%s", boostId))
	// btn2 := rm.URL("New User", "t.me/FrenlyRobot")

	rm.Inline(
		rm.Row(btn1),
	)

	return rm
}

func getClaimButton() *telebot.ReplyMarkup {
	rm := &telebot.ReplyMarkup{}
	btn := rm.URL("Claim The Reward 🚀", "t.me/FrenlyRobot?start=claim")

	rm.Inline(
		rm.Row(btn),
	)

	return rm
}

func getGroupButton(link string) *telebot.ReplyMarkup {
	rm := &telebot.ReplyMarkup{}
	btn := rm.URL("Claim In Group", link)

	rm.Inline(
		rm.Row(btn),
	)

	return rm
}

func getButtonLink(name string, link string) *telebot.ReplyMarkup {
	rm := &telebot.ReplyMarkup{}
	btn := rm.URL(name, link)

	rm.Inline(
		rm.Row(btn),
	)

	return rm
}

func getButtonsBoost(name string, link string) *telebot.ReplyMarkup {
	rm := &telebot.ReplyMarkup{}
	btn1 := rm.URL(name, link)
	btn2 := rm.Data("Missing Boost", link)

	rm.Inline(
		rm.Row(btn1, btn2),
	)

	return rm
}

func notifyEnd(u *User) {
	var rb *telebot.ReplyMarkup
	msg := lCycleFinished
	unb := u.getUnboosted()
	if len(unb) > 0 {
		rb = getRestartButtonBoost(unb[0].Link)
		msg += fmt.Sprintf("\n\n<b><u>Your miner's health is at %d%%!</u></b>\n\n<b><u>Boost your miner by clicking the button bellow and then boost button under each post that bot leads you to! This needs to be done to collect full reward.</u></b>", u.health())
	} else {
		rb = getRestartButton()
	}

	rec := &telebot.Chat{
		ID: u.TelegramId,
	}

	_, err := b.Send(rec, msg, rb, telebot.NoPreview)
	if err != nil {
		if strings.Contains(err.Error(), "blocked") {
			u.BotBlocked = true
			if err := db.Save(u).Error; err != nil {
				loge(err)
			}
		} else {
			loge(err)
		}
	}
}

func notifyPrize(u *User) *telebot.Message {
	cb := getClaimButton()
	msg := fmt.Sprintf(lWonPrize, u.Name)

	rec := &telebot.Chat{
		ID: News,
	}

	recGroup := &telebot.Chat{
		ID: getGroup(),
	}

	mc, err := b.Send(rec, msg, telebot.NoPreview)
	if err != nil {
		loge(err)
	}

	time.Sleep(time.Second * 5)

	claimMsg := fmt.Sprintf(lClaimPrize, u.Name)

	mg, err := b.Send(recGroup, claimMsg, cb, telebot.NoPreview)
	if err != nil {
		loge(err)
	}

	gb := getGroupButton(fmt.Sprintf("https://t.me/FrenlyCoin/%d", mg.ID))

	_, err = b.Edit(mc, gb)
	if err != nil {
		loge(err)
	}

	return mc
}

func postExists(postId string, chId int) bool {
	// c := getChannel(chId)

	sm := telebot.StoredMessage{
		MessageID: postId,
		ChatID:    int64(chId),
	}

	ch, err := b.ChatByID(int64(chId))
	if err != nil {
		log.Println(err)
		return false
	}

	err = b.React(ch, sm, telebot.ReactionOptions{})
	if err != nil {
		if strings.Contains(err.Error(), "MESSAGE_ID_INVALID") {
			return false
		} else {
			return true
		}
	}

	return true
}

func logTelegramSilent(message string) {
	message = getCallerInfo() + message
	rec := &telebot.Chat{
		ID: int64(7422140567),
	}
	b.Send(rec, message, telebot.Silent)
}

func notifyCashout(u *User, amount int64, tgid int64) {
	username := u.Name
	if len(u.Code) > 0 && u.Code != u.AddressWithdraw {
		username = fmt.Sprintf("@%s", u.Code)
	} else if u.Code == u.AddressWithdraw {
		username = u.Name
	}

	createdAt := "unknown"
	if !u.CreatedAt.IsZero() {
		createdAt = u.CreatedAt.Format("Jan 2, 2006")
	}

	frenAmount := float64(u.TMU) / float64(Mul9)
	compounds := u.CompoundCount
	cashoutAmount := float64(amount) / float64(Mul9)
	depositAddress := u.AddressDeposit
	if len(depositAddress) == 0 {
		depositAddress = u.AddressWithdraw
	}
	if len(depositAddress) == 0 {
		depositAddress = "unknown"
	}

	msg := fmt.Sprintf(lCashOut, username, formatNumber(cashoutAmount), createdAt, u.CycleCountTotal, compounds, formatNumber(frenAmount), depositAddress, depositAddress)

	rm := &telebot.ReplyMarkup{}
	payURL := fmt.Sprintf("https://app.tonkeeper.com/transfer/%s?amount=%d", u.AddressWithdraw, amount)
	payBtn := rm.URL("Pay", payURL)
	doneBtn := rm.Data("Done", "done")
	cancelBtn := rm.Data("Cancel", "cancel")

	rm.Inline(
		rm.Row(payBtn, doneBtn, cancelBtn),
	)

	rec := &telebot.Chat{ID: tgid}
	_, err := b.Send(rec, msg, rm, telebot.NoPreview)
	if err != nil {
		loge(err)
	}
}

func notifyRestartInactive(msg string, tgid int64) {
	rb := getRestartButton()
	rec := &telebot.Chat{
		ID: tgid,
	}
	_, err := b.Send(rec, msg, rb, telebot.NoPreview)
	if err != nil {
		loge(err)
	}
}

func notifyInactive() {
	var users []*User
	if err := db.Find(&users).Error; err != nil {
		loge(err)
		return
	}

	for _, u := range users {
		time.Sleep(time.Second)
		if u.isActive() || u.BotBlocked {
			continue
		}
		notifyRestartInactive(lPayoutsEnabled, u.TelegramId)
		notify(fmt.Sprintf("User %s is inactive, sent restart notification.", u.Name), Frenly)
	}
}
