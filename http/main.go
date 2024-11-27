package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/ethdb/memorydb"
	"github.com/fhenixprotocol/fheos/precompiles"
	fhedriver "github.com/fhenixprotocol/warp-drive/fhe-driver"
)

const (
	maxRetries = 3
)

// Struct to parse the incoming JSON request
type HashRequest struct {
	UType        byte   `json:"utype"`
	LhsHash      string `json:"lhsHash"`
	RhsHash      string `json:"rhsHash"`
	RequesterUrl string `json:"requesterUrl"`
}

type DecryptRequest struct {
	UType        byte   `json:"utype"`
	Hash         string `json:"hash"`
	RequesterUrl string `json:"requesterUrl"`
}

type SealOutputRequest struct {
	UType        byte   `json:"utype"`
	Hash         string `json:"hash"`
	PKey         string `json:"pkey"`
	RequesterUrl string `json:"requesterUrl"`
}

type HashResultUpdate struct {
	TempKey    []byte `json:"tempKey"`
	ActualHash []byte `json:"actualHash"`
}

type DecryptResultUpdate struct {
	CtKey     []byte `json:"ctKey"`
	Plaintext string `json:"plaintext"`
}

type SealOutputResultUpdate struct {
	CtKey []byte `json:"ctKey"`
	Value string `json:"value"`
}

var tp precompiles.TxParams

func doWithRetry(operation func() error) error {
	var lastErr error
	for i := 0; i < maxRetries; i++ {
		if err := operation(); err != nil {
			lastErr = err
			// Small exponential backoff before retrying
			time.Sleep(time.Duration(i+1) * time.Millisecond * 50)
			fmt.Printf("Retrying operation (attempt %d/%d)\n", i+1, maxRetries)
		} else {
			return nil
		}
	}
	return fmt.Errorf("operation failed after %d attempts: %v", maxRetries, lastErr)
}

func responseToServer(url string, tempKey []byte, json []byte) {
	err := doWithRetry(func() error {
		// Create a new request using http.NewRequest
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(json))
		if err != nil {
			return fmt.Errorf("error creating request: %v", err)
		}

		// Set the request content-type to application/json
		req.Header.Set("Content-Type", "application/json")

		// Send the request using http.Client
		client := &http.Client{
			// TODO : Adjust this timeout after gathering some real data
			Timeout: 5 * time.Second,
		}

		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("error sending request: %v", err)
		}
		defer resp.Body.Close()

		// Read and print the response body
		_, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("error reading response: %v", err)
		}

		// Check if the status code indicates success
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return fmt.Errorf("server returned non-success status code: %d for %s", resp.StatusCode, url)
		}

		return nil
	})

	if err != nil {
		log.Printf("Failed to send response after all retries: %v", err)
		return
	}

	fmt.Printf("Update requester %s with the result of %+v\n", url, tempKey)
}

func handleResult(url string, tempKey []byte, actualHash []byte) {
	fmt.Printf("Got result for %s : %s\n", hex.EncodeToString(tempKey), hex.EncodeToString(actualHash))
	// JSON data to be sent in the request body
	jsonData, err := json.Marshal(HashResultUpdate{TempKey: tempKey, ActualHash: actualHash})
	if err != nil {
		log.Printf("Failed to marshal update for requester %s with the result of %+v: %v", url, tempKey, err)
		return
	}

	responseToServer(url, tempKey, jsonData)
}

func handleDecryptResult(url string, ctKey []byte, plaintext *big.Int) {
	fmt.Printf("Got result for %s : %s\n", hex.EncodeToString(ctKey), plaintext)
	plaintextString := plaintext.Text(16)
	jsonData, err := json.Marshal(DecryptResultUpdate{CtKey: ctKey, Plaintext: plaintextString})
	if err != nil {
		log.Printf("Failed to marshal decrypt result for requester %s with the result of %+v: %v", url, ctKey, err)
		return
	}

	responseToServer(url, ctKey, jsonData)
}

