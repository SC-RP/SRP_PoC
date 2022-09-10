package main

import (
	"fmt"
	"os"
	"path/filepath"
	sm "selectionManager"
)

var (
	directory = getHomeDir() + "/SRP_Evaluation/"
	fileName  = directory + "TxRecords_0x007CcfFb7916F37F7AEEf05E8096ecFbe55AFc2f.json"
)

func main() {
	fmt.Println("********* Stats from Bootstrap Node ********* ")
	// unmarshal the data
	txTime := sm.ReadJSONFile(fileName)
	// get statistics of extracted data
	if txTime != nil {
		sm.GetStats(txTime)
	} else {
		fmt.Println("Unmarshaled map is empty!!")
	}
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
