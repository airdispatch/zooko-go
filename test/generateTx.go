package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"

	"getmelange.com/zooko/account"
	"github.com/melange-app/nmcd/btcec"
	"github.com/melange-app/nmcd/btcutil"
	"github.com/melange-app/nmcd/chaincfg"
)

var (
	privateKey = "a3ffd55805cf6be94c3b2b096c2ec5aa8ee2e0768ace200e15dd780d2321c918"
	utxo       = account.UTXO{
		TxID:     "0a73cc240a2faabf24c379f4182cef89957e5a608224a917dd6836a6703c97b6",
		Output:   1,
		Amount:   10e5,
		PkScript: nil,
	}
	pkScriptHex = "76a914cde96cbba3b67d01557c26be51a193939f71ec3b88ac"
	toAddress   = "NCvUPX465h7UypXJaUNy7u2FUxs1uGcKgr"
	// toAddress = "NFM8McCgXPMiCMCs37y4tza4pSMEEsXR1i"
)

func main() {
	pkScriptData, err := hex.DecodeString(pkScriptHex)
	if err != nil {
		log.Fatal(err)
	}

	data, err := hex.DecodeString(privateKey)
	if err != nil {
		log.Fatal(err)
	}

	key, _ := btcec.PrivKeyFromBytes(btcec.S256(), data)
	if err != nil {
		log.Fatal(err)
	}

	utxo.PkScript = pkScriptData

	acc := &account.Account{
		Keys:    key,
		Unspent: account.UTXOList([]*account.UTXO{&utxo}),
	}

	addr, err := btcutil.DecodeAddress(toAddress, &chaincfg.MainNetParams)
	if err != nil {
		fmt.Println("Addr error")
		log.Fatal(err)
	}

	// pkHashAddr := addr.(*btcutil.AddressPubKeyHash)
	// pkHashAddr.ScriptAddre

	pkHash := hex.EncodeToString(addr.ScriptAddress())

	tx, err := acc.TransferFunds(5e5, pkHash)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Here is your raw transaction:")

	buf := new(bytes.Buffer)
	err = tx.MsgTx.Serialize(buf)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(hex.EncodeToString(buf.Bytes()))
}
