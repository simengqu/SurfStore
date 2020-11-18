package surfstore

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	// "fmt"
)

type BlockStore struct {
	BlockMap map[string]Block
}

func (bs *BlockStore) GetBlock(blockHash string, blockData *Block) error {
	if _, ok := bs.BlockMap[blockHash]; ok {
		
		*blockData = bs.BlockMap[blockHash]
		return nil
	}
	return errors.New("Block data does not exist.")
}

func (bs *BlockStore) PutBlock(block Block, succ *bool) error {
	h := sha256.Sum256(block.BlockData)
	he := hex.EncodeToString(h[:])
	bs.BlockMap[he] = block
	return nil
}

func (bs *BlockStore) HasBlocks(blockHashesIn []string, blockHashesOut *[]string) error {
	for i := 0; i < len(blockHashesIn); i++ {
		if _, ok := bs.BlockMap[blockHashesIn[i]]; ok {
			// blockHashesOut = append(*blockHashesOut, blockHashesIn[i])
		}
	}

	return nil
}

// This line guarantees all method for BlockStore are implemented
var _ BlockStoreInterface = new(BlockStore)
