#!/bin/bash
# copy this bash to the root directory where the copro-aggregator, copro-task-manager-hardhat and nitro/fheos directories are present

AggregatorDir="$(pwd)/copro-aggregator"
TaskManagerDir="$(pwd)/copro-task-manager-hardhat"
FheosDir="$(pwd)/nitro/fheos/"
FheosImage="ghcr.io/fhenixprotocol/localfhenix:v0.2.4"
FheosServerImage="fheosserver"

BuildParam=$1

# Verify that directories exist
missing_dirs=()


RED='\e[4;31m'
GREEN='\e[4;32m'
YELLOW='\e[4;33m'

NC='\e[0m' # No Color
step_index=1

if [ ! -d "$AggregatorDir" ]; then
  missing_dirs+=("Aggregator Directory: $AggregatorDir")
fi

if [ ! -d "$TaskManagerDir" ]; then
  missing_dirs+=("Task Manager Directory: $TaskManagerDir")
fi

if [ ! -d "$FheosDir" ]; then
  missing_dirs+=("Fheos Directory: $FheosDir")
fi

# If any directories are missing, print them and exit
if [ ${#missing_dirs[@]} -ne 0 ]; then
  echo -e "\n${RED}Error: One or more directories do not exist. Please check the paths in the configuration file.${NC}"
  echo -e "${RED}Missing directories:${NC}"
  for dir in "${missing_dirs[@]}"; do
    echo "$dir"
  done
  exit 1
fi

function step() {
    echo -e "\n${YELLOW}[${step_index}] $1${NC}"
    step_index=$((step_index + 1))
}

function err(){
    echo -e "${RED}Error: $1${NC}"
    exit 1
}

function success(){
    echo -e "${GREEN}$1${NC}"
}

function listenDockerLogs {
    echo -e "${YELLOW}Listening to Docker container logs for 'localfhenix_copro'${NC}"
    docker logs -f localfhenix_copro | while read line; do
        if echo "$line" | grep -q "HTTP server started"; then
            echo -e "${GREEN}Started successfully${NC}"
            break
        fi
    done
}

# Placeholder functions for each task
function runningLocalFhenix {
    step "Running Local Fhenix"
    cd $TaskManagerDir
    pnpm install
    if docker ps | grep $FheosImage; then
            success "Local Fhenix is already running"
    else
        docker run -d --name localfhenix_copro -p 42069:8547 -p 42070:8548 -p 42000:3000 -it $FheosImage 
        if [ $? -ne 0 ]; then
            err "Error: Failed to run Local Fhenix"
        fi 
        success "Local Fhenix is started"
        listenDockerLogs
    fi
}

function buildFheosServer {
    if [ "$BuildParam" == "build" ]; then
        step "Building Fheos server"
        cd $FheosDir
        docker build -t $FheosServerImage .
        if [ $? -ne 0 ]; then
            err "Error: Failed to build Fheos server"
        fi 
        success "Fheos server is built"
    fi
}

function runningFheosServer {
    step "Running Fheos server"
    if docker ps | grep $FheosServerImage; then
            success "Fheos server is already running"
            return
    fi
    if docker ps -a | grep fheos_server; then
        step "Removing existing Fheos server container"
        docker container rm fheos_server
    fi
    docker run -d --network host --name fheos_server -e COPROCESSOR_MODE=1 -it -p 8448:8448 $FheosServerImage
    if [ $? -ne 0 ]; then
        err "Error: Failed to run Fheos server"
    fi 
    success "Fheos server is started"
}

function deployContracts {
    step "Running deploy TaskManager and Example contracts"
    cd $TaskManagerDir
    pnpm pnpm clean
    output=$(pnpm task:deploy)
    if [ $? -ne 0 ]; then
        err "Error: Failed to deploy contracts"
    fi
}

function copyDeployedTM {
    step "Copying the deployed TaskManager to the AggregatorDir: $AggregatorDir"
    cp $TaskManagerDir/deployments/localfhenix/TaskManager.json $AggregatorDir
}

function startAggregator {
    step "Starting the Aggregator"
    cd $AggregatorDir
    pnpm install
    node aggregator.js
}

function stopDocker {
    local name=$1
    step "Stopping ${name} container"
    docker stop $name
    step "Removing ${name} container"
    docker rm $name
}

if [ "$1" == "stop" ]; then
    stopDocker localfhenix_copro
    stopDocker fheos_server
    exit 0
fi

# Main Script Logic - Run tasks
echo "Starting the script with configured directories"
current_directory=$(pwd)
runningLocalFhenix
buildFheosServer
runningFheosServer
deployContracts
copyDeployedTM
if [ "$BuildParam" != "NA" ]; then
    startAggregator
fi
cd $current_directory