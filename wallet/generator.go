package wallet

import (
	"encoding/hex"
	"log"

	"github.com/decred/dcrd/bech32"
	"github.com/decred/dcrd/dcrec/secp256k1/v2"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

type Wallet struct {
	Mnemonic   string
	PrivateKey string
	PublicKey  string
	Address    string
}

// GenerateWalletV2 generates a wallet from a mnemonic phrase.
func GenerateWalletV2() (*Wallet, error) {
	var (
		entropy      []byte
		mnemonic     string
		masterKey    *bip32.Key
		radixPathKey *bip32.Key
		index0Key    *bip32.Key
		address      string
		err          error
	)

	// @TODO: support multiple bitSizes
	entropy, err = bip39.NewEntropy(128)
	if err != nil {
		log.Println("Unable to generate entropy")
		return nil, err
	}

	mnemonic, err = bip39.NewMnemonic(entropy)
	if err != nil {
		log.Println("Unable to generate mnemonic")
		return nil, err
	}

	masterKey, err = bip32.NewMasterKey(seedFromMnemonic(mnemonic))
	if err != nil {
		log.Println("Unable to generate masterKey")
		return nil, err
	}

	radixPathKey, err = deriveRadixPath(masterKey)
	if err != nil {
		log.Println("Unable to derive Radix Path")
		return nil, err
	}

	// @TODO: Check if its faster to check the first 1-5 keys or generate a new one.
	index0Key, err = getKeyForIndex(radixPathKey, 0)
	if err != nil {
		log.Println("Unable to generate masterKey")
		return nil, err
	}

	address, err = deriveAddressFromPubKey(index0Key.PublicKey().Key, "rdx")
	if err != nil {
		log.Println("Unable to generate address from pubkey")
		return nil, err
	}

	return &Wallet{
		Mnemonic:   mnemonic,
		PrivateKey: masterKey.String(),
		PublicKey:  index0Key.PublicKey().String(),
		Address:    address,
	}, nil
}

// GenerateWallet generates wallets f
func GenerateWallet() (*Wallet, error) {
	var (
		privateKey, _    = secp256k1.GeneratePrivateKey()
		pubKeyCompressed = privateKey.PubKey().SerializeCompressed()
		address          string
		err              error
	)

	address, err = deriveAddressFromPubKey(pubKeyCompressed, "rdx")
	if err != nil {
		return nil, err
	}

	return &Wallet{
		PrivateKey: hex.EncodeToString(privateKey.Serialize()),
		PublicKey:  hex.EncodeToString(pubKeyCompressed),
		Address:    address,
	}, nil
}

func deriveAddressFromPubKey(pubKey []byte, hrp string) (string, error) {
	byteSlice, err := bech32.ConvertBits(append([]byte("\x04"), pubKey...), 8, 5, true)
	if err != nil {
		log.Println("ConvertBits err")
		return "", err
	}

	return bech32.Encode(hrp, byteSlice)
}

func seedFromMnemonic(mnemonic string) []byte {
	return bip39.NewSeed(mnemonic, "")
}

func deriveRadixPath(key *bip32.Key) (*bip32.Key, error) {
	var (
		count  uint32
		keyErr error
	)

	for _, s := range []struct {
		identifier uint32
		isHardened bool
	}{
		{
			identifier: 44,
			isHardened: true,
		},
		{
			identifier: 1022,
			isHardened: true,
		},
		{
			identifier: 0,
			isHardened: true,
		},
		{
			identifier: 0,
			isHardened: false,
		},
	} {
		if s.isHardened {
			count += bip32.FirstHardenedChild + s.identifier
			key, keyErr = key.NewChildKey(bip32.FirstHardenedChild + s.identifier)
			if keyErr != nil {
				return nil, keyErr
			}
		} else {
			count += s.identifier
			key, keyErr = key.NewChildKey(s.identifier)
		}
	}

	return key, nil
}

func getKeyForIndex(key *bip32.Key, index uint32) (*bip32.Key, error) {
	return key.NewChildKey(bip32.FirstHardenedChild + index)
}
