const Web3 = require("web3");
var web3 = new Web3('ws://localhost:8546');
const fs = require('fs');
const path = require('path');


// JSON file's path in which addresses of parameters to be used in experiments
const MainDir = getPath()
const parametersFilePath = MainDir+'files/parameter_configuration.json'; 
const resultsFilePath = MainDir+'Results/results.csv';
const NodesRecordDirectory = MainDir+"SRP_Evaluation/"; // directory of records produced by nodes (docker containers)
const NUM_OF_BLOCKS_MINED_PRIOR =3


//---- extract parameters from parameters_configuration file
try {
    const jsonString = fs.readFileSync(parametersFilePath)
    var parameters = JSON.parse(jsonString)
  } catch(err) {
    console.log(err)
    return
  }
  const victimAddressesFilePath = MainDir+parameters.registeredSCs_FilePath;
  const sentTxsDir = MainDir+parameters.sentTxsRecords_Directory;
  const experimentID = parameters.experiment_id;
  // sending parameters
  const TotalNumOfTxs = parameters.numOftxs;
  const RATE = parameters.rate; // number of txs per second
  
 

//----- Start Collecting
var blocks = [];
var totalNumberOfTxs = 0;
var numberOfConfirmedTxs = 0;


// extract addresses of victim contracts from the JSON file & store them in the array 'victimAddresses'
try {
    const jsonString = fs.readFileSync(victimAddressesFilePath)
    var victimAddresses = JSON.parse(jsonString)
} catch(err) {
    console.log(err)
    return
}
// get number of senders
var numberOfSenders = victimAddresses.length;
var numberOfConfirmedTXsList = new Array(numberOfSenders);
var numberOfNodes = parameters.numOfNodes;

//compute number of nodes within a subset and get the factor to compute the total success rate
var ValidatorSubsetSIZE = parseInt(numberOfNodes)/parseInt(numberOfSenders);
var NumberOFShouldBeSentTxs_Factor = ValidatorSubsetSIZE + 1;

// ---------------
// create directory for records of sent txs 
var dirName = numberOfNodes+'Nodes_'+numberOfSenders+'Victims_'+RATE+'tps_'+TotalNumOfTxs+'txs_'+parameters.rv_time+'msRV_'
var sentTxsDirectory = sentTxsDir+dirName+experimentID+'/';
console.log("CollectBlock:::SentTxs Records Directory = "+sentTxsDirectory);

// get cmd arguments
cmdArgs = process.argv.slice(2);
//senderIndex = cmdArgs[0];

// get all past blocks up to latest block number then update tx records
getBlocks().then(updateConfTime_ofSentTxRecords);


///////// ///////// /////// /// /// Functions // // // // / / / / / / / / / / // / 

async function getBlocks(){
    // get latest published block number
    var latestblock = await web3.eth.getBlockNumber();
    numberOfBlocks = latestblock;
    console.log("\nCollectBlock:::latest block number is "+latestblock+"\n");
    // iterate over all block numbers starting from 3 (since first 3 blocks include genesis initialization & registering txs)
    for (var i = NUM_OF_BLOCKS_MINED_PRIOR; i <= latestblock; i++){
        //index = i + blocks[0].number;
        await web3.eth.getBlock(i, function(error, block){
            if(!error){
                //if (block.gasUsed>=0){
                // check if block is not empty
                if (block.transactions.length>=0){
                    console.log("CollectBlock:::block#"+block.number+" has "+block.transactions.length+" txs");
                    blocks.push({
                        number: block.number,
                        numberOfTxs: block.transactions.length,
                        txs: block.transactions,
                        timestamp: block.timestamp,
                    });
                }
                //}
            }  
        });
    }
    
    
}


