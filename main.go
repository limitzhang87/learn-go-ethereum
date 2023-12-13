package main

import (
	"fmt"
	"github.com/limitzhang87/learn-go-ethereum/chain"
	"log"
)

func main() {

	bc := chain.NewBlockchain()
	fmt.Println(bc.GetBalance("limit"))

	tx, err := chain.NewUTXOTransaction("limit", "limit2", 3, bc)
	if err != nil {
		log.Fatal(err)
	}
	bc.MinedBlock([]*chain.Transaction{tx}, "aa")

	tx, err = chain.NewUTXOTransaction("limit", "limit3", 4, bc)
	if err != nil {
		log.Fatal(err)
	}
	bc.MinedBlock([]*chain.Transaction{tx}, "bb")

	fmt.Println(bc.GetBalance("limit"))
	fmt.Println(bc.GetBalance("limit2"))
	fmt.Println(bc.GetBalance("limit3"))
}
