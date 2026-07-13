package bot

import (
	"log"
	"net/http"

	"github.com/go-macaron/binding"
	"github.com/go-macaron/cache"
	"github.com/go-macaron/csrf"
	"github.com/go-macaron/session"
	macaron "gopkg.in/macaron.v1"
)

func initMacaron() *macaron.Macaron {
	mac := macaron.Classic()

	mac.Use(macaron.Renderer())
	mac.Use(cache.Cacher())
	mac.Use(CustomHeaderMiddleware())
	mac.Use(session.Sessioner())
	mac.Use(csrf.Csrfer())

	mac.Options("/*", func(ctx *macaron.Context) {
		ctx.Resp.Header().Set("Access-Control-Allow-Origin", "*") // or your frontend origin
		ctx.Resp.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		ctx.Resp.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, X-Custom-Header, ngrok-skip-browser-warning")
		ctx.Resp.WriteHeader(http.StatusOK)
	})

	mac.Get("/data/:telegramid/:referral/:code/:name", viewData)
	mac.Get("/paid/:telegramid", viewPayment)
	mac.Get("/stats", viewStats)

	mac.Post("/save/:telegramid", binding.Bind(UserForm{}), viewSave)
	mac.Post("/compound/:telegramid", viewCompound)
	mac.Post("/withdraw/:telegramid", viewWithdraw)
	mac.Post("/restart/:telegramid", viewRestart)
	mac.Post("/boost/:telegramid/:boostid", viewBoost)

	log.Println(conf.Port)

	go mac.Run(conf.Port)

	return mac
}

func CustomHeaderMiddleware() macaron.Handler {
	return func(ctx *macaron.Context) {
		// Add header **after** the next handler has run
		ctx.Next()

		// Add your custom header(s) here
		ctx.Resp.Header().Set("ngrok-skip-browser-warning", "true")
		// You can add more:
		// ctx.Resp.Header().Set("X-Powered-By", "Go-Macaron")
		// ctx.Resp.Header().Add("Cache-Control", "no-cache")
	}
}
