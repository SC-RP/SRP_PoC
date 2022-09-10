package registration

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	sm "selectionManager"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	MainDir = getHomeDir()
	//SMAddress is the address of SelectionManager contract stored in genesis block
	SMAddress                = "0xDB571079aF66EDbB1a56d22809584d39C20001D9"
	genKeystoreDir           = MainDir + "/files/generatedKeystore_keystore/keystore"
	registeredSCsFilePath    = MainDir + "/files/registeredSCs.json"
	validatorKeystoreDir     = MainDir + "/files/keystore"
	etherTransferKeystoreDir = MainDir + "/files/etherTransferKeystore"
	GenesisPath              = MainDir + "/files/genesis.json"
	ParametersFilePath       = MainDir + "/files/parameter_configuration.json"

	sleepInterval time.Duration = 500 * time.Microsecond
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

// RegisterNodes registers validators by sending txs to selection manager contract
func RegisterNodes(NumOfAccounts int, NumOfValidators int, SMAddr common.Address, SM *sm.SelectionManager, client *ethclient.Client) {
	// get tx signers of stored accounts from which registration txs are sent
	ValtxOpts := GetTxOpts(0, NumOfAccounts, validatorKeystoreDir, NumOfAccounts+NumOfValidators)
	if ValtxOpts == nil {
		return
	}
	// get number of validators from extract list of tx signers
	lengthOfVals := len(ValtxOpts)

	// Send Registration tx, for each validator account
	var regtxs []*types.Transaction
	fmt.Printf("**RegisterCandidates::\n**** Extracted %d out of %d signers for Registration Txs ****\n", lengthOfVals, NumOfValidators)
	for i := 0; i < lengthOfVals; i++ {
		fmt.Printf("**RegisterCandidates::Registering Account %d: %s \n", i, ValtxOpts[i].From.Hex())
		tx, err := SM.Register(ValtxOpts[i])
		failedCall := true
		for err != nil {
			if failedCall {
				fmt.Println("**RegisterCandidates::Couldn't register candidates, Re-trying . . .")
				fmt.Printf("**RegisterCandidates::Calling Register() Failed for account:%d %s, err: %v\n", i, ValtxOpts[i].From.Hex(), err)
				fmt.Println("**RegisterCandidates::Re-trying ....")
				failedCall = false
			}
			time.Sleep(sleepInterval)
			tx, err = SM.Register(ValtxOpts[i])
			//return
		}
		//fmt.Printf("**RegisterCandidates::Registering Account %d: %s \n", i, ValtxOpts[i].From.Hex())
		regtxs = append(regtxs, tx)

	}
	fmt.Println()
	// wait until all nodes are registered as validators, if unsuccessful exit program
	for i := 0; i < len(regtxs); i++ {
		if regtxs[i] != nil {
			msg, err := regtxs[i].AsMessage(types.NewEIP155Signer(regtxs[i].ChainId()))
			if err != nil {
				fmt.Printf("**RegisterCandidates:: Error while getting message format of a reg-tx\nErr: %v\n", err)
				os.Exit(1)
			}
			receipt, err := bind.WaitMined(context.Background(), client, regtxs[i])
			if err != nil {
				fmt.Printf("**RegisterCandidates::Acc:%s - Error while waiting for reg-tx to get Mined\nErr: %v\n", msg.From().Hex(), err)
				os.Exit(1)
			}

			if receipt.Status == types.ReceiptStatusSuccessful {
				fmt.Printf("**RegisterCandidates::Acc:%s is registered as a candidate\n", msg.From().Hex())
			} else {
				fmt.Printf("**RegisterCandidates::Acc:%s - reg-tx failed\n", msg.From().Hex())
				os.Exit(1)
			}
		}
	}

}