// updates records of sent txs with confirmation time extracted from blocks
function updateConfTime_ofSentTxRecords() {  
    console.log("CollectBlock:::Number of extracted blocks = "+blocks.length);
    
    if( blocks.length > 0 ){
        console.log("\nCollectBlock:::updating "+numberOfSenders+" lists of sentTxs Records. . .\n");
        var sentTxs = extractRecords();
        if (sentTxs.length > 0){
        // iterate over number of senders
        for (var i = 0; i < numberOfSenders; i++){
            console.log("CollectBlock:::Updating conf_time of "+i+"th record list with "+sentTxs[i].length+" txs sent to victimAddress "+victimAddresses[i]);
            // iterate over i-th list of sentTx records
            for (var t = 0; t < sentTxs[i].length; t++){
                totalNumberOfTxs++; 
                // iterate over list of blocks
                for (var j = 0; j < blocks.length; j++) {
                    // iterate over list of txs in jth block
                    for (var k = 0; k < blocks[j].numberOfTxs; k++){
                        // update sentTx record if sentTx is found in current block then break this loop
                        if (sentTxs[i][t].hash == blocks[j].txs[k]) {
                            sentTxs[i][t].conf_time = blocks[j].timestamp*1000; // multiply by 1000 to convert it from secs into millisecs
                            numberOfConfirmedTXsList[i]++;
                            numberOfConfirmedTxs++;
                            break;
                        }  
                    }
                }
            }
            console.log("CollectBlock:::-- Number of confirmed txs for this sender = "+numberOfConfirmedTXsList[i]+", success-rate="+(numberOfConfirmedTXsList[i]/sentTxs[i].length*100).toFixed(2)+"%\n");
        }
    }
    } else {
        console.log("CollectBlock:::-- No blocks!!. . . Exiting "+"\n");
        process.exit(0);

    }
   
    console.log("\nCollectBlock:::Saving "+totalNumberOfTxs+" records of txs sent by "+numberOfSenders+" senders");
    // save updated records
    saveSentToFile(sentTxs);
    // get performance measures
    var blockMeasures = getBlockThroughput();
    let measures = getPerformanceMeasures(sentTxs);
    let block_tps = blockMeasures[0].toFixed(3);
    let sentDuration_tps = computeTPS(sentTxs).toFixed(3);
    let total_SuccessRate = blockMeasures[3]/(numberOfSenders*TotalNumOfTxs*NumberOFShouldBeSentTxs_Factor)*100
    console.log("\nCollectBlock:::TOTAL NUMBER OF CONFIRMED TXS =  "+numberOfConfirmedTxs+", Unconfirmed "+(totalNumberOfTxs-numberOfConfirmedTxs)+" sent txs\n");
    console.log("CollectBlock:::TPS via blocks-measurement is "+block_tps);
    console.log("CollectBlock:::TPS via sending-duration is "+sentDuration_tps);
    console.log("\nCollectBlock:::\nTotal Succes-Rate is "+total_SuccessRate+"%");
    // append results to csv file
    results = readCSV();
    results.push({
    experimentNumber: parameters.experiment_id,
    rv_time: parameters.rv_time,
    numOfNodes: parameters.numOfNodes,
    numOfVictims: numberOfSenders,
    numOfTxsPerVictim: TotalNumOfTxs,
    sendingRate: RATE,
    waitingTime: measures[0],
    blockTPS: block_tps,
    durationTPS: sentDuration_tps,
    successRate: measures[1],
    totalSuccessRate: total_SuccessRate,
    numOfConfirmedTxs: numberOfConfirmedTxs,
    numOfUnconfirmedTxs: (TotalNumOfTxs*numberOfSenders - numberOfConfirmedTxs),
    numberOfMinedTxs: blockMeasures[3],
    numberOfBlocks: blocks.length,
    avgBlockSize: blockMeasures[1].toFixed(3),
    avgBlockTime: blockMeasures[2].toFixed(3),
    });

    exportCSVFile(results);
    moveNodeRecords();
    //if (blocks.length == numberOfBlocks+1){
    console.log("CollectBlock:::Got "+blocks.length+" blocks up to latest block & saved updated records, Exiting . . ")
    process.exit(0)
    //}
}

