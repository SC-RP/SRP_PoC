package core

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"regexp"
	"strings"
	"sync"

	sm "selectionManager"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	gethMath "github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ethereum/go-ethereum/rpc"

	gethLog "github.com/ethereum/go-ethereum/log"

	//"github.com/ethereum/go-ethereum/crypto"
	//"runtime-protection/selectionManager"

	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	contractName                         = "SelectionManager"
	compiledContract                     = "/root//build" + contractName + ".json"
	configPath                           = "/root/config.json"
	parametersFilePath                   = "/root/files/parameter_configuration.json"
	InitialContractAddress               = "0xDB571079aF66EDbB1a56d22809584d39C20001D9"
	wordSize                             = 32
	txTimeout              time.Duration = 120 * time.Second
)

// RVTimeEmulation is the sleep time that emulates applying RV
var RVTimeEmulation time.Duration          //= 180 * time.Millisecond //(default is 180 milliseconds)
var NumOfValidators_within_subset *big.Int // this variable is the number of validators within a subset
var writingRecordsInterval int             // number of Tx Records at which json file of records is upated (originally every 500, now added from parameter_config)

// Backend wraps all methods required for getting state of  eth.Statedb.
type Backend interface {
	BlockChain() *BlockChain
	//Miner() *miner.Miner
	TxPool() *TxPool
}

var deployCounter uint = 0
var writeRecordCounter int = 1

// SCGuard class
type SCGuard struct {
	client            *ethclient.Client                 // client provider used to interact with on-chain execution (i.e. selection contract)
	candidateAddress  common.Address                    // publicKey/accountAddress of candidate, a potential validator
	selected          bool                              // true if candidate is currently selected to validate a given tx
	scs               map[common.Address]*SmartContract // list of txs assigned to candidate for runtime validation -- used to store candidate's index in selected list
	contract          *sm.SelectionManager              // instance bound to the deployed contract
	contractAddress   common.Address                    // address of current address object
	eth               Backend                           // an interface implementation of *EthAPIBackend (in package eth)
	contractSession   *sm.SelectionManagerSession       // an instance of the contract with pre-set call and transact optionss
	TxRecords         map[string]*sm.TxTime             // list of txs' timestamp-records processed by candidate
	TxRecordsMapMutex sync.RWMutex                      // a mutex lock used to avoid race conditions by goroutine writing to the map
	FileName          string                            // name of the JSON file that stores timestamp-records
	CSVlogName        string                            // name of the CSV file that stores selection event reocords
	//backend          *EthAPIBackend            // acts as the client without requiring an RPC/IPC connection (provides useful APIs)
}

// SCG is an accessible global pointer to the smart contract guard
var SCG *SCGuard

// Transaction represents tx selected for runtime verification/analysis
type Transaction struct {
	//txHash         [32]byte 	// hash of the transaction, which is used as its ID
	//isMalicious    bool   	// true if tx is classified as malicious by runtime analysis engines
	txState    string    // can either be 'processing', 'benign', 'maliciaous'
	totalVotes int64     // total number of votes received via pool
	results    int64     // total number of votes classifying tx as malicious received via pool
	timeout    time.Time // if the timeout elapses and tx is not processed yet, it is discarded
}

// SmartContract represents a sc that txs invoking it are processed by SCG's RV
type SmartContract struct {
	selectionIndex  *big.Int                // the index of candidate in selected validator list (this is needed for submitting verification result)
	numOfValidators int64                   // number of validators assigned to the tx
	txs             map[string]*Transaction // list of txs assigned to candidate for runtime validation -- used to store candidate's index in selected list
	timeout         time.Time               // once the timeout elapses a new set of nodes are selected
}

// NewSCGuard creates and returns a smart contract guard instance
func NewSCGuard(candidateAddress common.Address, eth Backend) *SCGuard {
	//func NewSCGuard(candidateAddress common.Address, backend *EthAPIBackend) *SCGuard {
	//var scg SCGuard
	scg := &SCGuard{
		candidateAddress:  candidateAddress,
		eth:               eth,
		contract:          &sm.SelectionManager{},
		contractSession:   &sm.SelectionManagerSession{},
		client:            &ethclient.Client{},
		TxRecordsMapMutex: sync.RWMutex{},
	}
	scg.scs = make(map[common.Address]*SmartContract)
	scg.TxRecords = make(map[string]*sm.TxTime)
	scg.FileName = "/root/go/bin/evaluation/TxRecords_" + candidateAddress.Hex() + ".json"
	scg.CSVlogName = "/root/go/bin/events/TxRecords_" + candidateAddress.Hex() + ".csv"
	titles := []string{time.Now().Format("2006-01-02 15:04:05.000"), "Txhash", "NumOfVals", "Timestamp", "Selected", "Selected_OnChain", "Unix_Timestamp"}
	// initialize CSV log file by writting field titles
	sm.CSVWriter(scg.CSVlogName, titles)
	// set experiment configuration parameters
	parameters := getConfigurationParameters()
	writingRecordsInterval = int(parameters["writingRecordsInterval"].(float64))
	gethLog.Info("SCG: Initialized & Set Candidate Address ", "address", candidateAddress)
	//
	return scg
}

///////////////////////////////Settings & Registration////////////////////////////////

// connectToClient() takes provider URL, connects to specified provider, & sets client object of scg
func (scg *SCGuard) connectToClient(provider string) {
	// Connect to local geth client
	client, err := ethclient.Dial(provider)
	if err != nil {
		gethLog.Warn("SCG: Couldn't connect to provider")
	} else {
		gethLog.Info("SCG: Connected to backend")
	}
	scg.client = client
}

