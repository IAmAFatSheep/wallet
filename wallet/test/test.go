package main

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"wallet/btcutil"
)

func main() {
	privKey:="8cafd75c347886f32275cbc638e7e9099a77649bf3b5b3dcef56cde3d644d2ce"
	fmt.Println(len(privKey))
	byte,_:=hex.DecodeString(privKey)
	s:=b58encode(byte)
	//fmt.Println()
	add, err := btcutil.DecodeWIF(s)
	if err!=nil {
		fmt.Println(err)
	}
	fmt.Println(add)
}


func b58encode(b []byte) (s string) {
	/* See https://en.bitcoin.it/wiki/Base58Check_encoding */

	const BITCOIN_BASE58_TABLE = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

	/* Convert big endian bytes to big int */
	x := new(big.Int).SetBytes(b)

	/* Initialize */
	r := new(big.Int)
	m := big.NewInt(58)
	zero := big.NewInt(0)
	s = ""

	/* Convert big int to string */
	for x.Cmp(zero) > 0 {
		/* x, r = (x / 58, x % 58) */
		x.QuoRem(x, m, r)
		/* Prepend ASCII character */
		s = string(BITCOIN_BASE58_TABLE[r.Int64()]) + s
	}

	return s
}