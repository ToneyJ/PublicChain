package BLC

import "fmt"

//查看余额
func (cli CLI) GetBalance(address string) {
	blockchain := BlockchainObject()
	defer blockchain.DB.Close()

	amount := blockchain.GetBalance(address)
	fmt.Println(amount)
}