// SetSelectionManagerBackend gets the deployed instance of SelectionManagerContract via eth.backendAPI
func (scg *SCGuard) SetSelectionManagerBackend(initial bool) {
	var ok, noState bool
	var err error
	contrAddr := common.HexToAddress(InitialContractAddress)
	blockNumber := rpc.BlockNumber(0)
	//timeout := 0
	gethLog.Info("SCG: Setting SelectionManager. . .")
	// connect to local IPC provider & set scg.client
	scg.connectToClient("/root/.ethereum/devchain/geth.ipc")
	//if it is initial then contract is found in genesis block (block 0), otherwise extract address from config file
	if initial {
		ok, noState, err = scg.verifyContractAddressBackend(contrAddr, blockNumber)
	} else {
		contrAddr, blockNumber = scg.getContractAddressBackend(configPath)
		gethLog.Info("SCG: Provided blockNumber is", "number", blockNumber.Int64())
		// verify retrieved address, if an error occurs then keep verifying until it works
		ok, noState, err = scg.verifyContractAddressBackend(contrAddr, blockNumber)
		if err != nil || noState { // if error or no state yet, then enter loop of address verification
			gethLog.Info("SCG: Error verifying contract address, re-verifying") //the if needed to log once
			for err != nil || noState {
				time.Sleep(10 * time.Second)
				ok, noState, err = scg.verifyContractAddressBackend(contrAddr, blockNumber)

			}
			gethLog.Info("SCG: Re-verification terminated (with no errors)")

		}
	}

	// create an instance only if the extracted address is a valid one of a deployed SelectionManager contract
	if ok {
		// set the contract address
		scg.contractAddress = contrAddr
		// create an instance of the deployed contract
		selectionContract, err := sm.NewSelectionManager(scg.contractAddress, scg.client)
		if err != nil {
			log.Fatalf("SCG: Couldn't create an instance of deployed contract: %v ", err)
		}
		// set the contract instance
		scg.contract = selectionContract
		// set the session (used to avoid adding bind options when using same account)
		// get transactOpts to create a session using candidate's account
		transactOpts := getTxOpts(scg.candidateAddress)
		// Wrap the contract instance into a session
		scg.contractSession = &sm.SelectionManagerSession{
			Contract: scg.contract,
			CallOpts: bind.CallOpts{
				Pending: false,
				Context: context.Background(),
			},
			TransactOpts: *transactOpts,
		}
		gethLog.Info("SCG: Created instance of deployed contract", "address", scg.contractAddress)
		go scg.WatchSelection()
	} else {
		gethLog.Info("SCG: Retrieved address is not of a deployed contract")
	}
}

// SetCandidateAddress sets the candidate's address to the provided addr
func (scg *SCGuard) SetCandidateAddress(addr common.Address) {
	scg.candidateAddress = addr
}

// GetCandidateAddress sets the candidate's address to the provided addr
func (scg *SCGuard) GetCandidateAddress() common.Address {
	return scg.candidateAddress
}

// AddTxLocally adds an new tx entry invoking given scAddr stored in local scs
// only if scs & a sc entry exist, returns false if it fails to do so
func (scg *SCGuard) AddTxLocally(txhash string, scAddr common.Address) bool {
	// return false if a sc entry with scAddr doesn't exists or scs is not set
	if scg.scs[scAddr] == nil || scg.scs == nil {
		return false
	}
	//return false if tx already exists
	if scg.scs[scAddr].txs[txhash] != nil {
		return false
	}
	if _, found := scg.scs[scAddr].txs[txhash]; !found {
		gethLog.Info("SCG: Received a new Tx, creating a new Tx entry. . .", "hash", txhash, "recipient", scAddr)
		// create & initialize a new tx entry
		scg.scs[scAddr].txs[txhash] = &Transaction{"processing", 0, 0, time.Now()}
	} else {
		gethLog.Info("SCG: Received a Tx with existing entry", "hash", txhash, "recipient", scAddr)
	}

	return true

}

///////////////////////////////Selection Management////////////////////////////////

// WatchSelection watches for SelectedSubset event emitted by contract
func (scg *SCGuard) WatchSelection() {
	// the retrieved TxHash from the SelectedSubset event
	var SCAddr []common.Address

	// Watch for a SelectedSubset event
	watchOpts := &bind.WatchOpts{Context: context.Background(), Start: nil}
	// Setup a channel for results
	channel := make(chan *sm.SelectionManagerSelectedSubset)
	gethLog.Info("SCG: Subscribing for SelectedSubset events. . .")
	// Start a goroutine which watches for new events
	go func(channel chan *sm.SelectionManagerSelectedSubset) {
		_, e := scg.contract.WatchSelectedSubset(watchOpts, channel, SCAddr)
		if e != nil {
			gethLog.Warn("Failed to Watch, Error:", "err", e)
		}
		gethLog.Info("SCG: Watching for SelectedSubset events. . . ")

	}(channel)
	for {
		// Receive events from the channel
		var selectedEvent *sm.SelectionManagerSelectedSubset = <-channel
		// timestamp selection event
		timestamp := time.Now()
		// get data of the logged event
		scAddr := selectedEvent.ScAddr
		numOfVals := selectedEvent.NumOfValidators.String()
		gethLog.Info("SCG: A subset of validators is selected for a new Tx ", "number", selectedEvent.NumOfValidators,
			"address", scAddr)
		selected := false
		selectedOnChain := false
		// check if node is selected as a validator for given sc address & initialize scs if not set
		if scg.checkSelectedList(scAddr) {
			// get index of candidate in list of selected nodes
			selectionIndex := scg.scs[scAddr].selectionIndex
			gethLog.Info("SCG: Selected as a validator for new Tx (contr addr, node addr)", "address", scAddr, "number", selectionIndex, "address", scg.candidateAddress)
			// only apply analysis after querying the contract to ensure that node is selected
			if scg.IsSelected(scAddr, scg.scs[scAddr].selectionIndex) {
				gethLog.Info("SCG: Double checked on-chain (contr addr, node addr)", "address", scAddr, "number", selectionIndex, "address", scg.candidateAddress)
				selectedOnChain = true
			}
			selected = true
		}

		////CSV LOG EVENTS
		go func(scAddr string, numOfVals string, timestamp time.Time, selected bool, selectedOnChain bool) {
			// create a new array of strings, where each slot is a field of the selection event along with a timestamp & a bool fields
			event := []string{scAddr, numOfVals, timestamp.Format("2006-01-02 15:04:05.000000000"), strconv.FormatBool(selected), strconv.FormatBool(selectedOnChain), strconv.FormatInt(timestamp.UnixNano(), 10)}
			// write created array to CSV file
			sm.CSVWriter(scg.CSVlogName, event)
		}(scAddr.Hex(), numOfVals, timestamp, selected, selectedOnChain)
		////

	}
}

