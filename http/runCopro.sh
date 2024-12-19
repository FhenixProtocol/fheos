#!/bin/bash
# copy this bash to the root directory where the copro-aggregator, copro-task-manager-hardhat and nitro/fheos directories are present

AggregatorDir="$(pwd)/copro-aggregator"
TaskManagerDir="$(pwd)/copro-task-manager-hardhat"
CtRegistryDir="$(pwd)/ct-registry"
FheosDir="$(pwd)/nitro/fheos/"
FheosImage="ghcr.io/fhenixprotocol/nitro/localfhenix:v0.3.2-alpha.17"
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

if [ ! -d "$CtRegistryDir" ]; then
  missing_dirs+=("Ct Registry Directory: $CtRegistryDir")
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
    local name=$1
    echo -e "${YELLOW}Listening to Docker container logs for 'localfhenix_copro'${NC}"
    local didStart=$(docker logs -f $name | while read line; do
        if echo "$line" | grep -q "HTTP server started"; then
            echo -e "${GREEN}Started successfully${NC}"
            break
        fi
    done)
    if [ -z "$didStart" ]; then
        err "Error: Failed to start Local Fhenix $name"
    fi
}

# Placeholder functions for each task
function runLocalFhenix {
    local name=$1
    local portRpc=$2
    local portWs=$3
    local portFaucet=$4

    cd $TaskManagerDir
    pnpm install
    if docker ps | grep $name; then
        success "$name is already running"
    else
        success "localfhenix is not running"
        docker run -d --name $name -p $portRpc:8547 -p $portWs:8548 -p $portFaucet:3000 -it $FheosImage
        if [ $? -ne 0 ]; then
            err "Error: Failed to run Local Fhenix: $name"
        fi 
        success "Local Fhenix ($name) is started"
        listenDockerLogs $name
    fi
}

function buildFheosServer {
    if [ "$BuildParam" == "build" ] || [ "$BuildParam" == "debugbuild" ]; then
        step "Building Fheos server"
        cd $FheosDir
        docker build -t $FheosServerImage .
        if [ $? -ne 0 ]; then
            err "Error: Failed to build Fheos server"
        fi 
        success "Fheos server is built"
    fi
}

function runFheosServer {
    step "Running Fheos server"
    if docker ps | grep $FheosServerImage; then
            success "Fheos server is already running"
            return
    fi
    if docker ps -a | grep fheos_server; then
        step "Removing existing Fheos server container"
        docker container rm fheos_server
    fi
    if [ "$BuildParam" == "debug" ] || [ "$BuildParam" == "debugbuild" ]; then
      docker run -d --network host --name fheos_server -e DEBUG_MODE=1 -it -p 8448:8448 -p 4002:4002  $FheosServerImage
    else
      docker run -d --network host --name fheos_server -it -p 8448:8448 -p 4002:4002  $FheosServerImage
    fi
    if [ $? -ne 0 ]; then
        err "Error: Failed to run Fheos server"
    fi 
    success "Fheos server is started"
}

function deployContracts {
    local contractRepo=$1

    step "Running deploy $(basename "$contractRepo")"
    cd $contractRepo
    pnpm pnpm clean
    output=$(pnpm task:deploy)
    if [ $? -ne 0 ]; then
        err "Error: Failed to deploy $(basename "$contractRepo") contracts"
    fi
}

function copyDeployedContract {
    local contractPath=$1
    step "Copying the deployed Contract to the AggregatorDir: $(basename $contractPath)"
    cp $contractPath $AggregatorDir
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
step "Running Local Fhenix (consumer chain)"
runLocalFhenix localfhenix_copro_da 42069 42070 42000
step "Running Local Fhenix (da chain)"
runLocalFhenix localfhenix_copro_consumer 42169 42170 42100
buildFheosServer
runFheosServer
deployContracts $TaskManagerDir
deployContracts $CtRegistryDir
copyDeployedContract $CtRegistryDir/deployments/localfhenix/CiphertextRegistry.json
copyDeployedContract $TaskManagerDir/deployments/localfhenix/TaskManager.json

copyDeployedCtRegistry
if [ "$BuildParam" != "NA" ]; then
    startAggregator
fi
cd $current_directory
