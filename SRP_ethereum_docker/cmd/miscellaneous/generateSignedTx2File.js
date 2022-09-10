//-- Module Imports
const performance = require('perf_hooks').performance;
const Web3 = require('web3')
const fs = require('fs')
const keythereum = require("keythereum"); // used to import keys
const { numberToHex } = require('web3-utils');
const provider = 'http://localhost:8545';

const 	MainDir = getPath()
// Json file's path in which addresses of victim contracts are stored
const victimAddressesFilePath = MainDir + 'files/victimAddresses.json';
// Directory path in which signed txs are stored
const signedTxsDirectory = MainDir + 'files/signedTxs/';
//const keysdir = "./files/generatedKeystore_keystore";
const keysdir = MainDir + 'files/etherTransferKeystore_keystore';

// connect to node via http using web provider
var web3 = new Web3(new Web3.providers.HttpProvider(provider))

const NUM_OF_TX = 2000;//number of accounts from which we send Txs, Max 12.000 accounts



const abi = [{"inputs":[],"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"string","name":"_name","type":"string"},{"indexed":false,"internalType":"uint256","name":"_numOfNames","type":"uint256"}],"name":"RegisteredNames","type":"event"},{"inputs":[{"internalType":"uint256","name":"","type":"uint256"}],"name":"Names","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"string","name":"_name","type":"string"}],"name":"addName","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"getNamesCount","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"}]

// exract addresses of victim contracts from the JSON file & store them in the array 'victimAddresses'
try {
    const jsonString = fs.readFileSync(victimAddressesFilePath)
    var victimAddresses = JSON.parse(jsonString)
  } catch(err) {
    console.log(err)
    return
}
// print extract addresses of the victim contracts
numOfVictims = victimAddresses.length
console.log('Extracted ' +numOfVictims+' victim addresses')
for (i = 0; i < victimAddresses.length; i++) {
    console.log('address['+i+'] = '+ victimAddresses[i])
}

invokeNameCollection(victimAddresses);




//var signedTxs = [];

// loop for number of exracted victim addresses (currently, its 10-1 = 9)
//for (var victimIndex = 0; victimIndex < victimAddresses.length - 1; victimIndex++){
// for (const victimAddr of victimAddresses) {
    
//     // get current victim contract's address
//     //victimAddr = victimAddresses[victimIndex];
//     // create victim contract instance 
//     victimContract = new web3.eth.Contract(abi, victimAddr);
//     // create, sign, and send registeration transactions
//     console.log("Creating & Sigining "+NUM_OF_TX+" txs invoking contract["+victimIndex+"]= "+victimAddr)
//     invokeNameCollection(victimContract, victimAddr, victimIndex);//.then(saveSignedToFile(victimAddr));
   
//     victimIndex++;
// }
/**
 * Creates, signs, & sends registeration txs to SelectionManager contract
 * for each account found in keystore, where private keys are retrieved 
 * using 
 * @param victimContract - Web3 Object of the deployed contract.
 */
