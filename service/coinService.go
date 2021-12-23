package service

import (
	"fmt"

	"github.com/seenark/coin-name/repository"
)

type coinService struct {
	Repo repository.CoinRepository
}

var AllCoins map[string]CoinResponse = map[string]CoinResponse{}

func NewCoinService(repo repository.CoinRepository) ICoinService {
	return &coinService{
		Repo: repo,
	}
}

func SetToCache(coins []CoinResponse) {
	AllCoins = make(map[string]CoinResponse)
	for _, coin := range coins {
		AllCoins[coin.Symbol] = coin
	}
}

func (cs coinService) GetFromCache(symbols []string) []CoinResponse {
	resCoin := []CoinResponse{}
	fmt.Printf("symbols: %v\n", len(symbols))
	if len(symbols) == 0 || (len(symbols) == 1 && symbols[0] == "") {
		for _, coin := range AllCoins {
			resCoin = append(resCoin, coin)
		}
		return resCoin
	}

	for _, coin := range symbols {
		newCoin := AllCoins[coin]
		resCoin = append(resCoin, newCoin)
	}
	return resCoin
}

func (cs coinService) FetchAllAndSetToCache() error {
	all, err := cs.GetAll([]string{}, "")
	if err != nil {
		return err
	}
	SetToCache(all)
	return nil
}

func (cs coinService) GetAll(symbol []string, name string) ([]CoinResponse, error) {
	coins, err := cs.Repo.GetAll(symbol, name)
	if err != nil {
		return nil, err
	}
	all := []CoinResponse{}
	for _, coin := range coins {
		all = append(all, CoinResponse{
			Symbol: coin.Symbol,
			Name:   coin.Name,
		})
	}
	return all, nil
}
func (cs coinService) GetNameBySymbol(symbol string) (*CoinResponse, error) {
	coin, err := cs.Repo.GetBySymbol(symbol)
	if err != nil {
		return nil, err
	}
	return &CoinResponse{
		Symbol: coin.Symbol,
		Name:   coin.Name,
	}, nil
}
func (cs coinService) Create(coin CoinResponse) (string, error) {
	newCoin := repository.Coin{
		Symbol: coin.Symbol,
		Name:   coin.Name,
	}
	id, err := cs.Repo.Create(newCoin)
	if err != nil {
		return "", err
	}
	go cs.FetchAllAndSetToCache()
	return id, nil
}

func (cs coinService) CreateMany(coins []CoinResponse) ([]string, error) {
	newCoins := []repository.Coin{}
	for _, c := range coins {
		newCoins = append(newCoins, repository.Coin{
			Symbol: c.Symbol,
			Name:   c.Name,
		})
	}
	ids, err := cs.Repo.CreateMany(newCoins)
	if err != nil {
		return nil, err
	}
	go cs.FetchAllAndSetToCache()
	return ids, nil
}

func (cs coinService) Update(symbol string, coinRes CoinResponse) (string, error) {
	find, err := cs.Repo.GetBySymbol(symbol)
	if err != nil {
		return "", err
	}
	newCoin := repository.Coin{
		Id:     find.Id,
		Symbol: coinRes.Symbol,
		Name:   coinRes.Name,
	}
	id, err := cs.Repo.Update(find.Id.Hex(), newCoin)
	if err != nil {
		return "", err
	}
	go cs.FetchAllAndSetToCache()
	return id, nil
}
func (cs coinService) Delete(symbol string) error {
	find, err := cs.Repo.GetBySymbol(symbol)
	if err != nil {
		return err
	}
	err = cs.Repo.Delete(find.Id.Hex())

	go cs.FetchAllAndSetToCache()
	return err
}
