// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package core

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/consensus/misc"

	//scguard "github.com/ethereum/go-ethereum/core/scguard"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
)

var evmStartTime time.Time

// StateProcessor is a basic Processor, which takes care of transitioning
// state from one point to another.
//
// StateProcessor implements Processor.
type StateProcessor struct {
	config *params.ChainConfig // Chain configuration options
	bc     *BlockChain         // Canonical block chain
	engine consensus.Engine    // Consensus engine used for block rewards
}

// NewStateProcessor initialises a new StateProcessor.
func NewStateProcessor(config *params.ChainConfig, bc *BlockChain, engine consensus.Engine) *StateProcessor {
	return &StateProcessor{
		config: config,
		bc:     bc,
		engine: engine,
	}
}

// Process processes the state changes according to the Ethereum rules by running
// the transaction messages using the statedb and applying any rewards to both
// the processor (coinbase) and any included uncles.
//
// Process returns the receipts and logs accumulated during the process and
// returns the amount of gas that was used in the process. If any of the
// transactions failed to execute due to insufficient gas it will return an error.
func (p *StateProcessor) Process(block *types.Block, statedb *state.StateDB, cfg vm.Config) (types.Receipts, []*types.Log, uint64, error) {
	var (
		receipts types.Receipts
		usedGas  = new(uint64)
		header   = block.Header()
		allLogs  []*types.Log
		gp       = new(GasPool).AddGas(block.GasLimit())
	)
	// Mutate the block and state according to any hard-fork specs
	if p.config.DAOForkSupport && p.config.DAOForkBlock != nil && p.config.DAOForkBlock.Cmp(block.Number()) == 0 {
		misc.ApplyDAOHardFork(statedb)
	}
	// Iterate over and process the individual transactions
	for i, tx := range block.Transactions() {
		statedb.Prepare(tx.Hash(), block.Hash(), i)
		log.Info("PROCESSOR: Applying Tx", "hash", tx.Hash())
		receipt, err := ApplyTransaction(p.config, p.bc, nil, gp, statedb, header, tx, usedGas, cfg)
		// //--------------------------------------------------------------> added SCGuard
		if err != nil { //&& err != ErrUnverifiedTx {
			return nil, nil, 0, err
		}
		// // //--------------------------------------------------------------> ended SCGuard
		receipts = append(receipts, receipt)
		allLogs = append(allLogs, receipt.Logs...)
	}
	// Finalize the block, applying any consensus engine specific extras (e.g. block rewards)
	p.engine.Finalize(p.bc, header, statedb, block.Transactions(), block.Uncles())

	return receipts, allLogs, *usedGas, nil
}

