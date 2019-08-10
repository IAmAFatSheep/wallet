package main

import (
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/pkg/errors"
	"log"
	"wallet/model"
)

type BtcService struct {
	client *rpcclient.Client
}

//根据address获取未花费的tx
func (btcSrv *BtcService) GetUnspentByAddress(address string) (unspents []btcjson.ListUnspentResult, err error) {
	btcAdd, err := btcutil.DecodeAddress(address, &chaincfg.RegressionNetParams)
	if err != nil {
		return nil, err
	}
	adds := [1]btcutil.Address{btcAdd}
	unspents, err = btcSrv.client.ListUnspentMinMaxAddresses(1, 999999, adds[:])
	if err != nil {
		return nil, err
	}
	return
}

//预估交易费用
//addrForm来源地址，addrTo去向地址
//transfer 转账金额
func (btcSrv *BtcService) CalculatorFee(addrFrom string, transfer float64, wallet *model.Wallet)( float64,float64,error)  {
	unspents, err := btcSrv.GetUnspentByAddress(addrFrom)
	if err != nil {
		return 0,0,err
	}

	//各种参数声明 可以构建为内部小对象
	num:=0								//需要的内部txInput
	outsu := float64(0)                 //unspent单子相加
	totalTran := transfer               //总共花费
	for _, v := range unspents {
		if v.Amount == 0 {
			continue
		}
		if outsu < totalTran {
			outsu += v.Amount
			num++
		} else {
			break
		}
	}
	min := (float64(num) * 180 + 1 * 34 + 10 - 40)/10e8		//max fee
	max := (float64(num) * 180 + 1 * 34 + 10 + 40)/10e8		//min fee
	return max,min,nil
}

//转账
//addrForm来源地址，addrTo去向地址
//transfer 转账金额
//fee 小费
//wallet from钱包地址+
func (btcSrv *BtcService) SendAddressToAddress(addrFrom, addrTo string, transfer, fee float64, wallet *model.Wallet) error {
	//数据库获取prv pub key等信息，便于调试--------START------
	//actf, err := dhSrv.GetAccountByAddress(addrFrom)
	//if err != nil {
	//	return err
	//}
	//----------------------------------------END-----------

	unspents, err := btcSrv.GetUnspentByAddress(addrFrom)
	if err != nil {
		return err
	}
	//各种参数声明 可以构建为内部小对象
	outsu := float64(0)                 //unspent单子相加
	feesum := fee                       //交易费总和
	totalTran := transfer + feesum      //总共花费
	var pkscripts [][]byte              //txin签名用script
	tx := wire.NewMsgTx(wire.TxVersion) //构造tx

	for _, v := range unspents {
		if v.Amount == 0 {
			continue
		}
		if outsu < totalTran {
			outsu += v.Amount
			{
				//txin输入-------start-----------------
				hash, _ := chainhash.NewHashFromStr(v.TxID)
				outPoint := wire.NewOutPoint(hash, v.Vout)
				txIn := wire.NewTxIn(outPoint, nil, nil)

				tx.AddTxIn(txIn)

				//设置签名用script
				txinPkScript, err := hex.DecodeString(v.ScriptPubKey)
				if err != nil {
					return err
				}
				pkscripts = append(pkscripts, txinPkScript)
			}
		} else {
			break
		}
	}
	//家里穷钱不够
	if outsu < totalTran {
		return errors.New("not enough money !")
	}
	// 输出1, 给form----------------找零-------------------
	addrf, err := btcutil.DecodeAddress(addrFrom, &chaincfg.RegressionNetParams)
	if err != nil {
		return err
	}
	pkScriptf, err := txscript.PayToAddrScript(addrf)
	if err != nil {
		return err
	}
	baf := int64((outsu - totalTran) * 1e8)
	tx.AddTxOut(wire.NewTxOut(baf, pkScriptf))
	//输出2，给to------------------付钱-----------------
	addrt, err := btcutil.DecodeAddress(addrTo, &chaincfg.RegressionNetParams)
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
	err = sign(tx, wallet.Prikey, pkscripts) //签名
	if err != nil {
		return err
	}
	//广播
	txHash, err := btcSrv.client.SendRawTransaction(tx, false)
	if err != nil {
		return err
	}
	//这里最好也记一下当前的block count,以便监听block count比此时高度
	//大6的时候去获取当前TX是否在公链有效
	//dhSrv.AddTx(txHash.String(), addrFrom, []string{addrFrom, addrTo})
	fmt.Println("Transaction successfully signed")
	fmt.Println(txHash.String())
	return nil
}

//ListAddressTransactions method not found;btcd NOTE: This is a btcwallet extension.
func (btcSrv *BtcService) GetTxByAddress(addrs []string, name string) (interface{}, error) {
	ct := len(addrs)
	addresses := make([]btcutil.Address, 0, ct)
	for _, v := range addrs {
		address, err := btcutil.DecodeAddress(v, &chaincfg.RegressionNetParams)
		if err != nil {
			log.Println("一个废物")
		} else {
			addresses = append(addresses, address)
		}
	}

	txs, err := btcSrv.client.ListAddressTransactions(addresses, name)
	if err != nil {
		return nil, err
	}
	return txs, nil
}

//验证交易是否被公链证实
//txid:交易id
func (btcSrv *BtcService) CheckTxMergerStatus(txId string) error {
	txHash, err := chainhash.NewHashFromStr(txId)
	if err != nil {
		return err
	}
	txResult, err := btcSrv.client.GetTransaction(txHash)
	if err != nil {
		return err
	}
	//pow共识机制当6个块确认后很难被修改
	if txResult.Confirmations < 6 {
		return errors.New("还未被足够确认！")
	}
	return nil
}

//签名
//privkey的compress方式需要与TxIn的
func sign(tx *wire.MsgTx, privKey string, pkScripts [][]byte) error {
	wif, err := btcutil.DecodeWIF(privKey)
	if err != nil {
		return err
	}
	/* lookupKey := func(a btcutil.Address) (*btcec.PrivateKey, bool, error) {
	    return wif.PrivKey, false, nil
	} */
	for i, _ := range tx.TxIn {
		script, err := txscript.SignatureScript(tx, i, pkScripts[i], txscript.SigHashAll, wif.PrivKey, false)
		//script, err := txscript.SignTxOutput(&chaincfg.RegressionNetParams, tx, i, pkScripts[i], txscript.SigHashAll, txscript.KeyClosure(lookupKey), nil, nil)
		if err != nil {
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
			return err
		}
		log.Println("Transaction successfully signed")
	}
	return nil
}
