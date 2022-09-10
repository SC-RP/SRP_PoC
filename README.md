# SRP: An Efficient Runtime Protection Framework for Smart Contracts
This repository contains a proof of concept implementation of the Smart contract Runtime Protection framework, SRP. It also contains the programs and scripts needed to execute experiments that test and evaluate SRP.


### Local Ethereum Network

We ran our experiments, reported in the paper, on a local Ethereum network deployed using docker-compose. 
The current network consists of three nodes (can be scaled) of which is a bootstrap node. 
Each node is a docker container. Mining is enabled for the bootstrap node only. The configuration of the network is set in 
[docker-compose.yml](https://github.com/SC-RP/SRP_PoC/blob/main/SRP_ethereum_docker/docker-compose.yml).

### SRP Code
In our proof-of-concept implementation, we extend Geth by modifying its source code to embed the main components of SRP and modify the consensus protocol to integrate the off-chain mechanism. We created a [smart contract](SRP_ethereum_docker/contracts/SelectionManager.sol), written in Solidity, that implements Selection Manager, which manages the list of registered candidates and randomly assigns validators to subscribed smart contracts that seek protection. We added a module to the [github.com/ethereum/go-ethereum/core](SRP_ethereum_docker/GO_PKGs/ethereum/go-ethereum/core) package that implements the main functionalities of SRP's protocol.

## Repository Structure

* [./GO_PKG/](https://github.com/SC-RP/SRP_PoC/tree/main/SRP_ethereum_docker/GO_PKGs) - contains all Go packages and modules needed to run SRP. This includes the modified source code of go-ethereum (Geth v 1.9.16).
* [./cmd/](https://github.com/SC-RP/SRP_PoC/tree/main/SRP_ethereum_docker/cmd) - has all scripts and programs that are executed to setup, run, and collect measurements of SRP.

* [./contracts/](https://github.com/SC-RP/SRP_PoC/tree/main/SRP_ethereum_docker/contracts) - contains the solidity source code of selection manager contract and a sample victim contract.

* [./files/](https://github.com/SC-RP/SRP_PoC/tree/main/SRP_ethereum_docker/files)- contains all files imported during the setup and execution of the experiments, including the genesis file of Ethereum blockchain, configuration parameters, account keystore directories, list of registered smart contracts, and the signed transactions that invoke sample smart contract. 

* [./monitored-geth-bootstrap/](https://github.com/SC-RP/SRP_PoC/tree/main/SRP_ethereum_docker/geth-bootstrap) - has the Dockerfile of the bootstrap node and the entry point script of the container. 

* [./monitored-geth-client/](https://github.com/SC-RP/SRP_PoC/tree/main/SRP_ethereum_docker/geth-bootstrap) - has the Dockerfile of eth node and the entry point script of the container. 

* [./Results](https://github.com/SC-RP/SRP_PoC/tree/main/SRP_ethereum_docker/Results) - has the collected measurements stored in result csv file.


## Execution Requirements and Dependencies
The main requirements of setting up the testing environment for SRP are: 
1. [Docker](https://docs.docker.com/engine/install/) and [Docker-compose](https://docs.docker.com/compose/install/) 
2. [Golang](https://go.dev/doc/install) - Part of running SRP is initializing the protocol by registering candidate nodes (offchain validators) and triggering the selection process. Both registration and selection are executed onchain by invoking the selection manager contract [SelectionManager.sol](https://github.com/SC-RP/SRP_PoC/blob/main/SRP_ethereum_docker/contracts/SelectionManager.sol). This is done in [initSRP.go](https://github.com/SC-RP/SRP_PoC/blob/main/SRP_ethereum_docker/cmd/runSRP/initSRP.go). To run initSRP.go, you need to install Golang and add  our local modules: [registration](https://github.com/SC-RP/SRP_PoC/tree/main/SRP_ethereum_docker/GO_PKGs/registration) and [selectionManager](https://github.com/SC-RP/SRP_PoC/tree/main/SRP_ethereum_docker/GO_PKGs/selectionManager) to your GOPATH.

3. [Node.js](https://nodejs.dev/en/learn/how-to-install-nodejs/) - To send transactions using our scripts, you need to install Node.js.  

4. Repository Structure - Scripts and programs read/write files from given paths that are set according to the current directory structure. 


## Configuration of SRP 
Using the bash script, [runExperiment.sh](https://github.com/SC-RP/SRP_PoC/blob/main/SRP_ethereum_docker/cmd/runExperiment.sh), to run the experiment, then you can change configuration parameters directly from the bash script. Otherwise,  you can modify [parameter_configuration.json](https://github.com/SC-RP/SRP_PoC/blob/main/SRP_ethereum_docker/files/parameter_configuration.json), or execute [changeParameters.js](https://github.com/SC-RP/SRP_PoC/blob/main/SRP_ethereum_docker/cmd/changeParameters.js) with desired parameters entered in defined order as: 
```bash 
node changeParameters.js $rv_time $numOfNodes $numOfVictims $rate $numOftxs $experiment_id
```


## Running SRP & Collecting Measurements 
You can run SRP by executing the bash script
[./cmd/runExperiment.sh](https://github.com/SC-RP/SRP_PoC/blob/main/SRP_ethereum_docker/cmd/runExperiment.sh): 
```bash 
./runExperiment.sh
```

This script executes all steps needed to run SRP (in the correct order), including:
1. Setting configuration parameters
2. Running docker-compose
2. Initializing SRP 
3. Sending SC Txs
4. Collecting measurements 

⚠️ NOTE: Scripts import/export files from/to other directories in this repository. To run this code successfully, you must execute programs within their directories. 

## Expected Output 
* JSON Files
    - [sendSCtxs.js](https://github.com/SC-RP/SRP_PoC/blob/main/SRP_ethereum_docker/cmd/runSRP/sendSCtxs.js) creates a directory, with a name according to configuration parameters, in [./SentTxs_Records/](https://github.com/SC-RP/SRP_PoC/tree/main/SRP_ethereum_docker/SentTxs_Records).  and stores records of sent Txs in a json file
    - Each Geth node running on docker-compose outputs time-records of received transactions in a JSON file, named with its public key address, that can be found in  [./SRP_Evaluation/](https://github.com/SC-RP/SRP_PoC/tree/main/SRP_ethereum_docker/SRP_Evaluation). 
* CSV Files
    - Each Geth node outputs a CSV file, named with its public key address, of received Selection Events that were emitted by the Selection Manger SC. The files are stored in [./SRP_SelectionEvents/](https://github.com/SC-RP/SRP_PoC/tree/main/SRP_ethereum_docker/SRP_SelectionEvents)
    - Final results of average times and throughput are stored in a CSV file in [./Results/](https://github.com/SC-RP/SRP_PoC/tree/main/SRP_ethereum_docker/Results)
* Log File
    - Nodes output log information that are stored in [./cmd/](https://github.com/SC-RP/SRP_PoC/tree/main/SRP_ethereum_docker/cmd)compose_output.log

