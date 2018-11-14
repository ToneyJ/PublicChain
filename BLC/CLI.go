package BLC

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type CLI struct{}

//运行命令行
func (cli *CLI) Run() {
	isVaildArgs()

	addresslistcmd :=flag.NewFlagSet("addresslist", flag.ExitOnError)
	createwalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)
	sendBlockCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printBlockCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	getbalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)

	flagFrom := sendBlockCmd.String("from", "", "转账源地址。。。。")
	flagTo := sendBlockCmd.String("to", "", "转账目的地址。。。。")
	flaAmount := sendBlockCmd.String("amount", "", "转账金额")
	flagCreateBlockchain := createBlockchainCmd.String("address", "", "创世区块交易数据......")
	getbalanceWithAdress := getbalanceCmd.String("address", "", "要查询某一个账号的余额.......")

	switch os.Args[1] {
	case "send":
		err := sendBlockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		err := printBlockCmd.Parse(os.Args[:2])
		if err != nil {
			log.Panic(err)
		}
	case "createblockchain":
		err := createBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "getbalance":
		err := getbalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createwallet":
		err := createwalletCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "addresslist":
		err := addresslistcmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		printUsage()
		os.Exit(1)
	}

	//交易
	if sendBlockCmd.Parsed() {
		if *flagFrom == "" || *flagTo == "" || *flaAmount == "" {
			printUsage()
			os.Exit(1)
		}

		from := JSONToArray(*flagFrom)
		to := JSONToArray(*flagTo)
		for index,fromAddress := range from{
			if IsValidForAddress(fromAddress) == false || IsValidForAddress(to[index]) == false{
				fmt.Println("地址无效.....")
				printUsage()
				os.Exit(1)
			}
		}

		amount := JSONToArray(*flaAmount)
		cli.Send(from, to, amount)
	}

	//打印区块信息
	if printBlockCmd.Parsed() {
		cli.PrintChain()
	}

	//创建创世区块
	if createBlockchainCmd.Parsed() {

		if *flagCreateBlockchain == "" {
			fmt.Println("交易地址为空......")
			printUsage()
			os.Exit(1)
		}

		if len(*flagCreateBlockchain) <= 4 ||IsValidForAddress(*flagCreateBlockchain) == false {
			fmt.Println("地址无效.....")
			printUsage()
			os.Exit(1)
		}

		cli.createGenesisBlockchain(*flagCreateBlockchain)
	}

	//查询余额
	if getbalanceCmd.Parsed() {


		if len(*getbalanceWithAdress) <= 4 ||IsValidForAddress(*getbalanceWithAdress) == false {
			fmt.Println("地址无效.....")
			printUsage()
			os.Exit(1)
		}

		cli.GetBalance(*getbalanceWithAdress)
	}

	//创建钱包
	if createwalletCmd.Parsed() {
		cli.CreateWallet()
	}

	//打印钱包地址
	if addresslistcmd.Parsed() {
		cli.AddressLists()
	}
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("\taddresslist --输出所有钱包地址")
	fmt.Println("\tcreatewallet --创建钱包")
	fmt.Println("\tcreateblockchain -address --交易数据")
	fmt.Println("\tsend -from FROM -to TO -amount AMOUNT --交易明细.")
	fmt.Println("\tprintchain --输出区块信息")
	fmt.Println("\tgetbalance -address --输出区块信息.")
}

func isVaildArgs() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}
}
