package main

import (
	"fmt"
	"github.com/limitzhang87/learn-go-ethereum/chain"
)

func main() {

	bc := chain.NewBlockchain()
	//bc.AddBlock("Sent 1 BTC to limit")
	//bc.AddBlock("Send 2 more BTC to limit2")

	bci := bc.Iterator()
	next := true
	for next {
		var block *chain.Block
		block, next = bci.PreBlock()
		fmt.Printf("Prev. hash: %x \n", block.PrevBlockHash)
		fmt.Printf("TxFromAddr: %s \n", block.Transactions[0].Vin[0].FromAddr)
		fmt.Printf("Hash: %x \n", block.Hash)
		fmt.Printf("Nonce: %d \n", block.Nonce)
		pow := chain.NewProofOfWork(block)
		fmt.Printf("Pow: %t\n", pow.Validate())
		fmt.Println()
	}
}
