module github.com/fhenixprotocol/fheos

go 1.20

require (
	github.com/ethereum/go-ethereum v1.10.26
	github.com/holiman/uint256 v1.2.3
	github.com/sirupsen/logrus v1.9.3
)

require github.com/deckarep/golang-set/v2 v2.3.1 // indirect

replace github.com/ethereum/go-ethereum => ./go-ethereum

replace github.com/fhenixprotocol/go-tfhe => ./go-tfhe

replace github.com/fhenixprotocol/decryption-oracle => ./go-tfhe/decryption-oracle

require (
	github.com/fhenixprotocol/decryption-oracle-proto v0.0.0-20231127135430-74f6a7e6df5c // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/gammazero/deque v0.2.1 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/stretchr/testify v1.8.1 // indirect
	golang.org/x/exp v0.0.0-20230905200255-921286631fa9 // indirect
	golang.org/x/net v0.15.0 // indirect
	golang.org/x/text v0.13.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230822172742-b8732ec3820d // indirect
	google.golang.org/grpc v1.59.0 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
	gopkg.in/natefinch/npipe.v2 v2.0.0-20160621034901-c1b8fa8bdcce // indirect
)

require (
	github.com/StackExchange/wmi v1.2.1 // indirect
	github.com/btcsuite/btcd/btcec/v2 v2.2.0 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.0.1 // indirect
	github.com/fhenixprotocol/go-tfhe v0.0.0-20230806123311-862edf6c6556
	github.com/go-ole/go-ole v1.2.5 // indirect
	github.com/go-stack/stack v1.8.1 // indirect
	github.com/google/uuid v1.3.1 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/shirou/gopsutil v3.21.4-0.20210419000835-c7a38de76ee5+incompatible // indirect
	github.com/tklauser/go-sysconf v0.3.12 // indirect
	github.com/tklauser/numcpus v0.6.1 // indirect
	golang.org/x/crypto v0.14.0
	golang.org/x/sys v0.14.0 // indirect

)
