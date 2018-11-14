package BLC

import (
	"fmt"
	"os"
)

//转账
func (cli *CLI) Send(from, to, amount []string) {
	if dbExists() == false {
		fmt.Println("数据不存在。。。")
		os.Exit(1)
	}

	blockchain := BlockchainObject()
	defer blockchain.DB.Close()

	blockchain.MineNewBlcok(from, to, amount)
}

