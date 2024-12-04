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
	"reflect"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/ethdb/memorydb"
	"github.com/fhenixprotocol/fheos/precompiles"
	fhedriver "github.com/fhenixprotocol/warp-drive/fhe-driver"
)

const (
	maxRetriesForRepost = 3
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

// type BigInt struct {
// 	Value *big.Int
// }

// func (b *BigInt) UnmarshalJSON(data []byte) error {
// 	// Unmarshal the string into a temporary variable
// 	var str string
// 	if err := json.Unmarshal(data, &str); err != nil {
// 		return err
// 	}

// 	// Convert the string to a big.Int
// 	b.Value = new(big.Int)
// 	_, ok := b.Value.SetString(str, 10) // base 10
// 	if !ok {
// 		return fmt.Errorf("invalid big.Int format: %s", str)
// 	}
// 	return nil
// }

type TrivialEncryptRequest struct {
	Value        *big.Int `json:"value"`
	ToType       byte     `json:"toType"`
	SecurityZone int32    `json:"securityZone"`
	RequesterUrl string   `json:"requesterUrl"`
}

func (r *TrivialEncryptRequest) UnmarshalJSON(data []byte) error {
	// Define a temporary struct to unmarshal JSON into
	var aux struct {
		Value        string `json:"value"`        // Temporary string to handle big.Int
		ToType       byte   `json:"toType"`       // As-is
		SecurityZone string `json:"securityZone"` // Temporary string for int32 conversion
		RequesterUrl string `json:"requesterUrl"` // As-is
	}

	// Unmarshal into the temporary struct
	if err := json.Unmarshal(data, &aux); err != nil {
		fmt.Printf("Failed to unmarshal request: %v", err)
		return err
	}

	// Parse `value` as *big.Int
	r.Value = new(big.Int)
	if _, ok := r.Value.SetString(aux.Value, 16); !ok {
		return fmt.Errorf("invalid big.Int format: %s", aux.Value)
	}

	// Parse `securityZone` as int32
	securityZone, err := strconv.ParseInt(aux.SecurityZone, 10, 32)
	if err != nil {
		return fmt.Errorf("invalid securityZone format: %s", aux.SecurityZone)
	}
	r.SecurityZone = int32(securityZone)

	// Assign the other fields
	r.ToType = aux.ToType
	r.RequesterUrl = aux.RequesterUrl

	return nil
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
	for i := 0; i < maxRetriesForRepost; i++ {
		if err := operation(); err != nil {
			lastErr = err
			// Small exponential backoff before retrying
			time.Sleep(time.Duration(i+1) * time.Millisecond * 50)
			fmt.Printf("Retrying operation (attempt %d/%d)\n", i+1, maxRetriesForRepost)
		} else {
			return nil
		}
	}
	return fmt.Errorf("operation failed after %d attempts: %v", maxRetriesForRepost, lastErr)
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

type HandlerFunc interface {
	func(byte, []byte, *precompiles.TxParams, *precompiles.CallbackFunc) ([]byte, uint64, error) | // 1 operand
		func(byte, []byte, []byte, *precompiles.TxParams, *precompiles.CallbackFunc) ([]byte, uint64, error) | // 2 operands
		func(byte, []byte, []byte, []byte, *precompiles.TxParams, *precompiles.CallbackFunc) ([]byte, uint64, error) | // 3 operands
		func([]byte, byte, int32, *precompiles.TxParams, *precompiles.CallbackFunc) ([]byte, uint64, error) | // TrivialEncrypt
		func(byte, []byte, byte, *precompiles.TxParams, *precompiles.CallbackFunc) ([]byte, uint64, error) | // Cast
		func(byte, uint64, int32, *precompiles.TxParams, *precompiles.CallbackFunc) ([]byte, uint64, error) // Random
}

type GenericHashRequest struct {
	UType        byte     `json:"utype"`
	Inputs       []string `json:"inputs"`
	RequesterUrl string   `json:"requesterUrl"`
}

func handleRequest[T HandlerFunc](w http.ResponseWriter, r *http.Request, handler T) {
	fmt.Printf("Got a request from %s\n", r.RemoteAddr)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Failed to read request body: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var req GenericHashRequest
	if err := json.Unmarshal(body, &req); err != nil {
		log.Printf("Failed to unmarshal request: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get the number of expected inputs based on handler type
	handlerType := reflect.TypeOf(handler)
	expectedInputs := handlerType.NumIn() - 3 // subtract utype, txParams, and callback

	// Convert all hex strings to byte arrays
	decodedInputs := [][]byte{}
	for i, hexStr := range req.Inputs {
		if hexStr == "" {
			continue
		}
		decoded, err := hex.DecodeString(hexStr)
		if err != nil {
			e := fmt.Sprintf("Invalid input hex string at position %d: %s %+v", i, hexStr, err)
			fmt.Println(e)
			http.Error(w, e, http.StatusBadRequest)
			return
		}
		decodedInputs = append(decodedInputs, decoded)
	}

	if len(decodedInputs) != expectedInputs {
		log.Printf("Handler expects %d inputs, got %d", expectedInputs, len(decodedInputs))
		http.Error(w, fmt.Sprintf("Handler expects %d inputs, got %d", expectedInputs, len(decodedInputs)), http.StatusBadRequest)
		return
	}

	callback := precompiles.CallbackFunc{
		CallbackUrl: req.RequesterUrl,
		Callback:    handleResult,
	}

	// Prepare the arguments for the handler call
	args := make([]reflect.Value, handlerType.NumIn())
	args[0] = reflect.ValueOf(req.UType)
	for i, input := range decodedInputs {
		args[i+1] = reflect.ValueOf(input)
	}
	args[len(args)-2] = reflect.ValueOf(&tp)
	args[len(args)-1] = reflect.ValueOf(&callback)

	// Call the handler with the appropriate number of inputs
	results := reflect.ValueOf(handler).Call(args)

	// TODO : handle gasUsed at index 1
	result, _, errInterface := results[0].Interface().([]byte), results[1].Interface().(uint64), results[2].Interface()
	if errInterface != nil {
		err = errInterface.(error)
		e := fmt.Sprintf("Operation failed: %+v", err)
		fmt.Println(e)
		http.Error(w, e, http.StatusInternalServerError)
		return
	}

	res := []byte(hex.EncodeToString(result))
	w.Write(res)
	fmt.Printf("Started processing the request for tempkey %s\n", hex.EncodeToString(result))
}

// Helper function to handle decoding the request and calling the respective function
func handleTwoRequest(w http.ResponseWriter, r *http.Request, handler func(byte, []byte, []byte, *precompiles.TxParams, *precompiles.CallbackFunc) ([]byte, uint64, error)) {
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

	// var trivialHash []byte
	// for i := 0; i <= 50; i++ {

	// 	// Create a byte slice of size 32
	// 	toEncrypt := make([]byte, 32)

	// 	// Convert the integer to bytes and store it in the byte slice
	// 	toEncrypt[31] = uint8(i)
	// 	trivialHash, _, err = precompiles.TrivialEncrypt(toEncrypt, 2, 0, &tp, nil)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("failed to generate trivial hash for %d: %v", i, err)
	// 	}
	// 	fmt.Printf("Trivial hash for %d: %x\n", i, trivialHash)
	// }
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

func TrivialEncryptHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Got a request from %s\n", r.RemoteAddr)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var req TrivialEncryptRequest
	if err := json.Unmarshal(body, &req); err != nil {
		fmt.Printf("Failed unmarsheling request: %+v body is %+v\n", err, string(body))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	callback := precompiles.CallbackFunc{
		CallbackUrl: req.RequesterUrl,
		Callback:    handleResult,
	}

	// Convert the value strings to byte arrays
	hexStr := fmt.Sprintf("%064x", req.Value)
	value, err := hex.DecodeString(hexStr)
	if err != nil {
		e := fmt.Sprintf("Invalid value: %s %+v", req.Value, err)
		fmt.Println(e)
		http.Error(w, e, http.StatusBadRequest)
		return
	}

	result, _, err := precompiles.TrivialEncrypt(value, req.ToType, req.SecurityZone, &tp, &callback)
	if err != nil {
		e := fmt.Sprintf("Operation failed: %+v", err)
		fmt.Println(e)
		http.Error(w, e, http.StatusBadRequest)
		return
	}

	// Respond with the result
	res := []byte(hex.EncodeToString(result))
	w.Write(res)
	fmt.Printf("Started processing the request for tempkey %s\n", hex.EncodeToString(result))
}

func main() {
	_, err := initFheos()
	if err != nil {
		log.Fatalf("Failed to initialize FHEOS: %v", err)
		os.Exit(1)
	}

	handlers := getHandlers()
	log.Printf("Got %d handlers", len(handlers))
	// iterate handlers
	for i, handler := range handlers {
		http.HandleFunc(handler.Name, handler.Handler)
		log.Printf("Added handler for %s in index %d", handler.Name, i)
	}

	http.HandleFunc("/Decrypt", DecryptHandler)
	log.Printf("Added handler for /Decrypt")
	http.HandleFunc("/SealOutput", SealOutputHandler)
	log.Printf("Added handler for /SealOutput")
	http.HandleFunc("/TrivialEncrypt", TrivialEncryptHandler)
	log.Printf("Added handler for /TrivialEncrypt")

	// Start the server
	log.Println("Server listening on port 8448...")
	if err := http.ListenAndServe(":8448", nil); err != nil {
		log.Fatalf("Server stopped: %v", err)
	}
}
