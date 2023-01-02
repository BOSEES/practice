#!/bin/bash
# Run this script to quickly install, setup, and run the current version of practice without docker.
# ./scripts/test_node.sh [clean|c]

KEY="practice1"
CHAINID="practice-t1"
MONIKER="localpractice"
KEYALGO="secp256k1"
KEYRING="test"
LOGLpracticeL="info"

practiced config keyring-backend $KEYRING
practiced config chain-id $CHAINID

command -v practiced > /dev/null 2>&1 || { echo >&2 "practiced command not found. Ensure this is setup / properly installed in your GOPATH."; exit 1; }
command -v jq > /dev/null 2>&1 || { echo >&2 "jq not installed. More info: https://stedolan.github.io/jq/download/"; exit 1; }

from_scratch () {

  make install

  # remove existing daemon.
  rm -rf ~/.practice/* 
  
  # practice1efd63aw40lxf3n4mhf7dzhjkr453axurv2zdzk
  echo "decorate bright ozone fork gallery riot bus exhaust worth way bone indoor calm squirrel merry zero scheme cotton until shop any excess stage laundry" | practiced keys add $KEY --keyring-backend $KEYRING --algo $KEYALGO --recover
  # practice1hj5fveer5cjtn4wd6wstzugjfdxzl0xps73ftl
  echo "wealth flavor believe regret funny network recall kiss grape useless pepper cram hint member few certain unveil rather brick bargain curious require crowd raise" | practiced keys add feeacc --keyring-backend $KEYRING --algo $KEYALGO --recover
  
  practiced init $MONIKER --chain-id $CHAINID 

  # Function updates the config based on a jq argument as a string
  update_test_genesis () {
    # update_test_genesis '.consensus_params["block"]["max_gas"]="100000000"'
    cat $HOME/.practice/config/genesis.json | jq "$1" > $HOME/.practice/config/tmp_genesis.json && mv $HOME/.practice/config/tmp_genesis.json $HOME/.practice/config/genesis.json
  }

  # Set gas limit in genesis
  update_test_genesis '.consensus_params["block"]["max_gas"]="100000000"'
  update_test_genesis '.app_state["gov"]["voting_params"]["voting_period"]="15s"'

  update_test_genesis '.app_state["staking"]["params"]["bond_denom"]="upractice"'  
  update_test_genesis '.app_state["bank"]["params"]["send_enabled"]=[{"denom": "upractice","enabled": true}]'
  # update_test_genesis '.app_state["staking"]["params"]["min_commission_rate"]="0.100000000000000000"' # sdk 46 only   

  update_test_genesis '.app_state["mint"]["params"]["mint_denom"]="upractice"'  
  update_test_genesis '.app_state["gov"]["deposit_params"]["min_deposit"]=[{"denom": "upractice","amount": "1000000"}]'
  update_test_genesis '.app_state["crisis"]["constant_fee"]={"denom": "upractice","amount": "1000"}'  

  update_test_genesis '.app_state["tokenfactory"]["params"]["denom_creation_fee"]=[{"denom":"upractice","amount":"100"}]'

  update_test_genesis '.app_state["feeshare"]["params"]["allowed_denoms"]=["upractice"]'

  # Allocate genesis accounts
  practiced add-genesis-account $KEY 10000000upractice,1000utest --keyring-backend $KEYRING
  practiced add-genesis-account feeacc 1000000upractice,1000utest --keyring-backend $KEYRING

  practiced gentx $KEY 1000000upractice --keyring-backend $KEYRING --chain-id $CHAINID

  # Collect genesis tx
  practiced collect-gentxs

  # Run this to ensure practicerything worked and that the genesis file is setup correctly
  practiced validate-genesis
}


if [ $# -eq 1 ] && [ $1 == "clean" ] || [ $1 == "c" ]; then
  echo "Starting from a clean state"
  from_scratch
fi

echo "Starting node..."

# Opens the RPC endpoint to outside connections
sed -i '/laddr = "tcp:\/\/127.0.0.1:26657"/c\laddr = "tcp:\/\/0.0.0.0:26657"' ~/.practice/config/config.toml
sed -i 's/cors_allowed_origins = \[\]/cors_allowed_origins = \["\*"\]/g' ~/.practice/config/config.toml
sed -i 's/enable = false/enable = true/g' ~/.practice/config/app.toml

practiced start --pruning=nothing  --minimum-gas-prices=0upractice  