// checkSelectedList checks if node is within selected list for a given sc
// it creates an entry for the sc in local scs if it does not exists
func (scg *SCGuard) checkSelectedList(scAddr common.Address) bool {
	var isSelected = false
	// set call options to call a view function
	callOpts := &bind.CallOpts{Context: context.Background(), Pending: false}
	// call the view function to get current validators stored in blockchain state
	validators, _ := scg.contract.GetValidators(callOpts, scAddr)
	// make a scs map if not initialized yet
	if scg.scs == nil {
		scg.scs = make(map[common.Address]*SmartContract)
	}
	if _, found := scg.scs[scAddr]; !found {
		scg.scs[scAddr] = &SmartContract{nil, int64(len(validators)), make(map[string]*Transaction), time.Now()}
	}
	// iterate over set of validators & break if candidate is a validator
	for i, v := range validators {
		// if candidate is validator then append new sc entry to scs map with the retrieved index
		if v == scg.candidateAddress { // _ is the index i
			isSelected = true
			// store sc in scs & set values
			scg.scs[scAddr].selectionIndex = big.NewInt(int64(i))
			break
		}
	}

	return isSelected
}

// TriggerSelection triggers the selection process executed on-chain by
// the SelectionManager contract if sc wasn't already added,
// //returns true if it triggered selection (not necessarily successfully), otherwise returns false
func (scg *SCGuard) TriggerSelection(scAddr common.Address) { //from common.Address) { //bool {
	scg.getNumberOfValidatorsWithinSubset()
	// use contract session to call IsTxExists()
	isExists, err := scg.contractSession.IsSCExists(scAddr)
	if err != nil {
		gethLog.Error("SCG: Error Calling isSCExists()", "err", err)
	}
	if !isExists {
		gethLog.Info("SCG: triggering-selection, Calling Select()", "recipient", scAddr, "address", scg.candidateAddress)

		// use contract session to call SelectFixed()
		tx, err := scg.contractSession.SelectFixed(scAddr, NumOfValidators_within_subset)
		if err != nil {
			gethLog.Warn("SCG: trigger-selection tx failed, SelectFixed() Error", "err", err, "address", scAddr)
			//fmt.Printf("Error is: %v", err)
			//return false
		} else {
			receipt, err := bind.WaitMined(context.Background(), scg.client, tx)
			if err != nil {
				gethLog.Warn("SCG: trigger-selection tx failed, WaitMined() Error", "err", err, "address", scAddr)
				//fmt.Printf("Error is: %v", err)
				//return false
			}

			if receipt.Status != types.ReceiptStatusSuccessful {
				//panic("Call failed")
				gethLog.Warn("SCG: trigger-selection tx mined but failed", "address", scAddr, "gas", receipt.GasUsed, "address", scg.candidateAddress)
				//return false
			} else {
				// Record Time
				//go scg.RecordTime(txhash, "Selected", time.Now(), from)
				gethLog.Info("SCG: trigger-selection Tx succeeded", "address", scAddr, "gas", receipt.GasUsed, "address", scg.candidateAddress)
			}

		}
	} else {
		gethLog.Info("SCG: Tx already exists, Didn't send trigger-selection tx", "address", scAddr, "address", scg.candidateAddress)
	}

	//return false

}

func (scg *SCGuard) checkSelectionIndex(scAddr common.Address) bool {
	// make a scs map if not initialized yet
	if scg.scs == nil {
		scg.scs = make(map[common.Address]*SmartContract)
	}
	if _, found := scg.scs[scAddr]; found {
		if scg.scs[scAddr].selectionIndex != nil {
			return true
		}

	}

	return false

}

////////////////////////////Runtime Verification////////////////////////////////

// getRVTimeoutValue extracts RV timeout value from parameters' JSON file
func (scg *SCGuard) getRVTimeoutValue() {
	// get number of validators/victim-contracts from parameter_configuration file
	// Open our jsonFile
	jsonFile, err := os.Open(parametersFilePath)
	// if os.Open returns an error then handle it
	if err != nil {
		fmt.Println("Error Opening parameter_config.json")
		fmt.Println(err)
	}
	// defer the closing of jsonFile
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var parameters map[string]interface{}
	json.Unmarshal([]byte(byteValue), &parameters)

	//set the timeout value of RV
	RVTimeEmulation = time.Duration(parameters["rv_time"].(float64)) * time.Microsecond //time.Millisecond
}

// getNumberOfValidatorsWithinSubset extracts number of nodes & victims from parameters' JSON file to compute
// number of nodes within each subset (NOTE: it is assumed that nodes are equally distributed)
func (scg *SCGuard) getNumberOfValidatorsWithinSubset() {
	// get number of validators/victim-contracts from parameter_configuration file
	// Open our jsonFile
	jsonFile, err := os.Open(parametersFilePath)
	// if os.Open returns an error then handle it
	if err != nil {
		fmt.Println("Error Opening parameter_config.json")
		fmt.Println(err)
	}
	// defer the closing of jsonFile
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var parameters map[string]interface{}
	json.Unmarshal([]byte(byteValue), &parameters)

	//get number of nodes and validators as int64
	numOfNodes := parameters["numOfNodes"].(int64)
	numOfVictims := parameters["numOfVictims"].(int64)

	// compute number of validators within a subset given number of (candidate) nodes & (victim) contracts
	NumOfValidators_within_subset = big.NewInt(numOfNodes / numOfVictims)

}