func handleSealOutputResult(url string, ctKey []byte, value string) {
	fmt.Printf("Got result for %s : %s\n", hex.EncodeToString(ctKey), value)
	jsonData, err := json.Marshal(SealOutputResultUpdate{CtKey: ctKey, Value: value})
	if err != nil {
		log.Printf("Failed to marshal seal output result for requester %s with the result of %+v: %v", url, ctKey, err)
		return
	}

	responseToServer(url, ctKey, jsonData)
}

// Helper function to handle decoding the request and calling the respective function
func handleRequest(w http.ResponseWriter, r *http.Request, handler func(byte, []byte, []byte, *precompiles.TxParams, *precompiles.CallbackFunc) ([]byte, uint64, error)) {
	fmt.Printf("Got a request from %s\n", r.RemoteAddr)
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
		e := fmt.Sprintf("Invalid lhsHash: %s %+v", req.LhsHash, err)
		fmt.Println(e)
		http.Error(w, e, http.StatusBadRequest)
		return
	}

	rhsHash, err := hex.DecodeString(req.RhsHash)
	if err != nil {
		e := fmt.Sprintf("Invalid lhsHash: %s %+v", req.LhsHash, err)
		fmt.Println(e)
		http.Error(w, e, http.StatusBadRequest)
		return
	}

	callback := precompiles.CallbackFunc{
		CallbackUrl: req.RequesterUrl,
		Callback:    handleResult,
	}

	result, _, err := handler(req.UType, lhsHash, rhsHash, &tp, &callback)

	if err != nil {
		e := fmt.Sprintf("Operation failed: %s %+v", req.LhsHash, err)
		fmt.Println(e)
		http.Error(w, e, http.StatusBadRequest)
		return
	}

	res := []byte(hex.EncodeToString(result))
	// Respond with the result
	w.Write(res)
	fmt.Printf("Started processing the request for tempkey %s\n", hex.EncodeToString(result))
}

