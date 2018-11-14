package BLC

//创建创世区块
func (cli *CLI) createGenesisBlockchain(address string) {
	blockchain := NewBlockhain(address)
	defer blockchain.DB.Close()
}

