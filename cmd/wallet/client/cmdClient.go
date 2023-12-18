package client

import (
	"context"
	"flag"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/limitzhang87/learn-go-ethereum/chain"
	"log"
	"math/big"
	"os"
)

const (
	CmdCreateWallet = "createWallet"
	CmdTransfer     = "transfer"
)

type CmdClient struct {
	NetWork string // 区块链地址
	DataDir string // 数据路径
}

func NewCmdClient(netWork, dataDir string) *CmdClient {
	return &CmdClient{
		NetWork: netWork,
		DataDir: dataDir,
	}
}

func (c *CmdClient) Help() {
	fmt.Print("./wallet createWallet -pass PASSWORD --for create new wallet")
	fmt.Println("./wallet transfer -from FROM -to TO -value VALUE --for transfer from acct to `to`")
}

// Run 运行方法
func (c *CmdClient) Run() {
	// 判断参数是否正确
	if len(os.Args) < 2 {
		c.Help()
		os.Exit(-1)
	}

	// 1. 立flag
	cwCmd := flag.NewFlagSet(CmdCreateWallet, flag.ExitOnError)

	tCmd := flag.NewFlagSet(CmdTransfer, flag.ExitOnError)

	// 2. 立flag参数
	cmCmdPass := cwCmd.String("pass", "", "PASSWORD")

	tCmdFrom := tCmd.String("from", "", "FROM")
	tCmdTo := tCmd.String("to", "", "TO")
	tCmdTransfer := tCmd.Int64("value", 0, "VALUE")

	// 3. 解析命令行参数
	switch os.Args[1] {
	case CmdCreateWallet:
		err := cwCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Failed to Parse cw_cmd", err)
			return
		}
	case CmdTransfer:
		err := tCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Failed to Parse transfer_cmd", err)
			return
		}
	}

	// 4. 确认flag参数出现
	if cwCmd.Parsed() {
		fmt.Println("params is ", *cmCmdPass)
		err := c.createWallet(*cmCmdPass)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (c *CmdClient) createWallet(pass string) error {
	w, err := chain.NewWallet(c.DataDir)
	if err != nil {
		fmt.Println("Failed to create wallet", err)
		return err
	}
	return w.StoreKey(pass)
}

func (c *CmdClient) transfer(from, to string, value int64) error {
	// 1. 钱包加载
	wallet, err := chain.LoadWallet(from, c.DataDir)
	if err != nil {
		fmt.Println("load wallet err", err)
		return err
	}
	// 2. 连接到以太坊
	cli, err := ethclient.Dial(c.NetWork)
	defer cli.Close()

	// 3. 获取nonce
	fromAddr := common.HexToAddress(from)
	nonce, err := cli.NonceAt(context.Background(), fromAddr, nil)
	if err != nil {
		fmt.Println("get transfer nonce err", err)
		return err
	}
	// 4. 创建交易
	gasLimit := uint64(300000)
	gasPrice := big.NewInt(21000000000)
	amount := big.NewInt(value)
	toAddress := common.HexToAddress(to)
	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		GasPrice: gasPrice,
		Gas:      gasLimit,
		To:       &toAddress,
		Value:    amount,
		Data:     []byte("Salary"),
	})

	return nil
}
