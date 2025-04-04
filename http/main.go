package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"flag"
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
	"github.com/fhenixprotocol/fheos/telemetry"
	fhedriver "github.com/fhenixprotocol/warp-drive/fhe-driver"
)

const (
	maxRetriesForRepost = 3
)

var tp precompiles.TxParams
var telemetryCollector *telemetry.TelemetryCollector

func getEventId(url string) string {
	parts := strings.Split(url, "_")
	if len(parts) == 0 {
		return ""
	}
	return parts[len(parts)-1]
}

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

func handleResult(url string, hadError bool, tempKey []byte, actualHash []byte) {
	fmt.Printf("Got hash result for %s : %s\n", hex.EncodeToString(tempKey), hex.EncodeToString(actualHash))
	// JSON data to be sent in the request body
	jsonData, err := json.Marshal(FheOperationResponse{TempKey: tempKey, ActualHash: actualHash})
	if err != nil {
		log.Printf("Failed to marshal update for requester %s with the result of %+v: %v", url, tempKey, err)
		return
	}

	responseToServer(url, tempKey, jsonData)
}

func handleDecryptResult(url string, ctHash []byte, plaintext *big.Int, transactionHash string, chainId uint64) {
	fmt.Printf("Got decrypt result for %s : %s\n", hex.EncodeToString(ctHash), plaintext)
	plaintextString := plaintext.Text(16)
	jsonData, err := json.Marshal(DecryptResponse{CtHash: ctHash, Plaintext: plaintextString, TransactionHash: transactionHash})
	if err != nil {
		log.Printf("Failed to marshal decrypt result for requester %s with the result of %+v: %v", url, ctHash, err)
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

func (g *FheOperationRequest) UnmarshalJSON(data []byte) error {
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

func inputsToTelemetry(inputs []fhedriver.CiphertextKey) string {
	var inputStrings []string
	for _, input := range inputs {
		inputStrings = append(inputStrings, hex.EncodeToString(input.Hash[:]))
	}
	return strings.Join(inputStrings, ", ")
}

func handleRequest[T HandlerFunc](w http.ResponseWriter, r *http.Request, handler T) {
	fmt.Printf("Got a FHE operation request from %s\n", r.RemoteAddr)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Failed to read request body: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var req FheOperationRequest
	if err := json.Unmarshal(body, &req); err != nil {
		log.Printf("Failed to unmarshal GenericHashRequest: %v, %+v", err, req)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	handlerType := reflect.TypeOf(handler)
	eventId := getEventId(req.RequesterUrl)
	telemetryCollector.AddTelemetry(telemetry.FheOperationRequestTelemetry{
		TelemetryType: "fhe_operation_request",
		OperationType: "generic_fhe_operation",
		ID:            eventId,
		Handle:        "",
		Inputs:        inputsToTelemetry(req.Inputs),
	})

	log.Printf("Request: %+v\n", req)

	// Get the number of expected inputs based on handler type
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
		EventId:     eventId,
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
		telemetryCollector.AddTelemetry(telemetry.FheOperationUpdateTelemetry{
			TelemetryType:  "fhe_operation_update",
			ID:             eventId,
			InternalHandle: hex.EncodeToString(result),
			Status:         fmt.Sprintf("failed err: %+v", err),
		})
		e := fmt.Sprintf("Operation failed: %+v", err)
		fmt.Println(e)
		http.Error(w, e, http.StatusInternalServerError)
		return
	}

	telemetryCollector.AddTelemetry(telemetry.FheOperationUpdateTelemetry{
		TelemetryType:  "fhe_operation_update",
		ID:             eventId,
		InternalHandle: hex.EncodeToString(result),
		Status:         "sync_part_done",
	})

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

func getEnvOrDefault(key string, defaultValue string) string {
	str := os.Getenv(key)
	if str == "" {
		return defaultValue
	}
	return str
}

func initFheos(configDir string) (*precompiles.TxParams, error) {
	if os.Getenv("FHEOS_DB_PATH") == "" {
		if err := os.Setenv("FHEOS_DB_PATH", "./fheosdb"); err != nil {
			return nil, fmt.Errorf("failed to set FHEOS_DB_PATH: %v", err)
		}
	}

	telemetryCollector = telemetry.NewTelemetryCollector(os.Getenv("FHEOS_TELEMETRY_PATH"))

	var config fhedriver.Config
	if configDir == "" {
		config = fhedriver.ConfigDefault
	} else {
		configData, err := os.ReadFile(configDir + "/fheos_config.json")
		if err != nil {
			return nil, fmt.Errorf("failed to read config file: %v", err)
		}

		if err := json.Unmarshal(configData, &config); err != nil {
			return nil, fmt.Errorf("failed to unmarshal config: %v", err)
		}
	}

	config.OracleType = getEnvOrDefault("FHEOS_ORACLE_TYPE", "local")
	if err := precompiles.InitFheConfig(&config, telemetryCollector); err != nil {
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
		UType           byte             `json:"utype"`
		Key             CiphertextKeyAux `json:"key"`
		RequesterUrl    string           `json:"requesterUrl"`
		TransactionHash string           `json:"transactionHash"`
		ChainId         uint64           `json:"chainId"`
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
	d.ChainId = uint64(aux.ChainId)

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

	eventId := getEventId(req.RequesterUrl)
	telemetryCollector.AddTelemetry(telemetry.FheOperationRequestTelemetry{
		TelemetryType: "fhe_operation_request",
		OperationType: "decrypt",
		ID:            eventId,
		Handle:        hex.EncodeToString(req.Key.Hash[:]),
		Inputs:        fmt.Sprintf("chain_id: %d", req.ChainId),
	})

	log.Printf("Decrypt Request: %+v\n", req)
	// Convert the hash strings to byte arrays
	callback := precompiles.DecryptCallbackFunc{
		CallbackUrl:     req.RequesterUrl,
		Callback:        handleDecryptResult,
		TransactionHash: req.TransactionHash,
		ChainId:         req.ChainId,
		EventId:         eventId,
	}

	_, _, err = precompiles.Decrypt(req.UType, fhedriver.SerializeCiphertextKey(req.Key), nil, &tp, &callback)
	if err != nil {
		telemetryCollector.AddTelemetry(telemetry.FheOperationUpdateTelemetry{
			TelemetryType:  "fhe_operation_update",
			ID:             eventId,
			InternalHandle: hex.EncodeToString(req.Key.Hash[:]),
			Status:         fmt.Sprintf("failed err: %+v", err),
		})
		e := fmt.Sprintf("Operation failed: %+v", err)
		fmt.Println(e)
		http.Error(w, e, http.StatusBadRequest)
		return
	}
	telemetryCollector.AddTelemetry(telemetry.FheOperationUpdateTelemetry{
		TelemetryType:  "fhe_operation_update",
		ID:             eventId,
		InternalHandle: hex.EncodeToString(req.Key.Hash[:]),
		Status:         "sync_part_done",
	})

	// Respond with the result
	w.Write(req.Key.Hash[:])
	fmt.Printf("Received decrypt request for %+v and type %+v\n", hex.EncodeToString(req.Key.Hash[:]), req.UType)
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
	eventId := getEventId(req.RequesterUrl)
	telemetryCollector.AddTelemetry(telemetry.FheOperationRequestTelemetry{
		TelemetryType: "fhe_operation_request",
		OperationType: "trivial_encrypt",
		ID:            eventId,
		Handle:        "",
		Inputs:        fmt.Sprintf("security_zone: %d, utype: %d, value: %s", req.SecurityZone, req.ToType, req.Value),
	})

	callback := precompiles.CallbackFunc{
		CallbackUrl: req.RequesterUrl,
		Callback:    handleResult,
		EventId:     eventId,
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
		telemetryCollector.AddTelemetry(telemetry.FheOperationUpdateTelemetry{
			TelemetryType:  "fhe_operation_update",
			ID:             eventId,
			InternalHandle: hex.EncodeToString(result),
			Status:         fmt.Sprintf("failed err: %+v", err),
		})
		e := fmt.Sprintf("Operation failed: %+v", err)
		fmt.Println(e)
		http.Error(w, e, http.StatusBadRequest)
		return
	}

	telemetryCollector.AddTelemetry(telemetry.FheOperationUpdateTelemetry{
		TelemetryType:  "fhe_operation_update",
		ID:             eventId,
		InternalHandle: hex.EncodeToString(result),
		Status:         "sync_part_done",
	})

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
	eventId := getEventId(req.RequesterUrl)

	telemetryCollector.AddTelemetry(telemetry.FheOperationRequestTelemetry{
		TelemetryType: "fhe_operation_request",
		OperationType: "cast",
		ID:            eventId,
		Handle:        hex.EncodeToString(req.Input.Hash[:]),
		Inputs:        fmt.Sprintf("security_zone: %d, utype: %d", req.Input.SecurityZone, req.UType),
	})

	callback := precompiles.CallbackFunc{
		CallbackUrl: req.RequesterUrl,
		Callback:    handleResult,
		EventId:     eventId,
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
		telemetryCollector.AddTelemetry(telemetry.FheOperationUpdateTelemetry{
			TelemetryType:  "fhe_operation_update",
			ID:             eventId,
			InternalHandle: hex.EncodeToString(result),
			Status:         fmt.Sprintf("failed err: %+v", err),
		})
		e := fmt.Sprintf("Operation failed: %+v", err)
		fmt.Println(e)
		http.Error(w, e, http.StatusBadRequest)
		return
	}

	telemetryCollector.AddTelemetry(telemetry.FheOperationUpdateTelemetry{
		TelemetryType:  "fhe_operation_update",
		ID:             eventId,
		InternalHandle: hex.EncodeToString(result),
		Status:         "sync_part_done",
	})

	// Respond with the result
	res := []byte(hex.EncodeToString(result))
	w.Write(res)
	fmt.Printf("Started processing the request for tempkey %s\n", hex.EncodeToString(result))
}

func createNetworkPublicKeyResponse(PublicKey []byte) ([]byte, error) {
	result := GetNetworkPublicKeyResponse{
		PublicKey: hex.EncodeToString(PublicKey),
	}

	responseData, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response: %+v", err)
	}
	return responseData, nil
}

func createCrsResponse(Crs []byte) ([]byte, error) {
	result := GetCrsResponse{
		Crs: hex.EncodeToString(Crs),
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

func StoreCtsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Got a verify request from %s\n", r.RemoteAddr)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var req StoreCtsRequest
	if err := json.Unmarshal(body, &req); err != nil {
		fmt.Printf("Failed unmarshaling request: %+v body is %+v\n", err, string(body))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// This is a placeholder for the actual signature verification
	if req.Signature != "toml" {
		http.Error(w, "Invalid signature", http.StatusBadRequest)
		return
	}

	hashes := []string{}

	for _, ct := range req.Cts {
		value, err := hex.DecodeString(hexOnly(ct.Value))
		if err != nil {
			e := fmt.Sprintf("Invalid Value: %s %+v", ct.Value, err)
			fmt.Println(e)
			http.Error(w, e, http.StatusBadRequest)
			return
		}

		hash, _, err := precompiles.StoreCt(ct.UType, value, int32(ct.SecurityZone), &tp, nil)
		if err != nil {
			e := fmt.Sprintf("Operation failed: %+v", err)
			fmt.Println(e)
			http.Error(w, e, http.StatusBadRequest)
			return
		}
		hashes = append(hashes, hex.EncodeToString(hash))
	}

	fmt.Printf("Stored %d cts: %+v\n", len(hashes), hashes)

	result := StoreCtsResponse{
		Hashes: hashes,
	}

	responseData, err := json.Marshal(result)
	if err != nil {
		fmt.Printf("Failed to marshal response: %+v\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = w.Write(responseData)
	if err != nil {
		fmt.Printf("Failed to write response: %v\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// If it's a preflight OPTIONS request, respond OK
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

func GetCrsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Got a getCrs request from %s\n", r.RemoteAddr)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var req GetCrsRequest
	if err := json.Unmarshal(body, &req); err != nil {
		fmt.Printf("Failed unmarshaling request: %+v body is %+v\n", err, string(body))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	expectedCrs, err := precompiles.GetCrs(int32(req.SecurityZone), &tp)
	if err != nil {
		e := fmt.Sprintf("Operation failed: %+v", err)
		fmt.Println(e)
		http.Error(w, e, http.StatusBadRequest)
		return
	}

	responseData, err := createCrsResponse(expectedCrs)
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

	var req GetCTRequest
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

	response := GetCTResponse{
		Data:         hex.EncodeToString(ct.Data),
		SecurityZone: uint8(ct.SecurityZone),
		UintType:     uint8(ct.Properties.EncryptionType),
		Compact:      ct.Properties.Compact,
		Gzipped:      ct.Properties.Gzipped,
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
	configDir := flag.String("config-dir", "", "Path to config directory")
	flag.Parse()

	if *configDir != "" {
		log.Printf("Using config directory: %s", *configDir)
	}

	_, err := initFheos(*configDir)
	if err != nil {
		log.Fatalf("Failed to initialize FHEOS: %v", err)
		os.Exit(1)
	}

	// Create two separate muxes for different ports
	publicMux := http.NewServeMux()
	privateMux := http.NewServeMux()

	handlers := getHandlers()
	log.Printf("Got %d handlers", len(handlers))

	// Private endpoints on port 8449
	for _, handler := range handlers {
		privateMux.HandleFunc(handler.Name, handler.Handler)
	}

	privateMux.HandleFunc("/Decrypt", DecryptHandler)
	privateMux.HandleFunc("/StoreCts", StoreCtsHandler)
	privateMux.HandleFunc("/TrivialEncrypt", TrivialEncryptHandler)
	privateMux.HandleFunc("/Cast", CastHandler)

	// Public endpoints on port 8448
	publicMux.HandleFunc("/GetNetworkPublicKey", GetNetworkPublicKeyHandler)
	publicMux.HandleFunc("/GetCrs", GetCrsHandler)
	publicMux.HandleFunc("/GetCT", GetCTHandler)
	publicMux.HandleFunc("/Health", HealthHandler)

	// Wrap both muxes in the CORS middleware
	wrappedPublicMux := corsMiddleware(publicMux)
	wrappedPrivateMux := corsMiddleware(privateMux)

	// Start both servers in separate goroutines
	go func() {
		log.Println("Public server listening on port 8448...")
		if err := http.ListenAndServe(":8448", wrappedPublicMux); err != nil {
			log.Fatalf("Public server stopped: %v", err)
		}
	}()

	log.Println("Private server listening on port 8449...")
	if err := http.ListenAndServe(":8449", wrappedPrivateMux); err != nil {
		log.Fatalf("Private server stopped: %v", err)
	}
}
