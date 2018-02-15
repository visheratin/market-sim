package chain

import (
	"crypto/sha256"

	"github.com/visheratin/market-sim/data"
)

type Block struct {
	prevBlockHash string
	hash          string
	providerID    string
	dataIDs       map[string]string
}

func createBlock(prevHash string, data []data.Data) Block {
	block := Block{}
	rootHash, dataHashes := getRootHash(data)
	h := sha256.New()
	h.Write([]byte(prevHash))
	h.Write(rootHash)
	blockHash := h.Sum(nil)
	block.hash = string(blockHash)
	block.dataIDs = map[string]string{}
	for idx, item := range data {
		block.dataIDs[item.ID] = dataHashes[idx]
	}
	return block
}

func getRootHash(data []data.Data) ([]byte, []string) {
	hashes := [][]byte{}
	dataHashes := []string{}
	for _, item := range data {
		h := sha256.New()
		hashes = append(hashes, item.Contents)
		_, err := h.Write(item.Contents)
		if err != nil {
			dataHashes = append(dataHashes, "")
		} else {
			dataHashes = append(dataHashes, string(h.Sum(nil)))
		}
	}
	for {
		newHashes := [][]byte{}
		for i := 0; i < len(hashes)-1; i = i + 2 {
			h := sha256.New()
			h.Write(hashes[i])
			h.Write(hashes[i+1])
			newHashes = append(newHashes, h.Sum(nil))
		}
		if len(hashes)%2 == 1 {
			h := sha256.New()
			h.Write(hashes[len(hashes)-1])
			newHashes = append(newHashes, h.Sum(nil))
		}
		hashes = newHashes
		if len(newHashes) == 1 {
			break
		}
	}
	return hashes[0], dataHashes
}

func (block *Block) Copy() Block {
	result := Block{}
	result.hash = block.hash
	result.providerID = block.providerID
	result.dataIDs = block.dataIDs
	return result
}
