package repository

import "go.mongodb.org/mongo-driver/bson/primitive"

type CoinKLine struct {
	Id          primitive.ObjectID `json:"id" bson:"_id"`
	Symbol      string             `json:"symbol" bson:"symbol"`
	ClosePrices []float64          `json:"closePrices" bson:"closePrices"`
}

type ICoinKLineRepository interface {
	GetMultiple([]string) ([]CoinKLine, error)
	GetBySymbol(string) (*CoinKLine, error)
	Create(CoinKLine) error
	Update(string, CoinKLine) error
	Delete(string) error
}
