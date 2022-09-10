package selectionManager

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"strconv"
	"time"

	gethLog "github.com/ethereum/go-ethereum/log"

	"github.com/ethereum/go-ethereum/common"
)

const (
	///Set of colors for println
	colorReset        = "\033[0m"
	colorRed          = "\033[31m"
	colorGreen        = "\033[32m"
	colorYellow       = "\033[33m"
	colorBrightYellow = "\033[93m"
	colorBlue         = "\033[34m"
	colorBrightBlue   = "\033[94m"
	colorPurple       = "\033[35m"
	colorBrightPurple = "\033[95m"
	colorCyan         = "\033[36m"
	colorBrightCyan   = "\033[96m"
	colorWhite        = "\033[37m"

	titleColor    = "\033[1;32m%s\033[0m"
	subtitleColor = "\033[1;36m%s\033[0m"
	avgColor      = "\033[1;93m%s\033[0m"
)

//JSONTxTime is a json object of TxTime
type JSONTxTime struct {
	Sender              string  `json:"sender"`
	Node                string  `json:"node"`
	Triggered           int64   `json:"triggered"`
	Selected            int64   `json:"selected"`
	Processed           int64   `json:"processed"`
	ResultSubmitted     int64   `json:"resultSubmitted"`
	Executed            int64   `json:"executed"`
	PoolExecuted        int64   `json:"poolExecuted"` // timpestamp of executing tx when From pool & not mining
	PoolDelay           float64 `json:"poolDelay"`    //total delay when From pool & not mining
	TotalDelay          float64 `json:"totalDelay"`   //total delay: executed - submitted
	SelectionDuration   float64 `json:"selection"`    //selection duration: Selected - Triggered
	ProccessingDuration float64 `json:"proccessing"`  //RV processing duration: ResultSubmitted - Processed
	EVMDuration         float64 `json:"evmDuration"`  // time needed for a tx to be executed by the EVM

}

//TxTime is an object used to record execution time of tx's lifecycle
type TxTime struct {
	//Txhash          common.Hash    `json:"id"`              // hash that identifies the tx
	//Txhash          string         `json:"id"`              // hash that identifies the tx
	Sender common.Address `json:"sender"` // address of account that sent tx
	To     common.Address `json:"to"`     // address of account that sent tx
	Node   common.Address `json:"node"`   // address of account assigned to the SCG candidate

	Triggered time.Time `json:"triggered"` // timestamp of triggering selection - SCG: when sending select() tx
	Selected  time.Time `json:"selected"`  // timestamp of receiving mined-ack - SCG: when select() tx gets mined
	//Submitted           time.Time      `json:"submitted"`       // timestamp of submitting tx for mining
	Processed           time.Time `json:"processed"`       // timestamp of submitting RV result - SCG: when sending submitResult() tx
	ResultSubmitted     time.Time `json:"resultSubmitted"` // timestamp of receiving mined-ack - SCG: when submitResult() tx gets mined
	Executed            time.Time `json:"executed"`        // timestamp of executing tx
	PoolExecuted        time.Time `json:"poolExecuted"`    // timpestamp of executing tx when From pool & not mining
	PoolDelay           float64   `json:"poolDelay"`       //total delay when From pool & not mining
	TotalDelay          float64   `json:"totalDelay"`      //total delay: executed - submitted
	SelectionDuration   float64   `json:"selection"`       //selection duration: Selected - Triggered
	ProccessingDuration float64   `json:"proccessing"`     //RV processing duration: ResultSubmitted - Processed
	EVMDuration         float64   `json:"evmDuration"`     // time needed for a tx to be executed by the EVM
}

/*type Unmarshaler interface {
	UnmarshalJSON([]byte) error
}*/

//NewTxTime creates a new record entry of a given tx
func NewTxTime(Sender common.Address, Node common.Address, To common.Address) *TxTime {
	t := &TxTime{
		//Txhash: Txhash,
		Sender: Sender,
		Node:   Node,
		To:     To,
	}
	return t
}

