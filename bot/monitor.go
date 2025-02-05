package bot

import (
	"log"
	"time"
)

type Monitor struct {
	Miners []*User
}

func (m *Monitor) loadMiners() {
	m.Miners = nil
	db.Preload("Referrer").Preload("Boosts").Find(&m.Miners)
}

func (m *Monitor) sendNotifications() {
	// counter := 1
	for _, miner := range m.Miners {
		if m.isSending(miner) {
			notifyEnd(miner)
			log.Printf("Notification: %s", miner.Name)
		}

		// if m.isSendingWeekly(miner, 10080) {
		// 	sendNotificationWeekly(miner)
		// 	log.Printf("Notification Weekly: %s %d", miner.Address, counter)
		// 	counter++
		// }
	}
}

func (m *Monitor) isSending(miner *User) bool {
	if miner.ID != 0 &&
		time.Since(miner.MiningTime).Minutes() > 1410 &&
		time.Since(miner.MiningTime).Minutes() < 1440 &&
		miner.LastNotification.Day() != time.Now().Day() &&
		miner.TelegramId != 0 {

		miner.LastNotification = time.Now()
		err := db.Save(miner).Error
		if err != nil {
			loge(err)
		}

		return true
	}

	return false
}

// func (m *Monitor) isSendingWeekly(miner *Miner, limit int64) bool {
// 	if miner.ID != 0 &&
// 		(int64(m.Height)-miner.MiningHeight) >= limit &&
// 		(miner.MiningTime.Hour() == time.Now().Hour() ||
// 			miner.MiningTime.IsZero()) &&
// 		time.Since(miner.LastNotificationWeekly) > time.Hour*168 &&
// 		miner.TelegramId != 0 {

// 		miner.LastNotificationWeekly = time.Now()
// 		err := db.Save(miner).Error
// 		if err != nil {
// 			log.Println(err)
// 			logTelegram(err.Error())
// 		}

// 		return true
// 	}
// 	return false
// }

func (m *Monitor) minerExists(telId int64) bool {
	for _, mnr := range m.Miners {
		if int64(mnr.TelegramId) == telId {
			return true
		}
	}

	return false
}

func (m *Monitor) isTriggeringChannelPost() bool {
	kv := &KeyValue{Key: "lastPostDay"}
	db.FirstOrCreate(kv, kv)

	if time.Now().Hour() == 16 && time.Now().Day() != int(kv.ValueInt) {
		kv.ValueInt = int64(time.Now().Day())
		db.Save(kv)
		return true
	}

	return false
}

func (m *Monitor) start() {
	for {
		m.loadMiners()

		m.sendNotifications()

		// log.Println("Monitor users loaded.")

		if m.isTriggeringChannelPost() {
			notifyRestart()
		}

		time.Sleep(time.Second * MonitorTick)
	}
}

func initMonitor() *Monitor {
	m := &Monitor{}
	go m.start()
	return m
}
