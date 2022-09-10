package core

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	gethLog "github.com/ethereum/go-ethereum/log"
)

/******** HELPER/UTILITY FUNCTIONS **********/
// watchFile watches for changes in a given file, check https://stackoverflow.com/questions/8270441/go-language-how-detect-file-changing
func watchFile(filePath string) error {
	initialStat, err := os.Stat(filePath)
	if err != nil {
		return err
	}

	for { // forever & break if change is detected
		stat, err := os.Stat(filePath)
		if err != nil {
			return err
		}
		// if size or modification time are not the same, then a change is detected
		if stat.Size() != initialStat.Size() || stat.ModTime() != initialStat.ModTime() {
			deployCounter++
			break
		}

		time.Sleep(1 * time.Second)
	}

	return nil
}

// importKey imports & returns the key of the first account in ./keystore,
// which is hardcoded below (eth.accounts[0], bootnode's assigned account)
func importKey() *keystore.Key {

	// NOTE: this is only suitable for the workspace of ether-docker-master
	passphrase := ""

	file := "/root/.ethereum/devchain/keystore/UTC--2016-02-29T14-52-47.838345060Z--794f74c8916310d6a0009bb8a43a5acab59a58ad"
	ks := keystore.NewKeyStore("/root/.ethereum/devchain/keystore", keystore.StandardScryptN, keystore.StandardScryptP)
	// Open the account key file
	jsonBytes, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	// Import the siging account
	signAcc, err := ks.Import(jsonBytes, passphrase, passphrase)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Printf("account found: signAcc.addr=%s; signAcc.url=%s\n", signAcc.Address.String(), signAcc.URL)
	//fmt.Println()

	// Unlock the signing account
	errUnlock := ks.Unlock(signAcc, passphrase)
	if errUnlock != nil {
		fmt.Println("account unlock error:")
		panic(err)
	}
	//fmt.Printf("account unlocked: signAcc.addr=%s; signAcc.url=%s\n", signAcc.Address.String(), signAcc.URL)
	//fmt.Println()

	// Get the private key
	keyWrapper, keyErr := keystore.DecryptKey(jsonBytes, passphrase)
	if keyErr != nil {
		fmt.Println("key decrypt error:")
		panic(keyErr)
	}

	if err := os.Remove(file); err != nil {
		log.Fatal(err)
	}

	return keyWrapper

}

//getAccount returns an account found in keystore given the address of the account
func getAccount(addr common.Address, ks *keystore.KeyStore) accounts.Account {
	var account accounts.Account
	var err error
	if ks.HasAddress(addr) {
		account, err = ks.Find(accounts.Account{Address: addr})
		if err != nil {
			gethLog.Warn("SCG: Method Call tx failed, Couldn't extract account", "address", addr)
			fmt.Printf("Error is: %v", err)
			//log.Panic(err)
		}
	} else {
		gethLog.Warn("SCG: Could not find account of given address", "address", addr)
		//fmt.Printf("Account not found given address:%s\n", addr.Hex())
	}

	return account
}

//getTxOpts creates and returns a transaction signer from a decrypted key of the given account stored in keystore
func getTxOpts(addr common.Address) *bind.TransactOpts {
	// create a keystore object for local keystore directory
	ks := keystore.NewKeyStore("/root/.ethereum/devchain/keystore", keystore.StandardScryptN, keystore.StandardScryptP)
	// get the account of the given address
	account := getAccount(addr, ks)
	// unlock the account
	errUnlock := ks.Unlock(account, "")
	if errUnlock != nil {
		gethLog.Warn("SCG: account unlock error:", "err", errUnlock)
		//panic(err)
	}
	// set the transaction options using the extracted variables
	transactOpts, err := bind.NewKeyStoreTransactor(ks, account)
	if err != nil {
		gethLog.Warn("SCG: result-submission tx failed, Couldn't create transaction signer", "err", err)
	}
	return transactOpts
}

// hex2int converts a long hex number to uint64
func hex2int(hexStr string) uint64 {
	// remove 0x suffix if found in the input string
	cleaned := strings.Replace(hexStr, "0x", "", -1)

	// base 16 for hexadecimal
	result, _ := strconv.ParseUint(cleaned, 16, 64)
	return uint64(result)
}

// Test prints/logs a test message to test scguard in geth
func Test() {
	fmt.Printf("Changed Testing SCGuard \n")
	gethLog.Info("SCG: Testing SCGuard..")
}
