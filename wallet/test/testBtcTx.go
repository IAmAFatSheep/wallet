package main

import (
	"fmt"
	"github.com/blockcypher/gobcy"
)


func main() {
	//btc := gobcy.API{"a32f202761be4948affc85fd1f0a7f93", "btc", "main"}
	//note the change to BlockCypher Testnet
	bcy := gobcy.API{"a32f202761be4948affc85fd1f0a7f93", "bcy", "test"}
	//generate two addresses
	addr1, err := bcy.GenAddrKeychain()
	addr2, err := bcy.GenAddrKeychain()
	fmt.Println(addr1.Address)
	//fmt.Println(addr2.Address)
	//fmt.Println(addr1.Private)
	//fmt.Println(addr1.Public)
	//use faucet to fund first
	_, err = bcy.Faucet(addr1, 3e5)
	if err != nil {
		fmt.Println(err)
	}
	//Post New TXSkeleton
	skel, err := bcy.NewTX(gobcy.TempNewTX(addr1.Address, addr2.Address, 2e5), false)
	//Sign it locally
	err = skel.Sign([]string{addr1.Private})
	if err != nil {
		fmt.Println(err)
	}
	//Send TXSkeleton
	skel, err = bcy.SendTX(skel)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", skel)
}
