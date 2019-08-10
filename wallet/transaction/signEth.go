package transaction

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"wallet/model"
)



//以太坊交易发送
//toAddress string,			   // receiver address
//amount *big.Int,			   // send amount
//gasLimit uint64,			   // gasLimit
//gasPrice *big.Int,	// gasProve
//httpUrl 			   // http服务地址
func SendTx(wallet *model.Wallet, toAddress string, amount *big.Int, gasLimit uint64, gasPrice *big.Int, httpUrl string) error {
	fromPrivkey:=wallet.EcdsaPrivateKey
	fromAddress:=wallet.EthAddress

	// 交易接收方
	toAddr := common.HexToAddress(toAddress)

	// 交易发送方
	fromAddr := common.HexToAddress(fromAddress)

	// 创建客户端
	client, err := ethclient.Dial(httpUrl)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// nonce获取
	nonce, err := client.PendingNonceAt(context.Background(), fromAddr)
	if err!=nil{
		return nil
	}

	// 认证信息组装
	auth := bind.NewKeyedTransactor(fromPrivkey)
	//auth,err := bind.NewTransactor(strings.NewReader(mykey),"111")
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = amount // in wei
	//auth.Value = big.NewInt(100000)     // in wei
	auth.GasLimit = gasLimit // in units
	//auth.GasLimit = uint64(0) // in units
	auth.GasPrice = gasPrice
	auth.From = fromAddr

	// 交易创建
	tx := types.NewTransaction(nonce, toAddr, amount, gasLimit, gasPrice, []byte{})

	// 交易签名
	signedTx, err := auth.Signer(types.HomesteadSigner{}, auth.From, tx)
	if err!=nil{
		return err
	}
	//signedTx ,err := types.SignTx(tx,types.HomesteadSigner{},fromPrivkey)

	// 交易发送
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// 等待挖矿完成
	bind.WaitMined(context.Background(), client, signedTx)
	return nil
}
