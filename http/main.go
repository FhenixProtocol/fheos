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
	"strings"
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

type DecryptRequest struct {
	UType        	byte                    `json:"utype"`
	Key          	fhedriver.CiphertextKey `json:"key"`
	RequesterUrl 	string                  `json:"requesterUrl"`
	TransactionHash string                	`json:"transactionHash"`
	ChainId         uint                   	`json:"chainId"`
}

type SealOutputRequest struct {
	UType        	byte                    `json:"utype"`
	Key          	fhedriver.CiphertextKey `json:"key"`
	PKey         	string                  `json:"pkey"`
	RequesterUrl 	string                  `json:"requesterUrl"`
	TransactionHash string                	`json:"transactionHash"`
	ChainId         uint                   	`json:"chainId"`
}

type MockDecryptRequest struct {
	CtHash string `json:"ctHash"`
	Permit string `json:"permit"`
}

type MockSealOutputRequest struct {
	CtHash    string `json:"ctHash"`
	Permit    string `json:"permit"`
	PublicKey string `json:"publicKey"`
}

type VerifyRequest struct {
	UType        byte   `json:"utype"`
	Value        string `json:"value"`
	SecurityZone byte   `json:"securityZone"`
}

type GetNetworkPublicKeyRequest struct {
	SecurityZone byte `json:"securityZone"`
}

type GetNetworkPublicKeyResult struct {
	PublicKey string `json:"publicKey"`
}

type VerifyResult struct {
	CtHash    string `json:"ctHash"`
	Signature string `json:"signature"`
}

type TrivialEncryptRequest struct {
	Value        *big.Int `json:"value"`
	ToType       byte     `json:"toType"`
	SecurityZone int32    `json:"securityZone"`
	RequesterUrl string   `json:"requesterUrl"`
}

type CastRequest struct {
	UType        byte                    `json:"utype"`
	Input        fhedriver.CiphertextKey `json:"input"`
	ToType       string                  `json:"toType"`
	RequesterUrl string                  `json:"requesterUrl"`
}
type HashResultUpdate struct {
	TempKey    []byte `json:"tempKey"`
	ActualHash []byte `json:"actualHash"`
}

type DecryptResultUpdate struct {
	CtHash    		[]byte `json:"ctHash"`
	Plaintext 		string `json:"plaintext"`
	TransactionHash string `json:"transactionHash"`
}

type SealOutputResultUpdate struct {
	CtHash 			[]byte `json:"ctHash"`
	PK     			string `json:"pk"`
	Value  			string `json:"value"`
	TransactionHash string `json:"transactionHash"`
}

type CTRequest struct {
	Hash string `json:"hash"`
}

type CTResponse struct {
	CiphertextData 			string `json:"ciphertextData"`
	SecurityZone   			uint8  `json:"securityZone"`
	IsTriviallyEncrypted 	bool   `json:"isTriviallyEncrypted"`
	UintType             	uint8  `json:"uintType"`
	Compact              	bool   `json:"compact"`
	Compressed           	bool   `json:"compressed"`
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

	fmt.Printf("Update requester %s with the result of %+v\n", url, hex.EncodeToString(tempKey))
}

func handleResult(url string, tempKey []byte, actualHash []byte) {
	fmt.Printf("Got hash result for %s : %s\n", hex.EncodeToString(tempKey), hex.EncodeToString(actualHash))
	// JSON data to be sent in the request body
	jsonData, err := json.Marshal(HashResultUpdate{TempKey: tempKey, ActualHash: actualHash})
	if err != nil {
		log.Printf("Failed to marshal update for requester %s with the result of %+v: %v", url, tempKey, err)
		return
	}

	responseToServer(url, tempKey, jsonData)
}

