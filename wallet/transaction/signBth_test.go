package transaction

import (
	"fmt"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
	"testing"
	"wallet/bip44"
	"wallet/model"
)

func TestSendTx(t *testing.T) {
	//total:= GetTotalUnspentFromHttp("1HT7xU2Ngenf7D4yocz2SAcnNLW7rK8d4E")
	//// Generate a mnemonic for memorization or user-friendly seeds
	//entropy, _ := bip39.NewEntropy(128)
	//mnemonic, _ := bip39.NewMnemonic(entropy)

	// Generate a Bip32 HD wallet for the mnemonic and a user supplied password
	seed := bip39.NewSeed("quiz theory exclude illness recall dynamic hub kit rocket domain destroy comfort", "Secret Passphrase")

	masterKey, _ := bip32.NewMasterKey(seed)
	//publicKey := masterKey.PublicKey()

	// Display mnemonic and keys
	//fmt.Println("Mnemonic: ", mnemonic)
	//fmt.Println("Master private key: ", masterKey)
	//fmt.Println("Master public key: ", publicKey)

	fkey, _ := bip44.NewKeyFromMasterKey(masterKey, bip44.TypeBitcoin, bip32.FirstHardenedChild, 0, 0)
	var w model.Wallet
	w.GenarateKey(fkey)
	fromAddress := w.GetBtcAddress()
	//fmt.Println(fromAddress)
	toAddress := "1HT7xU2Ngenf7D4yocz2SAcnNLW7rK8d4E"
	//mqPtk6HtBER5pdNciycXYKmPFq6NVZYiHw
	err:= SendAddressToAddress(fromAddress, toAddress, 0.0001, 0.00001, &w)
	if err!=nil {
		fmt.Println(err)
	}
}

