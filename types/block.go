package types

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"time"

	"github.com/Ragnar-BY/go-blockchain/utils"
)

type BlockHeader struct {
	PrevBlockHash [32]byte
	DataHash      [32]byte

	Time  int64
	Nonce [8]byte

	Hash [32]byte
}

func NewBlockHeader(prevBlockHash [32]byte, dataHash [32]byte) *BlockHeader {

	return &BlockHeader{PrevBlockHash: prevBlockHash, DataHash: dataHash}
}

func (bh *BlockHeader) FindNonce() error {

	bh.Time = time.Now().UnixNano()
	header, err := bh.HeaderNoNonce()
	if err != nil {
		return err
	}
	nonce, hash, err := pow.Run(header)
	if err != nil {
		return err
	}
	bh.Hash = hash
	bh.Nonce = nonce
	return nil

}

//check if blockHeader hash is under PoW target
func (bh *BlockHeader) Validate() (bool, error) {
	header, err := bh.HeaderNoNonce()
	if err != nil {
		return false, err
	}
	return pow.IsValid(header, bh.Nonce)
}

//prevBlockHash+dataHash
func (bh *BlockHeader) HeaderNoNonce() ([32]byte, error) {
	return utils.EncodeAndHash([]interface{}{
		bh.PrevBlockHash[:],
		bh.DataHash[:],
		bh.Time,
	})
}

type Block struct {
	Header *BlockHeader
	Data   []byte
}

func NewBlock(h *BlockHeader, data []byte) *Block {
	block := &Block{Header: h, Data: data}
	return block
}

func (b *Block) Hash() [32]byte {

	return b.Header.Hash
}

//gob serialize and deserialize
func (b *Block) Serialize() ([]byte, error) {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)
	return result.Bytes(), err
}
func DeserializeBlock(b []byte) (*Block, error) {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(b))
	err := decoder.Decode(&block)

	return &block, err
}

func (b *Block) ToString() string {
	t := time.Unix(0, b.Header.Time)
	str := fmt.Sprintf("Block :[PrevHash: %x, Data: [%s] , Hash %x, CreatedAt %v]",
		b.Header.PrevBlockHash, b.Data, b.Hash(), t.Format("2006-01-02 15:04:05.99"))
	return str
}
