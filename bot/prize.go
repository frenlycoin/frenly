package bot

import (
	"fmt"
	"log"
	"time"

	"gopkg.in/telebot.v3"
)

type PrizeManager struct {
	InactiveMiners []*User
	ActiveMiners   []*User
	WinMsg         *telebot.Message
}

func (pm *PrizeManager) loadMiners() {
	pm.InactiveMiners = nil
	pm.ActiveMiners = nil

	for _, u := range mon.Miners {
		if u.isActive() {
			pm.ActiveMiners = append(pm.ActiveMiners, u)
		} else {
			pm.InactiveMiners = append(pm.InactiveMiners, u)
		}
	}
}

func (pm *PrizeManager) isTriggering() bool {
	kv := &KeyValue{Key: "lastPrizeDay"}
	db.FirstOrCreate(kv, kv)

	if time.Now().Hour() == 9 && time.Now().Day() != int(kv.ValueInt) {
		kv.ValueInt = int64(time.Now().Day())
		db.Save(kv)
		return true
	}

	return false
}

func (pm *PrizeManager) executeLosers() {
	l := make(map[int]bool)
	c := len(pm.InactiveMiners)

	for range c / 10 {
		ui := generateRandNum(c - 1)
		l[ui] = true
	}

	for ui := range l {
		lu := pm.InactiveMiners[ui]
		msg := fmt.Sprintf(lNotWon, fmt.Sprintf("t.me/FrenlyNews/%d", pm.WinMsg.ID))
		log.Println(msg)
		notify(msg, lu.TelegramId)
		log.Printf("Loser: %s", lu.Name)
	}
}

func (pm *PrizeManager) executeWinner() *telebot.Message {
	wn := generateRandNum(len(pm.ActiveMiners))
	w := pm.ActiveMiners[wn]

	kv := &KeyValue{Key: "prizeWinner"}
	db.FirstOrCreate(kv, kv)

	kv.ValueInt = int64(w.ID)
	db.Save(kv)

	wm := notifyPrize(w)

	log.Printf("Winner: %s", w.Name)

	return wm
}

func (pm *PrizeManager) start() {
	for {
		if pm.isTriggering() {
			pm.loadMiners()

			pm.WinMsg = pm.executeWinner()

			pm.executeLosers()
		}

		time.Sleep(time.Second * PrizeTick)
	}
}

func initPrizeManager() *PrizeManager {
	pm := &PrizeManager{}
	go pm.start()
	return pm
}
