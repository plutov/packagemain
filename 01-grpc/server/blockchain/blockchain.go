package blockchain

import (
	"bytes"
	"crypto/sha256"
)

// Block keeps block headers
type Block struct {
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
}

// Blockchain keeps a sequence of Blocks
type Blockchain struct {
	Blocks []*Block
}

// SetHash calculates and sets block hash
func (b *Block) SetHash() {
	headers := bytes.Join([][]byte{b.PrevBlockHash, b.Data}, []byte{})
	hash := sha256.Sum256(headers)

	b.Hash = hash[:]
}

// NewBlock creates and returns Block
func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{[]byte(data), prevBlockHash, []byte{}}
	block.SetHash()
	return block
}

// NewGenesisBlock creates and returns genesis Block
func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}

// AddBlock saves provided data as a block in the blockchain
func (bc *Blockchain) AddBlock(data string) *Block {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := NewBlock(data, prevBlock.Hash)
	bc.Blocks = append(bc.Blocks, newBlock)

	return newBlock
}

// NewBlockchain creates a new Blockchain with genesis Block
func NewBlockchain() *Blockchain {
	return &Blockchain{[]*Block{NewGenesisBlock()}}
}
