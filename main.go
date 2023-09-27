package main

import (
	"publicChain/BLC"
)

func main() {
	blockchain := BLC.CreateBlockchainWithGenesis()
	blockchain.AddBlockToBlockchain("to 100rmb")
	defer blockchain.DB.Close()

	blockchain.PrintChain()
}
