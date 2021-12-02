package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/seenark/coin-name/config"
	"github.com/seenark/coin-name/handlers"
	"github.com/seenark/coin-name/routine"

	"github.com/seenark/coin-name/repository"
	"github.com/seenark/coin-name/service"
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
	makeSymbolAsIndexes(coinCollection)

	coinRepository := repository.NewCoinRepository(coinCollection, &ctx)
	coinService := service.NewCoinService(coinRepository)

	klineDb := mongoClient.Database(cf.Mongo.KlineDbName)
	hourCollection := klineDb.Collection(cf.Mongo.HourKlineCollection)
	makeSymbolAsIndexes(hourCollection)

	klineRepository := repository.NewKLineRepository(hourCollection, ctx)

	app := fiber.New()
	// Default config
	app.Use(cors.New())
	coinGroup := app.Group("/coin")
	klineGroup := app.Group("/kline")
	handlers.NewCoinHandler(coinService, coinGroup)
	handlers.NewKlineHandler(klineGroup, klineRepository)

	// fmt.Printf("cf: %v\n", cf)
	// port := os.Getenv("PORT")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	go routine.NewFetchKlineRoutine(klineRepository)
	app.Listen(fmt.Sprintf(":%v", cf.Port))
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

func makeSymbolAsIndexes(collection *mongo.Collection) {
	indexName, err := collection.Indexes().CreateOne(
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
}
