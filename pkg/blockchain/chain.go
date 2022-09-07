package blockchain

type Chain struct {
	Blocks []*Block
}

func (c *Chain) ChainBlock(data string) {
	lastBlock := c.Blocks[len(c.Blocks)-1]
	newblock := CreateBlock(data, lastBlock.Hash)
	c.Blocks = append(c.Blocks, newblock)
}

func InitChain() *Chain {
	return &Chain{[]*Block{Genesis()}}
}
