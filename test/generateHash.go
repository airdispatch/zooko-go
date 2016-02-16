package main

import (
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcutil"
)

func main() {
	name := "text/melange"
	rand, _ := hex.DecodeString("e4de61166713cf9e")

	hash := btcutil.Hash160(append(rand, []byte(name)...))
	fmt.Println(hex.EncodeToString(hash))
}
