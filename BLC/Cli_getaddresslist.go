package BLC

import "fmt"

//打印所有钱包地址
func (cli *CLI) AddressLists() {
	wallets,_:= NewWallets()

	for address,_ := range wallets.WalletsMap{
		fmt.Println(address)
	}
}
