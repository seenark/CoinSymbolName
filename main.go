package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/seenark/coin-name/config"
	"github.com/seenark/coin-name/handlers"
	"github.com/seenark/coin-name/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	initTimeZone()
	cf := config.GetConfig()

	ctx := context.TODO()

	mongoClient := connectMongo(cf.Mongo.Username, cf.Mongo.Password)
	db := mongoClient.Database(cf.Mongo.DbName)
	coinCollection := db.Collection(cf.Mongo.CoinNameCollection)
	indexName, err := coinCollection.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bson.D{{Key: "symbol", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	)
	if err != nil {
		panic(err)
	}
	fmt.Println("index name:", indexName)

	coinRepository := repository.NewCoinRepository(coinCollection, &ctx)
	handler := handlers.NewCoinHandler(coinRepository)

	app := fiber.New()
	app.Post("/coin/", handler.Create)
	app.Post("/coin/many", handler.CreateMany)
	app.Get("/coin", handler.GetAll)
	app.Get("/symbol/:symbol", handler.GetBySymbol)
	app.Get("/coin/:id", handler.GetById)
	app.Put("/coin/:id", handler.Update)
	app.Delete("/coin/:id", handler.Delete)
	fmt.Printf("cf: %v\n", cf)
	app.Listen(fmt.Sprintf("localhost:%d", cf.Port))
}

func connectMongo(username string, password string) *mongo.Client {
	clientOptions := options.Client().
		ApplyURI(fmt.Sprintf("mongodb+srv://%s:%s@hdgcluster.xmgsx.mongodb.net/myFirstDatabase?retryWrites=true&w=majority", username, password))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func initTimeZone() {
	ict, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		panic(err)
	}

	time.Local = ict
}