// triggers Selection for each NC account
func AssignValidatorsToContracts(NumOfAccounts int, NumOfValidators int, NumOfContracts int, SM *sm.SelectionManager, client *ethclient.Client) {
	// get addresses of all NCs
	NCAddrs := GetSCsAddresses(GenesisPath, common.HexToAddress(SMAddress))
	// set index for NC contracts
	NCIndex := 0
	// set the transaction options using the extracted variables
	txOpts := GetTxOpts(0, NumOfAccounts, validatorKeystoreDir, NumOfAccounts+NumOfValidators)
	var selectiontxs []*types.Transaction
	var registeredSCs []string
	NumOfValidators_within_subset := big.NewInt(int64(NumOfValidators / NumOfContracts))
	fmt.Printf("**RegisterCandidates::Sending %d SC-Selection txs (for %d Contracts),\n where number of nodes within a subset is %d \n\n", len(txOpts), NumOfContracts, NumOfValidators_within_subset)

	for i := 0; i < len(txOpts) && i < NumOfContracts; i++ {
		// the following loop is added to select n validators sequentially for each contract
		tx, err := SM.SelectFixed(txOpts[i], NCAddrs[NCIndex], NumOfValidators_within_subset)
		if err != nil {
			fmt.Printf("**RegisterCandidates::Calling SelectFixed() Failed for contract %s, sender account:%d %s,\n err: %v\n", NCAddrs[NCIndex].Hex(), i, txOpts[i].From.Hex(), err)
			os.Exit(1)
		}
		fmt.Printf("**RegisterCandidates::Triggered Selection for contract %s, Sender Account %d: %s \n", NCAddrs[NCIndex].Hex(), i, txOpts[i].From.Hex())
		selectiontxs = append(selectiontxs, tx)

		registeredSCs = append(registeredSCs, NCAddrs[NCIndex].Hex())
		NCIndex = (NCIndex + 1) % len(NCAddrs) // using modulus to avoid out-of-bound index
	}
	fmt.Println()
	failed := false
	// Note: sender account is used as an index to map to the contract selection
	for i := 0; i < len(selectiontxs); i++ {
		if selectiontxs[i] != nil {
			msg, err := selectiontxs[i].AsMessage(types.NewEIP155Signer(selectiontxs[i].ChainId()))
			if err != nil {
				fmt.Printf("**RegisterCandidates:: Error while getting message format of a selection-tx\nErr: %v\n", err)
				os.Exit(1)
			}
			receipt, err := bind.WaitMined(context.Background(), client, selectiontxs[i])
			if err != nil {
				fmt.Printf("**RegisterCandidates::Acc:%s - Error while waiting for select-tx to get Mined\nErr: %v\n", msg.From().Hex(), err)
				os.Exit(1)
			}

			if receipt.Status == types.ReceiptStatusSuccessful {
				fmt.Printf("**RegisterCandidates::Acc:%s triggered selection successfully\n", msg.From().Hex())
			} else {
				fmt.Printf("**RegisterCandidates::Acc:%s - select-tx failed,  tx-hash = %s\n", msg.From().Hex(), receipt.TxHash.Hex())
				failed = true
			}
		}
	}

	if failed {
		fmt.Println("**RegisterCandidates::Failed to send selection txs")
		os.Exit(1)
	}
	// write list of registered victim contracts to a json file
	WriteJSONFile(registeredSCsFilePath, registeredSCs)
	fmt.Println("Wrote list of registered victim contracts to Json File:", registeredSCsFilePath)
	fmt.Println("DONE registering candidates and victim contracts")
	//}
	// Exit successfully -- used to resolve issue appeared when running from bash
	os.Exit(0)

}

// ConnectToNode tries to connect to blockchain node (miner/bootstrap)
func ConnectToNode() (*ethclient.Client, error) {
	client, errCon := ethclient.Dial("http://127.0.0.1:8545")
	noConnection := true
	for errCon != nil {
		if noConnection {
			fmt.Println("**RegisterCandidates::Couldn't connect, Re-trying . . .")
			noConnection = false
		}
		client, errCon = ethclient.Dial("http://127.0.0.1:8545")
		//panic(errCon)
	}

	return client, errCon
}

