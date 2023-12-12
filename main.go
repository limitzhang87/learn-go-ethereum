package main

import (
	"encoding/hex"
	"fmt"
	"github.com/limitzhang87/learn-go-ethereum/chain"
)

func main() {
	bc := chain.NewBlockchain()

	txId, _ := hex.DecodeString("e43c8f4e5eae0a36b7cbd84e4962f5f8874c5d73055d70baae555d475e2bce55")
	tx := &chain.Transaction{
		VIn: []chain.TXInput{
			{
				TxId:     txId,
				VOutIdx:  0,
				FromAddr: "limit",
			},
		},
		VOut: []chain.TXOutput{
			{
				Value:  2,
				ToAddr: "limit2",
			},
			{
				Value:  8,
				ToAddr: "limit",
			},
		},
	}
	tx.SetId()
	bc.AddBlock([]*chain.Transaction{tx})

	fmt.Println(bc.GetBalance("limit"))
	fmt.Println(bc.GetBalance("limit2"))
}