// MarshalJSON converts TxTime struct to another format then marshals it
func (t *TxTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		//	Txhash          string `json:"id"`
		Sender    string `json:"sender"`
		To        string `json:"to"`
		Node      string `json:"node"`
		Triggered int64  `json:"triggered"`
		Selected  int64  `json:"selected"`
		//Submitted           int64   `json:"submitted"`
		Processed           int64   `json:"processed"`
		ResultSubmitted     int64   `json:"resultSubmitted"`
		Executed            int64   `json:"executed"`
		PoolExecuted        int64   `json:"poolExecuted"` // timpestamp of executing tx when From pool & not mining
		PoolDelay           float64 `json:"poolDelay"`    //total delay when From pool & not mining
		TotalDelay          float64 `json:"totalDelay"`   //total delay: executed - submitted
		SelectionDuration   float64 `json:"selection"`    //selection duration: Selected - Triggered
		ProccessingDuration float64 `json:"proccessing"`  //RV processing duration: ResultSubmitted - Processed
		EVMDuration         float64 `json:"evmDuration"`
	}{
		//Txhash:          t.Txhash.String(),
		//Txhash:          t.Txhash,
		Sender:    t.Sender.Hex(),
		To:        t.To.Hex(),
		Node:      t.Node.Hex(),
		Triggered: t.Triggered.UnixNano(),
		Selected:  t.Selected.UnixNano(),
		//Submitted:           t.Submitted.UnixNano(),
		Processed:           t.Processed.UnixNano(),
		ResultSubmitted:     t.ResultSubmitted.UnixNano(),
		Executed:            t.Executed.UnixNano(),
		PoolExecuted:        t.PoolExecuted.UnixNano(),
		PoolDelay:           t.PoolDelay,
		TotalDelay:          t.TotalDelay,
		SelectionDuration:   t.SelectionDuration,
		ProccessingDuration: t.ProccessingDuration,
		EVMDuration:         t.EVMDuration,
	})
}

//type TxTimes map[string]*TxTime

// UnmarshalTxTimeMap is called to unmarshal file into type JSONTxTime
func UnmarshalTxTimeMap(j []byte) (map[string]*JSONTxTime, error) {
	var rawStrings map[string]*JSONTxTime
	// t is the map of TxTime
	//t := make(map[string]*TxTime)
	// read it as a map of TxTime
	err := json.Unmarshal(j, &rawStrings)
	if err != nil {
		fmt.Println("Json Unmarshal error !")
		return nil, err
	}

	return rawStrings, nil
}

