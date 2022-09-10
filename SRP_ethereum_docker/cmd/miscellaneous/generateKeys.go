package main

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/core"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
)

var (
	MainDir                      = getHomeDir()
	generatedKeystoreDir         = MainDir + "/files/generatedKeystore_keystore/keystore"
	validatorsKeystoreDir        = MainDir + "/files/keystore"
	etherTransferKeystoreDir     = MainDir + "/files/etherTransferKeystore_keystore/etherTransferKeystore" //keystore"
	genesisPath                  = MainDir + "/files/genesis.json"
	genesisNewPath               = MainDir + "/files/newGenesis.json"
	txOptsPath                   = MainDir + "/files/txOpts.gob"
	numOfReqAccounts             = 15000
	numOfValidators          int = 0
	//numOfAccounts is the number of non-validator accounts
	numOfAccounts int = 10 - numOfValidators
)

// getHomeDir returns the main directory path of current user
func getHomeDir() string {
	path, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	if filepath.Base(path) == "runSRP" || filepath.Base(path) == "collectMeasurements" || filepath.Base(path) == "miscellaneous" {
		path = filepath.Dir(path)
	}
	if filepath.Base(path) == "cmd" {
		path = filepath.Dir(path)
	}
	fmt.Println("Main Directory: ", path)
	return path
}

func main() {

	//accs := getAccounts(generatedKeystoreDir)
	//addAccsToGenesis(accs, 0, 5000)

	/*txOpts := getTxOpts(0, 8271, generatedKeystoreDir)
	err := writeGob(txOptsPath, txOpts)
	if err != nil {
		fmt.Println(err)
	}
	//WriteTxOptsJSONFile(txOptsPath, txOpts)*/

	// print number of accounts currently stored in given keystore directory
	accs := getAccounts(etherTransferKeystoreDir)
	fmt.Printf("*** There are %d accounts ***\n", len(accs))
	//fmt.Printf("*** Generating %d new accounts to reach %d accounts ***\n", numOfReqAccounts, numOfReqAccounts+len(accs))
	// generate numOfReqAccounts accounts & store them in given keystore directory
	//generateAccounts(numOfReqAccounts, etherTransferKeystoreDir)
	// delete any account having same address as any of the accounts stored in validatorsKeystoreDir
	//compareKeystores(etherTransferKeystoreDir, validatorsKeystoreDir)
	//fmt.Printf("*** Adding %d accounts to genesis file ***\n", len(accs))
	//addAccsToGenesis(accs, 0, len(accs))
}

//getTxOpts creates and returns a list of transaction signers obtained from accounts stored in keystore
//start is the starting index of the account list (counting from left), while cut is number of accounts to be
//execluded (counting from the right), it is set to zero if all accounts following the account at start are included
func getTxOptsOfGeneratedKeys(start int, cut int) []*bind.TransactOpts {
	var transactOpts []*bind.TransactOpts
	var errUnlock error
	// create a keystore object for local keystore directory
	ks := keystore.NewKeyStore(generatedKeystoreDir, keystore.StandardScryptN, keystore.StandardScryptP)
	// get all accounts stored in the directory
	accounts := ks.Accounts()
	fmt.Printf("\n**** Extracted %d out of %d accounts ****\n", len(accounts), numOfReqAccounts+numOfValidators)
	if len(accounts) >= numOfReqAccounts+numOfValidators {
		for i := start; i < len(accounts)-cut; i++ {
			// unlock the account
			errUnlock = ks.Unlock(accounts[i], "")
			if errUnlock != nil {
				fmt.Printf("Account unlock error for account %s, err: %v\n", accounts[i].Address.Hex(), errUnlock)
				panic(errUnlock)
			}
			// set the transaction options using the extracted variables
			txOpts, err := bind.NewKeyStoreTransactor(ks, accounts[i])
			if err != nil {
				fmt.Printf("Couldn't create transaction signer for account %s, err: %v\n", accounts[i].Address.Hex(), err)
				panic(err)
			}
			transactOpts = append(transactOpts, txOpts)
			fmt.Printf("i:%d - New Signer TxOpts for account %s\n", i, accounts[i].Address.Hex())

		}
	} else {
		fmt.Printf("Number of accounts (%d) is less than the expected (%d)!!\n", len(accounts), numOfReqAccounts+numOfValidators)
		return nil
	}

	return transactOpts
}

