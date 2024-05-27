#!/usr/bin/env bash

# Function to display usage
usage() {
    echo "Usage: $0 [--debug]"
    echo "  --debug    Run in debug mode"
    exit 1
}

# Check for the --debug flag
DEBUG_MODE=${DEBUG_MODE:-0}
for arg in "$@"
do
    case $arg in
        --debug)
        DEBUG_MODE=1
        shift # Remove --debug from processing
        ;;
        *)
        usage
        ;;
    esac
done

if [[ "${DEBUG_MODE}" -eq 0 ]]; then
    echo "Starting in normal mode"
else
    echo "Starting in debug mode"
fi

# Start the FHE engine server
fhe-engine-server -c /home/user/fhenix/fhe_engine.toml &

# Wait for the server to start
sleep 2

# Path to the keys directory
KEYS_PATH="/home/fhenix/keys/0/tfhe/sks"

# Check if keys exist in the KEYS_PATH
if [ "$(ls -A $KEYS_PATH 2>/dev/null)" ]; then
    echo "Keys already exist in $KEYS_PATH, no need to initialize state."
else
    echo "No keys found in $KEYS_PATH. Initializing state..."
    # Initialize the FHE state
    fheos init-state

    # Wait for the keys to be loaded
    sleep 3
    # Command to initialize state goes here
fi

# Start the faucet server
node faucet/server.js &

# Nitro service is started only in "normal mode"
if [[ "${DEBUG_MODE}" -eq 0 ]]; then
    echo "Starting Nitro service in normal mode"
    nitro --conf.file /config/sequencer_config.json \
          --metrics \
          --node.feed.output.enable \
          --node.feed.output.port 9642 \
          --http.api net,web3,eth,txpool,debug \
          --node.seq-coordinator.my-url ws://sequencer:8548 \
          --graphql.enable \
          --graphql.vhosts "*" \
          --graphql.corsdomain "*" \
          --conf.env-prefix "NITRO" \
          --conf.fhenix.log-level 4
fi

# Start in debug mode if requested
if [[ "${DEBUG_MODE}" -eq 1 ]]; then
    echo "Starting in debug mode"
    /go/bin/dlv --listen=:4001 --headless=true --log=true --accept-multiclient --api-version=2 exec /usr/local/bin/nitro -- --conf.file /config/sequencer_config.json --node.dangerous.no-l1-listener --node.feed.output.enable --node.feed.output.port 9642 --http.api net,web3,eth,txpool,debug --node.seq-coordinator.my-url ws://sequencer:8548 --graphql.enable --graphql.vhosts "*" --graphql.corsdomain "*"
fi
