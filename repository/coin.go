package repository

import "go.mongodb.org/mongo-driver/bson/primitive"

type Coin struct {
	Id     primitive.ObjectID `json:"id" bson:"_id"`
	Symbol string             `json:"symbol" bson:"symbol"`
	Name   string             `json:"name" bson:"name"`
}

type CoinRepository interface {
	GetAll([]string, string) ([]Coin, error)
	GetById(string) (*Coin, error)
	GetBySymbol(string) (*Coin, error)
	Create(Coin) (string, error)
	CreateMany([]Coin) ([]string, error)
	Update(string, Coin) (string, error)
	Delete(string) error
}
