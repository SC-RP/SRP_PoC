#!/bin/bash
set -e
#cd /root/eth-net-intelligence-api
#perl -pi -e "s/XXX/$(hostname)/g" app.json
#/usr/bin/pm2 start ./app.json
sleep 3
cd /root/go/src/github.com/ethereum/go-ethereum/
#go install ./cmd/geth
make geth
#cd /root/go/bin/
./build/bin/geth --datadir=~/.ethereum/devchain init "/root/files/genesis.json"
sleep 3
BOOTSTRAP_IP=`getent hosts bootstrap | cut -d" " -f1`
GETH_OPTS=${@/XXX/$BOOTSTRAP_IP}
#ID=$(python3 /root/setMiner.py)
#geth $GETH_OPTS
#cd /root/go/bin/
./build/bin/geth --datadir=~/.ethereum/devchain --nodekeyhex=091bd6067cb4612df85d9c1ff85cc47f259ced4d4cd99816b14f35650f59c322 --rpcapi "db,personal,eth,net,web3,txpool" --rpccorsdomain="*" --networkid=456719 --rpc --rpcaddr="0.0.0.0" --ws --wsport "8546" --wsapi "db,eth,net,web3,personal,debug,txpool" --wsaddr "0.0.0.0" --wsorigins "*" --allow-insecure-unlock --unlock 0,1,2,3,4 --password /root/files/password/pass.txt --etherbase 0 --vmdebug --mine --minerthreads=1 --metrics --pprof --metrics.addr 0.0.0.0 --pprof.addr 0.0.0.0