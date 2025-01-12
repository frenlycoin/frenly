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

	b.Handle(telebot.OnChannelPost, commandChannelPost)
	b.Handle(telebot.OnUserJoined, commandJoin)

	return b
}

func notify(msg string, tgid int64) {
	rec := &telebot.Chat{
		ID: tgid,
	}
	b.Send(rec, msg, telebot.NoPreview)
}

func notifyRestart() {
	mb := getRestartButton()
	rec := &telebot.Chat{
		ID: News,
	}

	msg := lRestartMining

	m, err := b.Send(rec, msg, mb, telebot.Silent)
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

func getStartButton() *telebot.ReplyMarkup {
	rm := &telebot.ReplyMarkup{}
	btn := rm.URL("Start Mining Now 🚀", "t.me/FrenlyRobot")

	rm.Inline(
		rm.Row(btn),
	)

	return rm
}

func getFrenlyButton() *telebot.ReplyMarkup {
	rm := &telebot.ReplyMarkup{}
	btn := rm.URL("Boost Frenly Miner 🚀", "https://t.me/FrenlyRobot?start=boost")

	rm.Inline(
		rm.Row(btn),
	)

	return rm
}

func getClaimButton() *telebot.ReplyMarkup {
	rm := &telebot.ReplyMarkup{}
	btn := rm.Data("Claim The Reward 🚀", "claim")

	rm.Inline(
		rm.Row(btn),
	)

	b.Handle(&btn, commandClaim)

	return rm
}

func getGroupButton(link string) *telebot.ReplyMarkup {
	rm := &telebot.ReplyMarkup{}
	btn := rm.URL("Claim In Group", link)

	rm.Inline(
		rm.Row(btn),
	)

	b.Handle(&btn, commandClaim)

	return rm
}

func notifyEnd(tgid int64) {
	rb := getRestartButtonChannel()

	rec := &telebot.Chat{
		ID: tgid,
	}

	_, err := b.Send(rec, lCycleFinished, rb, telebot.NoPreview)
	if err != nil {
		if strings.Contains(err.Error(), "blocked") {
			u := getUser(tgid)
			u.BotBlocked = true
			if err := db.Save(u).Error; err != nil {
				loge(err)
			}
		} else {
			loge(err)
		}
	}
}

func notifyPrize(u *User) {
	cb := getClaimButton()
	msg := fmt.Sprintf(lWonPrize, u.Name)

	rec := &telebot.Chat{
		ID: News,
	}

	recGroup := &telebot.Chat{
		ID: Group,
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
}
