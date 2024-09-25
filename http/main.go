package http

import (
	"encoding/hex"
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/fhenixprotocol/fheos/precompiles"
	"io/ioutil"
	"log"
	"net/http"
)

// Struct to parse the incoming JSON request
type HashRequest struct {
	UType   byte   `json:"utype"`
	LhsHash string `json:"lhsHash"`
	RhsHash string `json:"rhsHash"`
}

func handleResult(tempKey []byte, newHash []byte) {

}

// Helper function to handle decoding the request and calling the respective function
func handleRequest(w http.ResponseWriter, r *http.Request, handler func(byte, []byte, []byte, *precompiles.TxParams, *precompiles.CallbackFunc) ([]byte, uint64, error)) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var req HashRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Convert the hash strings to byte arrays
	lhsHash, err := hex.DecodeString(req.LhsHash)
	if err != nil {
		http.Error(w, "Invalid lhsHash", http.StatusBadRequest)
		return
	}

	rhsHash, err := hex.DecodeString(req.RhsHash)
	if err != nil {
		http.Error(w, "Invalid rhsHash", http.StatusBadRequest)
		return
	}

	txParams := precompiles.TxParams{
		Commit:          true,
		EthCall:         false,
		CiphertextDb:    nil, // LIOR TODO
		ContractAddress: common.Address{},
		GetBlockHash: func(u uint64) common.Hash {
			return common.Hash{}
		},
		BlockNumber: nil,
		ErrChannel:  nil,
	}
	callback := precompiles.CallbackFunc{
		Callback: handleResult,
	}
	// Call the handler function (either Shr or Add)
	result, _, err := handler(req.UType, lhsHash, rhsHash, &txParams, &callback)

	if err != nil {
		http.Error(w, "Operation failed", http.StatusBadRequest)
		return
	}

	// Respond with the result
	w.Write(result)
}

// Handler for the /shr endpoint
func ShrHandler(w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, precompiles.Shr)
}

// Handler for the /add endpoint
func AddHandler(w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, precompiles.Add)
}

func main() {
	// Set up the HTTP server and routes
	http.HandleFunc("/shr", ShrHandler)
	http.HandleFunc("/add", AddHandler)

	// Start the server
	log.Println("Server listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
