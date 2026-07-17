package bot

import (
	"fmt"
	"log"

	"gopkg.in/telebot.v3"
)

func commandStats(c telebot.Context) error {
	priceKv := &KeyValue{Key: "dexLastPrice"}
	price := float64(0)
	if err := db.Where("key = ?", priceKv.Key).FirstOrCreate(priceKv).Error; err == nil {
		price = float64(priceKv.ValueInt) / float64(Mul9)
	} else {
		loge(err)
	}

	msg := fmt.Sprintf(lStats, cch.StatsCache.Miners, cch.StatsCache.ActiveMiners, cch.StatsCache.TMU, cch.StatsCache.RewardTMU, price)

	log.Println(c.Chat().ID)

	_, err := b.Send(c.Chat(), msg, telebot.NoPreview)
	return err
}
