package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/seenark/coin-name/service"
)

func (h CoinHandler) Create(c *fiber.Ctx) error {
	coin := new(service.CoinResponse)
	err := c.BodyParser(&coin)
	if err != nil {
		return err
	}

	id, err := h.CoinService.Create(*coin)
	if err != nil {
		return err
	}
	return c.Status(http.StatusCreated).JSON(fiber.Map{"Created": id})
}

func (h CoinHandler) CreateMany(c *fiber.Ctx) error {
	coins := []service.CoinResponse{}
	err := c.BodyParser(&coins)
	if err != nil {
		return err
	}
	fmt.Printf("coins: %v\n", coins)
	ids, err := h.CoinService.CreateMany(coins)
	if err != nil {
		return err
	}
	return c.Status(http.StatusCreated).JSON(ids)
}

func (h CoinHandler) GetAll(c *fiber.Ctx) error {

	symbol := c.Query("symbol")
	name := c.Query("name")

	split := strings.Split(symbol, ",")

	for index, v := range split {
		if v == "" {
			continue
		}
		split[index] = strings.Trim(v, " ")

	}

	all, err := h.CoinService.GetAll(split, name)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Somthing went wrong"})
	}
	return c.Status(http.StatusOK).JSON(all)
}

// func (h CoinHandler) GetById(c *fiber.Ctx) error {
// 	id := c.Params("id")
// 	if id == "" {
// 		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "please specified ID"})
// 	}
// 	coin, err := h.CoinService.GetById(id)
// 	if err != nil {
// 		return err
// 	}
// 	return c.Status(http.StatusOK).JSON(coin)
// }

func (h CoinHandler) GetBySymbol(c *fiber.Ctx) error {
	symbol := c.Params("symbol")
	if symbol == "" {
		return c.SendStatus(http.StatusBadRequest)
	}
	coin, err := h.CoinService.GetNameBySymbol(symbol)
	if err != nil {
		return c.SendStatus(http.StatusNotFound)
	}
	return c.JSON(coin)
}

func (h CoinHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.SendStatus(http.StatusBadRequest)
	}
	coin := new(service.CoinResponse)
	err := c.BodyParser(&coin)
	if err != nil {
		return err
	}
	id, err = h.CoinService.Update(id, *coin)
	if err != nil {
		return err
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{"updated": id})

}

func (h CoinHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.SendStatus(http.StatusBadRequest)
	}

	err := h.CoinService.Delete(id)
	if err != nil {
		return err
	}
	return c.SendStatus(http.StatusOK)
}