// ApplyTransactionAsWorker attempts to apply a transaction to the given state database
// and uses the input parameters for its environment. It returns the receipt
// for the transaction, gas used and an error if the transaction failed,
// indicating the block was invalid. (It's modified to apply SCGuard)
func ApplyTransactionAsWorker(config *params.ChainConfig, bc ChainContext, author *common.Address, gp *GasPool, statedb *state.StateDB, header *types.Header, tx *types.Transaction, usedGas *uint64, cfg vm.Config) (*types.Receipt, error) {
	evmStartTime = time.Now()
	msg, err := tx.AsMessage(types.MakeSigner(config, header.Number))
	if err != nil {
		return nil, err
	}
	// Create a new context to be used in the EVM environment
	context := NewEVMContext(msg, header, bc, author)
	// Create a new environment which holds all relevant information
	// about the transaction and calling mechanisms.
	vmenv := vm.NewEVM(context, statedb, config, cfg)

	// //--------------------------------------------------------------> added SCGuard
	if SCG.isContractCreated() { //only apply SCGuard logic, if contract has been initialized
		// checking if recipient is a contract
		to := msg.To()
		if to != nil && SCG.GetContractAddress() != *to { //if it is not a contract creation tx & it is not invoking SelectionManager
			codeAtTo := statedb.GetCode(*to)
			if codeAtTo != nil { // if no errors resulted from fetching code at the recipient address
				if len(codeAtTo) > 0 { //if recipient is a contract
					log.Info("SCG, wr: Verifying Tx invoking a contract", "address", *msg.To(), "hash", tx.Hash())
					//status := SCG.CheckStatus(tx.Hash(), *to)
					status := SCG.GetLocalTxStatus(tx.Hash().Hex(), *to)
					if status == "benign" { //apply tx only if its status is benign
						//if execute { // apply tx only if it is benign
						log.Info("SCG, wr: Applying Tx invoking a contract", "address", *msg.To(), "hash", tx.Hash())
						// //--------------------------------------------------------------> partially-ended SCGuard

						// Apply the transaction to the current state (included in the env)
						result, err := ApplyMessage(vmenv, msg, gp)
						if err != nil {
							return nil, err
						}
						fmt.Printf("Worker Onchain: RESULT: TX:%v (sent to:%v) FailedStatus=%v, UsedGas=%d, GP=%v\n", tx.Hash().Hex(), tx.To().Hash().TerminalString(), result.Failed(), result.UsedGas, gp.Gas())
						// Update the state with pending changes
						var root []byte
						if config.IsByzantium(header.Number) {
							statedb.Finalise(true)
						} else {
							root = statedb.IntermediateRoot(config.IsEIP158(header.Number)).Bytes()
						}
						*usedGas += result.UsedGas

						// Create a new receipt for the transaction, storing the intermediate root and gas used by the tx
						// based on the eip phase, we're passing whether the root touch-delete accounts.
						receipt := types.NewReceipt(root, result.Failed(), *usedGas)
						receipt.TxHash = tx.Hash()
						receipt.GasUsed = result.UsedGas
						// if the transaction created a contract, store the creation address in the receipt.
						if msg.To() == nil {
							receipt.ContractAddress = crypto.CreateAddress(vmenv.Context.Origin, tx.Nonce())
						}
						// Set the receipt logs and create a bloom for filtering
						receipt.Logs = statedb.GetLogs(tx.Hash())
						receipt.Bloom = types.CreateBloom(types.Receipts{receipt})
						receipt.BlockHash = statedb.BlockHash()
						receipt.BlockNumber = header.Number
						receipt.TransactionIndex = uint(statedb.TxIndex())

						////RECORD FOR EVALUATION
						txhash := tx.Hash().Hex()
						evmDuration := time.Now().Sub(evmStartTime).Seconds()
						// Record Time
						go SCG.RecordTime(txhash, "Executed", time.Now(), evmDuration, msg.From(), *to)
						/////
						return receipt, err
					} else if status == "malicious" {
						log.Info("SCG, wr: Rejected Malicious Tx invoking a contract", "address", *to, "hash", tx.Hash())
						return nil, ErrMaliciousTx
					} else if status == "processing" { // if status is processing, then it has reached chain
						log.Info("SCG, wr: Tx reached chain but not processed yet", "address", *to, "hash", tx.Hash())
						return nil, ErrUnverifiedTx
					} else { // if status is none of the above, then it is not processed yet
						log.Info("SCG, wr: Tx have not reached chain yet", "address", *to, "hash", tx.Hash())
						return nil, ErrUnverifiedTx
					}
				}
			}
		}
	} // //--------------------------------------------------------------> ended SCGuard

	// otherwise apply transaction normally

	// Apply the transaction to the current state (included in the env)
	result, err := ApplyMessage(vmenv, msg, gp)
	if err != nil {
		return nil, err
	}
	// Update the state with pending changes
	var root []byte
	if config.IsByzantium(header.Number) {
		statedb.Finalise(true)
	} else {
		root = statedb.IntermediateRoot(config.IsEIP158(header.Number)).Bytes()
	}
	*usedGas += result.UsedGas

	// Create a new receipt for the transaction, storing the intermediate root and gas used by the tx
	// based on the eip phase, we're passing whether the root touch-delete accounts.
	receipt := types.NewReceipt(root, result.Failed(), *usedGas)
	receipt.TxHash = tx.Hash()
	receipt.GasUsed = result.UsedGas
	// if the transaction created a contract, store the creation address in the receipt.
	if msg.To() == nil {
		receipt.ContractAddress = crypto.CreateAddress(vmenv.Context.Origin, tx.Nonce())
		// SCG trigger selection for newly created contracts
		go SCG.TriggerSelection(receipt.ContractAddress)
	}
	// Set the receipt logs and create a bloom for filtering
	receipt.Logs = statedb.GetLogs(tx.Hash())
	receipt.Bloom = types.CreateBloom(types.Receipts{receipt})
	receipt.BlockHash = statedb.BlockHash()
	receipt.BlockNumber = header.Number
	receipt.TransactionIndex = uint(statedb.TxIndex())

	return receipt, err

}

