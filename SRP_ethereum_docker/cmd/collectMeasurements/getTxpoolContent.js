const Web3 = require('web3')
var web3 = new Web3(new Web3.providers.WebsocketProvider('https://localhost:8546'));//Web3(new Web3.providers.HttpProvider(provider))

const TIMEOUT = 60*1000;
const SelectionManager_Addr = "0xDB571079aF66EDbB1a56d22809584d39C20001D9";

var start_time = new Date().getTime();
// Extend web3
web3.eth.extend({
    property: 'txpool',
    methods: [{
      name: 'content',
      call: 'txpool_content'
    },{
      name: 'inspect',
      call: 'txpool_inspect'
    },{
      name: 'status',
      call: 'txpool_status'
    }]
});

console.log("\nPoolContent:::*** Checking TxPool until it is empty then collecting records ****\n")
getTxpoolStatus();


async function getTxpoolStatus() {
    var end = true;
    var queuedTracker = [-1,-2,-3,-4,-5];
    var pendingTracker = [-1,-2,-3,-4,-5];
    let queueTimeoutCtr = 0;
    let pendingTimeoutCtr = 0;
    while (end){
    // check the txpool status
    txpoolstatus = await web3.eth.txpool.status();
    let pending = web3.utils.hexToNumber(txpoolstatus.pending)
    let queued = web3.utils.hexToNumber(txpoolstatus.queued)
    queuedTracker[queueTimeoutCtr] = queued;
    pendingTracker[pendingTimeoutCtr] = queued;
    console.log(`PoolContent:::Pending tx count in txpool: ${web3.utils.hexToNumber(txpoolstatus.pending)}`)
    console.log(`PoolContent:::Queued tx count in txpool: ${web3.utils.hexToNumber(txpoolstatus.queued)}`)
    if (pending == 0 && queued == 0){
        console.log("\nPoolContent:::-- No more Txs in pool...Exiting --\n")
        var consumedTime = new Date().getTime() - start_time;
        console.log("PoolContent:::Time consumed is "+(consumedTime/1000)+" seconds")
        break;
    } else {
        const allEqual = arr => arr.every( v => v === arr[0] )
        if (allEqual(queuedTracker) && pending == 0){
            console.log("\nPoolContent:::-- Number of Queued Txs did not change ("+queuedTracker.length+" loops)...Exiting --\n")
            var consumedTime = new Date().getTime() - start_time;
            console.log("PoolContent:::Time consumed is "+(consumedTime/1000)+" seconds")
            break;
        } else if (allEqual(pendingTracker)){
          console.log("\nPoolContent:::-- Number of PENDING Txs did not change ("+pendingTracker.length+" loops)...Exiting --\n")
          var consumedTime = new Date().getTime() - start_time;
          console.log("PoolContent:::Time consumed is "+(consumedTime/1000)+" seconds")
          break;
        }
    }
    queueTimeoutCtr = (queueTimeoutCtr+1)%(queuedTracker.length)
    pendingTimeoutCtr = (pendingTimeoutCtr+1)%(pendingTracker.length)
    await sleep(TIMEOUT);
    }
    process.exit(0);


}

function sleep(ms) {
    console.log("PoolContent:::Sleeping for "+(ms/1000)+" seconds")
    return new Promise(resolve => setTimeout(resolve, ms));
}

// var subscription = web3.eth.subscribe('pendingTransactions', async function(error, txHash){
//     if (!error)
//       console.log(`PoolContent:::New tx logged: ${txHash}`);

//       // check the txpool status
//       txpoolstatus = await web3.eth.txpool.status()
//       console.log(`PoolContent:::Pending tx count in txpool: ${web3.utils.hexToNumber(txpoolstatus.pending)}`)

//       // get all the pending transactions
//       pendingTransactions = await web3.eth.getPendingTransactions()

//     }
//   )

// web3.eth.extend({
//     property: 'txpool',
//     methods: [{
//       name: 'content',
//       call: 'txpool_content'
//     },{
//       name: 'inspect',
//       call: 'txpool_inspect'
//     },{
//       name: 'status',
//       call: 'txpool_status'
//     }]
//   });









// const subscription = web3.eth.subscribe('pendingTransactions', (err, txHash) => {
//     if (err) {
//         console.error(err)
//     } else {
//         console.log("PoolContent:::Result of pendingTxs Subscription: \n"+txHash);
//     }
// });

// subscription.on('data', (txHash) => {
// setTimeout(async () => {
// try {
// let tx = await web3.eth.getTransaction(txHash);
// if (tx && tx.to) {// This is the point you might be looking for to filter the address
// if (tx.to.toLowerCase() == SelectionManager_Addr.toLocaleLowerCase()) {
// console.log('TX hash: ',txHash ); // transaction hash
// console.log('TX To Address: ',tx.to ); // transaction hash
// console.log('TX confirmation: ',tx.transactionIndex ); // "null" when transaction is pending
// console.log('TX nonce: ',tx.nonce ); // number of transactions made by the sender prior to this one
// console.log('TX block hash: ',tx.blockHash ); // hash of the block where this transaction was in. "null" when transaction is pending
// console.log('TX block number: ',tx.blockNumber ); // number of the block where this transaction was in. "null" when transaction is pending
// console.log('TX sender address: ',tx.from ); // address of the sender
// console.log('TX amount(in Ether): ',web3.utils.fromWei(tx.value, 'ether')); // value transferred in ether
// console.log('TX date: ',new Date()); // transaction date
// console.log('TX gas price: ',tx.gasPrice ); // gas price provided by the sender in wei
// console.log('TX gas: ',tx.gas ); // gas provided by the sender.
// console.log('TX input: ',tx.input ); // the data sent along with the transaction.
// console.log('=====================================') // a visual separator
// }
// }
// } catch (err) {
// console.error(err);
// }
// })
// });