func handleDecryptResult(url string, ctHash []byte, plaintext *big.Int, transactionHash string, chainId uint) {
	fmt.Printf("Got decrypt result for %s : %s\n", hex.EncodeToString(ctHash), plaintext)
	plaintextString := plaintext.Text(16)
	jsonData, err := json.Marshal(DecryptResultUpdate{CtHash: ctHash, Plaintext: plaintextString, TransactionHash: transactionHash})
	if err != nil {
		log.Printf("Failed to marshal decrypt result for requester %s with the result of %+v: %v", url, ctHash, err)
		return
	}

	responseToServer(url, ctHash, jsonData)
}

func handleSealOutputResult(url string, ctHash []byte, pk []byte, value string, transactionHash string, chainId uint) {
	fmt.Printf("Got sealoutput result for %s : %s\n", hex.EncodeToString(ctHash), value)
	jsonData, err := json.Marshal(SealOutputResultUpdate{CtHash: ctHash, PK: hex.EncodeToString(pk), Value: value, TransactionHash: transactionHash})
	if err != nil {
		log.Printf("Failed to marshal seal output result for requester %s with the result of %+v: %v", url, ctHash, err)
		return
	}

	responseToServer(url, ctHash, jsonData)
}

type HandlerFunc interface {
	func(byte, []byte, *precompiles.TxParams, *precompiles.CallbackFunc) ([]byte, uint64, error) | // 1 operand
		func(byte, []byte, []byte, *precompiles.TxParams, *precompiles.CallbackFunc) ([]byte, uint64, error) | // 2 operands
		func(byte, []byte, []byte, []byte, *precompiles.TxParams, *precompiles.CallbackFunc) ([]byte, uint64, error) | // 3 operands
		func([]byte, byte, int32, *precompiles.TxParams, *precompiles.CallbackFunc) ([]byte, uint64, error) | // TrivialEncrypt
		func(byte, uint64, int32, *precompiles.TxParams, *precompiles.CallbackFunc) ([]byte, uint64, error) // Random
}

type CiphertextKeyAux struct {
	IsTriviallyEncrypted bool   `json:"IsTriviallyEncrypted"`
	UintType             int    `json:"UintType"` // Assuming EncryptionType is an int or compatible
	SecurityZone         int32  `json:"SecurityZone"`
	Hash                 string `json:"Hash"`
}
type GenericHashRequest struct {
	UType        byte                      `json:"uType"`
	Inputs       []fhedriver.CiphertextKey `json:"inputs"`
	RequesterUrl string                    `json:"requesterUrl"`
}

func convertInput(input CiphertextKeyAux) (*fhedriver.CiphertextKey, error) {
	decoded, err := hex.DecodeString(hexOnly(input.Hash[2:]))
	if err != nil {
		e := fmt.Sprintf("Invalid input hex string at position %s %+v", input.Hash, err)
		return nil, fmt.Errorf(e)
	}

	return &fhedriver.CiphertextKey{
		IsTriviallyEncrypted: input.IsTriviallyEncrypted,
		UintType:             fhedriver.EncryptionType(input.UintType),
		SecurityZone:         input.SecurityZone,
		Hash:                 [32]byte(decoded),
	}, nil
}

func convertInputs(inputs []CiphertextKeyAux) ([]fhedriver.CiphertextKey, error) {
	var convertedInputs []fhedriver.CiphertextKey
	for _, input := range inputs {
		// First two bytes are -x
		if len(input.Hash) == 66 {
			convertedInput, err := convertInput(input)
			if err != nil {
				return nil, err
			}

			convertedInputs = append(convertedInputs, *convertedInput)
		}
	}

	return convertedInputs, nil
}

func (g *GenericHashRequest) UnmarshalJSON(data []byte) error {
	var aux struct {
		UType        byte               `json:"UType"`
		Inputs       []CiphertextKeyAux `json:"Inputs"`
		RequesterUrl string             `json:"RequesterUrl"`
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		log.Printf("Failed to unmarshal GenericHashRequestAux: %v, %+v", err, aux)
		return err
	}

	convertedInputs, err := convertInputs(aux.Inputs)
	if err != nil {
		return err
	}

	g.UType = aux.UType
	g.Inputs = convertedInputs
	g.RequesterUrl = aux.RequesterUrl
	return nil
}

