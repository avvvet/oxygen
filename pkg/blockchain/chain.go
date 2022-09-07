package blockchain

import "fmt"

type Chain struct {
	Blocks []*Block
}

func InitChain() (*Chain, error) {
	block, err := Genesis()
	return &Chain{[]*Block{block}}, err
}

func (c *Chain) ChainBlock(data string) {
	lastBlock := c.Blocks[len(c.Blocks)-1]
	newblock, err := CreateBlock(data, lastBlock.Hash)
	if err != nil {
		fmt.Print(err)
	} else {
		c.Blocks = append(c.Blocks, newblock)
	}
}