func getConfigurationParameters() map[string]interface{} {
	// get number of validators/victim-contracts from parameter_configuration file
	// Open our jsonFile
	jsonFile, err := os.Open(parametersFilePath)
	// if os.Open returns an error then handle it
	if err != nil {
		fmt.Println("Error Opening parameter_config.json")
		fmt.Println(err)
	}
	// defer the closing of jsonFile
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var parameters map[string]interface{}
	json.Unmarshal([]byte(byteValue), &parameters)

	return parameters
}

// ApplyRuntimeVerification performs runtime analysis to determine if given tx
// leads to an exploitation
//func (scg *SCGuard) ApplyRuntimeVerification(txHash [32]byte, scAddr common.Address) {
func (scg *SCGuard) ApplyRuntimeVerification(tx *types.Transaction, sender common.Address, scAddr common.Address) {
	txHash := tx.Hash()
	// apply RV & submit result only if index is not nil
	if scg.checkSelectionIndex(scAddr) {
		// invoke runtime analysis engine
		scg.ExecuteTransactionOffchain(tx)
		// .....
		// get RVTimeEmulation from JSON file
		scg.getRVTimeoutValue()
		time.Sleep(RVTimeEmulation)
		gethLog.Info("SCG: Applied RV to given tx", "hash", common.BytesToHash(txHash[:]).Hex(), "recipient", scAddr)
		// submit result
		go scg.SubmitResult(txHash, scAddr, scg.scs[scAddr].selectionIndex, false, sender, *tx.To())
	} else {
		gethLog.Info("SCG: Skipping Result-Submission", "hash", common.BytesToHash(txHash[:]).Hex(), "recipient", scAddr)
	}

}

///////////////////////////////Result Enforcement////////////////////////////////

// SubmitResult attempts to submit result, if tx got mined but failed,
// it re-attempts only if tx is still unprocessed and selection state hasn't changed
func (scg *SCGuard) SubmitResult(txHash [32]byte, scAddr common.Address, index *big.Int, result bool, sender common.Address, To common.Address) {
	//	numOfTxsSent := 0
	//numOfTxsSent_ptr := &numOfTxsSent
	if scg.submitResult(txHash, scAddr, index, result, sender, To) {
		status := scg.CheckStatus(txHash, scAddr)
		txhash := common.BytesToHash(txHash[:]).Hex()
		//resubmit just to avoid un-updated blockchain state issue
		if status != "benign" && status != "malicious" {
			//if scg.CheckStatus(txHash, scAddr) == "processing" { //&& scg.IsSelected(txHash, index) {
			gethLog.Info("SCG: re-submitting result after prev-mined with failure", "hash", txhash, "address", scg.candidateAddress)
			scg.submitResult(txHash, scAddr, index, result, sender, To)
		} else {
			gethLog.Info("SCG: state-change detected, will not re-submit result", "hash", txhash, "address", scg.candidateAddress)
		}

	}

}

// submitResult allows validators to submit their runtime-verification result by
// performing a contract call to the func SubmitResult() of SelectionManager contract
// (signed by their candidate account) -- result is set to true, if tx is malicious
// Returns true if transaction got mined but failed
func (scg *SCGuard) submitResult(txHash [32]byte, scAddr common.Address, index *big.Int, result bool, sender common.Address, To common.Address) bool {
	txhash := common.BytesToHash(txHash[:]).Hex()
	gethLog.Info("SCG: Submitting Result. . .", "hash", txhash, "recipient", scAddr, "address", scg.candidateAddress)
	// use session to call SubmitResult()
	tx, err := scg.contractSession.SubmitResult(txHash, scAddr, index, result)
	if err != nil {
		gethLog.Warn("SCG: result-submission tx failed, Calling SubmitResult Error", "err", err, "hash", txhash, "address", scg.candidateAddress)
		//gethLog.Warn("SCG: Error", "err", err)
		status := scg.CheckStatus(txHash, scAddr)
		//resubmit
		for status != "benign" && status != "malicious" {
			gethLog.Info("SCG: re-submitting result after calling error", "hash", txhash, "address", scg.candidateAddress)
			tx, err = scg.contractSession.SubmitResult(txHash, scAddr, index, result)
			if err != nil {
				gethLog.Warn("SCG: result-Re-submission tx failed, Calling SubmitResult Error", "err", err, "hash", txhash, "address", scg.candidateAddress)
			}
			if tx != nil {
				break
			}
			status = scg.CheckStatus(txHash, scAddr)

		}

		//fmt.Printf("Error is: %v", err)
	}
	if tx != nil {
		//RECORD FOR EVALUATION
		go scg.RecordTime(txhash, "Processed", time.Now(), 0, sender, To)
		/////
		receipt, err := bind.WaitMined(context.Background(), scg.client, tx)
		if err != nil {
			gethLog.Warn("SCG: result-submission tx failed, Waiting for tx to get mined Error", "err", err, "hash", txhash, "address", scg.candidateAddress)
			//gethLog.Warn("SCG: Error", "err", err)
			//fmt.Printf("Error is: %v", err)
			return false
		}
		if receipt.Status != types.ReceiptStatusSuccessful {
			//panic("Call failed")
			gethLog.Warn("SCG: result-submission tx mined but failed", "hash", txhash, "recipient", scAddr, "gas", receipt.GasUsed, "address", scg.candidateAddress)
			return true

		}
		////RECORD FOR EVALUATION
		// Record Time
		go scg.RecordTime(txhash, "ResultSubmitted", time.Now(), 0, scg.contractAddress, To)
		////
		//gethLog.Info("SCG: Result-submission Tx Succeeded", "hash", receipt.TxHash, "gas", receipt.GasUsed, "address", receipt.ContractAddress)
		gethLog.Info("SCG: Result-submission Tx Succeeded", "hash", txhash, "recipient", scAddr, "gas", receipt.GasUsed, "address", scg.candidateAddress)

		//}
	}

	return false
}

