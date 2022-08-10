package main

import (
	"bufio"
	"fmt"
	"github.com/radogs/radix-vanity-address-generator/config"
	"github.com/radogs/radix-vanity-address-generator/validator"
	"github.com/radogs/radix-vanity-address-generator/wallet"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

func main() {
	fmt.Print(`
--------------------------------------------
------ Radix Vanity address generator ------
------ Use at your own risk!          ------
------ Powered by radogs.tech, woof! -------
--------------------------------------------
`)
	fmt.Println("Enter your desired wallet suffix")
	fmt.Println("comma delimited, use the following format: d0gs,w00f,nfts")
	fmt.Println("For a single match just type the word: sweet")
	fmt.Println("The longer the word the longer it will take. Try not to use 'o'")

	wordsInput := bufio.NewScanner(os.Stdin)
	wordsInput.Scan()

	lookFor, err := validator.Validate(wordsInput.Text())
	if err != nil {
		panic(err)
	}

	c := config.NewConfig()
	c.SetWords(lookFor)

	err = work(c)
	if err != nil {
		panic(err)
	}
}

func work(c *config.Config) error {
	var (
		startedAt    = time.Now()
		mutex        = sync.Mutex{}
		stop         = make(chan bool)
		lookFor      = c.GetWords()
		logAt        time.Time
		totalChecked = 0
		resError     error
	)

	fmt.Println("Woof! Let's go!")
	logAt = startedAt.Add(30 * time.Second)

	go func() {
		for {
			generatedWallet, err := wallet.GenerateWallet()
			if err != nil {
				mutex.Lock()
				resError = err
				break
			}

			match := matches(lookFor, generatedWallet.Address)
			if match != nil {
				writeErr := writeFile(generatedWallet)
				if writeErr != nil {
					mutex.Lock()
					resError = err
					break
				}
				c.Match(*match)
				lookFor = c.GetWords()

				fmt.Printf("Matched %s! %d out of %d matched\n", *match, c.Found(), c.Total())
				if len(lookFor) == 0 {
					mutex.Lock()
					break
				}
			} else {
				now := time.Now()
				totalChecked = totalChecked + 1

				if now.After(logAt) {
					timeRan := now.Sub(startedAt)
					fmt.Printf("Generated %d wallets in %.1f minutes. %d/%d matches\n", totalChecked, timeRan.Minutes(), c.Found(), c.Total())
					logAt = now.Add(30 * time.Second)
				}
			}
		}

		stop <- true
		mutex.Unlock()
	}()
	<-stop

	if resError != nil {
		return resError
	}

	timeRan := time.Now().Sub(startedAt)

	fmt.Printf("Woof, done! Matched a total of %d address \n", c.Found())
	fmt.Printf("Use them at your own risk and have fun! \n")
	fmt.Printf("Time ran %s", timeRan.String())

	return nil
}

func matches(lookFor []string, address string) *string {
	for _, suffix := range lookFor {
		if strings.HasSuffix(address, suffix) {
			return &suffix
		}
	}

	return nil
}

func writeFile(wallet *wallet.Wallet) error {
	ex, err := os.Executable()
	if err != nil {
		return err
	}
	exPath := filepath.Dir(ex)

	result := fmt.Sprintf("Address:\n%s\nPublic key:\n%s\nPrivate key:\n%s",
		wallet.Address, wallet.PublicKey, wallet.PrivateKey)

	err = ioutil.WriteFile(fmt.Sprintf("%s/%s.txt", exPath, wallet.Address), []byte(result), 0644)
	if err != nil {
		return err
	}

	return nil
}
