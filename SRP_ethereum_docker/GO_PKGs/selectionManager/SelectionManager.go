// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package selectionManager

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// SelectionManagerABI is the input ABI used to generate the binding from.
const SelectionManagerABI = "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"_txHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"_scAddr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_numOfValidators\",\"type\":\"uint256\"}],\"name\":\"AddedTx\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"_txHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"_result\",\"type\":\"bool\"}],\"name\":\"ProcessedTransaction\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_numOfCandidates\",\"type\":\"uint256\"}],\"name\":\"RegisteredNodes\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_scAddr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_numOfValidators\",\"type\":\"uint256\"}],\"name\":\"SelectedSubset\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"getCandidates\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNumOfSCs\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNumOfTxs\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"txHash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"scAddr\",\"type\":\"address\"}],\"name\":\"getStatus\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"scAddr\",\"type\":\"address\"}],\"name\":\"getValidators\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"scAddr\",\"type\":\"address\"}],\"name\":\"isSCExists\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"scAddr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"isSelected\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"txHash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"scAddr\",\"type\":\"address\"}],\"name\":\"isTxExists\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"register\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"scAddr\",\"type\":\"address\"}],\"name\":\"select\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"scAddr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"numOfVals\",\"type\":\"uint256\"}],\"name\":\"selectFixed\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"txHash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"scAddr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"malicious\",\"type\":\"bool\"}],\"name\":\"submitResult\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// SelectionManagerBin is the compiled bytecode used for deploying new contracts.
var SelectionManagerBin = "0x608060405234801561001057600080fd5b5060006002819055506000600381905550611c72806100306000396000f3fe608060405234801561001057600080fd5b50600436106100b45760003560e01c80636b9bdf34116100715780636b9bdf34146102505780638d1ead58146102b4578063a40d5bb3146102d2578063db5f9f3614610320578063e160e120146103e7578063ff8744a614610443576100b4565b806306a49fce146100b95780631aa3a00814610118578063337299eb146101225780634b4dd27d146101885780634f49f01c146101ee5780636110fb7514610232575b600080fd5b6100c16104dc565b6040518080602001828103825283818151815260200191508051906020019060200280838360005b838110156101045780820151818401526020810190506100e9565b505050509050019250505060405180910390f35b6101206105c8565b005b61016e6004803603604081101561013857600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803590602001909291905050506107bc565b604051808215151515815260200191505060405180910390f35b6101d46004803603604081101561019e57600080fd5b8101908080359060200190929190803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610868565b604051808215151515815260200191505060405180910390f35b6102306004803603602081101561020457600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291905050506108e7565b005b61023a610d6a565b6040518082815260200191505060405180910390f35b6102b26004803603608081101561026657600080fd5b8101908080359060200190929190803573ffffffffffffffffffffffffffffffffffffffff16906020019092919080359060200190929190803515159060200190929190505050610d74565b005b6102bc61143b565b6040518082815260200191505060405180910390f35b61031e600480360360408110156102e857600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919080359060200190929190505050611445565b005b61036c6004803603604081101561033657600080fd5b8101908080359060200190929190803573ffffffffffffffffffffffffffffffffffffffff1690602001909291905050506117a2565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156103ac578082015181840152602081019050610391565b50505050905090810190601f1680156103d95780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b610429600480360360208110156103fd57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291905050506118a0565b604051808215151515815260200191505060405180910390f35b6104856004803603602081101561045957600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919050505061190a565b6040518080602001828103825283818151815260200191508051906020019060200280838360005b838110156104c85780820151818401526020810190506104ad565b505050509050019250505060405180910390f35b60608060025467ffffffffffffffff811180156104f857600080fd5b506040519080825280602002602001820160405280156105275781602001602082028036833780820191505090505b50905060008090505b6002548110156105c05760008082815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1682828151811061057957fe5b602002602001019073ffffffffffffffffffffffffffffffffffffffff16908173ffffffffffffffffffffffffffffffffffffffff16815250508080600101915050610530565b508091505090565b600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060010160019054906101000a900460ff161561066e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602d815260200180611b9c602d913960400191505060405180910390fd5b33600080600254815260200190815260200160002060006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550600254600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000018190555060018060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060010160016101000a81548160ff021916908315150217905550600260008154809291906001019190505550600a600254106107ba577fadedc13440f3716a50214d07efec5c66588ad2249ceb6b61bb88774c499de8b96002546040518082815260200191505060405180910390a15b565b60003373ffffffffffffffffffffffffffffffffffffffff16600560008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600001600084815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1614905092915050565b6000600560008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600301600084815260200190815260200160002060000160009054906101000a900460ff16156108dc57600190506108e1565b600090505b92915050565b600280541015610942576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526028815260200180611bc96028913960400191505060405180910390fd5b600560008273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060010160009054906101000a900460ff16156109e8576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602a815260200180611c13602a913960400191505060405180910390fd5b6000606060025467ffffffffffffffff81118015610a0557600080fd5b50604051908082528060200260200182016040528015610a345781602001602082028036833780820191505090505b50905060008090505b60016002805481610a4a57fe5b0401811015610c67576002543342834302604051602001808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1660601b815260140183815260200182815260200193505050506040516020818303038152906040528051906020012060001c81610ac857fe5b069250600083905060008090505b600254811015610c58576000848381518110610aee57fe5b60200260200101511415610c3157600560008773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206002016000815480929190600101919050555060008083815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16600560008873ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600001600085815260200190815260200160002060006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506001848381518110610c2057fe5b602002602001018181525050610c58565b600160025403821415610c4357600091505b81806001019250508080600101915050610ad6565b50508080600101915050610a3d565b506001600560008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060010160006101000a81548160ff0219169083151502179055506004600081548092919060010191905055508273ffffffffffffffffffffffffffffffffffffffff167ff19127c2d25627f3ce43a78995d3a4f2e1f344c6ded175aaf766cd705672c684600560008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600201546040518082815260200191505060405180910390a2505050565b6000600354905090565b3373ffffffffffffffffffffffffffffffffffffffff16600560008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600001600084815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1614610e6b576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526022815260200180611bf16022913960400191505060405180910390fd5b600560008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600301600085815260200190815260200160002060000160009054906101000a900460ff16611035576001600560008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600301600086815260200190815260200160002060000160006101000a81548160ff021916908315150217905550600560008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600401600081548092919060010191905055506040518060400160405280600a81526020017f70726f63657373696e6700000000000000000000000000000000000000000000815250600560008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060030160008681526020019081526020016000206001019080519060200190611033929190611af6565b505b80156110a257600560008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206003016000858152602001908152602001600020600301600081548092919060010191905055505b600560008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600301600085815260200190815260200160002060020160008154809291906001019190505550600560008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060020154600560008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206003016000868152602001908152602001600020600201541415611435576002600560008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206003016000868152602001908152602001600020600201548161120857fe5b04600560008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600301600086815260200190815260200160002060030154111561134f576040518060400160405280600981526020017f6d616c6963696f75730000000000000000000000000000000000000000000000815250600560008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060030160008681526020019081526020016000206001019080519060200190611305929190611af6565b507f97cb40ddc2b324466029fb7b34af563c866d2fc02512c622674dd57b1e30625484600160405180838152602001821515151581526020019250505060405180910390a1611434565b6040518060400160405280600681526020017f62656e69676e0000000000000000000000000000000000000000000000000000815250600560008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600301600086815260200190815260200160002060010190805190602001906113ee929190611af6565b507f97cb40ddc2b324466029fb7b34af563c866d2fc02512c622674dd57b1e30625484600060405180838152602001821515151581526020019250505060405180910390a15b5b50505050565b6000600454905090565b600360025410156114a1576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526028815260200180611bc96028913960400191505060405180910390fd5b600560008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060010160009054906101000a900460ff1615611547576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602a815260200180611c13602a913960400191505060405180910390fd5b600080600090505b828110156116fa5760025481600454600202018161156957fe5b06915060008083815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16600560008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600001600083815260200190815260200160002060006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506000600182016002548161164057fe5b0614156116a3576001600560008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060010160006101000a81548160ff0219169083151502179055505b60018101600560008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060020181905550808060010191505061154f565b506004600081548092919060010191905055508273ffffffffffffffffffffffffffffffffffffffff167ff19127c2d25627f3ce43a78995d3a4f2e1f344c6ded175aaf766cd705672c684600560008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600201546040518082815260200191505060405180910390a2505050565b606080600560008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060030160008581526020019081526020016000206001018054600181600116156101000203166002900480601f01602080910402602001604051908101604052809291908181526020018280546001816001161561010002031660029004801561188f5780601f106118645761010080835404028352916020019161188f565b820191906000526020600020905b81548152906001019060200180831161187257829003601f168201915b505050505090508091505092915050565b6000600560008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060010160009054906101000a900460ff16156119005760019050611905565b600090505b919050565b6060600060025411611984576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601c8152602001807f4e756d626572206f662063616e64696461746573206973207a65726f0000000081525060200191505060405180910390fd5b6000600560008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060020154905060608167ffffffffffffffff811180156119e457600080fd5b50604051908082528060200260200182016040528015611a135781602001602082028036833780820191505090505b50905060008090505b82811015611aeb57600560008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600001600082815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16828281518110611aa457fe5b602002602001019073ffffffffffffffffffffffffffffffffffffffff16908173ffffffffffffffffffffffffffffffffffffffff16815250508080600101915050611a1c565b508092505050919050565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f10611b3757805160ff1916838001178555611b65565b82800160010185558215611b65579182015b82811115611b64578251825591602001919060010190611b49565b5b509050611b729190611b76565b5090565b611b9891905b80821115611b94576000816000905550600101611b7c565b5090565b9056fe63616e6e6f742072656769737465723a2063616e64696461746520616c726561647920726567697374657265644e756d626572206f662063616e6469646174657320646964206e6f7420726561636820322079657453656e646572206973206e6f7420612073656c65637465642076616c696461746f7243616e6e6f742070726f636573733a207472616e73616374696f6e20616c726561647920657869737473a2646970667358221220987d393761185f7e72939889ba931702108ebef2552a9ea2cdc2a72276b56e4964736f6c634300060b0033"

