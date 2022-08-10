package wallet

import (
	"encoding/hex"
	"log"

	"github.com/decred/dcrd/bech32"
	"github.com/decred/dcrd/dcrec/secp256k1"
)

type Wallet struct {
	PrivateKey string
	PublicKey  string
	Address    string
}

func GenerateWallet() (*Wallet, error) {
	var (
		privateKey, _    = secp256k1.GeneratePrivateKey()
		pubKeyCompressed = privateKey.PubKey().SerializeCompressed()
		prependedPubKey  = append([]byte("\x04"), pubKeyCompressed...)
		address          string
		byteSlice        []byte
		err              error
	)

	byteSlice, err = bech32.ConvertBits(prependedPubKey, 8, 5, true)
	if err != nil {
		log.Println("ConvertBits err")
		return nil, err
	}

	address, err = bech32.Encode("rdx", byteSlice)
	if err != nil {
		log.Println("bech32 Encode err")
		return nil, err
	}

	return &Wallet{
		PrivateKey: hex.EncodeToString(privateKey.Serialize()),
		PublicKey:  hex.EncodeToString(pubKeyCompressed),
		Address:    address,
	}, nil
}
