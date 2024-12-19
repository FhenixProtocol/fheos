#!/usr/bin/env bash

# Function to display usage
usage() {
    echo "Usage: $0 [--debug] [--detach]"
    echo "  --debug    Run in debug mode"
    exit 1
}

echo "$@"
# Check for the --debug flag
DEBUG_MODE=${DEBUG_MODE:-0}
DETACH_MODE=${DETACH_MODE:-0}
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
    echo "Starting fheos-server in normal mode"
else
    echo "Starting fheos-server in debug mode"
fi

# renault-server -c /home/user/fhenix/renault-server.toml &
# Start the FHE engine server
fhe-engine-server -c /home/user/fhenix/fhe_engine.toml &

# Wait for the server to start
sleep 2

if [[ "${DETACH_MODE}" -eq 1 ]]; then
  exec /bin/bash
fi

if [[ "${DEBUG_MODE}" -eq 1 ]]; then
  /go/bin/dlv --listen=:4002 --headless=true --log=true --accept-multiclient --api-version=2 exec /usr/local/bin/coprocessor
else
  coprocessor
fi