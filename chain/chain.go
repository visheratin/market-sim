package chain

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/visheratin/market-sim/data"
)

type Chain struct {
	blocks    []Block
	instances []*Chain
}

func (chain *Chain) Init(parent *Chain) error {
	if parent != nil {
		chain.blocks = parent.blocks
		chain.instances = parent.instances
		chain.instances = append(chain.instances, parent)
		for _, instance := range parent.instances {
			instance.instances = append(instance.instances, chain)
		}
		parent.instances = append(parent.instances, chain)
	} else {
		chain.blocks = []Block{}
		chain.instances = []*Chain{}
		chain.blocks = append(chain.blocks, Block{
			hash: "init block",
		})
	}
	return nil
}

func (chain *Chain) Add(data []data.Data) {
	newBlock := createBlock(chain.blocks[len(chain.blocks)-1].hash, data)
	newBlock.prevBlockHash = chain.blocks[len(chain.blocks)-1].hash
	chain.blocks = append(chain.blocks, newBlock)
	err := chain.Sync()
	if err != nil {
		fmt.Println(err)
		chain.recover()
	}
}

func (chain *Chain) recover() {
	for idx, block := range chain.blocks {
		blocksCounter := map[string]int{}
		blocks := map[string]Block{}
		chains := map[string][]*Chain{}
		blocksCounter[string(block.hash)] = 1
		blocks[block.hash] = block
		for _, instance := range chain.instances {
			if block.hash != instance.blocks[idx].hash {
				if _, ok := blocksCounter[instance.blocks[idx].hash]; !ok {
					blocksCounter[instance.blocks[idx].hash] = 1
					blocks[instance.blocks[idx].hash] = instance.blocks[idx]
					chains[instance.blocks[idx].hash] = []*Chain{}
				} else {
					blocksCounter[instance.blocks[idx].hash] = blocksCounter[instance.blocks[idx].hash] + 1
				}
				chains[instance.blocks[idx].hash] = append(chains[instance.blocks[idx].hash], instance)
			}
		}
		if len(blocksCounter) > 1 {
			max := 0
			var maxBlock Block
			for k, v := range blocksCounter {
				if v > max {
					max = v
					maxBlock = blocks[k]
				}
			}
			if block.hash != maxBlock.hash {
				chain.blocks = chain.blocks[:idx]
				newBlock := maxBlock.Copy()
				newBlock.prevBlockHash = chain.blocks[len(chain.blocks)-1].hash
				chain.blocks = append(chain.blocks, newBlock)
			}
			for _, instance := range chain.instances {
				if instance.blocks[idx].hash != maxBlock.hash {
					instance.blocks = instance.blocks[:idx]
					newBlock := maxBlock.Copy()
					newBlock.prevBlockHash = instance.blocks[len(instance.blocks)-1].hash
					instance.blocks = append(instance.blocks, newBlock)
				}
			}
		}
	}
}

func (chain *Chain) Sync() error {
	for _, instance := range chain.instances {
		newBlock := chain.blocks[len(chain.blocks)-1]
		prevBlock := &instance.blocks[len(instance.blocks)-1]
		if prevBlock.hash != newBlock.prevBlockHash {
			return fmt.Errorf("previous blocks do not match")
		}
		newBlock.prevBlockHash = prevBlock.hash
		instance.blocks = append(instance.blocks, newBlock)
	}
	return nil
}

func (chain *Chain) ValidateBlock(data []data.Data) error {
	if len(data) == 0 {
		return fmt.Errorf("data slice is empty")
	}
	dataID := data[0].ID
	for idx, block := range chain.blocks {
		for _, id := range block.dataIDs {
			if id == dataID {
				testBlock := createBlock(chain.blocks[idx-1].hash, data)
				if testBlock.hash == block.hash {
					return nil
				}
				return fmt.Errorf("data do not match")
			}
		}
	}
	return fmt.Errorf("data was not found")
}

func (chain *Chain) ValidateData(checkData data.Data) error {
	totalScore := 0
	output := make(chan error)
	go func(d data.Data) {
		output <- chain.checkData(d)
	}(checkData)
	for _, instance := range chain.instances {
		go func(d data.Data) {
			output <- instance.checkData(d)
		}(checkData)
	}
	counter := 0
	for err := range output {
		if err == nil {
			totalScore++
			if totalScore > (len(chain.instances)+1)/2.0 {
				return nil
			}
		} else {
			if (counter - totalScore) > (len(chain.instances)+1)/2.0 {
				return fmt.Errorf("data does not match blockchain hash")
			}
		}
		counter++
	}
	return fmt.Errorf("data does not match blockchain hash")
}

func (chain *Chain) checkData(data data.Data) error {
	latency := rand.Intn(10000)
	sleepDuration, _ := time.ParseDuration(strconv.Itoa(latency) + "ms")
	time.Sleep(sleepDuration)
	for _, block := range chain.blocks {
		if chainHash, ok := block.dataIDs[data.ID]; ok {
			h := sha256.New()
			h.Write(data.Contents)
			hash := string(h.Sum(nil))
			if chainHash == hash {
				return nil
			}
			return fmt.Errorf("data does not match blockchain hash")
		}
	}
	return fmt.Errorf("data was not found")
}

func (chain *Chain) GetRelatedData(dataID string) (map[string]string, error) {
	for _, block := range chain.blocks {
		for _, id := range block.dataIDs {
			if id == dataID {
				return block.dataIDs, nil
			}
		}
	}
	return nil, fmt.Errorf("data was not found")
}
