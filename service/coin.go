package service

type CoinResponse struct {
	Symbol string `json:"symbol"`
	Name   string `json:"name"`
}

type ICoinService interface {
	GetAll([]string, string) ([]CoinResponse, error)
	GetNameBySymbol(string) (*CoinResponse, error)
	Create(CoinResponse) (string, error)
	CreateMany([]CoinResponse) ([]string, error)
	Update(string, CoinResponse) (string, error)
	Delete(string) error
	FetchAllAndSetToCache() error
	GetFromCache(symbols []string) []CoinResponse
}
