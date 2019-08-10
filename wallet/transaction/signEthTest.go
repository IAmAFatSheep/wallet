package transaction

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"math/big"
	"testing"
)

// 交易发起方keystore文件地址
var fromKeyStoreFile = "";

// keystore文件对应的密码
var password = "";

//交易发送方地址
var fromAddress = ""

// 交易接收方地址
var toAddress = ""

// http服务地址, 例:http://localhost:8545
var httpUrl = "http://ip:port"

/*
	以太坊交易发送
*/
func TestSign(t *testing.T) {
	// 交易发送方
	// 获取私钥方式一，通过keystore文件

	fromKeystore, err := ioutil.ReadFile(fromKeyStoreFile)
	require.NoError(t, err)
	fromKey, err := keystore.DecryptKey(fromKeystore, password)
	fromPrivkey := fromKey.PrivateKey
	fromPubkey := fromPrivkey.PublicKey
	fromAddr := crypto.PubkeyToAddress(fromPubkey)

	// 获取私钥方式二，通过私钥字符串
	//privateKey, err := crypto.HexToECDSA("私钥字符串")

	// 交易接收方
	toAddr := common.HexToAddress(toAddress)

	// 交易发送方
	fromAddr = common.HexToAddress(fromAddress)
	// 数量
	amount := big.NewInt(14)

	// gasLimit
	var gasLimit uint64 = 300000

	// gasPrice
	var gasPrice *big.Int = big.NewInt(200)

	// 创建客户端
	client, err := ethclient.Dial(httpUrl)
	require.NoError(t, err)

	// nonce获取
	nonce, err := client.PendingNonceAt(context.Background(), fromAddr)

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
	//signedTx ,err := types.SignTx(tx,types.HomesteadSigner{},fromPrivkey)
	require.NoError(t, err)

	// 交易发送
	serr := client.SendTransaction(context.Background(), signedTx)
	if serr != nil {
		fmt.Println(serr)
	}

	// 等待挖矿完成
	bind.WaitMined(context.Background(), client, signedTx)

}

