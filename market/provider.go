package market

import (
	"fmt"
	"math/rand"

	"github.com/visheratin/market-sim/chain"
	"github.com/visheratin/market-sim/data"
)

type Provider struct {
	ID            string
	Title         string
	data          map[string]data.Data
	chainInstance *chain.Chain
}

func (provider *Provider) init(chainInstance *chain.Chain) {
	provider.data = map[string]data.Data{}
	provider.chainInstance = &chain.Chain{}
	provider.chainInstance.Init(chainInstance)
}

func (provider *Provider) GetData(id string) []byte {
	if data, ok := provider.data[id]; ok {
		return data.Contents
	}
	return nil
}

func (provider *Provider) UploadData(data []data.Data) {
	provider.chainInstance.Add(data)
	for _, item := range data {
		provider.data[item.ID] = item
	}
}

func (provider *Provider) TestData(id string) error {
	for _, item := range provider.data {
		if item.ID == id {
			err := provider.chainInstance.ValidateData(item)
			return err
		}
	}
	return fmt.Errorf("data with specified ID was not found")
}

func (provider *Provider) CorruptData(dataID string) error {
	if data, ok := provider.data[dataID]; ok {
		for idx := range data.Contents {
			data.Contents[idx] += byte(rand.Intn(255))
		}
		provider.data[dataID] = data
		return nil
	}
	return fmt.Errorf("data was not found")
}
