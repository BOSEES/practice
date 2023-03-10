#!/bin/bash

if [ -z "$CONTRACT_ID" ]; then
  CONTRACT_ID=1
fi

CHAIN_ID="${CHAIN_ID:-practice-t1}"
NODE="${NODE:-http://localhost:26657}"

# NODE="http://localhost:26657"
TX_FLAGS="--from practice1 --keyring-backend test --chain-id $CHAIN_ID --gas 10000000 --fees 20000upractice --node $NODE --output json --yes -b block"

# upload the contract
practiced tx wasm store ./scripts/oracle_querier.wasm $TX_FLAGS

# Get the transaction upload hash to query in the next step
TX_HASH=$(practiced tx wasm instantiate $CONTRACT_ID "{}" --label "ORACLE QUERIER" $TX_FLAGS --no-admin  | jq -r '.txhash') && echo $TX_HASH

# Query the logs for the contract address
CONTRACT_ADDR=$(practiced query tx $TX_HASH --output json --node $NODE | jq -r '.logs[0].events[0].attributes[0].value') && echo "Address: $CONTRACT_ADDR"

practiced q wasm contract-state smart $CONTRACT_ADDR '{"exchange_rate": {"denom":"upractice"}}' --node $NODE