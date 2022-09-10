package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	reg "registration"
	sm "selectionManager"

	"github.com/ethereum/go-ethereum/common"
)

// This program initializes SRP, aka SCGuard, by registering nodes as candidates
// and assigning each candidate to a victim contract (currently set as NameCollection)
/////

func main() {
	// convert selection manager address to hex
	SMAddr := common.HexToAddress(reg.SMAddress)
	// Get number of validators/victim-contracts from parameter_configuration file
	jsonFile, err := os.Open(reg.ParametersFilePath)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println("**RegisterCandidates::Error Opening parameter_config.json")
		fmt.Println(err)
	}
	// defer the closing of json File
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	// decode json file
	var parameters map[string]interface{}
	json.Unmarshal([]byte(byteValue), &parameters)

	// set the number of validator nodes/accounts
	NumOfValidators := int(parameters["numOfNodes"].(float64))
	// set the number of non-validator accounts
	NumOfAccounts := 10 - NumOfValidators
	// NumOfContracts is the number of contracts other than SelectionManager
	NumOfContracts := int(parameters["numOfVictims"].(float64))

	// connect to node
	client, _ := reg.ConnectToNode()
	// set selection manager instance
	SM, errContract := sm.NewSelectionManager(SMAddr, client)
	if errContract != nil {
		panic(errContract)
	}
	// Register Nodes as Candidates
	reg.RegisterNodes(NumOfAccounts, NumOfValidators, SMAddr, SM, client)
	// Trigger Selection, assigning candidates to smart contracts
	reg.AssignValidatorsToContracts(NumOfAccounts, NumOfValidators, NumOfContracts, SM, client)

}