//returns throughput obtained from blocks, measured as (avg block size)/(avg inter-block time)
function getBlockThroughput(){
    var numberOfBlocks = blocks.length;
    var avgInterBlockTime = 0.0; // average time between blocks
    var avgBlockSize = blocks[0].numberOfTxs; // as average number of txs included in a block
    var tps = 0.0;
    //console.log("CollectBlock:::First block's size = "+ avgBlockSize)
    for (var j=numberOfBlocks-1; j > 0; j--){
        avgInterBlockTime += (blocks[j].timestamp - blocks[j-1].timestamp);
        avgBlockSize += blocks[j].numberOfTxs;
        //console.log("CollectBlock:::Avg block's size = "+ avgBlockSize)
    }
    let numberOfMinedTxs = avgBlockSize;
    avgBlockSize = avgBlockSize/numberOfBlocks;
    avgInterBlockTime = avgInterBlockTime/numberOfBlocks;
    // measure throughput
    console.log("\nCollectBlock:::Avg block's size = "+ avgBlockSize)
    console.log("CollectBlock:::Avg inter-block time = "+ avgInterBlockTime)
    console.log("\nCollectBlock:::Total Number of Blocks = "+blocks.length)
    console.log("CollectBlock:::Total Number Of mined txs = "+numberOfMinedTxs)
    tps = avgBlockSize/avgInterBlockTime;

    return [tps, avgBlockSize, avgInterBlockTime, numberOfMinedTxs];
}

//returns throughput obtained from sent txs, measured as (number of confirmed txs)/(last_block's timestamp - first_sent_tx's timestamp)
function computeTPS(allTxs) {
    let duration = blocks[blocks.length-1].timestamp- allTxs[0][0].send_time/1000;
    // console.log("CollectBlock:::-- Last Block's Timestamp = "+(blocks[blocks.length-1].timestamp)+" secs");
    // console.log("CollectBlock:::-- First sent Tx's Sending Timestamp = "+(allTxs[0][0].send_time/1000)+" secs");
    // console.log("CollectBlock:::-- Duration = "+duration+" secs");
    return numberOfConfirmedTxs/duration;
}

// returns an array of arrays that stores list of (latest) sent txs extracted from given sentTxsDirectory
function extractRecords(){
    // get list of files in which sent tx records are stored 
    var listOfFiles = getListofFiles(sentTxsDirectory);
    // define array of arrays for all txs, each inner array contains records of one sender
    var allTxs = new Array(numberOfSenders);
    // extract sent txs and store them in an array of arrays
    for (var i=0; i < numberOfSenders; i++){
        for (var j=0; j < listOfFiles.length; j++){
            let txsFileName = sentTxsDirectory+listOfFiles[j];
            console.log("\n***CollectBlock:::current file is "+ txsFileName)
            // only read non-empty json files
            if (fs.statSync(txsFileName).size > 0 && path.extname(txsFileName) == '.json') {
                // only extract files sent to latest victim sc addresses
                if (txsFileName == sentTxsDirectory+'txs_'+i+'_'+victimAddresses[i].toLowerCase()+'.json')
                console.log("\n***CollectBlock:::extracting records from file "+ txsFileName)
                let rawdata = fs.readFileSync(txsFileName);  
                allTxs[i] = JSON.parse(rawdata);
                // initialize counters of confirmed txs
                numberOfConfirmedTXsList[i] = 0;

            }
        }
    }
    return allTxs;
}

// returns list of files (as names) found in provided directory
function getListofFiles(Direc) {
    // list of files
    FilesList = [];
    //requiring path and fs modules
    const path = require('path');
    const fs = require('fs');
    //joining path of directory 
    const directoryPath = Direc;
    console.log("CollectBlock:::Listing All files found in: ", Direc);
    let i = 0;
    fs.readdirSync(directoryPath).forEach(file => {
        //let isFile = fs.statSync(sentTxsDirectory+file).isFile(); --> raises an error: Error: ENOENT: no such file or directory, stat 'nodeRecords'
    if (file != "nodesRecords"){
      FilesList.push(file)
      console.log("CollectBlock::: file #" +i+ " is " + file); 
      i++
    }
    });
    console.log();
  
    return FilesList;
  }

  // saves updated records of sent txs into their files by overwriting them
  // save sent txs as records to JSON files, each saved with a different name based on worker's index & victim's address
function saveSentToFile(sentTxs) { 
    var fs = require('fs');
    // extract records of sent txs from double array and save them to corresponding files
    for (var i=0; i<victimAddresses.length; i++){
        var stringify = JSON.stringify(sentTxs[i]);
        fs.writeFileSync(sentTxsDirectory+'txs_'+i+'_'+victimAddresses[i].toLowerCase()+'.json', stringify , 'utf-8'); 
    }
    
}