// GetTxOpts creates & returns a list of transaction signers obtained from accounts stored in keystore
//start is the starting index of the account list (counting from left), while cut is number of accounts to be
//execluded (counting from the right), it is set to zero if all accounts following the account at start are included
func GetTxOpts(start int, cut int, keystoreDir string, totalNumOfAccounts int) []*bind.TransactOpts {
	var transactOpts []*bind.TransactOpts
	var errUnlock error
	// create a keystore object for local keystore directory
	ks := keystore.NewKeyStore(keystoreDir, keystore.StandardScryptN, keystore.StandardScryptP)
	// get all accounts stored in the directory
	accounts := ks.Accounts()
	fmt.Printf("**RegisterCandidates::\n**** Extracted %d accounts ****\n", len(accounts))
	if len(accounts) >= totalNumOfAccounts {
		for i := start; i < len(accounts)-cut; i++ {
			// unlock the account
			errUnlock = ks.Unlock(accounts[i], "")
			if errUnlock != nil {
				fmt.Printf("**RegisterCandidates::Account unlock error for account %s, err: %v\n", accounts[i].Address.Hex(), errUnlock)
				panic(errUnlock)
			}
			// set the transaction options using the extracted variables
			txOpts, err := bind.NewKeyStoreTransactor(ks, accounts[i])
			if err != nil {
				fmt.Printf("**RegisterCandidates::Couldn't create transaction signer for account %s, err: %v\n", accounts[i].Address.Hex(), err)
				panic(err)
			}
			transactOpts = append(transactOpts, txOpts)
			fmt.Printf("**RegisterCandidates::i:%d - New Signer TxOpts for account %s\n", i, accounts[i].Address.Hex())

		}
	} else {
		fmt.Printf("**RegisterCandidates::Number of accounts (%d) is less than the minimum expected (%d)!!\n", len(accounts), totalNumOfAccounts)
		return nil
	}

	return transactOpts
}

///////////////Read & Write SCs Addresses////////////////

// GetSCsAddresses returns an array of all hard-coded addresses of SCs in
// the given genesis file, excluding the SelectionManager contract
func GetSCsAddresses(genesisPath string, SelectionManagerAddr common.Address) []common.Address {
	var SCsAddrs []common.Address
	//// Read genesis
	// create a genesis struct
	genesis := new(core.Genesis)
	// open genesis file
	genesisFile, err := os.Open(genesisPath)
	if err != nil {
		fmt.Printf("**RegisterCandidates::Failed to read genesis file: %v", err)
		panic(err)
	}
	defer genesisFile.Close()
	// decode genesis file
	if err := json.NewDecoder(genesisFile).Decode(genesis); err != nil {
		fmt.Printf("**RegisterCandidates::invalid genesis file: %v\n", err)
		panic(err)
	}

	//numOfSCs := 0
	//// Extract all hard-coded SCs' addresses (i.e. accounts with code)
	// iterate over all accounts in genesis
	for accAddr, acc := range genesis.Alloc {
		// if account contains code then store accAddr in list of SCs' addresses
		if acc.Code != nil && len(acc.Code) > 0 && accAddr != SelectionManagerAddr {
			SCsAddrs = append(SCsAddrs, accAddr)
			//numOfSCs++
		}
	}

	fmt.Printf("**RegisterCandidates::\nExtracted %d smart-contracts' addresses\n", len(SCsAddrs))

	return SCsAddrs
}

//WriteJSONFile writes given data (as a map of TxTimes) to the given file
func WriteJSONFile(filename string, addresses []string) bool {
	// Preparing the data to be marshalled and written.
	dataBytes, err := json.MarshalIndent(addresses, "", "  ")
	if err != nil {
		fmt.Println(err)
		return false
	}
	err = ioutil.WriteFile(filename, dataBytes, 0644)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}
