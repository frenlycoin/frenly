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

func getRestartButtonBoost(boostLink string) *telebot.ReplyMarkup {
	kv := &KeyValue{Key: "restartPostId"}
	db.FirstOrCreate(kv, kv)
	link := fmt.Sprintf("https://t.me/FrenlyNews/%d", kv.ValueInt)

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
	btn1 := rm.Data("Boost Frenly Miner 🚀", boostId)
	btn2 := rm.URL("New User", "t.me/FrenlyRobot")

	rm.Inline(
		rm.Row(btn1, btn2),
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
		rb = getRestartButtonChannel()
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
