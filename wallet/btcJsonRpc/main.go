package main

import (
	"fmt"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"io/ioutil"
	"log"
	"path/filepath"
)


func main()  {
	ntfnHandlers := rpcclient.NotificationHandlers{
		OnFilteredBlockConnected: func(height int32, header *wire.BlockHeader, txns []*btcutil.Tx) {
			log.Printf("Block connected: %v (%d) %v",
				header.BlockHash(), height, header.Timestamp)
		},
		OnFilteredBlockDisconnected: func(height int32, header *wire.BlockHeader) {
			log.Printf("Block disconnected: %v (%d) %v",
				header.BlockHash(), height, header.Timestamp)
		},
	}

	// Connect to local btcd RPC server using websockets.
	btcdHomeDir := btcutil.AppDataDir("btcd", false)
	certs, err := ioutil.ReadFile(filepath.Join(btcdHomeDir, "rpc.cert"))
	if err != nil {
		log.Fatal(err)
	}
	connCfg := &rpcclient.ConnConfig{
		Host:         "localhost:8332",
		Endpoint:     "ws",
		User:         "rpcuser",
		Pass:         "rpcpass",
		Certificates: certs,
	}
	client, err := rpcclient.New(connCfg, &ntfnHandlers)
	if err!=nil {
		fmt.Println(err)
		return
	}
	fmt.Println(111)
	//bsv:=BtcService{client:client}
	//1HT7xU2Ngenf7D4yocz2SAcnNLW7rK8d4E
	//muyDoehpBExCbRRXLtDUpw5DaTb33UZeyG
	//unspend,err:= bsv.GetUnspentByAddress("muyDoehpBExCbRRXLtDUpw5DaTb33UZeyG")
	money,_:= client.GetBalance("12c6DSiU4Rq3P4ZxziKxzrL5LmMBrzjrJX")
	//client.SendRawTransaction()
	fmt.Println(money)
	num,_:= client.GetBlockCount()
	fmt.Println("block height:",num)
}
//310245