package BLC

type Blockchain struct {
  Blocks []*Block
}

//add new block
func (blockchain *Blockchain) AddBlock(data string) {
  //create new block
  prevBlock := blockchain.Blocks[len(blockchain.Blocks) - 1]
  block := NewBlock(data, prevBlock.Hash)
  //add block to blocks
  blockchain.Blocks = append(blockchain.Blocks, block)
}
//Create a blockchain with genesis block
func NewBlockChain() *Blockchain {
  return &Blockchain{[]*Block{NewGenesisBlock()}}
}
