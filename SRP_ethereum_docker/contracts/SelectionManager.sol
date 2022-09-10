//@SPDX-License-Identifier: UNLICENSED
pragma solidity >=0.4.22 <=0.7.0;

//0xDB571079aF66EDbB1a56d22809584d39C20001D9

contract SelectionManager {
    struct candidate {
        uint256 id;
        bool selected;
        bool isExist;
    }
    
    // represents a contract
    struct smartContract {
        mapping(uint256 => address) validators; // list of currently selected validators of a given tx
        bool isExist; 
        uint256 numOfValidators; // counts the number of selected validators
        mapping(bytes32 => transaction) txs;
        uint256 numOfTxs; // counts the number of txs invoking the SC  & were assigned to selected validators
        
    }
    struct transaction {
        bool isExist; // indicates whether a subset of candidates was already selected for the tx
        string status; // describes the status of the transaction, can be 'processing', 'malicious', or 'benign'
        uint256 totalVotes; // counts the number of validators that have submitted the runtime-verification result
        uint256 results; // counts the nunber of results classifying tx as malicious
    }

    mapping(uint256 => address) candidates; // used for direct access during selection
    mapping(address => candidate) candidatesInfo; // used for direct access during registration
    uint256 numOfCandidates;
    uint256 numOfSCs;
    mapping (address => smartContract) scs;

    // event that is emitted once a new tx is added 
    event AddedTx(bytes32 indexed _txHash, address _scAddr, uint256 _numOfValidators);
    
    // event that is emitted once the selection process terminates successfully
    event SelectedSubset(address indexed _scAddr, uint256 _numOfValidators);
    // event that emits upon reaching 10 candidates or higher
    event RegisteredNodes(uint256 _numOfCandidates);

    // event that emits upon receiving runtime verification result submitted by validators
    event ProcessedTransaction(bytes32 _txHash, bool _result);

    constructor() public {
        numOfCandidates = 0;
        numOfSCs = 0;
    }
    
   //******************* FUNCTIONS *******************//

    // registers sender's address as a candidate if it is qualified based on certain rules
    function register() public {

        require(
            !candidatesInfo[msg.sender].isExist,
            "cannot register: candidate already registered"
        );
        //register sender by adding it to list of candidates and storing its info
        candidates[numOfCandidates] = msg.sender;
        candidatesInfo[msg.sender].id = numOfCandidates;
        candidatesInfo[msg.sender].isExist = true;
        numOfCandidates++;

        if (numOfCandidates >= 10) {
            emit RegisteredNodes(numOfCandidates);
        }
    }

    // pseudorandomly selects a subset of nodes as validators of the given smart contract (identified by its address),
    // and emits an event once the selection process terminates successfully
    function select(address scAddr) public {
        //Only start selection if there are sufficient number of candidates
        require(
            numOfCandidates >= 3,
            "Number of candidates did not reach 2 yet"
        ); //isra: to be changed
        
        uint256 _id;
        uint256[] memory _selected = new uint[](numOfCandidates); //acts like a mapping, where index is _id & val is 0/1
        //select validators from the pool of candidates, currently selecting 50% of pool
        for (uint256 i = 0; i < numOfCandidates / 2 + 1; i++) {
            _id =
                uint256(
                    keccak256(
                        abi.encodePacked(
                            msg.sender,
                            block.timestamp,
                            block.number * i
                        )
                    )
                ) %
                numOfCandidates;
            // avoid selecting an already selected validator 
            uint256 _tempID = _id;
            for (uint256 j = 0; j < numOfCandidates; j++ ){
                if (_selected[_tempID] == 0){
                    scs[scAddr].numOfValidators++;
                    scs[scAddr].validators[i] = candidates[_tempID];
                    // set validator's local selected-flag
                    _selected[_tempID] = 1;
                    break;
                } else {
                    if (_tempID == numOfCandidates - 1){
                        _tempID = 0;
                    }
                    _tempID++;  
                }
            }
            
        }
        // if SC was not registered, add it to list of subscribed SCs
        if (!scs[scAddr].isExist) {
            // set the smart contract's exist flag
            scs[scAddr].isExist = true;
            // increment number of subscribed 
            numOfSCs++;
        }
        
        // triggering the selection event, passing the txhash that can be used to get the list of validators
        emit SelectedSubset(scAddr, scs[scAddr].numOfValidators); 
    }

    // selectFixed() is used for testing-purposes
    function selectFixed(address scAddr, uint256 numOfVals) public {
        //Only start selection if there are sufficient number of candidates
        require(
            numOfCandidates >= 3,
            "Number of candidates did not reach 3 yet"
        ); 

        uint256 _id;
        for (uint256 i = 0; i < numOfVals; i++) {
            _id = (2*numOfSCs+i) % numOfCandidates;
            // assign selected candidate to list of validators for the SC
            scs[scAddr].validators[i] = candidates[_id];
            if (numOfCandidates % (i+1) == 0){
                // set the transaction's exist flag
                scs[scAddr].isExist = true;
            }
            scs[scAddr].numOfValidators = i+1;
        }
        
        // increment number of processed txs
        numOfSCs++;
        
        // triggering the selection event, passing the txhash that can be used to get the list of validators
        emit SelectedSubset(scAddr, scs[scAddr].numOfValidators); 
    }




    // checks & verifies the result submitted by selected validators for the given tx
    function submitResult(
        bytes32 txHash,
        address scAddr,
        uint256 index,
        bool malicious
    ) public {
       require(scs[scAddr].validators[index] == msg.sender,"Sender is not a selected validator");
        //add tx if it doesn't exists
        if (!scs[scAddr].txs[txHash].isExist) {
             // set the transaction's exist flag
            scs[scAddr].txs[txHash].isExist = true;
            // increment number of processed txs
            scs[scAddr].numOfTxs++;
            // set the transaction's status as processing
            scs[scAddr].txs[txHash].status = "processing";
        }
        if (malicious) {
            scs[scAddr].txs[txHash].results++; 
        }
        scs[scAddr].txs[txHash].totalVotes++; 

        if (scs[scAddr].txs[txHash].totalVotes == scs[scAddr].numOfValidators) {
            if (scs[scAddr].txs[txHash].results > scs[scAddr].txs[txHash].totalVotes / 2) {
                scs[scAddr].txs[txHash].status = "malicious";
                emit ProcessedTransaction(txHash, true);
            } else {
                scs[scAddr].txs[txHash].status = "benign";
                emit ProcessedTransaction(txHash, false);
            }
        }
    }
    
  
    //returns true if transaction exists, otherwise false
    function isSCExists(address scAddr) public view returns (bool) {
       if (scs[scAddr].isExist)
            return true;
        else 
            return false;
    }
    //returns true if transaction exists, otherwise false
    function isTxExists(bytes32 txHash, address scAddr) public view returns (bool) {
       if (scs[scAddr].txs[txHash].isExist)
            return true;
        else 
            return false;
    }
    // returns current status of given tx, can be  processed, otherwise returns 'processing'
    function getStatus(bytes32 txHash, address scAddr) public view returns (string memory){
        string memory _status = scs[scAddr].txs[txHash].status;
        return _status;
    }

    // returns the list of validators for a given transaction, which is identified by tx hash
    function getValidators(address scAddr)
        public
        view
        returns (address[] memory)
    {
        require(numOfCandidates > 0, "Number of candidates is zero");
        uint256 _nVals = scs[scAddr].numOfValidators;
        address[] memory _validators = new address[](_nVals);
        for (uint256 i = 0; i < _nVals; i++) {
            _validators[i] = scs[scAddr].validators[i];
        }

        return _validators;
    }

    function getCandidates() public view returns (address[] memory) {
        address[] memory _candidates = new address[](numOfCandidates);
        for (uint256 i = 0; i < numOfCandidates; i++) {
            _candidates[i] = candidates[i];
        }
        return _candidates;
    }
    
    function getNumOfSCs() public view returns (uint256) {
        return numOfSCs;
    }
    
    function isSelected(address scAddr, uint256 index) public view returns (bool){
        return scs[scAddr].validators[index] == msg.sender;
    }
    
    
}