// WatchTransaction watches for ProcessedTransaction event emitted by contract
// returns true if analysis result for given transaction is benign,
// otherwise returns false
func (scg *SCGuard) WatchTransaction(txHash [32]byte) (bool, error) {
	//Solidity: ProcessedTransaction(bytes32 _txHash, bool _result);

	// Watch for a ProcessedTransaction event
	watchOpts := &bind.WatchOpts{Context: context.Background(), Start: nil}
	// Setup a channel for results
	channel := make(chan *sm.SelectionManagerProcessedTransaction)
	//fmt.Printf("Subscribing for SelectedSubset events. . .\n")
	gethLog.Info("SCG: Subscribing for ProcessedTransaction events. . .")
	// Start a goroutine which watches for new events
	go func() {
		_, e := scg.contract.WatchProcessedTransaction(watchOpts, channel)
		if e != nil {
			fmt.Printf("Failed to Watch ProcessedTx Error: %v", e)
		}
		gethLog.Info("SCG: Watching for ProcessedTransaction events. . . ")
		//defer sub.Unsubscribe()

	}()
	gethLog.Info("SCG: Waiting for Runtime-verification Result. . .", "hash", common.BytesToHash(txHash[:]).Hex())
	for { //keep waiting for ProcessedTx event of the given txhash (terminate if timeout elapses)
		select {
		//var processedTxEvent *sm.SelectionManagerProcessedTransaction = <-channel
		// Receive events from the channel
		case processedTxEvent := <-channel:
			gethLog.Info("SCG: Runtime-verification Result submission detected", "hash", common.BytesToHash(processedTxEvent.TxHash[:]).Hex())
			// if desired txhash's event received then log & break
			if processedTxEvent.TxHash == txHash {
				// if result is true, then log a warning that tx is malicious, otherwise log info
				if processedTxEvent.Result {
					gethLog.Warn("SCG: Tx is Malicious ", "address", scg.candidateAddress, "hash", common.BytesToHash(processedTxEvent.TxHash[:]).Hex())
					return false, ErrMaliciousTx
				}
				gethLog.Info("SCG: Tx is Benign", "address", scg.candidateAddress, "hash", common.BytesToHash(processedTxEvent.TxHash[:]).Hex())
				return true, nil

			}
		case <-time.After(180 * time.Second):
			gethLog.Info("SCG: WatchTransaction Timed-out", "address", scg.candidateAddress)
			return false, ErrWatchTimeout
		}
		return false, ErrWatchTimeout
	}
}

// CheckStatus reads the status of the given tx from the contract's state
func (scg *SCGuard) CheckStatus(txhash [32]byte, scAddr common.Address) string {
	// calling a view function, set call options
	callOpts := &bind.CallOpts{Context: context.Background(), Pending: false}
	status, _ := scg.contract.GetStatus(callOpts, txhash, scAddr)
	return status
}

// IsSelected returns true if candidate is selected for given sc at given index of list
func (scg *SCGuard) IsSelected(scAddr common.Address, index *big.Int) bool {
	// calling a view function, set call options
	callOpts := &bind.CallOpts{Context: context.Background(), Pending: false}
	selected, _ := scg.contract.IsSelected(callOpts, scAddr, index)
	return selected
}

/////////////////////////////// Relying On Pool Logic    //////////////////////////////////

// ResultSubmissionInput represents the arguments of result-submission tx
type ResultSubmissionInput struct {
	txHash    common.Hash
	scAddr    common.Address
	index     *big.Int
	malicious bool
}

// IsResultSubmission checks whether given tx is a result-submission tx for the given txhash,
// & keeps track of votes received so far, returns true if all required votes indicating that
// tx has been processed by SCg & hence ready for execution, otherwise returns false
func (scg *SCGuard) IsResultSubmission(tx *types.Transaction, from common.Address) {
	//txinput: 0x|SubmitResult()-signature '9eef8f1a'|32-bytes txhash|big.Int index|bool var malicious|
	gethLog.Info("SCG: Received a SelectionManager tx", "hash", tx.Hash())
	// decode tx input and extract method arguments
	arguments := decodePayload(tx.Data())

	if arguments != nil {
		// discard tx to avoid mining it -- but broadcast it too
		txhash := arguments.txHash
		scAddr := arguments.scAddr
		txHash := common.BytesToHash(txhash[:]).Hex()

		// check if scAddr extracted from input exists
		if _, found := scg.scs[scAddr]; !found {
			gethLog.Warn("SCG: Received a result-submission tx for non-existing SC entry", "hash", txHash, "recipient", scAddr, "address", from)
			isExists, err := scg.contractSession.IsSCExists(scAddr)
			if err != nil {
				gethLog.Error("SCG: Error Calling isSCExists()", "err", err)
			}

			if isExists { //if sc exists on-chain then
				// set & initialize a local entry for sc
				scg.checkSelectedList(scAddr)
			}

		} else {
			gethLog.Info("SCG: Received a result-submission tx", "hash", txHash, "recipient", scAddr, "address", from)
		}

		// if a tx entry of txhash extracted from input doesn't exists, create one
		if _, found := scg.scs[scAddr].txs[txHash]; !found {
			gethLog.Info("SCG: Received a result-submission tx for non-existing Tx entry, creating one. . .", "hash", txHash, "recipient", scAddr, "address", from)
			// create & initialize a tx entry
			scg.scs[scAddr].txs[txHash] = &Transaction{"processing", 0, 0, time.Now()}
		}

		if scg.isValidator(scAddr, from) {
			gethLog.Info("SCG: Sender of tx is a validator ", "hash", txHash, "recipient", scAddr, "address", from)
			// if sender is a validator, then add vote
			if arguments.malicious {
				scg.scs[scAddr].txs[txHash].results++
			}
			// keep track of total number of votes
			scg.scs[scAddr].txs[txHash].totalVotes++

			// check if all votes are received then calculate result & return true
			if scg.scs[scAddr].txs[txHash].totalVotes == scg.scs[scAddr].numOfValidators {
				if scg.scs[scAddr].txs[txHash].results > scg.scs[scAddr].txs[txHash].totalVotes/2 {
					scg.scs[scAddr].txs[txHash].txState = "malicious"
					gethLog.Info("SCG: Tx processed as malicious when relying on pool", "hash", txHash, "recipient", scAddr, "address", from)
				} else {
					scg.scs[scAddr].txs[txHash].txState = "benign"
					gethLog.Info("SCG: Tx processed as benign when relying on pool", "hash", txHash, "recipient", scAddr, "address", from)
				}

				////RECORD FOR EVALUATION
				// Record Time
				go scg.RecordTime(txHash, "PoolExecuted", time.Now(), 0, from, *tx.To())
				/////
			}
		}

	}

}

