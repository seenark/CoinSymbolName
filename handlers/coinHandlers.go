package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/seenark/coin-name/service"
)

type CoinHandler struct {
	CoinService service.ICoinService
}

func NewCoinHandler(cs service.ICoinService, app fiber.Router) {

	handler := CoinHandler{
		CoinService: cs,
	}

	app.Post("/", handler.Create)
	app.Post("/many", handler.CreateMany)
	app.Get("/", handler.GetAll)
	app.Get("/:symbol", handler.GetBySymbol)
	// app.Get("/coin/:id", handler.GetById)
	app.Put("/:id", handler.Update)
	app.Delete("/:id", handler.Delete)

	// return CoinHandler{
	// 	CoinRepository: cr,
	// }
}
