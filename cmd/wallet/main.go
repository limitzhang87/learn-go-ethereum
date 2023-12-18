package main

import "github.com/limitzhang87/learn-go-ethereum/cmd/wallet/client"

func main() {
	c := client.NewCmdClient("http:://localhost:8545", "./keystore")
	c.Run()
}