//WriteJSONFile writes given data (as a map of TxTimes) to the given file
func WriteJSONFile(filename string, data map[string]*TxTime) bool {
	// Preparing the data to be marshalled and written.
	dataBytes, err := json.MarshalIndent(data, "", "  ")
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

//ReadJSONFile checks given file, extracts all time records & returns them as a map of TxTimes
func ReadJSONFile(filename string) map[string]*JSONTxTime {
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
	//data := make(map[string]*TxTime) //{}
	//json.Unmarshal(file, &data)
	data, err := UnmarshalTxTimeMap(file)
	if err != nil {
		fmt.Println(err)
	}

	return data
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

// Seconds returns the given duration d as a floating point number of seconds
func Seconds(d int64) float64 {
	sec := d / 1000000000
	nsec := d % 1000000000
	return float64(sec) + float64(nsec)/1e9
}

// GetStats extracts stats of txs stored in the map of TxTimes,
// it returns the number of of sucessfully processed txs & thier average total delay,
// the number of un-executed txs, number of txs under processing (i.e. reached chain),
// and the number of txs that did not reach chain
// returns: numOfTxs, numOfSuccTxs, avgDelay, unknownDelays, failedTxs, reachedTxs, unreachedTxs
func GetStats(txTime map[string]*JSONTxTime) (int, int, float64, int, int, int, int) {
	numOfTxs := len(txTime)
	sumOfDelays := 0.0
	avgDelay := 0.0
	sumOfEVMTime := 0.0
	avgEVMTime := 0.0
	numOfSuccTxs := 0
	unknownDelays := 0
	failedTxs := 0
	unreachedTxs := 0
	reachedTxs := 0
	// From pool instead of mining
	poolSuccTxs := 0
	poolSumOfDelays := 0.0
	poolAvgDelay := 0.0
	poolUnknownDelays := 0
	poolFailedTxs := 0
	poolSuccUnexecuted := 0
	poolUnexecutedSumOfDelays := 0.0
	poolUnexecutedAvgDelay := 0.0
	// iterate over the map of tx records
	for _, txRecord := range txTime {
		// if tx is executed, then count tx as successful
		if txRecord.TotalDelay > 0.0 && txRecord.TotalDelay < 1000.0 && txRecord.Triggered != -6795364578871345152 {
			numOfSuccTxs++
			sumOfDelays = sumOfDelays + txRecord.TotalDelay
			sumOfEVMTime = sumOfEVMTime + txRecord.EVMDuration
			//fmt.Printf("tx: %s   - successfully processed\n", txhash)
		} else if txRecord.TotalDelay > 0.0 && txRecord.Triggered == -6795364578871345152 {
			unknownDelays++
		} else { // otherwise, count tx as failed & determine whether it reached chain or not
			failedTxs++
			// if selected & processed timestamps are not recorded, then tx didn't reach chain
			if txRecord.SelectionDuration == 0.0 && txRecord.ProccessingDuration == 0.0 {
				unreachedTxs++
				//	fmt.Printf("tx: %s   - did not reach chain\n", txhash)
			} else { // if any of the 2 timestamps are recorded then tx did reach chain
				reachedTxs++
				//	fmt.Printf("tx: %s   - did reach chain\n", txhash)
			}
		}

		// Compute values when From pool
		if txRecord.PoolDelay > 0.0 && txRecord.PoolDelay < 1000.0 && txRecord.Triggered != -6795364578871345152 {
			if txRecord.Executed != -6795364578871345152 {
				poolSuccTxs++
				poolSumOfDelays = poolSumOfDelays + txRecord.PoolDelay
			} else {
				poolSuccUnexecuted++
				poolUnexecutedSumOfDelays = poolUnexecutedSumOfDelays + txRecord.PoolDelay
			}

		} else if txRecord.TotalDelay > 0.0 && txRecord.Triggered == -6795364578871345152 {
			poolUnknownDelays++
		} else {
			poolFailedTxs++
		}

	}
	// compute average delay
	if numOfSuccTxs > 0 {
		avgDelay = sumOfDelays / float64(numOfSuccTxs)
		avgEVMTime = sumOfEVMTime / float64(numOfSuccTxs)
	}
	// compute average delay From pool (have been processed onchain)
	if poolSuccTxs > 0 {
		poolAvgDelay = poolSumOfDelays / float64(poolSuccTxs)
	}
	// compute average delay From pool (have not been processed onchain)
	if poolSuccUnexecuted > 0 {
		poolUnexecutedAvgDelay = poolUnexecutedSumOfDelays / float64(poolSuccUnexecuted)
	}

	fmt.Printf(titleColor, "\n####   "+strconv.Itoa(numOfTxs)+" txs were sent   ####")

	fmt.Printf(subtitleColor, "\n--- From Mined Blocks ---")
	strAvgDelay := fmt.Sprintf("%.3f", avgDelay)
	strEVMTime := fmt.Sprintf("%.3f", avgEVMTime*1000.0)
	fmt.Printf("\n * ")
	fmt.Printf(avgColor, strconv.Itoa(numOfSuccTxs))
	fmt.Printf(" successfully processed with ")
	fmt.Printf(avgColor, strAvgDelay+"s (EVM: "+strEVMTime+"ms)")
	fmt.Printf(" average delay,")
	fmt.Printf("\n * %d with unknown delay", unknownDelays)
	fmt.Printf("\n * %d unprocessed, %d reached chain & %d didn't reach\n", failedTxs, reachedTxs, unreachedTxs)
	strSuccRate := fmt.Sprintf("%.2f", float64(numOfSuccTxs)/float64(numOfTxs)*100.0)
	fmt.Printf("-----------------------\nOn-chain Success-rate = ")
	fmt.Printf(titleColor, strSuccRate)
	fmt.Printf(" \n-----------------------\n")

	fmt.Printf(subtitleColor, "--- From Pool ---")
	strpoolAvgDelay := fmt.Sprintf("%.3f", poolAvgDelay)
	fmt.Printf("\n * ")
	fmt.Printf(avgColor, strconv.Itoa(poolSuccTxs))
	fmt.Printf(" successfully processed with (onchain also) ")
	fmt.Printf(avgColor, strpoolAvgDelay+"s")
	fmt.Printf(" average delay,")
	fmt.Printf("\n * %d successfully processed (but not onchain) with %.3f average delay", poolSuccUnexecuted, poolUnexecutedAvgDelay)
	fmt.Printf("\n * %d with unknown delay", poolUnknownDelays)
	fmt.Printf("\n * %d unprocessed \n", poolFailedTxs)

	strPoolSuccRate := fmt.Sprintf("%.2f", float64((poolSuccTxs+poolSuccUnexecuted))/float64(numOfTxs)*100.0)
	fmt.Printf("-----------------------\nPool Success-rate = ")
	fmt.Printf(titleColor, strPoolSuccRate)
	fmt.Printf(" \n-----------------------\n")
	return numOfTxs, numOfSuccTxs, avgDelay, unknownDelays, failedTxs, reachedTxs, unreachedTxs

}

// GetCommunicationDelay computes the average communication delay of a given record,
// the avg delay is computed by getting the avg poolDelay excluding those processed
// by the node itself (i.e. those with "resultSubmitted": -6795364578871345152 ),
//  it returns both the average delay and an array of the poolDelays included in the computation
func GetCommunicationDelay(txTime map[string]*JSONTxTime) (float64, []float64) {
	numOfTxs := len(txTime)
	poolSuccTxs := 0
	poolSumOfDelays := 0.0
	poolAvgDelay := 0.0
	var countedDelays []float64

	// get sum of poolDelays of txs that were not validated by the candidate & were successfully executed on-chain
	for _, txRecord := range txTime {
		if txRecord.ResultSubmitted == -6795364578871345152 && txRecord.PoolDelay > 0.0 && txRecord.PoolDelay < 1000.0 {
			poolSuccTxs++
			poolSumOfDelays = poolSumOfDelays + txRecord.PoolDelay
			countedDelays = append(countedDelays, txRecord.PoolDelay)
		}
	}

	// get average of pool delay
	// compute average delay From pool (have been processed onchain)
	if poolSuccTxs > 0 {
		poolAvgDelay = poolSumOfDelays / float64(poolSuccTxs)
	}
	fmt.Printf(subtitleColor, "--- Communication Delay ---")
	strAvgDelay := fmt.Sprintf("%.5f", poolAvgDelay*1000.0)
	fmt.Printf("\n * ")
	fmt.Printf(avgColor, strconv.Itoa(poolSuccTxs))
	fmt.Printf(" counted for computing T_comm, which is ")
	fmt.Printf(avgColor, strAvgDelay+"ms\n")
	// get the percentage of counted txs
	strPoolSuccRate := fmt.Sprintf("%.2f", float64((poolSuccTxs))/float64(numOfTxs)*100.0)
	fmt.Printf("-----------------------\n Percentage of Counted = ")
	fmt.Printf(titleColor, strPoolSuccRate)
	fmt.Printf(" \n-----------------------\n")

	return poolAvgDelay, countedDelays

}

///////////////SELECTION LOGS////////////////

// SelectionEventLog stores the info of a selection event to be logged in ajson file
type SelectionEventLog struct {
	Txhash    string         `json:"txhash"`            // hash that identifies the tx
	NumOfVals *big.Int       `json:"numOfVals"`         // number of validators that were selected
	Timestamp time.Time      `json:"timestamp"`         // the time at which the event was received
	Selected  common.Address `json:"selectedCandidate"` // indicate whether the node was selected or not, (default is zero means wasn't selected)
}

// MarshalJSON converts SelectionEventLog struct to another format then marshals it
func (sl *SelectionEventLog) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Txhash    string `json:"txhash"`
		NumOfVals string `json:"numOfVals"`
		Timestamp int64  `json:"timestamp"`
		Selected  string `json:"selectedCandidate"`
	}{
		Txhash:    sl.Txhash,
		NumOfVals: sl.NumOfVals.String(),
		Timestamp: sl.Timestamp.UnixNano(),
		Selected:  sl.Selected.Hex(),
	})
}

