package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/radogs/radix-vanity-address-generator/config"
	"github.com/radogs/radix-vanity-address-generator/validator"
	"github.com/radogs/radix-vanity-address-generator/wallet"
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
	fmt.Println("The longer the word the longer it will take!")

	wordsInput := bufio.NewScanner(os.Stdin)
	wordsInput.Scan()

	lookFor, err := validator.Validate(wordsInput.Text())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Would you like a mnemonic phrase ( 12 words for your wallet seed)? Y/n (mnemonic mode is currently a lot slower..)")
	wordsInput = bufio.NewScanner(os.Stdin)
	wordsInput.Scan()

	c := config.NewConfig()
	c.SetWords(lookFor)
	c.SetMnemonicMode(validator.ParseYn(wordsInput.Text()))

	err = work(c)
	if err != nil {
		log.Fatal(err)
	}
}

// @TODO move this to the generator package
func generate(isMnemonicMode bool) func() (*wallet.Wallet, error) {
	if isMnemonicMode {
		fmt.Println("Entering MnemonicMode")
		return wallet.GenerateWalletV2
	}

	fmt.Println("Entering quickMode")
	return wallet.GenerateWallet
}

func work(c *config.Config) error {
	var (
		generateWallet = generate(c.GetMnemonicMode())
		startedAt      = time.Now()
		lookFor        = c.GetWords()
		logAt          = startedAt.Add(30 * time.Second)
		totalChecked   = 0
	)

	fmt.Println("Woof! Let's go!")
	for {
		generatedWallet, err := generateWallet()
		if err != nil {
			return err
		}

		match := matches(lookFor, generatedWallet.Address)
		if match != nil {
			writeErr := writeFile(generatedWallet)
			if writeErr != nil {
				return writeErr
			}
			c.Match(*match)
			lookFor = c.GetWords()

			fmt.Printf("Matched %s! %d out of %d matched\n", *match, c.Found(), c.Total())
			if len(lookFor) == 0 {
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

	result := fmt.Sprintf(`
Address: %s
Public key: %s
Private key: %s
mnemonic phrase: %s`,
		wallet.Address, wallet.PublicKey, wallet.PrivateKey, wallet.Mnemonic)

	return ioutil.WriteFile(fmt.Sprintf("%s/%s.txt", exPath, wallet.Address), []byte(result), 0644)
}