async function invokeNameCollection(victimAddresses) {
    var printed = false;
    // get addresses of accounst used to sign txs for victim contracts
    var addresses = getAddresses();
    console.log('Extracted ' +addresses.length+' account addresses for sending txs')
    var keyObjects = new Array(addresses.length);//[];
    var privateKeys = new Array(addresses.length);//[];
    var buffer = Buffer.from("", "utf-8");
    // list of signed txs sent to vicitimAddress & saved in JSON file 
    var signedTxs = [];
    // compute index based on
    //startIndex = victimIndex*NUM_OF_TX
    //endIndex = startIndex+NUM_OF_TX
    //console.log("StartIndex = "+ startIndex+" & Ending at = "+ endIndex + "-1")
    // if starting index exceeds the number of accounts in generated keystore then exit
    // if ((startIndex + NUM_OF_TX) >= addresses.length) {
    //     console.log("Number of accounts ("+addresses.length+") is insufficient !!")
    //     return
    // }
    var victimIndex = 0;
    var numbOfVictimContracts = victimAddresses.length;
    // timestamp t0
    var t0 = performance.now()
    console.log('There are ' +numbOfVictimContracts+' victim contracts')
    console.log('Signing ' +NUM_OF_TX+' txs for every SC account')
    console.log('Total of ' +NUM_OF_TX*numbOfVictimContracts+' txs to be signed')
    // import accounts as objects, for each retrieved account sign txs & store them
    for (var i = 0; i < NUM_OF_TX*numbOfVictimContracts; i++) {
        // get current victim contract's address
        victimAddress = victimAddresses[victimIndex];
        //let accountIndex = i+startIndex;
       

    if (victimAddress.toLowerCase() != "0xff37a57b8d373518abe222db1077ed9a968a5fd5" && victimAddress.toLowerCase() != "0xff37a57b8d373518abe222db1077ed9a968a5fd4" && victimAddress.toLowerCase() != "0xff37a57b8d373518abe222db1077ed9a968a5fd7" ) {
        // create victim contract instance 
        victimContract = new web3.eth.Contract(abi, victimAddress);
        keyObjects[i] = keythereum.importFromFile(addresses[i], keysdir);
        console.log("TX %d: %s", i, keyObjects[i].address);
        // get the private key of current account, note: keythereum returns buffer not hex string
        privateKeys[i] = keythereum.recover(buffer, keyObjects[i]).toString('hex');
        console.log("PrivateKey: ", privateKeys[i]);
        // get the latest tx nonce of current account using a helper function
        //let txNonce = getNonce(keyObjects[i].address);
        console.log("Constructing Tx . . .");
        // create and sign transactions
        let tx_builder = victimContract.methods.addName("HBKU");
        let encoded_tx = tx_builder.encodeABI();
        let transactionObject = {
            gas: 4000000,
            data: encoded_tx,
            from: keyObjects[i].address,
            to: victimAddress,
            //nonce: txNonce
        };
        // time = new Date().getTime();
        console.log("Signing Tx . . .");
        
        let signedTx = await web3.eth.accounts.signTransaction(transactionObject, privateKeys[i])
        signedTxs.push(signedTx)
        console.log("SignedTx["+i+"]= ")
        } else {
            if(printed == false){
                console.log("i="+i+": Skipping since it is already signed");
                printed = true;
            }
        }
        // update victimIndex and save txs to JSON file named with current victimAddress 
        if ((i+1)%NUM_OF_TX == 0 && i !==0){ //|| i == NUM_OF_TX*numbOfVictimContracts -1){
            console.log("From VictimAddress["+ victimIndex+ "] ="+victimAddress);
            victimIndex = (victimIndex + 1)%numbOfVictimContracts;
            console.log("To VictimAddress["+ victimIndex+ "] ="+victimAddresses[victimIndex]);
            console.log("Incremented victimIndex = "+ victimIndex+ ", numOfVictims = "+numbOfVictimContracts);
            // save signed txs to a JSON file named with the victim's address
            //console.log("SignedTxs[] = ", signedTxs)
            //signed.rawTransaction
            if (signedTxs.length > 0){
                saveSignedToFile(victimAddress, signedTxs);
            }
                
            //clear up list of txs
            signedTxs = [];
            printed = false;
        }
        
        

    }

    var t1 = performance.now()
    exec_time = t1-t0;
    console.log('There are ' +numbOfVictimContracts+' victim contracts')
    console.log('Signed ' +NUM_OF_TX+' txs for every SC account')
    console.log('Total of ' +NUM_OF_TX*numbOfVictimContracts+' txs were signed')
    console.log('Extracted ' +addresses.length+' account addresses for sending txs')
    console.log("The generation "+NUM_OF_TX+" Signed Txs Took " + exec_time + " milliseconds.")
    console.log("In minutes: "+millisToMinutesAndSeconds(exec_time))

}

function saveSignedToFile(scAddress, signedTxs) {
    // get name file based on victim's address
    var outFileName = signedTxsDirectory+'signedTxs_'+scAddress.toLowerCase();
    // read JSON file to append the array of signed txs to list of already saved signed txs
    let already_signedtxsJSON = fs.readFileSync(outFileName+'.json','utf-8');
    // transform JSON string to array
    let already_signedtxs = JSON.parse(already_signedtxsJSON);
    console.log("Read "+already_signedtxs.length+" already signed txs from "+outFileName)
    // concatenate array of the new signed txs with already signed txs
    let all_signedtxs = already_signedtxs.concat(signedTxs)
    // transform array to JSON string
    var stringify = JSON.stringify(all_signedtxs);
    console.log("Saving additional "+signedTxs.length+" signed txs. . .")
    console.log("Saving total of "+all_signedtxs.length+" signed txs to "+outFileName)
    // write JSON string to file
    fs.writeFileSync(outFileName+'.json', stringify , 'utf-8'); 
}

/************* Helper Functions **************/
function getAddresses() {
    var addresses = [];
    //requiring path and fs modules
    const fs = require('fs');
    //joining path of directory 
    //const directoryPath = path.join(__dirname, 'keystore');
    const directoryPath = keysdir + '/keystore'

    //passing directoryPath to read the directory synchronously
    files = fs.readdirSync(directoryPath);

    //listing all files using forEach
    files.forEach(function (file) {
        // get the raw data of each file and extract the address field 
        let rawdata = fs.readFileSync(directoryPath + `/${file}`);
        let account = JSON.parse(rawdata);
        // store the extracted address in the array
        addresses.push(account.address);
    });
    return addresses;
}

function millisToMinutesAndSeconds(millis) {
    var minutes = Math.floor(millis / 60000);
    var seconds = ((millis % 60000) / 1000).toFixed(0);
    return minutes + ":" + (seconds < 10 ? '0' : '') + seconds;
}




function getPath(){
    mainDir = __dirname.substring(0, __dirname.lastIndexOf('/')); // remove runSRP directory
    mainDir = mainDir.substring(0, mainDir.lastIndexOf('/'));// remove cmd directory
    return mainDir + '/'
    
    }