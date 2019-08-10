package main

import (
	"fmt"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
	"wallet/bip44"
	"wallet/model"
)

func main(){
	// Generate a mnemonic for memorization or user-friendly seeds
	entropy, _ := bip39.NewEntropy(128)
	mnemonic, _ := bip39.NewMnemonic(entropy)

	//"quiz theory exclude illness recall dynamic hub kit rocket domain destroy comfort"
	// Generate a Bip32 HD wallet for the mnemonic and a user supplied password
	seed := bip39.NewSeed("quiz theory exclude illness recall dynamic hub kit rocket domain destroy comfort", "Secret Passphrase")

	masterKey, _ := bip32.NewMasterKey(seed)
	publicKey := masterKey.PublicKey()

	// Display mnemonic and keys
	fmt.Println("Mnemonic: ", mnemonic)
	fmt.Println("Master private key: ", masterKey)
	fmt.Println("Master public key: ", publicKey)

	fkey, err := bip44.NewKeyFromMasterKey(masterKey, bip44.TypeBitcoin, bip32.FirstHardenedChild, 0, 0)
	if err!=nil {
		fmt.Println(err)
	}
	//privateKey, err := fkey.ECPrivKey()
	//privateKeyECDSA := privateKey.ToECDSA()
	//if err != nil {
	//	return nil, err
	//}
	//fpubkey:=fkey.PublicKey()
	//fmt.Println("btc private key: ", fkey)
	//fmt.Println("btc public key: ", GenaratePublicKey(fkey))
	var w model.Wallet
	pri,pub:=w.GenarateKey(fkey)
	fmt.Println("ecdsa",*w.EcdsaPrivateKey)
	fmt.Println("private:  ",pri)
	fmt.Println("public:  ",pub)
	fmt.Println(1)
	fmt.Println(w.GetBtcAddress())
}


//private  8cafd75c347886f32275cbc638e7e9099a77649bf3b5b3dcef56cde3d644d2ce
//address  mqPtk6HtBER5pdNciycXYKmPFq6NVZYiHw

//Mnemonic:  quiz theory exclude illness recall dynamic hub kit rocket domain destroy comfort
//address	 mzPdLAqLyN9n8Vx9zvrPLETCyex5zHY5CN

//private:   7bed4e771d732698ebe24f721b13a47cb61e4ecdf2e749305d5518bd515098df
