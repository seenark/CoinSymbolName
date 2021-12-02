package service

import "github.com/seenark/coin-name/repository"

type coinService struct {
	Repo repository.CoinRepository
}

func NewCoinService(repo repository.CoinRepository) ICoinService {
	return &coinService{
		Repo: repo,
	}
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
	return id, nil
}
func (cs coinService) Delete(symbol string) error {
	find, err := cs.Repo.GetBySymbol(symbol)
	if err != nil {
		return err
	}
	err = cs.Repo.Delete(find.Id.Hex())
	return err
}
