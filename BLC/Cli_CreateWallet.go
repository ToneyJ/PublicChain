package BLC

func (cli *CLI)CreateWallet(){
	wallets,_ := NewWallets()
	wallets.CreateNewWallet()
}