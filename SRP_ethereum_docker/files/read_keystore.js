
function getAddresses() {
    var addresses = [];
    //requiring path and fs modules
    const path = require('path');
    const fs = require('fs');
    //joining path of directory 
    const directoryPath = path.join(__dirname, 'keystore');
   
    //passsing directoryPath to read the directory synchronously
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

var addresses = getAddresses();
// print addresses
for (var i = 0; i < addresses.length; i++) {
    console.log("\nAccount %d: %s", i, addresses[i]);
}