function getPerformanceMeasures(allTxs) {
    // average inter-sending of all txs to blockchain 
    var averageInterSending = 0.0
    // average waiting time of all txs (= confirmationTime - sendingTime)
    var averageWaitingTime = 0.0

    var totalNumOfTxs = 0
    var waitingTimeSum = 0.0
    var intersendingTimeSum = 0.0
    var intersendingCtr = 0
    var waitingCtr = 0
    var numberOfUnconfirmedTxs = 0
    // exract all Txs and store them in an array of arrays
    for (var i = 0; i < allTxs.length; i++) {
        
        let sentTxs = allTxs[i];
        let numOfTxs = sentTxs.length;
        console.log("CollectBlock:::worker"+i+": Sent "+numOfTxs+" txs to victimAddress "+victimAddresses[i]);

        // compute average waiting time & average intersending time
        var s = numOfTxs - 1//index for intersending
        for (var j=0; j < numOfTxs; j++){
            if (sentTxs[j].conf_time == -111){
            numberOfUnconfirmedTxs++;
            } else {
            let wt = sentTxs[j].conf_time-sentTxs[j].send_time;
           // console.log("CollectBlock:::Tx["+i+"]["+j+"] --> conf=" + sentTxs[j].conf_time+" - send="+sentTxs[j].send_time+"= "+wt);
            if (wt > 0){
                waitingTimeSum = waitingTimeSum + (wt);
                waitingCtr++;
            }
            }
            
            if (s > 0){
            intersendingTimeSum = intersendingTimeSum + (sentTxs[s].send_time - sentTxs[s-1].send_time);
            intersendingCtr++;
            }
            s--;
        }
        totalNumOfTxs = totalNumOfTxs + numOfTxs;
        allTxs[i] = sentTxs;
    }


    averageWaitingTime = waitingTimeSum/waitingCtr;
    averageInterSending = intersendingTimeSum/intersendingCtr;


    // average inter-sending time between workers
    var averageParallelInterSending = 0.0
    var intersendingParallelTimeSum = 0.0
    var parallelCtr = 0

    // get the minimum length of arrays
    let minimumLength = allTxs[0].length;
    let minimIndex = 0;
    console.log("CollectBlock:::length of allTxs = "+allTxs.length)
    for (var j=0; j < allTxs.length; j++){
    //console.log("CollectBlock:::Record ["+j+"] of victim "+victimAddresses[j]+" has a length of "+allTxs[j].length)
    if (allTxs[j].length < minimumLength && allTxs[j].length != 0) {
        minimumLength = allTxs[j].length
        minimIndex = j
    }
    if (allTxs[j].length == 0){
        console.log("CollectBlock:::Removing an empty record of txs sent to victim address " + victimAddresses[j])
        allTxs.splice(j,1);
        j--;
    }
    }
    if (minimumLength == 0) {
    console.log("\nCollectBlock:::There are zero records of sent txs")
    // } else {
    //   console.log("\nCollectBlock:::Minimum number of txs record is "+minimumLength+" of victim "+victimAddresses[minimIndex])
    }
    // get intersending time between workers (i.e. parallelisim)
    for (var i=0; i < minimumLength; i++){
    for (var j=allTxs.length-1; j > 0; j--){
        intersendingParallelTimeSum = intersendingParallelTimeSum + Math.abs(allTxs[j][i].send_time-allTxs[j-1][i].send_time);
        parallelCtr++;
        // if (i==0){
        //   console.log("CollectBlock:::tx["+j+"]["+i+"] hash = "+allTxs[j][i].hash)
        //   console.log("CollectBlock::: tx["+(j-1)+"]["+i+"] hash = "+allTxs[j-1][i].hash)
        //   console.log("CollectBlock:::Inter-sending Time = tx["+j+"]["+i+"]:"+allTxs[j][i].send_time+" - "+ "tx["+(j-1)+"]["+i+"]:"+allTxs[j-1][i].send_time+" = "+Math.abs(allTxs[j][i].send_time-allTxs[j-1][i].send_time))
        // }
    }
    }
    console.log()
    averageParallelInterSending = intersendingParallelTimeSum/parallelCtr;

    let successRate = (((numberOfConfirmedTxs)/totalNumOfTxs)*100).toFixed(2);
    console.log("CollectBlock:::***Measurements computed for "+totalNumOfTxs+" txs sent by "+numberOfSenders+" validators***\n")
    console.log("CollectBlock:::-- Avg waiting time is "+averageWaitingTime.toFixed(3)+" ms for "+waitingCtr+" values,")
    console.log((numberOfConfirmedTxs)+" txs are confirmed of which "+(numberOfConfirmedTxs-waitingCtr)+" have negative waiting time, \nthe remaining "+numberOfUnconfirmedTxs+" txs are unconfirmed")
    console.log("\nCollectBlock:::-- Confirmation Success-Rate = "+successRate+"%\n")
    console.log("CollectBlock:::-- Avg intersending time is "+averageInterSending.toFixed(3)+" ms for "+intersendingCtr+" values")
    console.log("CollectBlock:::-- Avg parallel intersending time is "+averageParallelInterSending+" ms for "+parallelCtr+" values\n")

    return [averageWaitingTime.toFixed(3), successRate]
}

