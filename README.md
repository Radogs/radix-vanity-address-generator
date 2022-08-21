# Radogs Radix DLT vanity address generator
This CLI program attempts to generate a vanity wallet address with keys or mnemonic phrase and writes them to .txt files in the current folder or outputs them to the console when ran via Docker. 

**DO NOT SHARE THOSE WITH ANYONE!!!**

#### Powered by Radogs
Woof! Check us out at https://radogs.tech!

## Is this safe?
Always take care when working with private keys / mnemonic's! This application is meant to run on your own computer / Docker. 
So this application is as safe as your PC. As the license says its use at your own risk. If you're part of the Radix community and familiar with Go then please leave a comment!

## How to use
1. Navigate to the [releases page](https://github.com/Radogs/radix-vanity-address-generator/releases).
2. Select your OS
3. Download the executable
4. Run it on your system

## How to use with Docker
1. Check us out on [Docker](https://hub.docker.com/repository/docker/radogs/radix-vanity-address-generator)
2. Simply run ```docker run -i radogs/radix-vanity-address-generator```


## What is the difference between mnemonic and quick mode?
If you plan to use your address with a (hd) wallet like Radix wallet. You will need to use the mnemonic phrase to use this address in your wallet. Quick mode is for developers that for example only require the secret to sign transactions.

## How long will it take to generate my address?
The shorter the lookup word the quicker it will be. We recommend a max of 4-5 characters. Benchmarks will be added here later.

## Questions?
Join our Telegram! https://t.me/radogsNFT

## Work in progress
Todo list:
- Improve speed for mnemonic mode
- Move some logic out of main.go
- Add more test coverage
