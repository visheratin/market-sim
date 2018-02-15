package test

import (
	"bufio"
	"fmt"
	"os"

	uuid "github.com/satori/go.uuid"
	"github.com/visheratin/market-sim/data"
	"github.com/visheratin/market-sim/market"
)

func UploadData(market *market.Market, pathToData string) ([]string, []string, error) {
	file, err := os.Open(pathToData)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	correctDataIDs := []string{}
	corruptDataIDs := []string{}
	providerCounter := 0
	linesCounter := 0
	input := []data.Data{}
	for scanner.Scan() {
		line := scanner.Text()
		bytes := []byte(line)
		dataPart := data.Data{}
		id := uuid.NewV4()
		dataPart.ID = id.String()
		dataPart.Contents = bytes
		dataPart.Metadata = map[string]string{}
		dataPart.Metadata["id"] = dataPart.ID
		provider, err := market.GetProvider(providerCounter)
		if err == nil {
			dataPart.ProviderID = provider.ID
		} else {
			fmt.Println(err)
		}
		if linesCounter < 500 {
			input = append(input, dataPart)
			linesCounter++
		} else {
			correctDatas, corruptDatas, err := market.UploadData(input, provider.ID)
			input = []data.Data{}
			if err != nil {
				return nil, nil, fmt.Errorf("error on loading data")
			}
			correctDataIDs = append(correctDataIDs, correctDatas...)
			corruptDataIDs = append(corruptDataIDs, corruptDatas...)
			linesCounter = 0
			if providerCounter == (market.ProvidersNum() - 1) {
				providerCounter = 0
			} else {
				providerCounter++
			}
		}
	}
	return correctDataIDs, corruptDataIDs, nil
}

func SearchData(market *market.Market, id string) ([]data.Data, error) {
	criteria := []data.Criterion{}
	criteria = append(criteria, data.Criterion{Type: "id", Value: id})
	searchResult, err := market.SearchData(criteria)
	if err != nil {
		return nil, err
	}
	data, _ := market.ProvideData(searchResult)
	return data, nil
}