func handleRequest[T HandlerFunc](w http.ResponseWriter, r *http.Request, handler T) {
	fmt.Printf("Got a FHE operation request from %s\n", r.RemoteAddr)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Failed to read request body: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var req GenericHashRequest
	if err := json.Unmarshal(body, &req); err != nil {
		log.Printf("Failed to unmarshal GenericHashRequest: %v, %+v", err, req)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Request: %+v\n", req)

	// Get the number of expected inputs based on handler type
	handlerType := reflect.TypeOf(handler)
	expectedInputs := handlerType.NumIn() - 3 // subtract utype, txParams, and callback

	// Convert all hex strings to byte arrays
	decodedInputs := [][]byte{}
	for _, input := range req.Inputs {
		decodedInputs = append(decodedInputs, fhedriver.SerializeCiphertextKey(input))
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

	return &tp, nil
}

func (d *DecryptRequest) UnmarshalJSON(data []byte) error {
	var aux struct {
		UType        	byte             	`json:"utype"`
		Key          	CiphertextKeyAux 	`json:"key"`
		RequesterUrl 	string           	`json:"requesterUrl"`
		TransactionHash string				`json:"transactionHash"`
		ChainId 		uint				`json:"chainId"`
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		log.Printf("Failed to unmarshal DecryptRequestAux: %v, %+v", err, aux)
		return err
	}

	convertedInput, err := convertInput(aux.Key)
	if err != nil {
		return err
	}

	d.UType = aux.UType
	d.Key = *convertedInput
	d.RequesterUrl = aux.RequesterUrl
	d.TransactionHash = aux.TransactionHash
	d.ChainId = aux.ChainId

	return nil
}

func DecryptHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Got a decrypt request from %s\n", r.RemoteAddr)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var req DecryptRequest
	if err := json.Unmarshal(body, &req); err != nil {
		fmt.Printf("Failed unmarshaling request: %+v body is %+v\n", err, string(body))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Decrypt Request: %+v\n", req)
	// Convert the hash strings to byte arrays
	callback := precompiles.DecryptCallbackFunc{
		CallbackUrl: req.RequesterUrl,
		Callback:    handleDecryptResult,
		TransactionHash: req.TransactionHash,
		ChainId: req.ChainId,
	}

	_, _, err = precompiles.Decrypt(req.UType, fhedriver.SerializeCiphertextKey(req.Key), nil, &tp, &callback)
	if err != nil {
		e := fmt.Sprintf("Operation failed: %+v", err)
		fmt.Println(e)
		http.Error(w, e, http.StatusBadRequest)
		return
	}
	// Respond with the result
	w.Write(req.Key.Hash[:])
	fmt.Printf("Received decrypt request for %+v and type %+v\n", hex.EncodeToString(req.Key.Hash[:]), req.UType)
}

func (d *MockDecryptRequest) UnmarshalJSON(data []byte) error {
	var aux struct {
		CtHash string `json:"ctHash"`
		Permit string `json:"permit"`
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		log.Printf("Failed to unmarshal MockDecryptRequestAux: %v, %+v", err, aux)
		return err
	}

	d.CtHash = aux.CtHash
	d.Permit = aux.Permit

	return nil
}

func DecryptHandlerMock(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Got a decrypt request (currently mocked) from %s\n", r.RemoteAddr)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var req MockDecryptRequest
	if err := json.Unmarshal(body, &req); err != nil {
		fmt.Printf("Failed unmarshaling request: %+v body is %+v\n", err, string(body))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := []byte("42")
	_, err = w.Write(response)
	if err != nil {
		fmt.Printf("Failed to write response: %v\n", err)
		return
	}

	fmt.Printf("Received mock decrypt request for %+v and permit %+v\n", req.CtHash, req.Permit)
}

func (d *MockSealOutputRequest) UnmarshalJSON(data []byte) error {
	var aux struct {
		CtHash    string `json:"ctHash"`
		Permit    string `json:"permit"`
		PublicKey string `json:"publicKey"`
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		log.Printf("Failed to unmarshal MockSealOutputRequestAux: %v, %+v", err, aux)
		return err
	}

	d.CtHash = aux.CtHash
	d.Permit = aux.Permit
	d.PublicKey = aux.PublicKey

	return nil
}

func SealOutputHandlerMock(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Got a sealoutput request (currently mocked) from %s\n", r.RemoteAddr)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var req MockSealOutputRequest
	if err := json.Unmarshal(body, &req); err != nil {
		fmt.Printf("Failed unmarshaling request: %+v body is %+v\n", err, string(body))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := []byte("")
	_, err = w.Write(response)
	if err != nil {
		fmt.Printf("Failed to write response: %v\n", err)
		return
	}

	fmt.Printf("Received mock seal output request for %+v and permit %+v and pubkey %+v\n", req.CtHash, req.Permit, req.PublicKey)
}

func (s *SealOutputRequest) UnmarshalJSON(data []byte) error {
	var aux struct {
		UType        	byte             	`json:"utype"`
		Key          	CiphertextKeyAux 	`json:"key"`
		PKey         	string           	`json:"pkey"`
		RequesterUrl 	string           	`json:"requesterUrl"`
		TransactionHash string 				`json:"transactionHash"`
		ChainId 		uint  				`json:"chainId"`
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		log.Printf("Failed to unmarshal SealOutputRequestAux: %v, %+v", err, aux)
		return err
	}

	convertedInput, err := convertInput(aux.Key)
	if err != nil {
		return err
	}

	s.UType = aux.UType
	s.Key = *convertedInput
	s.PKey = aux.PKey
	s.RequesterUrl = aux.RequesterUrl
	s.TransactionHash = aux.TransactionHash
	s.ChainId = aux.ChainId

	return nil
}
func SealOutputHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Got a sealoutput request from %s\n", r.RemoteAddr)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var req SealOutputRequest
	if err := json.Unmarshal(body, &req); err != nil {
		fmt.Printf("Failed unmarshaling request: %+v body is %+v\n", err, string(body))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("SealOutput Request: %+v\n", req)
	callback := precompiles.SealOutputCallbackFunc{
		CallbackUrl: req.RequesterUrl,
		Callback:    handleSealOutputResult,
		TransactionHash: req.TransactionHash,
		ChainId: req.ChainId,
	}

	pkey, err := hex.DecodeString(hexOnly(req.PKey))
	if err != nil {
		e := fmt.Sprintf("Invalid pkey: %s %+v", req.PKey, err)
		fmt.Println(e)
		http.Error(w, e, http.StatusBadRequest)
		return
	}

	_, _, err = precompiles.SealOutput(req.UType, fhedriver.SerializeCiphertextKey(req.Key), pkey, &tp, &callback)
	if err != nil {
		e := fmt.Sprintf("Operation failed: %+v", err)
		fmt.Println(e)
		http.Error(w, e, http.StatusBadRequest)
		return
	}
	// Respond with the result
	w.Write(req.Key.Hash[:])
	fmt.Printf("Received seal output request for %+v and type %+v\n", hex.EncodeToString(req.Key.Hash[:]), req.UType)
}

func (t *TrivialEncryptRequest) UnmarshalJSON(data []byte) error {
	var aux struct {
		Value        string `json:"value"`
		ToType       byte   `json:"toType"`
		SecurityZone int32  `json:"securityZone"`
		RequesterUrl string `json:"requesterUrl"`
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		log.Printf("Failed to unmarshal TrivialEncryptRequestAux: %v, %+v", err, aux)
		return err
	}

	val, ok := new(big.Int).SetString(aux.Value[2:], 16)
	if !ok {
		return fmt.Errorf("invalid hex value: %s", aux.Value)
	}

	t.Value = val
	t.ToType = aux.ToType
	t.SecurityZone = aux.SecurityZone
	t.RequesterUrl = aux.RequesterUrl

	return nil
}

func TrivialEncryptHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Got a trivial encrypt request from %s\n", r.RemoteAddr)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var req TrivialEncryptRequest
	if err := json.Unmarshal(body, &req); err != nil {
		fmt.Printf("Failed unmarshaling TrivialEncryptRequest: %+v body is %+v\n", err, string(body))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("TrivialEncrypt Request: %+v\n", req)
	callback := precompiles.CallbackFunc{
		CallbackUrl: req.RequesterUrl,
		Callback:    handleResult,
	}

	// Convert the value strings to byte arrays
	hexStr := fmt.Sprintf("%064x", req.Value)
	value, err := hex.DecodeString(hexOnly(hexStr))
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

func (s *CastRequest) UnmarshalJSON(data []byte) error {
	var aux struct {
		UType        byte             `json:"utype"`
		Input        CiphertextKeyAux `json:"input"`
		ToType       string           `json:"toType"`
		RequesterUrl string           `json:"requesterUrl"`
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		log.Printf("Failed to unmarshal CastRequestAux: %v, %+v", err, aux)
		return err
	}

	convertedInput, err := convertInput(aux.Input)
	if err != nil {
		return err
	}

	s.UType = aux.UType
	s.Input = *convertedInput
	s.ToType = aux.ToType
	s.RequesterUrl = aux.RequesterUrl

	return nil
}

func CastHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Got a cast request from %s\n", r.RemoteAddr)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var req CastRequest
	if err := json.Unmarshal(body, &req); err != nil {
		fmt.Printf("Failed unmarshaling request: %+v body is %+v\n", err, string(body))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Cast Request: %+v\n", req)
	callback := precompiles.CallbackFunc{
		CallbackUrl: req.RequesterUrl,
		Callback:    handleResult,
	}

	// Convert the value strings to byte arrays
	toTypeInt, err := strconv.Atoi(req.ToType)
	if err != nil {
		e := fmt.Sprintf("Invalid toType: %s %+v", req.ToType, err)
		fmt.Println(e)
		http.Error(w, e, http.StatusBadRequest)
		return
	}

	result, _, err := precompiles.Cast(req.UType, fhedriver.SerializeCiphertextKey(req.Input), byte(toTypeInt), &tp, &callback)
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

func createVerifyResponse(ctHash []byte) ([]byte, error) {
	verifyResult := VerifyResult{
		CtHash:    hex.EncodeToString(ctHash),
		Signature: "Haim",
	}

	responseData, err := json.Marshal(verifyResult)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response: %+v", err)
	}
	return responseData, nil
}

func createNetworkPublicKeyResponse(PublicKey []byte) ([]byte, error) {
	result := GetNetworkPublicKeyResult{
		PublicKey: hex.EncodeToString(PublicKey),
	}

	responseData, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response: %+v", err)
	}
	return responseData, nil
}

func hasHexPrefix(str string) bool {
	return len(str) >= 2 && strings.HasPrefix(strings.ToLower(str), "0x")
}

func hexOnly(value string) string {
	if hasHexPrefix(value) {
		return value[2:]
	}
	return value
}

func UpdateCTHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Got a verify request from %s\n", r.RemoteAddr)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var req VerifyRequest
	if err := json.Unmarshal(body, &req); err != nil {
		fmt.Printf("Failed unmarshaling request: %+v body is %+v\n", err, string(body))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// fmt.Printf("Verify Request Value: %s\n", req.Value)
	value, err := hex.DecodeString(hexOnly(req.Value))
	if err != nil {
		e := fmt.Sprintf("Invalid Value: %s %+v", req.Value, err)
		fmt.Println(e)
		http.Error(w, e, http.StatusBadRequest)
		return
	}

	ctHash, _, err := precompiles.Verify(req.UType, value, int32(req.SecurityZone), &tp, nil)
	if err != nil {
		e := fmt.Sprintf("Operation failed: %+v", err)
		fmt.Println(e)
		http.Error(w, e, http.StatusBadRequest)
		return
	}

	responseData, err := createVerifyResponse(ctHash)
	if err != nil {
		e := fmt.Sprintf("Failed to marshal response: %+v", err)
		fmt.Println(e)
		http.Error(w, e, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseData)
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// If it’s a preflight OPTIONS request, respond OK
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Proceed with the next handler
		next.ServeHTTP(w, r)
	})
}

func GetNetworkPublicKeyHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Got a getNetworkPublicKey request from %s\n", r.RemoteAddr)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var req GetNetworkPublicKeyRequest
	if err := json.Unmarshal(body, &req); err != nil {
		fmt.Printf("Failed unmarshaling request: %+v body is %+v\n", err, string(body))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	expectedPk, err := precompiles.GetNetworkPublicKey(int32(req.SecurityZone), &tp)
	if err != nil {
		e := fmt.Sprintf("Operation failed: %+v", err)
		fmt.Println(e)
		http.Error(w, e, http.StatusBadRequest)
		return
	}

	responseData, err := createNetworkPublicKeyResponse(expectedPk)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseData)
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Got a health request from %s\n", r.RemoteAddr)
	response := map[string]bool{"success": true}

	value, err := json.Marshal(response)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(value)
}

func GetCTHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Got a GetCT request from %s\n", r.RemoteAddr)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var req CTRequest
	if err := json.Unmarshal(body, &req); err != nil {
		fmt.Printf("Failed unmarshaling request: %+v body is %+v\n", err, string(body))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("GetCT Request: %+v\n", req)

	if req.Hash == "" {
		http.Error(w, "Hash is not supported yet", http.StatusBadRequest)
		return
	}

	// Decode the hash from hex
	hash, err := hex.DecodeString(hexOnly(req.Hash))
	if err != nil {
		e := fmt.Sprintf("Invalid hash: %s %+v", req.Hash, err)
		fmt.Println(e)
		http.Error(w, e, http.StatusBadRequest)
		return
	}

	// Get the ciphertext from the ciphertext database
	ct, err := precompiles.GetCT(hash, &tp)
	if err != nil {
		e := fmt.Sprintf("Failed to get ciphertext: %+v", err)
		fmt.Println(e)
		if strings.Contains(err.Error(), "placeholder") {
			http.Error(w, e, http.StatusPreconditionRequired)
		} else {
			http.Error(w, e, http.StatusBadRequest)
		}
		return
	}

	response := CTResponse{
		CiphertextData: hex.EncodeToString(ct.Data),
		SecurityZone:   uint8(ct.Key.SecurityZone),
		IsTriviallyEncrypted: ct.Key.IsTriviallyEncrypted,
		UintType:           uint8(ct.Key.UintType),
		Compact:            ct.Compact,
		Compressed:         ct.Compressed,
	}

	responseData, err := json.Marshal(response)
	if err != nil {
		e := fmt.Sprintf("Failed to marshal response: %+v", err)
		fmt.Println(e)
		http.Error(w, e, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseData)
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
	for _, handler := range handlers {
		http.HandleFunc(handler.Name, handler.Handler)
	}

	http.HandleFunc("/Decrypt", DecryptHandler)
	http.HandleFunc("/SealOutput", SealOutputHandler)
	http.HandleFunc("/QueryDecrypt", DecryptHandlerMock)
	http.HandleFunc("/QuerySealOutput", SealOutputHandlerMock)
	http.HandleFunc("/UpdateCT", UpdateCTHandler)
	http.HandleFunc("/TrivialEncrypt", TrivialEncryptHandler)
	http.HandleFunc("/Cast", CastHandler)
	http.HandleFunc("/GetNetworkPublicKey", GetNetworkPublicKeyHandler)
	http.HandleFunc("/Health", HealthHandler)
	http.HandleFunc("/GetCT", GetCTHandler)

	// Wrap the default mux in the CORS middleware
	wrappedMux := corsMiddleware(http.DefaultServeMux)

	// Start the server
	log.Println("Server listening on port 8448...")
	if err := http.ListenAndServe(":8448", wrappedMux); err != nil {
		log.Fatalf("Server stopped: %v", err)
	}
}