// GetLocalTxStatus returns the status of given tx if tx exists locally
// otherwise, it return an empty string
func (scg *SCGuard) GetLocalTxStatus(txhash string, scAddr common.Address) string {
	if scg.IsTxExistsLocally(txhash, scAddr) {
		return scg.scs[scAddr].txs[txhash].txState
	}

	return ""
}

// IsTxExistsLocally returns true if given tx invoking SC is already stored locally
func (scg *SCGuard) IsTxExistsLocally(txhash string, scAddr common.Address) bool {
	if scg.scs == nil || scg.scs[scAddr] == nil {
		return false
	}
	if scg.scs[scAddr].txs[txhash] == nil {
		return false
	}

	return true
}

// IsTxTimeElapsed returns true if timeout of given tx has elapsed, otherwise returns false
func (scg *SCGuard) IsTxTimeElapsed(txhash string, scAddr common.Address, timestamp time.Time) bool {
	if scg.IsTxExistsLocally(txhash, scAddr) {
		if timestamp.Sub((scg.scs[scAddr].txs[txhash].timeout)) > txTimeout {
			return true
		}
	}
	return false
}

///////////////////////////////////////////////////////////////////////////////////////

// isValidator checks if given candidate address is within selected list for a given tx
func (scg *SCGuard) isValidator(scAddr common.Address, candidateAddress common.Address) bool {

	var isValidator = false
	// calling a view function, set call options
	callOpts := &bind.CallOpts{Context: context.Background(), Pending: false}
	validators, _ := scg.contract.GetValidators(callOpts, scAddr)

	for _, v := range validators {
		if v == candidateAddress { // _ is the index i
			isValidator = true
			break
		}
	}
	return isValidator

}

func decodePayload(txInput []byte) *ResultSubmissionInput {
	// example of transaction input data
	//txInput := "0x6b9bdf34fe70a7c4572f229edffb38b2be528d58410cb634c810ff5b941b37ccbc13df0d000000000000000000000000e4ee26c6ec935541cf58ab70ef323633f515124b00000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000001"

	// use ResultSubmissionInput struct to unpack arguments
	var data *ResultSubmissionInput

	if len(txInput) < 128 {
		gethLog.Info("SCG: Decoding tx, bytes less than 128")
		return nil
	}

	// get txInput method signature (which is of size: 4 bytes)
	decodedSig := txInput[0:4]

	// load contract ABI
	abi, err := abi.JSON(strings.NewReader(sm.SelectionManagerABI))
	if err != nil {
		gethLog.Error("SCG: Decoding tx Error, abi.Json()", "err", err)
		return nil
	}

	// recover Method from signature and ABI
	method, err := abi.MethodById(decodedSig)
	if err != nil {
		gethLog.Error("SCG: Decoding tx Error, abi.MethodById()", "err", err)
		return nil
	}
	if method == nil {
		gethLog.Warn("SCG: Decoding tx, sc method not found")
		return nil
	} else if method.Name == "submitResult" { // proceed only if input method's name is submitResult
		// initiliaze data
		data = &ResultSubmissionInput{}
		// extract arguments txInput Payload (after method signature, which is of 4 bytes)
		decodedData := txInput[4:]
		// set start & end
		start := 0
		end := start + wordSize
		for i := 0; i < len(decodedData)/wordSize; i++ {
			// extract arguments, where each argument is placed in a # bytes = wordSize
			if start < len(decodedData) && end <= len(decodedData) {
				if i == 0 {
					// extract first argument, which is txhash of size 32 bytes
					txh := decodedData[start:end]
					// get its hex representation
					encodedtxhash := hex.EncodeToString(txh)
					// get ([32]byte) Hash from the hex
					data.txHash = common.HexToHash(encodedtxhash)

				} else if i == 1 {
					// extract second argument, which is scAddr of size 20 bytes
					scBytes := decodedData[start:end]
					// convert bytes to common.Address
					data.scAddr = common.BytesToAddress(scBytes)
				} else if i == 2 {
					// extract third argument, which is a *big.Int of size 32 bytes
					indexBytes := decodedData[start:end]
					// get its hex representation
					encodedIndex := hex.EncodeToString(indexBytes)
					// parse index as a big.Int
					data.index = gethMath.MustParseBig256(encodedIndex)
				} else if i == 3 {
					// extract fourth argument, which is a boolean of size 1 byte (placed at the end)
					resultByte := decodedData[len(decodedData)-2:]
					encodedResult := hex.EncodeToString(resultByte)
					// if encoded Result is 0
					if encodedResult == "0000" {
						data.malicious = false
					} else if encodedResult == "0001" {
						data.malicious = true
					}

				}
			} else {
				gethLog.Warn("SCG: Decoding tx failed, start/end pointers exceeded length of data", "number", big.NewInt(int64(len(decodedData))))
				return nil
			}

			// update start & end pointers/indicies
			start = end
			end = start + wordSize

		}

		if data.malicious {
			gethLog.Info("SCG: Decoding tx, extracted data (malicious)", "hash", data.txHash, "address", data.scAddr, "number", data.index)
		} else {
			gethLog.Info("SCG: Decoding tx, extracted data (benign)", "hash", data.txHash, "address", data.scAddr, "number", data.index)
		}

	} else { //else (if retrieved method isn't submitResult return)
		gethLog.Info("SCG: Decoding tx, method not submitResult")
		return nil
	}

	return data
}

