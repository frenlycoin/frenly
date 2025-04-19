package bot

import (
	"log"
	"strconv"
	"strings"

	"gopkg.in/macaron.v1"
)

func viewBoost(ctx *macaron.Context) {
	br := &BoostResponse{Success: true}
	tgid := getTgId(ctx)
	p := ctx.Params("boostid")
	u := getUser(tgid)

	if strings.HasPrefix(p, "b-") {
		pids := strings.Split(p, "-")[1]
		pid, err := strconv.Atoi(pids)
		if err != nil {
			loge(err)
		}

		po := getPost(pid)
		boosted := false

		if po.ID != 0 {
			if u.MiningTime.Before(po.CreatedAt) {
				for _, p := range u.Boosts {
					if p.ID == po.ID {
						boosted = true
					}
				}

				if !boosted {
					err := db.Model(u).Association("Boosts").Append(po)
					if err != nil {
						loge(err)
					}

					u.MiningTime = po.CreatedAt
					if err := db.Save(u).Error; err != nil {
						loge(err)
					}
					log.Println("Saved mining time.")
				}
			}
		}
	}

	br.Health = u.health()

	ctx.Header().Add("Access-Control-Allow-Origin", "*")

	ctx.JSON(200, br)
}

type BoostResponse struct {
	Success bool  `json:"success"`
	Health  int64 `json:"health"`
}