// DeploySelectionManager deploys a new Ethereum contract, binding an instance of SelectionManager to it.
func DeploySelectionManager(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *SelectionManager, error) {
	parsed, err := abi.JSON(strings.NewReader(SelectionManagerABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(SelectionManagerBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SelectionManager{SelectionManagerCaller: SelectionManagerCaller{contract: contract}, SelectionManagerTransactor: SelectionManagerTransactor{contract: contract}, SelectionManagerFilterer: SelectionManagerFilterer{contract: contract}}, nil
}

// SelectionManager is an auto generated Go binding around an Ethereum contract.
type SelectionManager struct {
	SelectionManagerCaller     // Read-only binding to the contract
	SelectionManagerTransactor // Write-only binding to the contract
	SelectionManagerFilterer   // Log filterer for contract events
}

// SelectionManagerCaller is an auto generated read-only Go binding around an Ethereum contract.
type SelectionManagerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SelectionManagerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SelectionManagerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SelectionManagerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SelectionManagerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SelectionManagerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SelectionManagerSession struct {
	Contract     *SelectionManager // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SelectionManagerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SelectionManagerCallerSession struct {
	Contract *SelectionManagerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// SelectionManagerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SelectionManagerTransactorSession struct {
	Contract     *SelectionManagerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// SelectionManagerRaw is an auto generated low-level Go binding around an Ethereum contract.
type SelectionManagerRaw struct {
	Contract *SelectionManager // Generic contract binding to access the raw methods on
}

// SelectionManagerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SelectionManagerCallerRaw struct {
	Contract *SelectionManagerCaller // Generic read-only contract binding to access the raw methods on
}

// SelectionManagerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SelectionManagerTransactorRaw struct {
	Contract *SelectionManagerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSelectionManager creates a new instance of SelectionManager, bound to a specific deployed contract.
func NewSelectionManager(address common.Address, backend bind.ContractBackend) (*SelectionManager, error) {
	contract, err := bindSelectionManager(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SelectionManager{SelectionManagerCaller: SelectionManagerCaller{contract: contract}, SelectionManagerTransactor: SelectionManagerTransactor{contract: contract}, SelectionManagerFilterer: SelectionManagerFilterer{contract: contract}}, nil
}

// NewSelectionManagerCaller creates a new read-only instance of SelectionManager, bound to a specific deployed contract.
func NewSelectionManagerCaller(address common.Address, caller bind.ContractCaller) (*SelectionManagerCaller, error) {
	contract, err := bindSelectionManager(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SelectionManagerCaller{contract: contract}, nil
}

// NewSelectionManagerTransactor creates a new write-only instance of SelectionManager, bound to a specific deployed contract.
func NewSelectionManagerTransactor(address common.Address, transactor bind.ContractTransactor) (*SelectionManagerTransactor, error) {
	contract, err := bindSelectionManager(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SelectionManagerTransactor{contract: contract}, nil
}

// NewSelectionManagerFilterer creates a new log filterer instance of SelectionManager, bound to a specific deployed contract.
func NewSelectionManagerFilterer(address common.Address, filterer bind.ContractFilterer) (*SelectionManagerFilterer, error) {
	contract, err := bindSelectionManager(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SelectionManagerFilterer{contract: contract}, nil
}

// bindSelectionManager binds a generic wrapper to an already deployed contract.
func bindSelectionManager(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SelectionManagerABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SelectionManager *SelectionManagerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _SelectionManager.Contract.SelectionManagerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SelectionManager *SelectionManagerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SelectionManager.Contract.SelectionManagerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SelectionManager *SelectionManagerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SelectionManager.Contract.SelectionManagerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SelectionManager *SelectionManagerCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _SelectionManager.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SelectionManager *SelectionManagerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SelectionManager.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SelectionManager *SelectionManagerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SelectionManager.Contract.contract.Transact(opts, method, params...)
}

// GetCandidates is a free data retrieval call binding the contract method 0x06a49fce.
//
// Solidity: function getCandidates() view returns(address[])
func (_SelectionManager *SelectionManagerCaller) GetCandidates(opts *bind.CallOpts) ([]common.Address, error) {
	var (
		ret0 = new([]common.Address)
	)
	out := ret0
	err := _SelectionManager.contract.Call(opts, out, "getCandidates")
	return *ret0, err
}

// GetCandidates is a free data retrieval call binding the contract method 0x06a49fce.
//
// Solidity: function getCandidates() view returns(address[])
func (_SelectionManager *SelectionManagerSession) GetCandidates() ([]common.Address, error) {
	return _SelectionManager.Contract.GetCandidates(&_SelectionManager.CallOpts)
}

// GetCandidates is a free data retrieval call binding the contract method 0x06a49fce.
//
// Solidity: function getCandidates() view returns(address[])
func (_SelectionManager *SelectionManagerCallerSession) GetCandidates() ([]common.Address, error) {
	return _SelectionManager.Contract.GetCandidates(&_SelectionManager.CallOpts)
}

// GetNumOfSCs is a free data retrieval call binding the contract method 0x8d1ead58.
//
// Solidity: function getNumOfSCs() view returns(uint256)
func (_SelectionManager *SelectionManagerCaller) GetNumOfSCs(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _SelectionManager.contract.Call(opts, out, "getNumOfSCs")
	return *ret0, err
}

// GetNumOfSCs is a free data retrieval call binding the contract method 0x8d1ead58.
//
// Solidity: function getNumOfSCs() view returns(uint256)
func (_SelectionManager *SelectionManagerSession) GetNumOfSCs() (*big.Int, error) {
	return _SelectionManager.Contract.GetNumOfSCs(&_SelectionManager.CallOpts)
}

// GetNumOfSCs is a free data retrieval call binding the contract method 0x8d1ead58.
//
// Solidity: function getNumOfSCs() view returns(uint256)
func (_SelectionManager *SelectionManagerCallerSession) GetNumOfSCs() (*big.Int, error) {
	return _SelectionManager.Contract.GetNumOfSCs(&_SelectionManager.CallOpts)
}

// GetNumOfTxs is a free data retrieval call binding the contract method 0x6110fb75.
//
// Solidity: function getNumOfTxs() view returns(uint256)
func (_SelectionManager *SelectionManagerCaller) GetNumOfTxs(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _SelectionManager.contract.Call(opts, out, "getNumOfTxs")
	return *ret0, err
}

// GetNumOfTxs is a free data retrieval call binding the contract method 0x6110fb75.
//
// Solidity: function getNumOfTxs() view returns(uint256)
func (_SelectionManager *SelectionManagerSession) GetNumOfTxs() (*big.Int, error) {
	return _SelectionManager.Contract.GetNumOfTxs(&_SelectionManager.CallOpts)
}

// GetNumOfTxs is a free data retrieval call binding the contract method 0x6110fb75.
//
// Solidity: function getNumOfTxs() view returns(uint256)
func (_SelectionManager *SelectionManagerCallerSession) GetNumOfTxs() (*big.Int, error) {
	return _SelectionManager.Contract.GetNumOfTxs(&_SelectionManager.CallOpts)
}

// GetStatus is a free data retrieval call binding the contract method 0xdb5f9f36.
//
// Solidity: function getStatus(bytes32 txHash, address scAddr) view returns(string)
func (_SelectionManager *SelectionManagerCaller) GetStatus(opts *bind.CallOpts, txHash [32]byte, scAddr common.Address) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _SelectionManager.contract.Call(opts, out, "getStatus", txHash, scAddr)
	return *ret0, err
}

// GetStatus is a free data retrieval call binding the contract method 0xdb5f9f36.
//
// Solidity: function getStatus(bytes32 txHash, address scAddr) view returns(string)
func (_SelectionManager *SelectionManagerSession) GetStatus(txHash [32]byte, scAddr common.Address) (string, error) {
	return _SelectionManager.Contract.GetStatus(&_SelectionManager.CallOpts, txHash, scAddr)
}

// GetStatus is a free data retrieval call binding the contract method 0xdb5f9f36.
//
// Solidity: function getStatus(bytes32 txHash, address scAddr) view returns(string)
func (_SelectionManager *SelectionManagerCallerSession) GetStatus(txHash [32]byte, scAddr common.Address) (string, error) {
	return _SelectionManager.Contract.GetStatus(&_SelectionManager.CallOpts, txHash, scAddr)
}

// GetValidators is a free data retrieval call binding the contract method 0xff8744a6.
//
// Solidity: function getValidators(address scAddr) view returns(address[])
func (_SelectionManager *SelectionManagerCaller) GetValidators(opts *bind.CallOpts, scAddr common.Address) ([]common.Address, error) {
	var (
		ret0 = new([]common.Address)
	)
	out := ret0
	err := _SelectionManager.contract.Call(opts, out, "getValidators", scAddr)
	return *ret0, err
}

// GetValidators is a free data retrieval call binding the contract method 0xff8744a6.
//
// Solidity: function getValidators(address scAddr) view returns(address[])
func (_SelectionManager *SelectionManagerSession) GetValidators(scAddr common.Address) ([]common.Address, error) {
	return _SelectionManager.Contract.GetValidators(&_SelectionManager.CallOpts, scAddr)
}

// GetValidators is a free data retrieval call binding the contract method 0xff8744a6.
//
// Solidity: function getValidators(address scAddr) view returns(address[])
func (_SelectionManager *SelectionManagerCallerSession) GetValidators(scAddr common.Address) ([]common.Address, error) {
	return _SelectionManager.Contract.GetValidators(&_SelectionManager.CallOpts, scAddr)
}

// IsSCExists is a free data retrieval call binding the contract method 0xe160e120.
//
// Solidity: function isSCExists(address scAddr) view returns(bool)
func (_SelectionManager *SelectionManagerCaller) IsSCExists(opts *bind.CallOpts, scAddr common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _SelectionManager.contract.Call(opts, out, "isSCExists", scAddr)
	return *ret0, err
}

// IsSCExists is a free data retrieval call binding the contract method 0xe160e120.
//
// Solidity: function isSCExists(address scAddr) view returns(bool)
func (_SelectionManager *SelectionManagerSession) IsSCExists(scAddr common.Address) (bool, error) {
	return _SelectionManager.Contract.IsSCExists(&_SelectionManager.CallOpts, scAddr)
}

// IsSCExists is a free data retrieval call binding the contract method 0xe160e120.
//
// Solidity: function isSCExists(address scAddr) view returns(bool)
func (_SelectionManager *SelectionManagerCallerSession) IsSCExists(scAddr common.Address) (bool, error) {
	return _SelectionManager.Contract.IsSCExists(&_SelectionManager.CallOpts, scAddr)
}

// IsSelected is a free data retrieval call binding the contract method 0x337299eb.
//
// Solidity: function isSelected(address scAddr, uint256 index) view returns(bool)
func (_SelectionManager *SelectionManagerCaller) IsSelected(opts *bind.CallOpts, scAddr common.Address, index *big.Int) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _SelectionManager.contract.Call(opts, out, "isSelected", scAddr, index)
	return *ret0, err
}

// IsSelected is a free data retrieval call binding the contract method 0x337299eb.
//
// Solidity: function isSelected(address scAddr, uint256 index) view returns(bool)
func (_SelectionManager *SelectionManagerSession) IsSelected(scAddr common.Address, index *big.Int) (bool, error) {
	return _SelectionManager.Contract.IsSelected(&_SelectionManager.CallOpts, scAddr, index)
}

// IsSelected is a free data retrieval call binding the contract method 0x337299eb.
//
// Solidity: function isSelected(address scAddr, uint256 index) view returns(bool)
func (_SelectionManager *SelectionManagerCallerSession) IsSelected(scAddr common.Address, index *big.Int) (bool, error) {
	return _SelectionManager.Contract.IsSelected(&_SelectionManager.CallOpts, scAddr, index)
}

// IsTxExists is a free data retrieval call binding the contract method 0x4b4dd27d.
//
// Solidity: function isTxExists(bytes32 txHash, address scAddr) view returns(bool)
func (_SelectionManager *SelectionManagerCaller) IsTxExists(opts *bind.CallOpts, txHash [32]byte, scAddr common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _SelectionManager.contract.Call(opts, out, "isTxExists", txHash, scAddr)
	return *ret0, err
}

// IsTxExists is a free data retrieval call binding the contract method 0x4b4dd27d.
//
// Solidity: function isTxExists(bytes32 txHash, address scAddr) view returns(bool)
func (_SelectionManager *SelectionManagerSession) IsTxExists(txHash [32]byte, scAddr common.Address) (bool, error) {
	return _SelectionManager.Contract.IsTxExists(&_SelectionManager.CallOpts, txHash, scAddr)
}

// IsTxExists is a free data retrieval call binding the contract method 0x4b4dd27d.
//
// Solidity: function isTxExists(bytes32 txHash, address scAddr) view returns(bool)
func (_SelectionManager *SelectionManagerCallerSession) IsTxExists(txHash [32]byte, scAddr common.Address) (bool, error) {
	return _SelectionManager.Contract.IsTxExists(&_SelectionManager.CallOpts, txHash, scAddr)
}

// Register is a paid mutator transaction binding the contract method 0x1aa3a008.
//
// Solidity: function register() returns()
func (_SelectionManager *SelectionManagerTransactor) Register(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SelectionManager.contract.Transact(opts, "register")
}

// Register is a paid mutator transaction binding the contract method 0x1aa3a008.
//
// Solidity: function register() returns()
func (_SelectionManager *SelectionManagerSession) Register() (*types.Transaction, error) {
	return _SelectionManager.Contract.Register(&_SelectionManager.TransactOpts)
}

// Register is a paid mutator transaction binding the contract method 0x1aa3a008.
//
// Solidity: function register() returns()
func (_SelectionManager *SelectionManagerTransactorSession) Register() (*types.Transaction, error) {
	return _SelectionManager.Contract.Register(&_SelectionManager.TransactOpts)
}

// Select is a paid mutator transaction binding the contract method 0x4f49f01c.
//
// Solidity: function select(address scAddr) returns()
func (_SelectionManager *SelectionManagerTransactor) Select(opts *bind.TransactOpts, scAddr common.Address) (*types.Transaction, error) {
	return _SelectionManager.contract.Transact(opts, "select", scAddr)
}

// Select is a paid mutator transaction binding the contract method 0x4f49f01c.
//
// Solidity: function select(address scAddr) returns()
func (_SelectionManager *SelectionManagerSession) Select(scAddr common.Address) (*types.Transaction, error) {
	return _SelectionManager.Contract.Select(&_SelectionManager.TransactOpts, scAddr)
}

// Select is a paid mutator transaction binding the contract method 0x4f49f01c.
//
// Solidity: function select(address scAddr) returns()
func (_SelectionManager *SelectionManagerTransactorSession) Select(scAddr common.Address) (*types.Transaction, error) {
	return _SelectionManager.Contract.Select(&_SelectionManager.TransactOpts, scAddr)
}

// SelectFixed is a paid mutator transaction binding the contract method 0xa40d5bb3.
//
// Solidity: function selectFixed(address scAddr, uint256 numOfVals) returns()
func (_SelectionManager *SelectionManagerTransactor) SelectFixed(opts *bind.TransactOpts, scAddr common.Address, numOfVals *big.Int) (*types.Transaction, error) {
	return _SelectionManager.contract.Transact(opts, "selectFixed", scAddr, numOfVals)
}

// SelectFixed is a paid mutator transaction binding the contract method 0xa40d5bb3.
//
// Solidity: function selectFixed(address scAddr, uint256 numOfVals) returns()
func (_SelectionManager *SelectionManagerSession) SelectFixed(scAddr common.Address, numOfVals *big.Int) (*types.Transaction, error) {
	return _SelectionManager.Contract.SelectFixed(&_SelectionManager.TransactOpts, scAddr, numOfVals)
}

// SelectFixed is a paid mutator transaction binding the contract method 0xa40d5bb3.
//
// Solidity: function selectFixed(address scAddr, uint256 numOfVals) returns()
func (_SelectionManager *SelectionManagerTransactorSession) SelectFixed(scAddr common.Address, numOfVals *big.Int) (*types.Transaction, error) {
	return _SelectionManager.Contract.SelectFixed(&_SelectionManager.TransactOpts, scAddr, numOfVals)
}

// SubmitResult is a paid mutator transaction binding the contract method 0x6b9bdf34.
//
// Solidity: function submitResult(bytes32 txHash, address scAddr, uint256 index, bool malicious) returns()
func (_SelectionManager *SelectionManagerTransactor) SubmitResult(opts *bind.TransactOpts, txHash [32]byte, scAddr common.Address, index *big.Int, malicious bool) (*types.Transaction, error) {
	return _SelectionManager.contract.Transact(opts, "submitResult", txHash, scAddr, index, malicious)
}

// SubmitResult is a paid mutator transaction binding the contract method 0x6b9bdf34.
//
// Solidity: function submitResult(bytes32 txHash, address scAddr, uint256 index, bool malicious) returns()
func (_SelectionManager *SelectionManagerSession) SubmitResult(txHash [32]byte, scAddr common.Address, index *big.Int, malicious bool) (*types.Transaction, error) {
	return _SelectionManager.Contract.SubmitResult(&_SelectionManager.TransactOpts, txHash, scAddr, index, malicious)
}

// SubmitResult is a paid mutator transaction binding the contract method 0x6b9bdf34.
//
// Solidity: function submitResult(bytes32 txHash, address scAddr, uint256 index, bool malicious) returns()
func (_SelectionManager *SelectionManagerTransactorSession) SubmitResult(txHash [32]byte, scAddr common.Address, index *big.Int, malicious bool) (*types.Transaction, error) {
	return _SelectionManager.Contract.SubmitResult(&_SelectionManager.TransactOpts, txHash, scAddr, index, malicious)
}

// SelectionManagerAddedTxIterator is returned from FilterAddedTx and is used to iterate over the raw logs and unpacked data for AddedTx events raised by the SelectionManager contract.
type SelectionManagerAddedTxIterator struct {
	Event *SelectionManagerAddedTx // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SelectionManagerAddedTxIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SelectionManagerAddedTx)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SelectionManagerAddedTx)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SelectionManagerAddedTxIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SelectionManagerAddedTxIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SelectionManagerAddedTx represents a AddedTx event raised by the SelectionManager contract.
type SelectionManagerAddedTx struct {
	TxHash          [32]byte
	ScAddr          common.Address
	NumOfValidators *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterAddedTx is a free log retrieval operation binding the contract event 0x22c33fe72b667706dfe728273e99df11a541f14729723692fa6676520a1f3595.
//
// Solidity: event AddedTx(bytes32 indexed _txHash, address _scAddr, uint256 _numOfValidators)
func (_SelectionManager *SelectionManagerFilterer) FilterAddedTx(opts *bind.FilterOpts, _txHash [][32]byte) (*SelectionManagerAddedTxIterator, error) {

	var _txHashRule []interface{}
	for _, _txHashItem := range _txHash {
		_txHashRule = append(_txHashRule, _txHashItem)
	}

	logs, sub, err := _SelectionManager.contract.FilterLogs(opts, "AddedTx", _txHashRule)
	if err != nil {
		return nil, err
	}
	return &SelectionManagerAddedTxIterator{contract: _SelectionManager.contract, event: "AddedTx", logs: logs, sub: sub}, nil
}

// WatchAddedTx is a free log subscription operation binding the contract event 0x22c33fe72b667706dfe728273e99df11a541f14729723692fa6676520a1f3595.
//
// Solidity: event AddedTx(bytes32 indexed _txHash, address _scAddr, uint256 _numOfValidators)
func (_SelectionManager *SelectionManagerFilterer) WatchAddedTx(opts *bind.WatchOpts, sink chan<- *SelectionManagerAddedTx, _txHash [][32]byte) (event.Subscription, error) {

	var _txHashRule []interface{}
	for _, _txHashItem := range _txHash {
		_txHashRule = append(_txHashRule, _txHashItem)
	}

	logs, sub, err := _SelectionManager.contract.WatchLogs(opts, "AddedTx", _txHashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SelectionManagerAddedTx)
				if err := _SelectionManager.contract.UnpackLog(event, "AddedTx", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseAddedTx is a log parse operation binding the contract event 0x22c33fe72b667706dfe728273e99df11a541f14729723692fa6676520a1f3595.
//
// Solidity: event AddedTx(bytes32 indexed _txHash, address _scAddr, uint256 _numOfValidators)
func (_SelectionManager *SelectionManagerFilterer) ParseAddedTx(log types.Log) (*SelectionManagerAddedTx, error) {
	event := new(SelectionManagerAddedTx)
	if err := _SelectionManager.contract.UnpackLog(event, "AddedTx", log); err != nil {
		return nil, err
	}
	return event, nil
}

// SelectionManagerProcessedTransactionIterator is returned from FilterProcessedTransaction and is used to iterate over the raw logs and unpacked data for ProcessedTransaction events raised by the SelectionManager contract.
type SelectionManagerProcessedTransactionIterator struct {
	Event *SelectionManagerProcessedTransaction // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SelectionManagerProcessedTransactionIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SelectionManagerProcessedTransaction)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SelectionManagerProcessedTransaction)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SelectionManagerProcessedTransactionIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SelectionManagerProcessedTransactionIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SelectionManagerProcessedTransaction represents a ProcessedTransaction event raised by the SelectionManager contract.
type SelectionManagerProcessedTransaction struct {
	TxHash [32]byte
	Result bool
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterProcessedTransaction is a free log retrieval operation binding the contract event 0x97cb40ddc2b324466029fb7b34af563c866d2fc02512c622674dd57b1e306254.
//
// Solidity: event ProcessedTransaction(bytes32 _txHash, bool _result)
func (_SelectionManager *SelectionManagerFilterer) FilterProcessedTransaction(opts *bind.FilterOpts) (*SelectionManagerProcessedTransactionIterator, error) {

	logs, sub, err := _SelectionManager.contract.FilterLogs(opts, "ProcessedTransaction")
	if err != nil {
		return nil, err
	}
	return &SelectionManagerProcessedTransactionIterator{contract: _SelectionManager.contract, event: "ProcessedTransaction", logs: logs, sub: sub}, nil
}

// WatchProcessedTransaction is a free log subscription operation binding the contract event 0x97cb40ddc2b324466029fb7b34af563c866d2fc02512c622674dd57b1e306254.
//
// Solidity: event ProcessedTransaction(bytes32 _txHash, bool _result)
func (_SelectionManager *SelectionManagerFilterer) WatchProcessedTransaction(opts *bind.WatchOpts, sink chan<- *SelectionManagerProcessedTransaction) (event.Subscription, error) {

	logs, sub, err := _SelectionManager.contract.WatchLogs(opts, "ProcessedTransaction")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SelectionManagerProcessedTransaction)
				if err := _SelectionManager.contract.UnpackLog(event, "ProcessedTransaction", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseProcessedTransaction is a log parse operation binding the contract event 0x97cb40ddc2b324466029fb7b34af563c866d2fc02512c622674dd57b1e306254.
//
// Solidity: event ProcessedTransaction(bytes32 _txHash, bool _result)
func (_SelectionManager *SelectionManagerFilterer) ParseProcessedTransaction(log types.Log) (*SelectionManagerProcessedTransaction, error) {
	event := new(SelectionManagerProcessedTransaction)
	if err := _SelectionManager.contract.UnpackLog(event, "ProcessedTransaction", log); err != nil {
		return nil, err
	}
	return event, nil
}

// SelectionManagerRegisteredNodesIterator is returned from FilterRegisteredNodes and is used to iterate over the raw logs and unpacked data for RegisteredNodes events raised by the SelectionManager contract.
type SelectionManagerRegisteredNodesIterator struct {
	Event *SelectionManagerRegisteredNodes // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SelectionManagerRegisteredNodesIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SelectionManagerRegisteredNodes)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SelectionManagerRegisteredNodes)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SelectionManagerRegisteredNodesIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SelectionManagerRegisteredNodesIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SelectionManagerRegisteredNodes represents a RegisteredNodes event raised by the SelectionManager contract.
type SelectionManagerRegisteredNodes struct {
	NumOfCandidates *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterRegisteredNodes is a free log retrieval operation binding the contract event 0xadedc13440f3716a50214d07efec5c66588ad2249ceb6b61bb88774c499de8b9.
//
// Solidity: event RegisteredNodes(uint256 _numOfCandidates)
func (_SelectionManager *SelectionManagerFilterer) FilterRegisteredNodes(opts *bind.FilterOpts) (*SelectionManagerRegisteredNodesIterator, error) {

	logs, sub, err := _SelectionManager.contract.FilterLogs(opts, "RegisteredNodes")
	if err != nil {
		return nil, err
	}
	return &SelectionManagerRegisteredNodesIterator{contract: _SelectionManager.contract, event: "RegisteredNodes", logs: logs, sub: sub}, nil
}

// WatchRegisteredNodes is a free log subscription operation binding the contract event 0xadedc13440f3716a50214d07efec5c66588ad2249ceb6b61bb88774c499de8b9.
//
// Solidity: event RegisteredNodes(uint256 _numOfCandidates)
func (_SelectionManager *SelectionManagerFilterer) WatchRegisteredNodes(opts *bind.WatchOpts, sink chan<- *SelectionManagerRegisteredNodes) (event.Subscription, error) {

	logs, sub, err := _SelectionManager.contract.WatchLogs(opts, "RegisteredNodes")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SelectionManagerRegisteredNodes)
				if err := _SelectionManager.contract.UnpackLog(event, "RegisteredNodes", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRegisteredNodes is a log parse operation binding the contract event 0xadedc13440f3716a50214d07efec5c66588ad2249ceb6b61bb88774c499de8b9.
//
// Solidity: event RegisteredNodes(uint256 _numOfCandidates)
func (_SelectionManager *SelectionManagerFilterer) ParseRegisteredNodes(log types.Log) (*SelectionManagerRegisteredNodes, error) {
	event := new(SelectionManagerRegisteredNodes)
	if err := _SelectionManager.contract.UnpackLog(event, "RegisteredNodes", log); err != nil {
		return nil, err
	}
	return event, nil
}

// SelectionManagerSelectedSubsetIterator is returned from FilterSelectedSubset and is used to iterate over the raw logs and unpacked data for SelectedSubset events raised by the SelectionManager contract.
type SelectionManagerSelectedSubsetIterator struct {
	Event *SelectionManagerSelectedSubset // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SelectionManagerSelectedSubsetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SelectionManagerSelectedSubset)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SelectionManagerSelectedSubset)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SelectionManagerSelectedSubsetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SelectionManagerSelectedSubsetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SelectionManagerSelectedSubset represents a SelectedSubset event raised by the SelectionManager contract.
type SelectionManagerSelectedSubset struct {
	ScAddr          common.Address
	NumOfValidators *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterSelectedSubset is a free log retrieval operation binding the contract event 0xf19127c2d25627f3ce43a78995d3a4f2e1f344c6ded175aaf766cd705672c684.
//
// Solidity: event SelectedSubset(address indexed _scAddr, uint256 _numOfValidators)
func (_SelectionManager *SelectionManagerFilterer) FilterSelectedSubset(opts *bind.FilterOpts, _scAddr []common.Address) (*SelectionManagerSelectedSubsetIterator, error) {

	var _scAddrRule []interface{}
	for _, _scAddrItem := range _scAddr {
		_scAddrRule = append(_scAddrRule, _scAddrItem)
	}

	logs, sub, err := _SelectionManager.contract.FilterLogs(opts, "SelectedSubset", _scAddrRule)
	if err != nil {
		return nil, err
	}
	return &SelectionManagerSelectedSubsetIterator{contract: _SelectionManager.contract, event: "SelectedSubset", logs: logs, sub: sub}, nil
}

// WatchSelectedSubset is a free log subscription operation binding the contract event 0xf19127c2d25627f3ce43a78995d3a4f2e1f344c6ded175aaf766cd705672c684.
//
// Solidity: event SelectedSubset(address indexed _scAddr, uint256 _numOfValidators)
func (_SelectionManager *SelectionManagerFilterer) WatchSelectedSubset(opts *bind.WatchOpts, sink chan<- *SelectionManagerSelectedSubset, _scAddr []common.Address) (event.Subscription, error) {

	var _scAddrRule []interface{}
	for _, _scAddrItem := range _scAddr {
		_scAddrRule = append(_scAddrRule, _scAddrItem)
	}

	logs, sub, err := _SelectionManager.contract.WatchLogs(opts, "SelectedSubset", _scAddrRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SelectionManagerSelectedSubset)
				if err := _SelectionManager.contract.UnpackLog(event, "SelectedSubset", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseSelectedSubset is a log parse operation binding the contract event 0xf19127c2d25627f3ce43a78995d3a4f2e1f344c6ded175aaf766cd705672c684.
//
// Solidity: event SelectedSubset(address indexed _scAddr, uint256 _numOfValidators)
func (_SelectionManager *SelectionManagerFilterer) ParseSelectedSubset(log types.Log) (*SelectionManagerSelectedSubset, error) {
	event := new(SelectionManagerSelectedSubset)
	if err := _SelectionManager.contract.UnpackLog(event, "SelectedSubset", log); err != nil {
		return nil, err
	}
	return event, nil
}
