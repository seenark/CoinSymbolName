package repository

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CoinRepositoryDb struct {
	collection *mongo.Collection
	ctx        *context.Context
}

func NewCoinRepository(cl *mongo.Collection, ctx *context.Context) CoinRepository {
	return &CoinRepositoryDb{
		collection: cl,
		ctx:        ctx,
	}
}

func (cr *CoinRepositoryDb) GetAll(symbols []string, name string) ([]Coin, error) {
	fmt.Println("Get All From Repo")
	filters := bson.M{}
	newSymbols := []string{}
	for _, v := range symbols {
		if v != "" {
			newSymbols = append(newSymbols, v)
		}
	}
	if len(newSymbols) > 0 {
		symbolArr := []string{}
		symbolArr = append(symbolArr, symbols...)
		filters["symbol"] = bson.M{"$in": symbolArr}
	}
	if name != "" {
		filters["name"] = name
	}
	options := options.Find()
	options.SetSort(bson.D{primitive.E{Key: "symbol", Value: 1}})
	all := []Coin{}
	cur, err := cr.collection.Find(*cr.ctx, filters, options)

	if err != nil {
		return nil, err
	}

	for cur.Next(*cr.ctx) {
		coin := Coin{}
		err = cur.Decode(&coin)
		if err != nil {
			continue
		}
		all = append(all, coin)
	}

	return all, nil
}
func (cr *CoinRepositoryDb) GetById(id string) (*Coin, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	coin := Coin{}
	err = cr.collection.FindOne(*cr.ctx, bson.D{primitive.E{Key: "_id", Value: _id}}).Decode(&coin)
	if err != nil {
		return nil, err
	}
	return &coin, nil
}

func (cr *CoinRepositoryDb) GetBySymbol(symbol string) (*Coin, error) {
	coin := new(Coin)
	err := cr.collection.FindOne(*cr.ctx, bson.D{primitive.E{Key: "symbol", Value: symbol}}).Decode(&coin)
	if err != nil {
		return nil, err
	}
	return coin, nil
}

func (cr *CoinRepositoryDb) Create(c Coin) (string, error) {

	_id := primitive.NewObjectID()
	coin := Coin{
		Id:     _id,
		Symbol: c.Symbol,
		Name:   c.Name,
	}

	inserted, err := cr.collection.InsertOne(*cr.ctx, coin)
	if err != nil {
		return "", err
	}
	idObj, ok := inserted.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", fmt.Errorf("cannot translate inserted_id to obj_id")
	}
	return idObj.Hex(), nil
}

func (cr *CoinRepositoryDb) CreateMany(cs []Coin) ([]string, error) {
	items := []interface{}{}
	for _, c := range cs {
		c.Id = primitive.NewObjectID()
		items = append(items, bson.D{
			primitive.E{Key: "symbol", Value: c.Symbol},
			primitive.E{Key: "name", Value: c.Name},
			primitive.E{Key: "_id", Value: c.Id},
		})
	}
	opts := options.InsertMany().SetOrdered(false)
	result, err := cr.collection.InsertMany(*cr.ctx, items, opts)
	if err != nil {
		return nil, err
	}
	ids := []string{}
	for _, id := range result.InsertedIDs {
		fmt.Printf("id: %v\n", id)
		newId, ok := id.(primitive.ObjectID)
		if !ok {
			continue
		}
		ids = append(ids, newId.Hex())
	}
	return ids, nil
}
func (cr *CoinRepositoryDb) Update(id string, c Coin) (string, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return "", err
	}
	update := bson.D{primitive.E{Key: "$set", Value: c}}
	_, err = cr.collection.UpdateOne(*cr.ctx, bson.D{primitive.E{Key: "_id", Value: _id}}, update)
	// _, err := cr.collection.UpdateByID(*cr.ctx, id, update)
	if err != nil {
		return "", err
	}
	return c.Id.Hex(), nil
}
func (cr *CoinRepositoryDb) Delete(id string) error {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = cr.collection.DeleteOne(*cr.ctx, bson.D{primitive.E{Key: "_id", Value: _id}})
	if err != nil {
		return err
	}
	return nil
}
