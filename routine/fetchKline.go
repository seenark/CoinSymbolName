package routine

import (
	"fmt"
	"net/http"
	"time"

	"github.com/seenark/coin-name/binance"
	"github.com/seenark/coin-name/helpers"
	"github.com/seenark/coin-name/repository"
)

func NewFetchKlineRoutine(klineCollection repository.ICoinKLineRepository) {

	for {

		now := time.Now()
		// mil := now.Second()
		// if mil%10 == 0 {
		min := now.Minute()
		if min == 0 {
			all, err := klineCollection.GetMultiple([]string{})
			if err != nil {
				fmt.Println(err)
			}
			for _, sb := range all {
				err = updateKLineForSymbol(sb.Symbol, klineCollection)
				if err != nil {
					fmt.Println("Error in Fetch routine", err)
					continue
				}
			}
			helpers.PrintMemUsage()
			// time.Sleep(5 * time.Second)
			time.Sleep(58 * time.Minute)
		}

	}

}

func StoreHourKLineForSymbol(symbol string, klineCollection repository.ICoinKLineRepository) (*repository.CoinKLine, error) {
	now := time.Now()
	min := int64(now.Minute())
	minTime := time.Duration(min)
	end := now.Add(-minTime * time.Minute)
	start := end.AddDate(0, 0, -1)
	bClient := new(binance.BinanceClient)
	bClient.HttpClient = &http.Client{}
	klines := bClient.GetKLine(symbol, "1h", ToMilliseconds(start), ToMilliseconds(end), 24)
	closePrices := []float64{}
	for _, k := range klines {
		closePrices = append(closePrices, k.Close)
	}
	coinKl := repository.CoinKLine{
		Symbol:      symbol,
		ClosePrices: closePrices,
	}
	err := klineCollection.Create(coinKl)
	if err != nil {
		return nil, err
	}
	return &coinKl, nil
}

func updateKLineForSymbol(symbol string, klineCollection repository.ICoinKLineRepository) error {
	now := time.Now()
	min := int64(now.Minute())
	minTime := time.Duration(min)
	end := now.Add(-minTime * time.Minute)
	start := end.AddDate(0, 0, -1)
	bClient := new(binance.BinanceClient)
	bClient.HttpClient = &http.Client{}
	klines := bClient.GetKLine(symbol, "1h", ToMilliseconds(start), ToMilliseconds(end), 24)
	closePrices := []float64{}
	for _, k := range klines {
		closePrices = append(closePrices, k.Close)
	}
	coinKl := repository.CoinKLine{
		Symbol:      symbol,
		ClosePrices: closePrices,
	}
	err := klineCollection.Update(symbol, coinKl)
	if err != nil {
		return err
	}
	return nil
}

// helpers
func ToMilliseconds(t time.Time) int {
	return int(t.UnixNano()) / 1e6
}