//////////////////////////////////////////////////////////////////////////////////

/******* CONTRACT RELATED FUNCTIONS ********/
// WatchContractAddress watches for changes to config.json file by which the contract address is being shared
// if a change is detected, it sets the SelectionManager
func (scg *SCGuard) WatchContractAddress() {
	scg.SetSelectionManagerBackend(true) // initial Set
	gethLog.Info("SCG: Watching config.json for changes to get SelectionManager deployed-address")
	doneChan := make(chan bool)

	go func(doneChan chan bool) {
		defer func() {
			doneChan <- true
		}()
		for {
			err := watchFile(configPath)
			if err != nil {
				fmt.Println(err)
			}
			scg.SetSelectionManagerBackend(false)

		}

	}(doneChan)

	<-doneChan
}

// verifyContractAddressBackend checks if given address is an address of a deployed contract
func (scg *SCGuard) verifyContractAddressBackend(address common.Address, blockNumber rpc.BlockNumber) (bool, bool, error) {

	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	// if provided address is invalid, return false
	if !re.MatchString(address.Hex()) {
		gethLog.Warn("SCG: Provided address is invalid", "address", address)
		return false, false, nil
	}
	//	gethLog.Info("SCG: Provided blockNumber is", "number", blockNumber.Int64())
	// retrieve bytecode at given address
	state, _, err := scg.StateAndHeaderByNumber(context.Background(), blockNumber)
	if state == nil { // if no state retrieved, return second bool as true
		gethLog.Info("SCG: State doesn't exists yet, waiting for block containing contract")
		return false, true, err
		//log.Fatal(err)
	} else if err != nil {
		gethLog.Warn("SCG: State couldn't be retrieved", "err", err)
		return false, false, err
	}
	bytecode := state.GetCode(address)
	gethLog.Info("SCG: GotCode :)", "address", address)
	// return true if length of bytecode is greater than zero (i.e. there is code hence it's a contract)
	return len(bytecode) > 0, false, err
}

// getContractAddressBackend extracts the block number & address from config.json
func (scg *SCGuard) getContractAddressBackend(filePath string) (common.Address, rpc.BlockNumber) {
	// Get deployed Contract's address
	gethLog.Info("SCG: Extracting SelectionManager Contract Address..")
	// read config.json to get smart-contract's address
	configData, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Print(err)
	}
	//var conf Config
	var conf map[string]string
	// unmarshall config.json data
	err = json.Unmarshal([]byte(configData), &conf)
	if err != nil {
		fmt.Println("error:", err)
	}

	// assign extracted address
	contractAddress := common.HexToAddress(conf[contractName])

	bn, err := strconv.ParseInt(conf["blockNumber"], 10, 64)
	if err != nil {
		log.Fatalf("SCG: Error: blockNumber.UnmarshalJSON failed when getting blockNumber from %s\nErr: %v\n", configPath, err)
	}
	blockNumber := rpc.BlockNumber(bn)

	gethLog.Info("SCG: Extracted Contract Address ", "address", contractAddress)
	return contractAddress, blockNumber
}

//////////////////////////////////////////////////////////////////////////////////

/******** GETTER/FETCHER FUNCTIONS **********/

func (scg *SCGuard) isContractCreated() bool {
	if scg.contract == nil || scg.contractSession == nil {
		gethLog.Warn("SCG: SelectionManager Contract is not initialized yet")
		return false
	} else if *scg.contract == (sm.SelectionManager{}) {
		gethLog.Warn("SCG: SelectionManager Contract is not set yet")
		return false
	}
	return true
}

// GetContractAddress returns the address of the deployed SelectionManager contract
func (scg *SCGuard) GetContractAddress() common.Address {
	return scg.contractAddress
}

// HeaderByNumber fetches the header of the given block number & returns it if found
func (scg *SCGuard) HeaderByNumber(ctx context.Context, number rpc.BlockNumber) (*types.Header, error) {
	if number == rpc.LatestBlockNumber {
		return scg.eth.BlockChain().CurrentBlock().Header(), nil
	}
	return scg.eth.BlockChain().GetHeaderByNumber(uint64(number)), nil
}

// StateAndHeaderByNumber fetches the state & header of the given block number & returns them if found
func (scg *SCGuard) StateAndHeaderByNumber(ctx context.Context, number rpc.BlockNumber) (*state.StateDB, *types.Header, error) {
	header, err := scg.HeaderByNumber(ctx, number)
	if err != nil {
		return nil, nil, err
	}
	if header == nil {
		return nil, nil, errors.New("header not found")
	}
	stateDb, err := scg.eth.BlockChain().StateAt(header.Root)
	return stateDb, header, err
}

//////////////////////////////////////////////////////////////////////////////////

//////////////////////////// Offchain Tx Execution ////////////////////////////////
func (scg *SCGuard) ExecuteTransactionOffchain(tx *types.Transaction) {
	// get parameters from backend
	config := scg.eth.BlockChain().Config()
	bc := scg.eth.BlockChain()
	gp := new(GasPool).AddGas(scg.eth.BlockChain().GasLimit())
	statedb, state_err := scg.eth.BlockChain().State()
	if state_err != nil {
		fmt.Printf("SCG::Offchain::Failed to get stateDB, Err=%v\n", state_err)
		return

	}
	usedGas := new(uint64)
	cfg := *scg.eth.BlockChain().GetVMConfig()
	header := scg.eth.BlockChain().CurrentHeader()

	gethLog.Info("SCG::Offchain::Applying Tx Offchain . . .", "hash", tx.Hash())
	result, txState, err := ApplyTransactionOffchain(config, bc, nil, gp, statedb, header, tx, usedGas, cfg)
	if err != nil {
		fmt.Printf("SCG::Offchain::Failed Applying Tx Offchain Error: %v\n", err)
	}
	//fmt.Printf("SCG::Offchain::After Execution SNAPSHOT_ID=%v\n", scg.eth.BlockChain().State().Snapshot())
	fmt.Printf("SCG::Offchain::RESULT: TX FailedStatus=%v ReultData=%s, UsedGas=%d, GP=%v\n", result.Failed(), string(result.ReturnData), result.UsedGas, gp.Gas())
	fmt.Printf("SCG::Offchain::RESULT: BlockNumber=%v, txState=%s, \n", header.Number, txState)

}

