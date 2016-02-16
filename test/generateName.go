package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"

	"getmelange.com/zooko/account"
	"github.com/melange-app/nmcd/btcec"
)

var (
	privateKey = "a3ffd55805cf6be94c3b2b096c2ec5aa8ee2e0768ace200e15dd780d2321c918"
	utxo       = account.UTXO{
		TxID: "1d220b9fdf1336ef60f299efb2e2c7bb8c83be2ea00603e9ee82956270bf91e8",
		// TxID:     "0a876e9710ffbe7b38266189b54e1d49568cdbfe80f15e80159bc9295f47592f",
		Output:   0,
		Amount:   10e5,
		PkScript: nil,
	}
	// pkScriptHex = "76a914cde96cbba3b67d01557c26be51a193939f71ec3b88ac"
	pkScriptHex = "5114be685eca1025e492218b749b3e8548d666d2a6096d76a914cde96cbba3b67d01557c26be51a193939f71ec3b88ac"
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

	// tx, rand, err := acc.CreateNameNew("test/melange")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	rand, err := hex.DecodeString("e4de61166713cf9e")
	if err != nil {
		log.Fatal(err)
	}

	tx, err := acc.CreateNameFirstUpdate(rand, "test/melange", "Hello from Golang!")
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

	// fmt.Println("Here is the random data:")
	// fmt.Println(hex.EncodeToString(rand))
}

// Generates the following raw transaction
// 00710000012f59475f29c99b15805ef180fedb8c56491d4eb5896126387bbeff10976e870a010000006a4730440220787fdbb3dbbe4d1c7c1d0c16fea1ea1abb9d22bd535c7d3509582b746d8bab9802207649dc3c27744602e49a300d18fea834b69a57b82e0616201d909c98faa8b295012103b5f3686b9333650664440d82f2e860ece1746e431f017aa4b85c53e695e919edffffffff0240420f0000000000305114be685eca1025e492218b749b3e8548d666d2a6096d76a914cde96cbba3b67d01557c26be51a193939f71ec3b88ace0673500000000001976a914cde96cbba3b67d01557c26be51a193939f71ec3b88ac00000000

// Name details:
// Name: test/melange
// Rand: e4de61166713cf9e
// 256263

// Name first update:
// 0071000001e891bf70629582eee90306a02ebe838cbbc7e2b2ef99f260ef3613df9f0b221d000000006b483045022100cee90e8b50de86d8eff8b75db5df4406bbe563a7b4a2d2c42196bc547e2df0d00220552b3c93c16bfbd3fb722af145d8626a790523f2c208624d8962e26e17c4ae6f012103b5f3686b9333650664440d82f2e860ece1746e431f017aa4b85c53e695e919edffffffff0140420f000000000045520c746573742f6d656c616e676508e4de61166713cf9e1248656c6c6f2066726f6d20476f6c616e67216d6d76a914cde96cbba3b67d01557c26be51a193939f71ec3b88ac00000000

// 0071000001e891bf70629582eee90306a02ebe838cbbc7e2b2ef99f260ef3613df9f0b221d000000006b483045022100c46c477f3d5cd3a65ad477433054f3c7eb7735d92130f2043500c6f1a0d2964702200671287072eb4bc8357582604d3368a1429b3cb905bb3c152e3a4782dcaa8e74012103b5f3686b9333650664440d82f2e860ece1746e431f017aa4b85c53e695e919edffffffff0140420f0000000000455208e4de61166713cf9e0c746573742f6d656c616e67651248656c6c6f2066726f6d20476f6c616e67216d6d76a914cde96cbba3b67d01557c26be51a193939f71ec3b88ac00000000