func WriteTxOptsJSONFile(filename string, txOpts []*bind.TransactOpts) {
	// Preparing the data to be marshalled and written.
	dataBytes, err := json.MarshalIndent(txOpts, "", "  ")
	if err != nil {
		fmt.Println(err)
	}
	err = ioutil.WriteFile(filename, dataBytes, 0644)
	if err != nil {
		fmt.Println(err)
	}
}
func ReadJSONFile(filename string) []*bind.TransactOpts {
	//check the existence of file & handle errors
	err := checkFile(filename)
	if err != nil {
		fmt.Println(err)
	}

	// read the file
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
	}

	// unmarshal the data
	var txOpts []*bind.TransactOpts
	errUnmarshal := json.Unmarshal([]byte(file), &txOpts)
	if err != nil {
		fmt.Println(errUnmarshal)
	}

	return txOpts
}

//checkFile checks if given file exists, otherwise creates a new one
func checkFile(filename string) error {
	//check file's status
	_, err := os.Stat(filename)
	//create file if it doesn't exists
	if os.IsNotExist(err) {
		_, err := os.Create(filename)
		if err != nil {
			return err
		}
	}
	return nil
}
func generateAccounts(numOfAccounts int, keystoreDir string) []accounts.Account {
	var accs []accounts.Account
	// create a keystore object for local keystore directory
	ks := keystore.NewKeyStore(keystoreDir, keystore.StandardScryptN, keystore.StandardScryptP)

	for i := 0; i < numOfAccounts; i++ {
		// create a new account with empty passphrase
		acc, err := ks.NewAccount("")
		if err != nil {
			fmt.Printf("Error while attempting to generate %dth account:\n  %v\n", i, err)
			return nil
		}
		// append created account to list of accounts
		//fmt.Printf(" %dth account with address: %s \n", i, acc.Address.Hex())
		accs = append(accs, acc)

		// export account
		ks.Export(acc, "", "")
		fmt.Printf("Created & Exported %dth account with address: %s \n", i, acc.Address.Hex())

	}

	fmt.Printf("Generated & exported %d accounts \n", len(accs))

	return accs

}

func getAccounts(keystoreDir string) []accounts.Account {
	// create a keystore object for local keystore directory
	ks := keystore.NewKeyStore(keystoreDir, keystore.StandardScryptN, keystore.StandardScryptP)
	// get all accounts stored in the directory
	accs := ks.Accounts()

	fmt.Printf("Extracted %d accounts from dir:%s\n", len(accs), keystoreDir)

	return accs

}

func compareKeystores(srcKeystoreDir string, oldKeystoreDir string) int {
	numOfReplicas := 0
	// create a keystore object for old keystore directory
	oldks := keystore.NewKeyStore(oldKeystoreDir, keystore.StandardScryptN, keystore.StandardScryptP)
	// get all accounts stored in the directory
	oldAccs := oldks.Accounts()

	// create a keystore object for old keystore directory
	genks := keystore.NewKeyStore(srcKeystoreDir, keystore.StandardScryptN, keystore.StandardScryptP)
	// get all accounts stored in the directory
	genAccs := genks.Accounts()

	for i := 0; i < len(genAccs); i++ {
		for j := 0; j < len(oldAccs); j++ {
			if genAccs[i].Address == oldAccs[j].Address {
				fmt.Printf("Found an old account within generated accounts!!\n OldAcc #:%d addr: %s \n Deleting it from genDir", j, oldAccs[j].Address)
				genks.Delete(genAccs[i], "")
				numOfReplicas++
			}
		}
	}

	fmt.Printf("Number of account replicas is %d \n", numOfReplicas)

	return numOfReplicas

}

