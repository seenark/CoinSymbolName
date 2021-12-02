package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/seenark/coin-name/helpers"
	"github.com/seenark/coin-name/repository"
	"github.com/seenark/coin-name/routine"
)

type KLineHandler struct {
	Repo repository.ICoinKLineRepository
}

func NewKlineHandler(app fiber.Router, klineRepo repository.ICoinKLineRepository) {
	handler := KLineHandler{
		Repo: klineRepo,
	}

	app.Get("/", handler.getMultiple)
	app.Get("/:symbol", handler.getKline)
	app.Post("/", handler.create)
	app.Put("/:symbol", handler.update)
	app.Delete("/:symbol", handler.delete)
}

func (kh KLineHandler) create(c *fiber.Ctx) error {
	kline := repository.CoinKLine{}
	err := c.BodyParser(&kline)
	if err != nil {
		return err
	}
	err = kh.Repo.Create(kline)
	if err != nil {
		return err
	}
	return c.SendStatus(http.StatusCreated)

}

func (kh KLineHandler) getMultiple(ctx *fiber.Ctx) error {
	symbols := ctx.Query("symbols")
	split := strings.Split(symbols, ",")

	for index, v := range split {
		if v == "" {
			continue
		}
		split[index] = strings.Trim(v, " ")
	}

	all, err := kh.Repo.GetMultiple(split)
	if err != nil {
		ctx.SendStatus(http.StatusNotFound)
	}
	notFoundList := []string{}
	symbolMap := map[string]bool{}
	for _, s := range all {
		symbolMap[s.Symbol] = true
	}
	fmt.Printf("symbolMap: %v\n", symbolMap)
	for _, s := range split {
		if _, ok := symbolMap[s]; !ok {
			notFoundList = append(notFoundList, s)
		}
	}

	fmt.Printf("notFoundList: %v\n", notFoundList)
	fmt.Printf("notFoundList: %v\n", len(notFoundList))
	for _, sb := range notFoundList {
		if sb == "" {
			continue
		}
		ck, err := routine.StoreHourKLineForSymbol(sb, kh.Repo)
		if err != nil {
			fmt.Println("some error", err.Error())
			continue
		}
		all = append(all, *ck)
	}
	helpers.PrintMemUsage()
	return ctx.Status(http.StatusOK).JSON(all)
}

func (kh KLineHandler) getKline(c *fiber.Ctx) error {
	symbol := c.Params("symbol")
	kl, err := kh.Repo.GetBySymbol(symbol)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return c.SendStatus(http.StatusNotFound)
		}

		return err
	}
	return c.JSON(kl)
}

func (kh KLineHandler) update(c *fiber.Ctx) error {
	symbol := c.Params("symbol")
	kline := repository.CoinKLine{}
	c.BodyParser(&kline)
	err := kh.Repo.Update(symbol, kline)
	if err != nil {
		return err
	}
	return c.SendStatus(http.StatusOK)
}

func (kh KLineHandler) delete(c *fiber.Ctx) error {
	symbol := c.Params("symbol")
	err := kh.Repo.Delete(symbol)
	if err != nil {
		return err
	}
	return c.SendStatus(http.StatusOK)
}
