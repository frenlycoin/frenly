package bot

import (
	"fmt"
	"sort"

	"gopkg.in/macaron.v1"
)

func viewStats(ctx *macaron.Context) string {
	stats := "STAKED\n\n"
	var users []*User

	db.Preload("Referrer").Preload("Boosts").Order("tmu desc").Find(&users)

	for i, u := range users {
		stats += fmt.Sprintf("%d - %s: %.9f (%t)\n", i+1, u.Name, float64(u.TMU)/float64(Mul9), u.isActive())
	}

	rewards := make(map[string]float64)

	for _, u := range users {
		rf := float64(u.rewards(false)) / float64(Mul9)
		rewards[fmt.Sprintf("%s (%s)", u.Name, u.Code)] = rf
	}

	// Create slice of key-value pairs
	pairs := make([][2]interface{}, 0, len(rewards))
	for k, v := range rewards {
		pairs = append(pairs, [2]interface{}{k, v})
	}

	// Sort slice based on values
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i][1].(float64) > pairs[j][1].(float64)
	})

	// Extract sorted keys
	keys := make([]string, len(pairs))
	for i, p := range pairs {
		keys[i] = p[0].(string)
	}

	stats += "\nREWARD\n\n"

	// Print sorted map
	for _, k := range keys {
		stats += fmt.Sprintf("%s: %.9f\n", k, rewards[k])
	}

	return stats
}
