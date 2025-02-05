package bot

import (
	"log"

	"github.com/go-macaron/binding"
	"github.com/go-macaron/cache"
	macaron "gopkg.in/macaron.v1"
)

func initMacaron() *macaron.Macaron {
	mac := macaron.Classic()

	mac.Use(macaron.Renderer())
	mac.Use(cache.Cacher())

	mac.Get("/data/:telegramid/:referral/:code/:name", viewData)
	mac.Get("/paid/:telegramid", viewPayment)
	mac.Get("/stats", viewStats)

	mac.Post("/save/:telegramid", binding.Bind(UserForm{}), viewSave)
	mac.Post("/compound/:telegramid", viewCompound)
	mac.Post("/withdraw/:telegramid", viewWithdraw)
	mac.Post("/restart/:telegramid", viewRestart)

	log.Println(conf.Port)

	go mac.Run(conf.Port)

	return mac
}
