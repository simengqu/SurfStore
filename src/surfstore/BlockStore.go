package surfstore

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
)

type BlockStore struct {
	BlockMap map[string]Block
}

func (bs *BlockStore) GetBlock(blockHash string, blockData *Block) error {
	// panic("todo")
	if _, ok := bs.BlockMap[blockHash]; ok {
		*blockData = bs.BlockMap[blockHash]
		return nil
	}
	return errors.New("Block data does not exist.")
}

func (bs *BlockStore) PutBlock(block Block, succ *bool) error {
	// panic("todo")
	h := sha256.Sum256(block.BlockData)
	he := hex.EncodeToString(h[:])
	bs.BlockMap[he] = block
	*succ = true
	return nil
}

func (bs *BlockStore) HasBlocks(blockHashesIn []string, blockHashesOut *[]string) error {
	// panic("todo")
	for i := 0; i < len(blockHashesIn); i++ {
		if _, ok := bs.BlockMap[blockHashesIn[i]]; ok {
			// blockHashesOut = append(*blockHashesOut, blockHashesIn[i])
		}
	}

	return nil
}

// This line guarantees all method for BlockStore are implemented
var _ BlockStoreInterface = new(BlockStore)
