package chain

import (
	"fmt"
	"testing"
)

func TestNewBlockchain(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			"aaa",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBlockchain(); got != nil {
				got.AddBlock([]*Transaction{
					{
						ID:   nil,
						Vin:  nil,
						Vout: nil,
					},
				})
				//got.AddBlock("Send 2 more BTC to limit2")

				bci := got.Iterator()

				for block, next := bci.PreBlock(); next; {
					fmt.Printf("Prev. hash: %x \n", block.PrevBlockHash)
					fmt.Printf("DataHash: %s \n", block.HashTransactions())
					fmt.Printf("Hash: %x \n", block.Hash)
					fmt.Printf("Nonce: %d \n", block.Nonce)
					pow := NewProofOfWork(block)
					fmt.Printf("Pow: %t\n", pow.Validate())
					fmt.Println()
				}
			}
		})
	}
}
