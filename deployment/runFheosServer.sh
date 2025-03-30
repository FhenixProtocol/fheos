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
if [[ "${FHE_ENGINE_CONFIG_DIR}" != "" ]]; then
    fhe-engine-server -c "${FHE_ENGINE_CONFIG_DIR}/fhe_engine.toml" &
else
    fhe-engine-server -c /home/user/fhenix/fhe_engine.toml &
fi

echo "Waiting for connection to 127.0.0.1:50051..."
timeout=120
start_time=$(date +%s)
while true; do
    if nc -z 127.0.0.1 50051 2>/dev/null; then
        echo "Connected to 127.0.0.1:50051"
        echo "connected to engine"
        break
    fi
    
    current_time=$(date +%s)
    elapsed_time=$((current_time - start_time))
    
    if [ $((elapsed_time % 20)) -eq 0 ]; then
        echo "Elapsed time: $elapsed_time seconds"
    fi
    
    if [ $elapsed_time -ge $timeout ]; then
        echo "Timeout: Port 50051 is not available"
        exit 1
    fi
    
    sleep 1
done

if [[ "${DETACH_MODE}" -eq 1 ]]; then
  coprocessor &
  exec /bin/bash
fi

if [[ "${DEBUG_MODE}" -eq 1 ]]; then
  /go/bin/dlv --listen=:4002 --headless=true --log=true --accept-multiclient --api-version=2 exec /usr/local/bin/coprocessor
else
  coprocessor
fi