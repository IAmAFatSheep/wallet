package transaction

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"io/ioutil"
	"log"
	"net/http"
	"wallet/model"
)

type TotalUnspent struct {
	UnspentOutputs []UnspentOutput `json:"unspent_outputs"`
}

type UnspentOutput struct {
	TxHash          string `json:"tx_hash"`
	TxHashBigEndian string `json:"tx_hash_big_endian"`
	TxOutputN       int    `json:"tx_output_n"`
	Script          string `json:"script"`
	Value           int    `json:"value"`
	ValueHex        string `json:"value_hex"`
	Confirmations   int    `json:"confirmations"`
	TxIndex         int    `json:"tx_index"`
}

//交易广播后，响应结构
type TxResponse struct {
	Network string `json:"network"`
	Txid string `json:"txid"`
}

//从 https://blockchain.info 查询钱包地址的utxo
func GetTotalUnspentFromHttp(fromAddress string) (*TotalUnspent,error) {
	url := fmt.Sprintf("https://blockchain.info/unspent?active=%s", fromAddress)
	req, _ := http.NewRequest("GET", url, nil)
	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	if res.StatusCode!=200 {
		return nil,errors.New(string(body))
	}
	var total TotalUnspent
	if err := json.Unmarshal(body, &total); err != nil {
		fmt.Println(err)
		return nil,err
	}

	return &total,nil
}

//通过https://chain.so/api 广播交易
func Broadcast(txid string) (string,error)  {

	url := fmt.Sprintf("https://chain.so/api/v2/send_tx/DOGE?tx_hex=%s",txid)
	req, _ := http.NewRequest("POST", url, nil)
	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	if res.StatusCode!=200{
		return "",errors.New(string(body))
	}
	var tx TxResponse
	json.Unmarshal(body,&tx)
	return tx.Txid,nil
}



//转账
//addrForm来源地址，addrTo去向地址
//transfer 转账金额
//fee 小费
//wallet from钱包地址+
func SendAddressToAddress(addrFrom, addrTo string, transfer, fee float64, wallet *model.Wallet) error {

	unspents,err := GetTotalUnspentFromHttp(addrFrom)
	if err!=nil {
		return err
	}
	//各种参数声明 可以构建为内部小对象
	outsu := float64(0)                 //unspent单子相加
	feesum := fee                       //交易费总和
	totalTran := transfer + feesum      //总共花费
	var pkscripts [][]byte              //txin签名用script
	tx := wire.NewMsgTx(wire.TxVersion) //构造tx

	for _, v := range unspents.UnspentOutputs {
		if v.Value == 0 {
			continue
		}

		if outsu < totalTran {
			amount := float64(v.Value)
			outsu += amount
			{
				//txin输入-------start-----------------
				hash, _ := chainhash.NewHashFromStr(v.TxHash)
				outPoint := wire.NewOutPoint(hash, uint32(v.TxOutputN))
				txIn := wire.NewTxIn(outPoint, nil, nil)

				tx.AddTxIn(txIn)

				//设置签名用script
				txinPkScript, err := hex.DecodeString(v.Script)
				if err != nil {
					return err
				}
				pkscripts = append(pkscripts, txinPkScript)
			}
		} else {
			break
		}
	}

	if outsu < totalTran {
		return errors.New("not enough money !")
	}
	// 输出1, 给form----------------找零-------------------
	addrf, err := btcutil.DecodeAddress(addrFrom, &chaincfg.MainNetParams)
	if err != nil {
		fmt.Println(err)
		return err
	}

	pkScriptf, err := txscript.PayToAddrScript(addrf)
	if err != nil {
		fmt.Println(err)
		return err
	}
	baf := int64((outsu - totalTran) * 1e8)
	tx.AddTxOut(wire.NewTxOut(baf, pkScriptf))

	//输出2，给to------------------付钱-----------------
	addrt, err := btcutil.DecodeAddress(addrTo, &chaincfg.MainNetParams)
	if err != nil {
		return err
	}


	pkScriptt, err := txscript.PayToAddrScript(addrt)
	if err != nil {
		return err
	}
	bat := int64(transfer * 1e8)
	tx.AddTxOut(wire.NewTxOut(bat, pkScriptt))
	//-------------------输出填充end------------------------------
	err = sign(tx, wallet, pkscripts) //签名
	fmt.Println(2222)
	if err != nil {
		fmt.Println(333333)
		return err
	}
	fmt.Println(tx)

	//广播
	//txHash, err := btcSrv.client.SendRawTransaction(tx, false)
	txHash,err:=Broadcast(txSerializeString(tx))
	if err != nil {
		return err
	}
	////这里最好也记一下当前的block count,以便监听block count比此时高度
	////大6的时候去获取当前TX是否在公链有效
	fmt.Println("Transaction successfully signed and broadcast !")
	fmt.Println(txHash)
	return nil
}


//签名
//privkey的compress方式需要与TxIn的
func sign(tx *wire.MsgTx, wallet *model.Wallet, pkScripts [][]byte) error {
	var wkey btcec.PrivateKey
	wkey = btcec.PrivateKey(*wallet.EcdsaPrivateKey)
	//var ecdsa btcec.PrivateKey=w.EcdsaPrivateKey.(btcec.PrivateKey)
	wif, err := btcutil.NewWIF(&wkey, &chaincfg.MainNetParams, false)
	//add, err := btcutil.DecodeWIF(privKey)
	if err != nil {
		fmt.Println(err)
		return err
	}
	/* lookupKey := func(a btcutil.Address) (*btcec.PrivateKey, bool, error) {
	    return wif.PrivKey, false, nil
	} */
	for i, _ := range tx.TxIn {
		script, err := txscript.SignatureScript(tx, i, pkScripts[i], txscript.SigHashAll, wif.PrivKey, false)
		//script, err := txscript.SignTxOutput(&chaincfg.RegressionNetParams, tx, i, pkScripts[i], txscript.SigHashAll, txscript.KeyClosure(lookupKey), nil, nil)
		if err != nil {
			fmt.Println(err)
			return err
		}
		tx.TxIn[i].SignatureScript = script
		vm, err := txscript.NewEngine(pkScripts[i], tx, i,
			txscript.StandardVerifyFlags, nil, nil, -1)
		if err != nil {
			return err
		}
		err = vm.Execute()
		if err != nil {
			fmt.Println(err)
			return err
		}
		log.Println("Transaction successfully signed")
	}
	return nil
}


//交易序列化为16进制字符串
func txSerializeString(tx *wire.MsgTx) string {
	txHex := ""
	if tx != nil {
		// Serialize the transaction and convert to hex string.
		buf := bytes.NewBuffer(make([]byte, 0, tx.SerializeSize()))
		if err := tx.Serialize(buf); err != nil {
			log.Fatal("交易序列化失败!")
		}
		txHex = hex.EncodeToString(buf.Bytes())
	}
	return txHex
}