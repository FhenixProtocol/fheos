package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/ethdb/memorydb"
	"github.com/fhenixprotocol/fheos/precompiles"
	fhedriver "github.com/fhenixprotocol/warp-drive/fhe-driver"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"os"
	"strconv"
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

type HashResultUpdate struct {
	TempKey    []byte `json:"tempKey"`
	ActualHash []byte `json:"actualHash"`
}

type DecryptResultUpdate struct {
	CtKey     []byte   `json:"ctKey"`
	Plaintext *big.Int `json:"plaintext"`
}

var tp precompiles.TxParams

func responseToServer(url string, tempKey []byte, json []byte) {
	// Create a new request using http.NewRequest
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(json))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	// Set the request content-type to application/json
	req.Header.Set("Content-Type", "application/json")

	// Send the request using http.Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Read and print the response body
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response: %v", err)
	}

	fmt.Printf("Update requester %s with the result of %+v\n", url, tempKey)
}

func handleResult(url string, tempKey []byte, actualHash []byte) {
	fmt.Printf("Got result for %s : %s\n", hex.EncodeToString(tempKey), hex.EncodeToString(actualHash))
	// JSON data to be sent in the request body
	jsonData, err := json.Marshal(HashResultUpdate{TempKey: tempKey, ActualHash: actualHash})
	if err != nil {
		log.Fatalf("Failed to update requester %s with the result of %+v", url, tempKey)
	}

	responseToServer(url, tempKey, jsonData)
}

func handleDecryptResult(url string, ctKey []byte, plaintext *big.Int) {
	fmt.Printf("Got result for %s : %s\n", hex.EncodeToString(ctKey), plaintext)
	jsonData, err := json.Marshal(DecryptResultUpdate{CtKey: ctKey, Plaintext: plaintext})
	if err != nil {
		log.Fatalf("Failed to update requester %s with the result of %+v", url, ctKey)
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
		err := os.Mkdir("./keys/", 0755)
		if err != nil {
			return err
		}
	}

	if _, err := os.Stat("./keys/tfhe/"); os.IsNotExist(err) {
		err := os.Mkdir("./keys/tfhe/", 0755)
		if err != nil {
			return err
		}
	}

	for i := int32(0); i < securityZones; i++ {
		err := fhedriver.GenerateFheKeys(i)
		if err != nil {
			return fmt.Errorf("error generating FheKeys for securityZone %d: %s", i, err)
		}
	}
	return nil
}

func initFheos() (*precompiles.TxParams, error) {
	if os.Getenv("FHEOS_DB_PATH") == "" {
		err := os.Setenv("FHEOS_DB_PATH", "./fheosdb")
		if err != nil {
			return nil, err
		}
	}

	err := precompiles.InitFheConfig(&fhedriver.ConfigDefault)
	if err != nil {
		return nil, err
	}

	securityZones, err := getenvInt("FHEOS_SECURITY_ZONES", 1)
	if err != nil {
		return nil, err
	}
	err = generateKeys(int32(securityZones))
	if err != nil {
		return nil, err
	}

	err = precompiles.InitializeFheosState()
	if err != nil {
		return nil, err
	}

	err = os.Setenv("FHEOS_DB_PATH", "")
	if err != nil {
		return nil, err
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
			return nil, err
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

	result, _, err := precompiles.Decrypt(req.UType, hash, nil, &tp, &callback)
	if err != nil {
		e := fmt.Sprintf("Operation failed: %+v", err)
		fmt.Println(e)
		http.Error(w, e, http.StatusBadRequest)
		return
	}
	resultString := result.Text(16)
	// Respond with the result
	w.Write([]byte(resultString))
	fmt.Printf("Done processing decrypt request %+v (%s: %+v)\n", result, resultString, []byte(resultString))
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

	// Start the server
	log.Println("Server listening on port 8448...")
	log.Fatal(http.ListenAndServe(":8448", nil))
}