// ApplyTransaction attempts to apply a transaction to the given state database
// and uses the input parameters for its environment. It returns the receipt
// for the transaction, gas used and an error if the transaction failed,
// indicating the block was invalid.
func ApplyTransaction(config *params.ChainConfig, bc ChainContext, author *common.Address, gp *GasPool, statedb *state.StateDB, header *types.Header, tx *types.Transaction, usedGas *uint64, cfg vm.Config) (*types.Receipt, error) {
	evmStartTime = time.Now()
	msg, err := tx.AsMessage(types.MakeSigner(config, header.Number))
	if err != nil {
		return nil, err
	}
	// Create a new context to be used in the EVM environment
	context := NewEVMContext(msg, header, bc, author)
	// Create a new environment which holds all relevant information
	// about the transaction and calling mechanisms.
	vmenv := vm.NewEVM(context, statedb, config, cfg)

	// Apply the transaction to the current state (included in the env)
	result, err := ApplyMessage(vmenv, msg, gp)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Processor Onchain: RESULT: TX:%v (sent to:%v) FailedStatus=%v, UsedGas=%d, GP=%v\n", tx.Hash().Hex(), tx.To().Hash().TerminalString(), result.Failed(), result.UsedGas, gp.Gas())
	//fmt.Printf("Processor Onchain: RESULT: TX:%v (sent to %v) FailedStatus=%v ResultData=%v, UsedGas=%d, GP=%v\n", tx.Hash().Hex(), tx.To().Hash().Hex(), result.Failed(), result.ReturnData, result.UsedGas, gp.Gas())
	// Update the state with pending changes
	var root []byte
	if config.IsByzantium(header.Number) {
		statedb.Finalise(true)
	} else {
		root = statedb.IntermediateRoot(config.IsEIP158(header.Number)).Bytes()
	}
	*usedGas += result.UsedGas

	// Create a new receipt for the transaction, storing the intermediate root and gas used by the tx
	// based on the eip phase, we're passing whether the root touch-delete accounts.
	receipt := types.NewReceipt(root, result.Failed(), *usedGas)
	receipt.TxHash = tx.Hash()
	receipt.GasUsed = result.UsedGas
	// if the transaction created a contract, store the creation address in the receipt.
	if msg.To() == nil {
		receipt.ContractAddress = crypto.CreateAddress(vmenv.Context.Origin, tx.Nonce())
	}
	// Set the receipt logs and create a bloom for filtering
	receipt.Logs = statedb.GetLogs(tx.Hash())
	receipt.Bloom = types.CreateBloom(types.Receipts{receipt})
	receipt.BlockHash = statedb.BlockHash()
	receipt.BlockNumber = header.Number
	receipt.TransactionIndex = uint(statedb.TxIndex())

	if SCG.isContractCreated() { //only apply SCGuard logic, if contract has been initialized
		// checking if recipient is a contract -- statedb.GetCode(msg.To())
		to := msg.To()
		if to != nil && SCG.GetContractAddress() != *to { //if it is not a contract creation tx & it is not invoking SelectionManager
			codeAtTo := statedb.GetCode(*to)
			if codeAtTo != nil { // if no errors resulted from fetching code at the recipient address
				if len(codeAtTo) > 0 { //if recipient is a contract
					////RECORD FOR EVALUATION
					evmDuration := time.Now().Sub(evmStartTime).Seconds()
					// Record Time
					go SCG.RecordTime(tx.Hash().Hex(), "Executed", time.Now(), evmDuration, msg.From(), *to)
					/////
				}
			}
		}
	}

	return receipt, err

}

// ApplyTransactionOffchain attempts to apply a transaction offchain (with RV) to the given state database
// based on the input parameters for its environment. It returns the execution result of the tx,
// an error if the transaction failed, and the result of applying RV (invalid, malicious, or benign tx)
func ApplyTransactionOffchain(config *params.ChainConfig, bc ChainContext, author *common.Address, gp *GasPool, statedb *state.StateDB, header *types.Header, tx *types.Transaction, usedGas *uint64, cfg vm.Config) (*ExecutionResult, string, error) {
	txState := "invalid"
	log.Info("OffChain: Applying Tx Offchain . . .", "hash", tx.Hash())
	snapShotID := statedb.Snapshot() // Snapshot returns an identifier for the current revision of the state.
	// we can then use the saved id to RevertToSnapshot(revid int) if needed

	fmt.Printf("Offchain: Prior Execution SNAPSHOT_ID=%v and GP=%v\n", snapShotID, gp.Gas())
	log.Info("OffChain: 1) Converting Tx as message . . .", "hash", tx.Hash())
	msg, err := tx.AsMessage(types.MakeSigner(config, header.Number))
	if err != nil {
		fmt.Printf("Offchain: Failed 1)Converting Error: %v\n", err)
		return nil, "invalid", err

	}

	// Create a new context to be used in the EVM environment
	context := NewEVMContext(msg, header, bc, author)
	log.Info("OffChain: 2) Created New EVM Context for tx", "hash", tx.Hash())

	// Create a new environment which holds all relevant information
	// about the transaction and calling mechanisms.
	vmenv := vm.NewEVM(context, statedb, config, cfg)
	log.Info("OffChain: 3) Created New EVM for tx", "hash", tx.Hash())

	log.Info("Offchain: 4) Applying Tx as Message . . . ", "hash", tx.Hash())
	// Apply the transaction to the current state (included in the env)
	result, err := ApplyMessage(vmenv, msg, gp)
	if err != nil {
		fmt.Printf("Offchain: Failed 4)Applying Message Error: %v\n", err)
		return nil, "invalid", err
	}
	log.Info("Offchain: Applied Tx as Message", "hash", tx.Hash())
	snapShotID = statedb.Snapshot()
	fmt.Printf("Offchain: After Execution SNAPSHOT_ID=%v\n", snapShotID)
	fmt.Printf("Offchain: 5) RESULT: TX FailedStatus=%v ResultData=%v, UsedGas=%d, GP=%v\n", result.Failed(), string(result.ReturnData), result.UsedGas, gp.Gas())
	fmt.Printf("Offchain: 5) RESULT: BlockNumber=%v, author=%v, \n", header.Number, author)

	txState = "benign"
	return result, txState, err

}