func getenvInt(key string, defaultValue int) (int, error) {
	s := os.Getenv(key)
	if s == "" {
		return defaultValue, nil
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return v, nil
}

func generateKeys(securityZones int32) error {
	if _, err := os.Stat("./keys/"); os.IsNotExist(err) {
		if err := os.Mkdir("./keys/", 0755); err != nil {
			return fmt.Errorf("failed to create keys directory: %v", err)
		}
	}

	if _, err := os.Stat("./keys/tfhe/"); os.IsNotExist(err) {
		if err := os.Mkdir("./keys/tfhe/", 0755); err != nil {
			return fmt.Errorf("failed to create tfhe directory: %v", err)
		}
	}

	for i := int32(0); i < securityZones; i++ {
		if err := fhedriver.GenerateFheKeys(i); err != nil {
			return fmt.Errorf("error generating FheKeys for securityZone %d: %v", i, err)
		}
	}
	return nil
}

func initFheos() (*precompiles.TxParams, error) {
	if os.Getenv("FHEOS_DB_PATH") == "" {
		if err := os.Setenv("FHEOS_DB_PATH", "./fheosdb"); err != nil {
			return nil, fmt.Errorf("failed to set FHEOS_DB_PATH: %v", err)
		}
	}

	if err := precompiles.InitFheConfig(&fhedriver.ConfigDefault); err != nil {
		return nil, fmt.Errorf("failed to init FHE config: %v", err)
	}

	securityZones, err := getenvInt("FHEOS_SECURITY_ZONES", 1)
	if err != nil {
		return nil, fmt.Errorf("failed to get security zones: %v", err)
	}

	if err := generateKeys(int32(securityZones)); err != nil {
		return nil, fmt.Errorf("failed to generate keys: %v", err)
	}

	if err := precompiles.InitializeFheosState(); err != nil {
		return nil, fmt.Errorf("failed to initialize FHEOS state: %v", err)
	}

	if err := os.Setenv("FHEOS_DB_PATH", ""); err != nil {
		return nil, fmt.Errorf("failed to clear FHEOS_DB_PATH: %v", err)
	}

	tp = precompiles.TxParams{
		Commit:          true,
		GasEstimation:   false,
		EthCall:         false,
		CiphertextDb:    memorydb.New(),
		ContractAddress: common.HexToAddress("0x0000000000000000000000000000000000000000"),
		GetBlockHash:    vm.GetHashFunc(nil),
		BlockNumber:     nil,
		ParallelTxHooks: nil,
	}

	var trivialHash []byte
	for i := 0; i <= 50; i++ {

		// Create a byte slice of size 32
		toEncrypt := make([]byte, 32)

		// Convert the integer to bytes and store it in the byte slice
		toEncrypt[31] = uint8(i)
		trivialHash, _, err = precompiles.TrivialEncrypt(toEncrypt, 2, 0, &tp, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to generate trivial hash for %d: %v", i, err)
		}
		fmt.Printf("Trivial hash for %d: %x\n", i, trivialHash)
	}
	return &tp, nil
}

func DecryptHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Got a request from %s\n", r.RemoteAddr)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var req DecryptRequest
	if err := json.Unmarshal(body, &req); err != nil {
		fmt.Printf("Failed unmarsheling request: %+v body is %+v\n", err, string(body))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Convert the hash strings to byte arrays
	hash, err := hex.DecodeString(req.Hash)
	if err != nil {
		e := fmt.Sprintf("Invalid hash: %s %+v", req.Hash, err)
		fmt.Println(e)
		http.Error(w, e, http.StatusBadRequest)
		return
	}

	callback := precompiles.DecryptCallbackFunc{
		CallbackUrl: req.RequesterUrl,
		Callback:    handleDecryptResult,
	}

	_, _, err = precompiles.Decrypt(req.UType, hash, nil, &tp, &callback)
	if err != nil {
		e := fmt.Sprintf("Operation failed: %+v", err)
		fmt.Println(e)
		http.Error(w, e, http.StatusBadRequest)
		return
	}
	// Respond with the result
	w.Write(hash)
	fmt.Printf("Received decrypt request for %+v and type %+v\n", hash, req.UType)
}

func SealOutputHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Got a request from %s\n", r.RemoteAddr)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var req SealOutputRequest
	if err := json.Unmarshal(body, &req); err != nil {
		fmt.Printf("Failed unmarsheling request: %+v body is %+v\n", err, string(body))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Convert the hash strings to byte arrays
	hash, err := hex.DecodeString(req.Hash)
	if err != nil {
		e := fmt.Sprintf("Invalid hash: %s %+v", req.Hash, err)
		fmt.Println(e)
		http.Error(w, e, http.StatusBadRequest)
		return
	}

	callback := precompiles.SealOutputCallbackFunc{
		CallbackUrl: req.RequesterUrl,
		Callback:    handleSealOutputResult,
	}

	pkey, err := hex.DecodeString(req.PKey)
	if err != nil {
		e := fmt.Sprintf("Invalid pkey: %s %+v", req.PKey, err)
		fmt.Println(e)
		http.Error(w, e, http.StatusBadRequest)
		return
	}

	_, _, err = precompiles.SealOutput(req.UType, hash, pkey, &tp, &callback)
	if err != nil {
		e := fmt.Sprintf("Operation failed: %+v", err)
		fmt.Println(e)
		http.Error(w, e, http.StatusBadRequest)
		return
	}
	// Respond with the result
	w.Write(hash)
	fmt.Printf("Received seal output request for %+v and type %+v\n", hash, req.UType)

}

func main() {
	_, err := initFheos()
	if err != nil {
		log.Fatalf("Failed to initialize FHEOS: %v", err)
		os.Exit(1)
	}

	handlers := getHandlers()
	// iterate handlers
	for _, handler := range handlers {
		http.HandleFunc(handler.Name, handler.Handler)
	}

	http.HandleFunc("/Decrypt", DecryptHandler)
	http.HandleFunc("/SealOutput", SealOutputHandler)

	// Start the server
	log.Println("Server listening on port 8448...")
	if err := http.ListenAndServe(":8448", nil); err != nil {
		log.Fatalf("Server stopped: %v", err)
	}
}