func addAccsToGenesis(accs []accounts.Account, start int, end int) {

	if end > len(accs) || start < 0 {
		fmt.Println("indices are out of bound !!")
		return
	}
	// create a genesis struct
	genesis := new(core.Genesis)

	// open genesis file
	genesisFile, err := os.Open(genesisPath)
	if err != nil {
		fmt.Printf("Failed to read genesis file: %v", err)
		return
	}
	defer genesisFile.Close()
	// decode genesis file
	if err := json.NewDecoder(genesisFile).Decode(genesis); err != nil {
		fmt.Printf("invalid genesis file: %v\n", err)
		return
	}
	// or use built-in Unmarshal
	//genesis.UnmarshalJSON(data)

	//genesis.Alloc[k]
	i := 0
	for k, v := range genesis.Alloc {
		i++
		fmt.Printf("Genesis Account %d: %s - balance: %s\n", i, k.Hex(), v.Balance.String())
	}

	// balance to be assigned to new account
	balance := new(big.Int)
	balance.SetString("120000000000000000000", 10)

	numOfAddedAccs := 0
	// assign each of the generated accounts to a Genesis Account & allocate balance to it
	for j := start; j < end; j++ {
		// address of extracted account
		addr := accs[j].Address
		if _, found := genesis.Alloc[addr]; found {
			fmt.Printf("Skipping Account %d: %s - already exists with balance: %s\n", j, addr.Hex(), genesis.Alloc[addr].Balance.String())
		} else {
			// create genesis account to be assigned to new account
			genAcc := new(core.GenesisAccount)
			genAcc.Balance = balance
			genesis.Alloc[addr] = *genAcc
			fmt.Printf("Added new Genesis Account %d: %s - balance: %s\n", j, addr.Hex(), genesis.Alloc[addr].Balance.String())
			numOfAddedAccs++
		}

	}

	// marshal Genesis
	dataBytes, err := genesis.MarshalJSON()
	if err != nil {
		fmt.Printf("Error marhasling genesis: %v\n", err)
		return
	}
	// write marshaled data to JSON files
	err = ioutil.WriteFile(genesisNewPath, dataBytes, 0644)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Added %d accounts to genesis file\n", numOfAddedAccs)

}

//getTxOpts creates and returns a list of transaction signers obtained from accounts stored in keystore
//start is the starting index of the account list (counting from left), while cut is number of accounts to be
//execluded (counting from the right), it is set to zero if all accounts following the account at start are included
func getTxOpts(start int, cut int, keystoreDir string) []*bind.TransactOpts {
	var transactOpts []*bind.TransactOpts
	var errUnlock error
	// create a keystore object for local keystore directory
	ks := keystore.NewKeyStore(keystoreDir, keystore.StandardScryptN, keystore.StandardScryptP)
	// get all accounts stored in the directory
	accounts := ks.Accounts()
	fmt.Printf("\n**** Extracted %d accounts ****\n", len(accounts))
	if len(accounts) >= numOfAccounts+numOfValidators {
		for i := start; i < len(accounts)-cut; i++ {
			// unlock the account
			errUnlock = ks.Unlock(accounts[i], "")
			if errUnlock != nil {
				fmt.Printf("Account unlock error for account %s, err: %v\n", accounts[i].Address.Hex(), errUnlock)
				panic(errUnlock)
			}
			// set the transaction options using the extracted variables
			txOpts, err := bind.NewKeyStoreTransactor(ks, accounts[i])
			if err != nil {
				fmt.Printf("Couldn't create transaction signer for account %s, err: %v\n", accounts[i].Address.Hex(), err)
				panic(err)
			}
			transactOpts = append(transactOpts, txOpts)
			fmt.Printf("i:%d - New Signer TxOpts for account %s\n", i, accounts[i].Address.Hex())

		}
	} else {
		fmt.Printf("Number of accounts (%d) is less than the minimum expected (%d)!!\n", len(accounts), numOfAccounts+numOfValidators)
		return nil
	}

	return transactOpts
}

func writeGob(filePath string, object interface{}) error {
	file, err := os.Create(filePath)
	if err == nil {
		encoder := gob.NewEncoder(file)
		encoder.Encode(object)
	}
	file.Close()
	return err
}

func readGob(filePath string, object interface{}) error {
	file, err := os.Open(filePath)
	if err == nil {
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(object)
	}
	file.Close()
	return err
}
