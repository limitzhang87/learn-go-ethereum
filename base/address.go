package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
)

// limit : 0x2c2059b05CfF9Fde6027298b2cABE365BcF74DE3
func main() {
	address := common.HexToAddress("0x2c2059b05CfF9Fde6027298b2cABE365BcF74DE3")

	fmt.Println(address.Hex())
	fmt.Println(address.Bytes())
	fmt.Println(address.Big())
}
