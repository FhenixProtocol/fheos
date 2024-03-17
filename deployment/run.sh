#!/usr/bin/env bash

fhe-engine-server -c /home/user/fhenix/fhe_engine.toml &

fheos init-state

node faucet/server.js &

nitro --conf.file /config/sequencer_config.json --metrics --node.feed.output.enable --node.feed.output.port 9642  --http.api net,web3,eth,txpool,debug --node.seq-coordinator.my-url  ws://sequencer:8548 --graphql.enable --graphql.vhosts "*" --graphql.corsdomain "*"