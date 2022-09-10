//-- Module Imports
const Web3 = require('web3')
const fs = require('fs');
const { mainModule } = require('process');

var web3 = new Web3(new Web3.providers.WebsocketProvider('ws://localhost:8546'));
const SLEEP = true

// JSON file's path in which addresses of parameters to be used in sending txs
const 	MainDir = getPath()
const parametersFilePath = MainDir + 'files/parameter_configuration.json'; 
// get number of txs to be sent by each virtual worker, which is passed as a cmd argument
cmdArgs = process.argv.slice(2);

const start_index = 0;
const TIMEOUT = 30*1000;

//---- extract parameters from parameters_configuration file
try {
    const jsonString = fs.readFileSync(parametersFilePath)
    var parameters = JSON.parse(jsonString)
} catch(err) {
    console.log(err)
    return
}
const victimAddressesFilePath = MainDir+parameters.registeredSCs_FilePath;
const signedTxsDirectory = MainDir+parameters.signedTxs_Directory;
const sentTxsDir = MainDir+parameters.sentTxsRecords_Directory;
const experimentID = parameters.experiment_id;
// sending parameters
const TotalNumOfTxs = parameters.numOftxs;
const RATE = parameters.rate; // number of txs per second
const delay = (1/RATE) *1000 //multiply by 1000 to convert from seconds to milliseconds


// extract addresses of victim contracts from the JSON file & store them in the array 'victimAddresses'
try {
    const jsonString = fs.readFileSync(victimAddressesFilePath)
    var victimAddresses = JSON.parse(jsonString)
} catch(err) {
    console.log(err)
    return
}

 var numberOfSenders = victimAddresses.length;
 var numberOfNodes = parameters.numOfNodes;
// ---------------
// create directory for records of sent txs 
var dirName = numberOfNodes+'Nodes_'+numberOfSenders+'Victims_'+RATE+'tps_'+TotalNumOfTxs+'txs_'+parameters.rv_time+'msRV_'
var sentTxsDirectory = sentTxsDir+dirName+experimentID+'/';
if (!fs.existsSync(sentTxsDirectory)){
    console.log("**SeqSending: Creating Directory: "+sentTxsDirectory)
    fs.mkdirSync(sentTxsDirectory);
    fs.mkdirSync(sentTxsDirectory+'nodesRecords/');
} else {
    let listOfFiles = getListofFiles(sentTxsDirectory);
    if (listOfFiles.length == 0){
        console.log("**SeqSending: SentTxs Directory: "+sentTxsDirectory+ " exists already but is Empty")
    } else if (cmdArgs.length >= 1){
        if (cmdArgs[0] == 'force'){
            for (const file of listOfFiles) {
                try {
                    fs.unlinkSync(sentTxsDirectory+file)
                    //file removed
                } catch(err) {
                    console.error(err)
                }
            }
            console.log("**SeqSending: NOTE: SentTxs Directory: "+sentTxsDirectory+ " exists already, stored files are REPLACED")
        }
    } else {
        console.log("**SeqSending: SentTxs Directory: "+sentTxsDirectory+ " exists already...rename it, exiting! ")
        process.exit(0)
    }
    
}
//-------------------------------------------
// create an array of signed txs arrays
var signedTxs = new Array(numberOfSenders);
// create an array of sent txs arrays
var sentTxs = [];//new Array(numberOfSenders);

// extract signed txs and store them in an array of arrays
for (var i=0; i < numberOfSenders; i++){
    let signedTxFilePath = signedTxsDirectory+'signedTxs_'+victimAddresses[i].toLowerCase()+'.json'
    let rawdata = fs.readFileSync(signedTxFilePath);  
    signedTxs[i] = JSON.parse(rawdata);
    // initialize records of sent txs
    sentTxs[i] = [];
    console.log("**SeqSending: Extracted, for worker"+i+", "+signedTxs[i].length+" txs of which "+ TotalNumOfTxs+ " tx to be sent to victimAddress "+victimAddresses[i]);
}
console.log("**SeqSending: Sending "+TotalNumOfTxs*numberOfSenders+" txs to "+numberOfSenders+ " victim Addresses");

res = sendTransactions().then(saveSentToFile);
//-------------------------------------------
function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}



// sends txs that are saved as rawdata (in signedTx.json)
async function sendTransactions() { 
    if (SLEEP){
        console.log("**SeqSending: Sleeping for "+(TIMEOUT/1000)+" seconds to allow all nodes to interconnect first")
        await sleep(TIMEOUT);
    }
    // send the extracted txs
    for (var j=0; j < TotalNumOfTxs; j++) { 
        let start_time = new Date().getTime();
        for (var i=0; i < numberOfSenders; i++){
            console.log("**SeqSending: worker"+i+": Sending "+j+"th tx with hash::"+signedTxs[i][j+start_index].transactionHash+ "to victimAddress "+victimAddresses[i]);
            let send_time = new Date().getTime();
            web3.eth.sendSignedTransaction(signedTxs[i][j+start_index].rawTransaction)

            console.log("**SeqSending: Saving Record of tx "+signedTxs[i][j+start_index].transactionHash);
           
            // the confirmation time value below is unreliable
            sentTxs[i].push({
                hash: signedTxs[i][j+start_index].transactionHash,
                send_time: send_time,
                conf_time: -111,
                to: victimAddresses[i],
            });
        }
        // sleep after sending N=numberOfSenders txs
        if (new Date().getTime()-start_time <= delay) {
            consumedTime = new Date().getTime() - start_time;
            console.log("**SeqSending: Tx Batch #"+j+": Sleeping for "+(delay-consumedTime)+" ms after Sending "+ numberOfSenders+" txs ");
            await sleep(delay - consumedTime);
        }
    }

}

// save sent txs as records to JSON files, each saved with a different name based on worker's index & victim's address
function saveSentToFile() { 
    let numOfSentTxs = 0;
    var fs = require('fs');
    // extract records of sent txs from double array and save them to corresponding files
    for (var i=0; i<numberOfSenders; i++){
        var stringify = JSON.stringify(sentTxs[i]);
        console.log("**SeqSending: worker"+i+": Saving records of "+sentTxs[i].length+" txs sent to victimAddress "+victimAddresses[i])
        fs.writeFileSync(sentTxsDirectory+'txs_'+i+'_'+victimAddresses[i].toLowerCase()+'.json', stringify , 'utf-8'); 
        numOfSentTxs+= sentTxs[i].length;
    }
    if (numOfSentTxs == (TotalNumOfTxs*numberOfSenders)){
        console.log("**SeqSending: Sent total of "+numOfSentTxs+" txs & saved their records, Exiting . . ")
        process.exit(0)
    }
    
}


// returns list of files stored in given directory
function getListofFiles(sentTxsDirectory) {
    // list of files
    FilesList = [];
    const fs = require('fs');
    const directoryPath = sentTxsDirectory;
    console.log("**SeqSending: Listing All files found in: ", sentTxsDirectory);
    let i = 0;
    fs.readdirSync(directoryPath).forEach(file => {
      FilesList.push(file)
      console.log("**SeqSending:  file #" +i+ " is " + file); 
      i++;
    });
  
  
    return FilesList;
  }



  function getPath(){
    mainDir = __dirname.substring(0, __dirname.lastIndexOf('/')); // remove runSRP directory
    mainDir = mainDir.substring(0, mainDir.lastIndexOf('/'));// remove cmd directory
    return mainDir + '/'
    
    }