//WriteJSONEvent writes given data (as a map of TxTimes) to the given file
func WriteJSONEvent(filename string, data map[string]*SelectionEventLog) {
	// Preparing the data to be marshalled and written.
	dataBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Println(err)
	}
	err = ioutil.WriteFile(filename, dataBytes, 0644)
	if err != nil {
		log.Println(err)
	}
}

// CSVWriter creates a CSV file, named fname,
// & appends a slice of strings (i.e.selection events),
// each stored in a new row
func CSVWriter(fname string, line []string) {
	f, err := os.OpenFile(fname, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		gethLog.Warn("SCG: CSV event-logging Failed", "err", err)
		return
	}
	w := csv.NewWriter(f)
	w.Write(line)
	w.Flush()
}

////

// UnmarshalJSON will be called when json.Unmarshal() is called
// instead of the standard unmarshal methods in the json package.
// It iterates over the keys and manually Unmarshal the data
// into TxTime struct depending on the type
func UnmarshalJSON(j []byte) (map[string]*TxTime, error) {
	var rawStrings map[string]map[string]string
	// t is the map of TxTime
	t := make(map[string]*TxTime)
	// read it as a map of TxTime
	err := json.Unmarshal(j, &rawStrings)
	if err != nil {
		fmt.Println("Json Unmarshal error !")
		return nil, err
	}

	for txhash, txtime := range rawStrings {
		for k, v := range txtime {
			if k == "sender" {
				t[txhash].Sender = common.HexToAddress(v)
			}
			if k == "triggered" {
				tt, err := strconv.ParseInt(v, 10, 64)
				if err != nil {
					fmt.Println("Parsing Unix time error !")
					return nil, err
				}
				t[txhash].Triggered = time.Unix(tt, 0)
			}
			if k == "selected" {
				tt, err := strconv.ParseInt(v, 10, 64)
				if err != nil {
					fmt.Println("Parsing Unix time error !")
					return nil, err
				}
				t[txhash].Selected = time.Unix(tt, 0)
			}
			if k == "processed" {
				tt, err := strconv.ParseInt(v, 10, 64)
				if err != nil {
					fmt.Println("Parsing Unix time error !")
					return nil, err
				}
				t[txhash].Processed = time.Unix(tt, 0)
			}
			if k == "resultSubmitted" {
				tt, err := strconv.ParseInt(v, 10, 64)
				if err != nil {
					fmt.Println("Parsing Unix time error !")
					return nil, err
				}
				t[txhash].ResultSubmitted = time.Unix(tt, 0)
			}
			if k == "executed" {
				tt, err := strconv.ParseInt(v, 10, 64)
				if err != nil {
					fmt.Println("Parsing Unix time error !")
					return nil, err
				}
				t[txhash].Executed = time.Unix(tt, 0)
			}
			if k == "poolExecuted" {
				tt, err := strconv.ParseInt(v, 10, 64)
				if err != nil {
					fmt.Println("Parsing Unix time error !")
					return nil, err
				}
				t[txhash].PoolExecuted = time.Unix(tt, 0)
			}
			if k == "poolDelay" {
				t[txhash].PoolDelay, err = strconv.ParseFloat(v, 32)
				if err != nil {
					fmt.Println("Parsing float error !")
					return nil, err
				}
			}
			if k == "totalDelay" {
				t[txhash].TotalDelay, err = strconv.ParseFloat(v, 32)
				if err != nil {
					fmt.Println("Parsing float error !")
					return nil, err
				}
			}
			if k == "selectionDuration" {
				t[txhash].SelectionDuration, err = strconv.ParseFloat(v, 32)
				if err != nil {
					fmt.Println("Parsing float error !")
					return nil, err
				}
			}
			if k == "proccessingDuration" {
				t[txhash].ProccessingDuration, err = strconv.ParseFloat(v, 32)
				if err != nil {
					fmt.Println("Parsing float error !")
					return nil, err
				}
			}

		}

	}

	return t, nil
}
