package main

import (
	"fmt"
	"strconv"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/visheratin/market-sim/market"
	"github.com/visheratin/market-sim/test"
)

func main() {
	m, _ := market.Create()
	for i := 0; i < 1000; i++ {
		provider := market.Provider{}
		provider.Title = "provider" + strconv.Itoa(i)
		id := uuid.NewV4()
		provider.ID = id.String()
		m.AddProvider(provider)
	}
	correctDataIDs, corruptDataIDs, err := test.UploadData(&m, "./input/All_GPUs.csv")
	if err != nil {
		fmt.Println(err)
		return
	}
	startTime := time.Now()
	test.SearchData(&m, correctDataIDs[0])
	fmt.Println(time.Since(startTime).String())
	startTime = time.Now()
	test.SearchData(&m, corruptDataIDs[0])
	fmt.Println(time.Since(startTime).String())
}
