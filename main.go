package main

import (
	"encoding/hex"
	"log"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/btcsuite/btcutil/bech32"
	"github.com/decred/dcrd/dcrec/secp256k1"
)

type wallet struct {
	privateKey string
	publicKey  string
	address    string
}

func maxParallelism() int {
	maxProcs := runtime.GOMAXPROCS(0)
	numCPU := runtime.NumCPU()
	if maxProcs < numCPU {
		return maxProcs
	}
	return numCPU
}

func generateWallet() (*wallet, error) {
	var (
		privateKey, _    = secp256k1.GeneratePrivateKey()
		pubKeyCompressed = privateKey.PubKey().SerializeCompressed()
		address          string
		byteSlice        []byte
		err              error
	)

	byteSlice, err = bech32.ConvertBits(pubKeyCompressed, 8, 5, true)
	if err != nil {
		return nil, err
	}

	address, err = bech32.Encode("rdx", byteSlice)
	if err != nil {
		return nil, err
	}

	return &wallet{
		privateKey: hex.EncodeToString(privateKey.Serialize()),
		publicKey:  hex.EncodeToString(pubKeyCompressed),
		address:    address,
	}, nil
}

func matches(address string) bool {
	lookFor := []string{"radix", "radogs"}

	for _, suffix := range lookFor {
		if strings.HasSuffix(address, suffix) {
			return true
		}
	}

	return false
}

func main() {
	var (
		matchedWallet *wallet
		startedAt     = time.Now()
		mutex         = sync.Mutex{}
		stop          = make(chan bool)
		totalChecked  = 0
		resError      error
	)

	for i := 0; i < maxParallelism(); i++ {
		go func() {
			for {
				generatedWallet, err := generateWallet()
				if err != nil {
					mutex.Lock()
					resError = err
					break
				}
				address := generatedWallet.address
				if matches(address) {
					mutex.Lock()
					if matchedWallet != nil {
						break
					}
					matchedWallet = generatedWallet
					break
				} else {
					totalChecked = totalChecked + 1
				}

				// Log
				log.Println(address[len(address)-5:])
			}

			stop <- true
			mutex.Unlock()
		}()
		<-stop
	}

	if resError != nil {
		log.Println(resError.Error())
		return
	}

	timeRan := time.Now().Sub(startedAt)

	log.Println("matchedPrivateKey")
	log.Println(matchedWallet.privateKey)
	log.Println("matchedPublicKey")
	log.Println(matchedWallet.publicKey)
	log.Println("matchedAddress")
	log.Println(matchedWallet.address)
	log.Println(timeRan.String())
}