////////////////// Tx TIME RECORDING  //////////////////

//RecordTime records the given timestamp and calculates delays storing records to a Json file
// A read/write mutex is used to avoid races and race conditions raising runtime panic
func (scg *SCGuard) RecordTime(txhash string, field string, timestamp time.Time, duration float64, from common.Address, to common.Address) {
	// get write lock
	scg.TxRecordsMapMutex.Lock()
	switch field {
	case "Triggered": // arrived at pool of nodes

		if _, found := scg.TxRecords[txhash]; found {
			// timestamp Triggered
			scg.TxRecords[txhash].Triggered = timestamp

		} else {
			//Create a new txRecord
			newTx := sm.NewTxTime(from, scg.candidateAddress, to)
			scg.TxRecords[txhash] = newTx
			// timestamp Triggered
			scg.TxRecords[txhash].Triggered = timestamp

		}
		// Add tx to local entries of txs
		scg.AddTxLocally(txhash, to)

	case "Selected": // triggered selection (in case of updating selection)

		// timestamp Selected
		scg.TxRecords[txhash].Selected = timestamp
		//	evaluate selection duration
		scg.TxRecords[txhash].SelectionDuration = timestamp.Sub((scg.TxRecords[txhash].Triggered)).Seconds()

	case "Processed": // when validator processes tx offchain & submits result

		// check if tx is already recorded, otherwise create new one & record timestamp
		if _, found := scg.TxRecords[txhash]; found {
			// timestamp Processed
			scg.TxRecords[txhash].Processed = timestamp

		} else {
			// Create a new txRecord
			newTx := sm.NewTxTime(from, scg.candidateAddress, to)
			scg.TxRecords[txhash] = newTx
			// timestamp Processed
			scg.TxRecords[txhash].Processed = timestamp

		}
		// realease write lock

	case "ResultSubmitted": // when result-submission tx reaches is mined (this was used is old-version)

		// timestamp ResultSubmitted
		scg.TxRecords[txhash].ResultSubmitted = timestamp
		//	evaluate processing duration
		scg.TxRecords[txhash].ProccessingDuration = timestamp.Sub((scg.TxRecords[txhash].Processed)).Seconds()

	case "PoolExecuted": //when nodes receive result submission

		// check if tx is already recorded, otherwise create new one & record timestamp
		if _, found := scg.TxRecords[txhash]; found {
			// check if if it wasn't already recorded, where timestamp is equal to default value below
			if scg.TxRecords[txhash].PoolExecuted.Unix() == -62135596800 {
				// timestamp ResultSubmitted
				scg.TxRecords[txhash].PoolExecuted = timestamp
				//	evaluate processing duration
				scg.TxRecords[txhash].PoolDelay = timestamp.Sub((scg.TxRecords[txhash].Triggered)).Seconds()

			}
		} else {
			// Create a new txRecord
			newTx := sm.NewTxTime(from, scg.candidateAddress, to)
			scg.TxRecords[txhash] = newTx
			// timestamp ResultSubmitted
			scg.TxRecords[txhash].PoolExecuted = timestamp
			//	evaluate processing duration
			scg.TxRecords[txhash].PoolDelay = timestamp.Sub((scg.TxRecords[txhash].Triggered)).Seconds()

		}

		/////

	case "Executed": // when tx is executed onchain, i.e., included in a block (note: this isn't equivalent to block construction time)

		// check if tx is already recorded, otherwise create new one & record timestamp
		if _, found := scg.TxRecords[txhash]; found {
			// check if if it wasn't already recorded, where timestamp is equal to default value below
			if scg.TxRecords[txhash].Executed.Unix() == -62135596800 {
				// timestamp ResultSubmitted
				scg.TxRecords[txhash].Executed = timestamp
				//	evaluate processing duration
				scg.TxRecords[txhash].TotalDelay = timestamp.Sub((scg.TxRecords[txhash].Triggered)).Seconds()
				//  record the EVM execution duration
				scg.TxRecords[txhash].EVMDuration = duration
			}
		} else {
			// Create a new txRecord
			newTx := sm.NewTxTime(from, scg.candidateAddress, to)
			scg.TxRecords[txhash] = newTx
			// timestamp ResultSubmitted
			scg.TxRecords[txhash].Executed = timestamp
			//	evaluate processing duration
			scg.TxRecords[txhash].TotalDelay = timestamp.Sub((scg.TxRecords[txhash].Triggered)).Seconds()
			//  record the EVM execution duration
			scg.TxRecords[txhash].EVMDuration = duration
		}

	}
	// write tx record to file every writingRecordsInterval only
	if len(scg.TxRecords) == writeRecordCounter*writingRecordsInterval {
		// write data to json file
		sm.WriteJSONFile(scg.FileName, scg.TxRecords)
		writeRecordCounter++
		gethLog.Info("Wrote TxRecords to file (number=counter", "number", writeRecordCounter)
	}
	// release write lock
	scg.TxRecordsMapMutex.Unlock()
}

// ExportLatestTxRecords writes the latest copy of TxRecords to predefined json file
// returns true if the records were written and false otherwise
func (scg *SCGuard) ExportLatestTxRecords() {

	scg.TxRecordsMapMutex.Lock()
	if len(scg.TxRecords) > 0 {
		gethLog.Info("No Tx records to write")
	}
	// write data to json file
	wrote := sm.WriteJSONFile(scg.FileName, scg.TxRecords)
	if wrote {
		gethLog.Info("TxRecord is written to Json File")
	} else {
		gethLog.Info("Write attempt of TxRecords failed")
	}
	scg.TxRecordsMapMutex.Unlock()
}