function exportCSVFile(values) {
    var headers = {
        experimentNumber: 'Experiment',
        rv_time: 'RVTime_(ms)',
        numOfNodes:'Nodes',
        numOfVictims: 'Senders',
        numOfTxsPerVictim: 'NumberOfTxs_perVictim',
        sendingRate: 'SendingRate',
        waitingTime: 'WaitingTime_(ms)',
        blockTPS: 'Block_tps',
        durationTPS: 'SentDuration_tps',
        successRate: 'SuccessRate',
        totalSuccessRate: 'TotalSuccessRate',
        numOfConfirmedTxs: 'ConfirmedTxs',
        numOfUnconfirmedTxs: 'UnconfirmedTxs',
        numberOfMinedTxs: 'TotalNumberOfMinedTxs',
        numberOfBlocks: 'Blocks',
        avgBlockSize: 'BlockSize_(txs)',
        avgBlockTime: 'BlockTime_(s)',
    };
    if (headers) {
        values.unshift(headers);
    }

    // Convert Object to JSON
    var jsonObject = JSON.stringify(values);

    var csv = convertToCSV(jsonObject);

    fs.writeFileSync(resultsFilePath, csv);

}

function convertToCSV(objArray) {
    var array = typeof objArray != 'object' ? JSON.parse(objArray) : objArray;
    var str = '';

    for (var i = 0; i < array.length; i++) {
        var line = '';
        for (var index in array[i]) {
            if (line != '') line += ','

            line += array[i][index];
        }

        str += line + '\r\n';
    }

    return str;
}

//var csv is the CSV file with headers
function readCSV(){
    csv = fs.readFileSync(resultsFilePath, 'utf8');
    var lines=csv.split("\n");
  
    const result = []
    const headers = lines[0].split(',')

    for (let i = 1; i < lines.length; i++) {        
        if (!lines[i])
            continue
        const obj = {}
        const currentline = lines[i].split(',')

        for (let j = 0; j < headers.length; j++) {
            obj[headers[j]] = currentline[j]
        }
        result.push(obj)
    }
    return result;
  
  }


function moveNodeRecords(){
    console.log("Moving File from SCG_Evaluation to nodesRecords ");
    
    var filesList = getListofFiles(NodesRecordDirectory);
    var destinationPath = sentTxsDirectory+'nodesRecords/';
    for (file of filesList){
        let pathToFile = NodesRecordDirectory + file;
        try {   fs.copyFileSync(pathToFile, destinationPath+file)
                console.log("CollectBlock:::Successfully copied and moved "+file+" to "+destinationPath);
            } catch(err) {
                throw err
            }
    }
    
}

function getPath(){
    mainDir = __dirname.substring(0, __dirname.lastIndexOf('/')); // remove runSRP directory
    mainDir = mainDir.substring(0, mainDir.lastIndexOf('/'));// remove cmd directory
    return mainDir + '/'
    
    }