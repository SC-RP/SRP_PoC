#!/bin/bash
parent_path=$( cd "$(dirname "${BASH_SOURCE[0]}")" ; pwd -P )

cd "$parent_path"
# set parameters
rv_time=1000
numOfNodes=3
numOfVictims=3
rate=5
numOftxs=100
experiment_id=1

numberOfEths=$(( $numOfNodes - 1))
init_timeout=2m

# kill docker compose & other backgrounds scripts/programs if any
docker-compose down
echo "Killing initSRP" 
pkill -f -9 initSRP

 echo "BASH: Updated experiment_$experiment_id & numOfVictims $numOfVictims & RV $rv_time microseconds" 
node changeParameters.js $rv_time $numOfNodes $numOfVictims $rate $numOftxs $experiment_id \
&& docker-compose up --build --scale eth=$numberOfEths > compose_output.log 2>&1 & # run docker-compose in background
echo "Sleeping for $init_timeout to allow nodes' initialization. ."
sleep $init_timeout 
# run scripts sequentially
go run ./runSRP/initSRP.go -release \
&& nodejs ./runSRP/sendSCtxs.js -release \
&& nodejs ./collectMeasurements/getTxpoolContent.js -release \
&& nodejs ./collectMeasurements/collectSentTxsFromBlocks.js  -release

docker_timeout=1m
echo "Sleeping for $docker_timeout and then shutting docker down . ."
sleep $docker_timeout 
docker-compose down


# to show compose log file on terminal run:     
    # tail -f compose_output.log 
