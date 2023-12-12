package main

import (
	"fmt"
	"github.com/limitzhang87/learn-go-ethereum/chain"
)

func main() {

	bc := chain.NewBlockchain()
	bc.AddBlock(nil)
	//bc.AddBlock("Send 2 more BTC to limit2")

	fmt.Println(bc.GetBalance("limit"))
}
