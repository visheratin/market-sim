package market

import (
	"fmt"
	"math/rand"

	"github.com/visheratin/market-sim/chain"
	"github.com/visheratin/market-sim/data"
)

type Market struct {
	providers     []*Provider
	chainInstance *chain.Chain
}

func Create() (Market, error) {
	market := Market{
		providers: []*Provider{},
	}
	market.chainInstance = &chain.Chain{}
	err := market.chainInstance.Init(nil)
	if err != nil {
		return Market{}, err
	}
	return market, nil
}

func (market *Market) AddProvider(provider Provider) {
	provider.init(market.chainInstance)
	provider.chainInstance = market.chainInstance
	market.providers = append(market.providers, &provider)
}

func (market *Market) GetProvider(idx int) (Provider, error) {
	if idx <= (len(market.providers) - 1) {
		return *market.providers[idx], nil
	}
	return Provider{}, fmt.Errorf("index is out of range")
}

func (market *Market) GetProviderByID(id string) (Provider, error) {
	for _, provider := range market.providers {
		if provider.ID == id {
			return *provider, nil
		}
	}
	return Provider{}, fmt.Errorf("provider was not found")
}

func (market *Market) ProvidersNum() int {
	return len(market.providers)
}

func (market *Market) UploadData(data []data.Data, providerID string) ([]string, []string, error) {
	correctDataIDs := []string{}
	corruptDataIDs := []string{}
	for idx, provider := range market.providers {
		if provider.ID == providerID {
			market.providers[idx].UploadData(data)
			for i := 0; i < len(data)-1; i++ {
				random := rand.Float32()
				if random > 0.99999 {
					market.providers[idx].CorruptData(data[i].ID)
					corruptDataIDs = append(corruptDataIDs, data[i].ID)
				} else {
					correctDataIDs = append(correctDataIDs, data[i].ID)
				}
			}
			return correctDataIDs, corruptDataIDs, nil
		}
	}
	return nil, nil, fmt.Errorf("provider was not found")
}

func (market *Market) SearchData(criteria []data.Criterion) ([]data.Data, error) {
	result := []data.Data{}
	for _, provider := range market.providers {
		for idx, data := range provider.data {
			isGood := true
			for key, value := range data.Metadata {
				for _, criterion := range criteria {
					if criterion.Type != key || criterion.Value != value {
						isGood = false
						break
					}
				}
				if !isGood {
					break
				}
			}
			if isGood {
				result = append(result, provider.data[idx])
				result[len(result)-1].Contents = nil
			}
		}
	}
	return result, nil
}

func (market *Market) validateData(data data.Data) error {
	for _, provider := range market.providers {
		if provider.ID == data.ProviderID {
			return provider.TestData(data.ID)
		}
	}
	return fmt.Errorf("provider with specified data was not found")
}

func (market *Market) ProvideData(data []data.Data) ([]data.Data, error) {
	for _, item := range data {
		err := market.validateData(item)
		if err != nil {
			fmt.Println(err)
			continue
		}
		var provider Provider
		for _, prov := range market.providers {
			if prov.ID == item.ProviderID {
				provider = *prov
				break
			}
		}
		item.Contents = provider.GetData(item.ID)
	}
	return data, nil
}
