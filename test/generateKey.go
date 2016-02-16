package main

import (
	"encoding/hex"
	"fmt"
	"log"

	"getmelange.com/zooko/account"
	"github.com/melange-app/nmcd/btcutil"
	"github.com/melange-app/nmcd/chaincfg"
)

func main() {
	acc, err := account.CreateAccount()
	if err != nil {
		log.Fatal(err)
	}

	k := acc.Keys.Serialize()
	fmt.Println("Your wallet private key is:")
	fmt.Println(hex.EncodeToString(k))
	fmt.Println("Your address is:")

	pkHash := btcutil.Hash160(acc.Keys.PubKey().SerializeCompressed())
	addr, err := btcutil.NewAddressPubKeyHash(pkHash, &chaincfg.MainNetParams)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(addr)
}
