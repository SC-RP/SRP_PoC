package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
)

func main() {
	var MainDir = getHomeDir()
	//SMAddress is the address of SelectionManager contract stored in genesis block
	SMAddress := "0xDB571079aF66EDbB1a56d22809584d39C20001D9"
	//GenesisPath is the path of the genesis file
	GenesisPath := MainDir + "/files/genesis.json"
	//OutFile Json file's path in which addresses of victim contracts are stored
	OutFile := MainDir + "/files/victimAddresses.json"
	// NCAddrs is the list of extracted victim addresses from genesis file
	NCAddrs := GetSCsAddresses(GenesisPath, common.HexToAddress(SMAddress))
	// write extracted addresses to json file
	WriteJSONFile(OutFile, NCAddrs)

}

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

// GetSCsAddresses returns an array of all hard-coded addresses of SCs in
// the given genesis file, excluding the SelectionManager contract
func GetSCsAddresses(genesisPath string, SelectionManagerAddr common.Address) []string {
	var SCsAddrs []string
	//// Read genesis
	// create a genesis struct
	genesis := new(core.Genesis)
	// open genesis file
	genesisFile, err := os.Open(genesisPath)
	if err != nil {
		fmt.Printf("Failed to read genesis file: %v", err)
		panic(err)
	}
	defer genesisFile.Close()
	// decode genesis file
	if err := json.NewDecoder(genesisFile).Decode(genesis); err != nil {
		fmt.Printf("invalid genesis file: %v\n", err)
		panic(err)
	}

	//// Extract all hard-coded SCs' addresses (i.e. accounts with code)
	// iterate over all accounts in genesis
	for accAddr, acc := range genesis.Alloc {
		// if account contains code then store accAddr in list of SCs' addresses
		if acc.Code != nil && len(acc.Code) > 0 && accAddr != SelectionManagerAddr {
			SCsAddrs = append(SCsAddrs, accAddr.Hex())

		}
	}

	fmt.Printf("\nExtracted %d smart-contracts' addresses\n", len(SCsAddrs))

	return SCsAddrs

}
