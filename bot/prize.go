package bot

import (
	"log"
	"time"
)

type PrizeManager struct {
}

func (pm *PrizeManager) start() {
	for {
		log.Println("Prize Tick")
		time.Sleep(time.Second * PrizeTick)
	}
}

func initPrizeManager() *PrizeManager {
	pm := &PrizeManager{}
	go pm.start()
	return pm
}
