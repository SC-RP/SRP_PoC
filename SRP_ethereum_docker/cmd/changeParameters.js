

// JSON file's path in which addresses of parameters to be used in experiments
const 	MainDir = getPath()
const parametersFilePath = MainDir + 'files/parameter_configuration.json'; 

const fs = require('fs')


// get parameters, which are passed as a cmd argument (order of paramters is like in the file)
cmdArgs = process.argv.slice(2);

//---- extract parameters from parameters_configuration file
try {
    const jsonString = fs.readFileSync(parametersFilePath)
    var parameters = JSON.parse(jsonString)
  } catch(err) {
    console.log(err)
    return
}


//---- update parameters_config.json file with passed parameters (multiply by 1 to convert them into integers)
// Structure: RV_Time numOfNodes numOfVictims sendingRate numOfTxs experimentID
parameters.rv_time = cmdArgs[0]* 1;
parameters.numOfNodes = cmdArgs[1]* 1;
parameters.numOfVictims = cmdArgs[2]* 1;
parameters.rate = cmdArgs[3]* 1;
parameters.numOftxs = cmdArgs[4]* 1;
parameters.experiment_id = cmdArgs[5]* 1;



// write updated json array into JSON file 
var stringify = JSON.stringify(parameters,null, "  ");
fs.writeFileSync(parametersFilePath, stringify , 'utf-8'); 





function getPath(){
  mainDir = __dirname.substring(0, __dirname.lastIndexOf('/')); // remove cmd directory
  return mainDir + '/'
  
  }