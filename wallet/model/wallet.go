package model

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"github.com/btcsuite/btcd/btcec"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip32"
	"golang.org/x/crypto/ripemd160"
	"wallet/utils"
)

const privKeyBytesLen = 32

type Wallet struct {
	EcdsaPrivateKey *ecdsa.PrivateKey
	Prikey          string
	Pubkey          string
	BtcAddress      string
	EthAddress      string
}

//parametor: bit32.key , GenarateKey return privateKey with publicKey
func(w *Wallet) GenarateKey(k *bip32.Key) (string, string) {
	privateKey, _ := btcec.PrivKeyFromBytes(btcec.S256(), k.Key)
	//if err != nil {
	//	fmt.Println("generate private key fail")
	//}

	privateKeyECDSA := privateKey.ToECDSA()
	w.EcdsaPrivateKey=privateKeyECDSA

	//publicKey := privateKeyECDSA.Public()
	//
	//publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	//
	//if !ok {
	//	fmt.Println("generate public key fail")
	//}

	d := privateKeyECDSA.D.Bytes()
	b := make([]byte, 0, privKeyBytesLen)
	priKey := utils.PaddedAppend(privKeyBytesLen, b, d)
	pubKey := append(privateKeyECDSA.PublicKey.X.Bytes(), privateKeyECDSA.PublicKey.Y.Bytes()...)
	w.Prikey=hex.EncodeToString(priKey)
	w.Pubkey=hex.EncodeToString(pubKey)
	return w.Prikey,w.Pubkey
}

// GetBtcAddress returns btc wallet address
func (w *Wallet) GetBtcAddress() (address string) {
	/* See https://en.bitcoin.it/wiki/Technical_background_of_Bitcoin_addresses */

	/* Convert the public key to bytes */
	pub_bytes,_ := hex.DecodeString(w.Pubkey)
	//pub_bytes:=[]byte{}
	//fmt.Println("pub_byte",pub_bytes)

	/* SHA256 Hash */
	//fmt.Println("2 - Perform SHA-256 hashing on the public key")
	sha256_h := sha256.New()
	sha256_h.Reset()
	sha256_h.Write(pub_bytes)
	pub_hash_1 := sha256_h.Sum(nil)
	//fmt.Println(byteString(pub_hash_1))
	//fmt.Println("=======================")
	//fmt.Println("hash1",pub_hash_1)
	/* RIPEMD-160 Hash */
	//fmt.Println("3 - Perform RIPEMD-160 hashing on the result of SHA-256")
	ripemd160_h := ripemd160.New()
	ripemd160_h.Reset()
	ripemd160_h.Write(pub_hash_1)
	pub_hash_2 := ripemd160_h.Sum(nil)
	//fmt.Println(byteString(pub_hash_2))
	//fmt.Println("=======================")
	/* Convert hash bytes to base58 check encoded sequence */
	//fmt.Println("hash2",pub_hash_2)
	//testNet:Ox6f  mainNet:0x00
	address = utils.B58checkencode(0x00, pub_hash_2)
	w.BtcAddress = address
	return address
}

// GetEthAddress returns btc wallet address
func (w *Wallet) GetEthAddress(key *ecdsa.PrivateKey) (address string) {

	address = crypto.PubkeyToAddress(key.PublicKey).Hex()
	w.EthAddress = address
	return
}
