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
#GETH_OPTS=${@/XXX/$BOOTSTRAP_IP}
ID=$(python3 /root/setMiner.py)
#cd /root/go/bin/
./build/bin/geth --datadir=~/.ethereum/devchain --rpccorsdomain="*" --networkid=456719 --rpc --bootnodes="enode://288b97262895b1c7ec61cf314c2e2004407d0a5dc77566877aad1f2a36659c8b698f4b56fd06c4a0c0bf007b4cfb3e7122d907da3b005fa90e724441902eb19e@$BOOTSTRAP_IP:30303" --allow-insecure-unlock --etherbase $ID 
#--mine --minerthreads=1
#./root/go/src/github.com/ethereum/go-ethereum/build/bin/geth attach ipc://root/.ethereum/devchain/geth.ipc