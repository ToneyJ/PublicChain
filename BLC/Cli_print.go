package BLC

//打印区块
func (cli *CLI) PrintChain() {

	BlockChain := BlockchainObject()
	defer BlockChain.DB.Close()
	BlockChain.Print()
}