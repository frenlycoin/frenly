package bot

import (
	"time"
)

type Cache struct {
	StatsCache *StatsCache
}

func (c *Cache) loadStatsCache() {
	tmu := float64(0)
	reward := uint64(0)
	var users []*User
	db.Preload("Referrer").Preload("Boosts").Find(&users)
	count := len(users)
	countActive := 0

	for _, u := range users {
		tmu += (float64(u.TMU) / float64(Mul9))
		reward += u.rewards(false)
		if u.isActive() {
			countActive++
		}
	}

	c.StatsCache.Miners = count
	c.StatsCache.ActiveMiners = countActive
	c.StatsCache.TMU = tmu
	c.StatsCache.RewardTMU = float64(reward) / float64(Mul9)
}

func (c *Cache) start() {
	for {
		c.loadStatsCache()

		// log.Println("Cache loaded.")

		time.Sleep(time.Second * CacheTick)
	}
}

func initCache() *Cache {
	c := &Cache{}
	c.StatsCache = &StatsCache{}
	go c.start()

	return c
}

type StatsCache struct {
	Miners       int
	ActiveMiners int
	TMU          float64
	RewardTMU    float64
}
