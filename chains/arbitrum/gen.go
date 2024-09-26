// Copyright 2021-2022, Offchain Labs, Inc.
// For license information, see https://github.com/nitro/blob/master/LICENSE

package arbitrum

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type HardHatArtifact struct {
	Format       string        `json:"_format"`
	ContractName string        `json:"contractName"`
	SourceName   string        `json:"sourceName"`
	Abi          []interface{} `json:"abi"`
	Bytecode     string        `json:"bytecode"`
}

type moduleInfo struct {
	contractNames []string
	abis          []string
	bytecodes     []string
}

func (m *moduleInfo) addArtifact(artifact HardHatArtifact) {
	abi, err := json.Marshal(artifact.Abi)
	if err != nil {
		log.Fatal(err)
	}
	m.contractNames = append(m.contractNames, artifact.ContractName)
	m.abis = append(m.abis, string(abi))
	m.bytecodes = append(m.bytecodes, artifact.Bytecode)
}

func (m *moduleInfo) exportABIs(dest string) {
	for i, name := range m.contractNames {
		path := filepath.Join(dest, name+".abi")
		abi := m.abis[i] + "\n"

		// #nosec G306
		err := os.WriteFile(path, []byte(abi), 0o644)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func capitalizeFirstLetter(input string) string {
	if input == "" {
		return input
	}

	// Split the input into the first letter and the rest of the string
	firstLetter := string(input[0])
	restOfInput := input[1:]

	// Capitalize the first letter and concatenate it with the rest of the string
	capitalizedFirstLetter := strings.ToUpper(firstLetter)
	capitalizedString := capitalizedFirstLetter + restOfInput

	return capitalizedString
}

func uncapitalizeFirstLetter(input string) string {
	if input == "" {
		return input
	}

	// Split the input into the first letter and the rest of the string
	firstLetter := string(input[0])
	restOfInput := input[1:]

	// Capitalize the first letter and concatenate it with the rest of the string
	capitalizedFirstLetter := strings.ToLower(firstLetter)
	capitalizedString := capitalizedFirstLetter + restOfInput

	return capitalizedString
}

type Param struct {
	Name string
	Type string
}

func StringInList(s string, list []string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}

func CreateTemplate(dirPath string) {
	path := filepath.Join(dirPath, "precompiles", "contracts.go")
	outPath := filepath.Join(dirPath, "FheOs_gen.sol")
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}

	outFile, err := os.Create(outPath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}

	defer file.Close()
	defer outFile.Close()

	outFile.WriteString(`// SPDX-License-Identifier: BSD-3-Clause-Clear
// solhint-disable one-contract-per-file
pragma solidity >=0.8.13 <0.9.0;

library Precompiles {
    //solhint-disable const-name-snakecase
    address public constant Fheos = address(128);
}
`)

	outFile.WriteString(`
interface FheOps {
`)

	// Create a bufio scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	// Iterate over each line in the file
	hadBanner := false
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "=====") {
			hadBanner = true
			continue
		}

		if !hadBanner {
			continue
		}

		Name := ""
		Ret := ""
		var params []Param
		if strings.Contains(line, "func ") && !strings.Contains(line, "go func ") {
			chunks := strings.Fields(line)
			inner := strings.Split(chunks[1], "(")
			Name = inner[0]

			param := Param{
				Name: "",
				Type: "",
			}
			Break := false
			if inner[1] != ")" {
				param.Name = inner[1]
				if strings.Contains(chunks[2], ")") {
					param.Type = strings.Replace(chunks[2], ")", "", 1)
					Break = true
				} else {
					param.Type = strings.Replace(chunks[2], ",", "", 1)
				}
				params = append(params, param)

			}

			i := 3
			for ; i < len(chunks); i += 2 {
				if Break {
					break
				}
				param.Name = strings.Replace(chunks[i], ",", "", 1)
				if strings.Contains(chunks[i+1], ")") {
					param.Type = strings.Replace(chunks[i+1], ")", "", 1)
					Break = true
				} else {
					param.Type = strings.Replace(chunks[i+1], ",", "", 1)
				}
				params = append(params, param)
			}

			if strings.Contains(chunks[i], "(") {
				Ret = strings.Replace(chunks[i], "(", "", 1)
				Ret = strings.Replace(Ret, ",", "", 1)
			}

			outLine := "    function " + uncapitalizeFirstLetter(Name) + "("
			for count, param := range params {
				if param.Type == "[]byte" {
					param.Type = "bytes memory"
				}

				if param.Type == "string" {
					param.Type = "string memory"
				}

				if param.Type == "byte" {
					param.Type = "uint8"
				}

				if param.Type == "*big.Int" {
					param.Type = "uint256"
				}

				if param.Type == "*TxParams" || param.Type == "*CallbackFunc" {
					continue
				}

				outLine += param.Type + " " + param.Name

				// Is it the last (Ignoring the TxParams)
				if count < len(params)-3 {
					outLine += ", "
				}
			}

			outLine += ") external pure"
			if Ret != "" && Name != "Log" {
				if Ret == "[]byte" {
					Ret = "bytes memory"
				}
				if Ret == "*big.Int" {
					Ret = "uint256"
				}
				if Ret == "string" {
					Ret = "string memory"
				}
				outLine += " returns (" + Ret + ")"
			}

			outLine += ";\n"

			_, err = outFile.WriteString(outLine)
			if err != nil {
				fmt.Println("Error writing to file:", err)
				return
			}

		}
	}

	outFile.WriteString("}")

	// Check for scanner errors
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}
}

func Gen(parent string, output string) {
	filePaths, err := filepath.Glob(filepath.Join(parent, "precompiles", "artifacts", "contracts", "*", "*.json"))
	fmt.Printf("%s, %+v\n", parent, filePaths)
	if err != nil {
		log.Fatal(err)
	}

	modules := make(map[string]*moduleInfo)

	for _, path := range filePaths {
		if strings.Contains(path, ".dbg.json") {
			continue
		}

		dir, file := filepath.Split(path)
		dir, _ = filepath.Split(dir[:len(dir)-1])
		_, module := filepath.Split(dir[:len(dir)-1])
		module = strings.ReplaceAll(module, "-", "_")
		module += "gen"

		name := file[:len(file)-5]

		data, err := os.ReadFile(path)
		if err != nil {
			log.Fatal("could not read", path, "for contract", name, err)
		}

		artifact := HardHatArtifact{}
		if err := json.Unmarshal(data, &artifact); err != nil {
			log.Fatal("failed to parse contract", name, err)
		}

		if strings.Contains(fmt.Sprintf("%+v", artifact), "BatchPoster") {
			fmt.Printf("File: %v, Module %v, Artifact %+v\n", path, module, artifact)
		}
		modInfo := modules[module]
		if modInfo == nil {
			modInfo = &moduleInfo{}
			modules[module] = modInfo
		}
		modInfo.addArtifact(artifact)
	}

	abi := ""

	for module, info := range modules {

		code, err := bind.Bind(
			info.contractNames,
			info.abis,
			info.bytecodes,
			nil,
			module,
			bind.LangGo,
			nil,
			nil,
		)

		code = "// This file is auto-generated by chains/arbitrum/gen.go\n" + code

		if err != nil {
			log.Fatal(err)
		}

		abi = info.abis[0]

		folder := filepath.Join(output, module)

		err = os.MkdirAll(folder, 0o755)
		if err != nil {
			log.Fatal(err)
		}

		/*
			#nosec G306
			This file contains no private information so the permissions can be lenient
		*/
		fmt.Printf("%s\n", filepath.Join(folder, module+".go"))
		err = os.WriteFile(filepath.Join(folder, module+".go"), []byte(code), 0o644)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("successfully generated go abi files")

	blockscout := filepath.Join(parent, "nitro-testnode", "blockscout", "init", "data")
	if _, err := os.Stat(blockscout); err != nil {
		fmt.Println("skipping abi export since blockscout is not present")
	} else {
		modules["precompilesgen"].exportABIs(blockscout)
		modules["node_interfacegen"].exportABIs(blockscout)
		fmt.Println("successfully exported abi files")
	}

	var functions []Function

	err = json.Unmarshal([]byte(abi), &functions)
	if err != nil {
		fmt.Printf("Error unmarshaling JSON: %v\n", err)
		return
	}

	var operations []Operation
	for _, f := range functions {
		parameters := ""
		innerParameters := ""

		if len(f.Inputs) > 0 {
			parameters += ", "
		}
		for count, arg := range f.Inputs {
			t := arg.Type
			if t == "bytes" {
				t = "[]byte"
			}

			if t == "uint8" {
				t = "byte"
			}

			if t == "uint256" {
				t = "*big.Int"
			}

			if t == "*TxParams" || t == "*CallbackFunc" {
				continue
			}

			parameters += fmt.Sprintf("%s %s", arg.Name, t)
			innerParameters += arg.Name
			if count < (len(f.Inputs) - 1) {
				parameters += ", "
			}
			if f.Name != "log" {
				innerParameters += ", "
			}
		}

		ret := "void"
		if len(f.Outputs) != 0 {
			ret = f.Outputs[0].Type
			if ret == "bytes" {
				ret = "[]byte"
			}

			if ret == "uint256" {
				ret = "*big.Int"
			}
		}

		operations = append(operations, Operation{
			Name:        capitalizeFirstLetter(f.Name),
			Inputs:      parameters,
			InnerInputs: innerParameters,
			ReturnType:  ret,
		})
	}

	f := "FheOps_gen.go"
	file, err := os.Create(f)
	if err != nil {
		log.Fatal(err)
	}

	file.WriteString(`// This file is auto-generated by chains/arbitrum/gen.go

package precompiles

import (
	"fmt"
	"github.com/ethereum/go-ethereum/metrics"
	fheos "github.com/fhenixprotocol/fheos/precompiles"
	"math/big"
	"time"
)

type FheOps struct {
	Address addr // 0x80
}
`)
	handlers, err := os.Create("http/handlers.auto.go")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()
	defer handlers.Close()

	handlers.WriteString(`package http
import (
	"github.com/fhenixprotocol/fheos/precompiles"
	"net/http"
)

type HandlerDef struct {
	Name    string
	Handler func(w http.ResponseWriter, r *http.Request)
}

`)

	getHandlersFunc := `func getHandlers() []HandlerDef {
		return []HandlerDef{
`

	logTemplate := GenerateLogTemplate()
	_, err = file.WriteString(logTemplate)
	if err != nil {
		fmt.Println("Error writing")
		return
	}

	for _, op := range operations {
		if op.Name == "Log" {
			continue
		}

		var (
			template     *template.Template
			funcTemplate string
			callTemplate string
		)
		if strings.Contains(op.Name, "GetNetworkPublicKey") {
			template = GenerateFHEOperationNoGasTemplate()
		} else {
			op.OperationTypeName = "toType"
			if strings.Contains(op.Inputs, "utype") {
				op.OperationTypeName = "utype"
			}

			template = GenerateFHEOperationTemplate()
			if len(op.Inputs) == 44 {
				funcTemplate, callTemplate = GenerateHandlerFunction(op.Name)
			}

		}
		err = template.Execute(file, op)
		if err != nil {
			fmt.Println("Error writing")
			return
		}

		handlers.WriteString(funcTemplate)
		getHandlersFunc += callTemplate
	}

	getHandlersFunc += `}
}`
	handlers.WriteString(getHandlersFunc)
}

type Operation struct {
	Name              string
	OperationTypeName string
	Inputs            string
	InnerInputs       string
	ReturnType        string
}

type Function struct {
	Inputs          []Argument `json:"inputs"`
	Name            string     `json:"name"`
	Outputs         []Argument `json:"outputs"`
	StateMutability string     `json:"stateMutability"`
	Type            string     `json:"type"`
}

// Argument represents a function argument
type Argument struct {
	InternalType string `json:"internalType"`
	Name         string `json:"name"`
	Type         string `json:"type"`
}

func GenerateFHEOperationNoGasTemplate() *template.Template {
	templateText := `
func (con FheOps) {{.Name}}(c ctx, evm mech{{.Inputs}}) ({{.ReturnType}}, error) {

	tp := fheos.TxParamsFromEVM(evm, c.caller)
	return fheos.{{.Name}}({{.InnerInputs}}&tp)
}
`

	tmpl, err := template.New("functionTemplate").Parse(templateText)
	if err != nil {
		fmt.Printf("Error parsing template: %s\n", err)
		return nil
	}

	return tmpl
}

func GenerateLogTemplate() string {
	templateText := `
func (con FheOps) Log(c ctx, evm mech, s string) error {
	tp := fheos.TxParamsFromEVM(evm, c.caller)
	gas, err := fheos.Log(s, &tp)

	if err != nil {
		return err
	}

	return c.Burn(gas)
}
`

	return templateText
}

func GenerateFHEOperationTemplate() *template.Template {
	templateText := `
func (con FheOps) {{.Name}}(c ctx, evm mech{{.Inputs}}) ({{.ReturnType}}, error) {
	tp := fheos.TxParamsFromEVM(evm, c.caller)
	if metrics.Enabled {
		h := fmt.Sprintf("%s/%s/%s", "fheos", "{{.Name}}", fheos.UtypeToString({{.OperationTypeName}}))
		defer func(start time.Time) {
			sampler := func() metrics.Sample {
				return metrics.NewBoundedHistogramSample()
			}
			metrics.GetOrRegisterHistogramLazy(h, nil, sampler).Update(time.Since(start).Microseconds())
		}(time.Now())
	}

	ret, gas, err := fheos.{{.Name}}({{.InnerInputs}}&tp, nil)

	if err != nil {
		if metrics.Enabled {
			c := fmt.Sprintf("%s/%s/%s/%s", "fheos", "{{.Name}}", fheos.UtypeToString({{.OperationTypeName}}), "error/fhe_failure")
			metrics.GetOrRegisterCounter(c, nil).Inc(1)
		}
		return ret, err
	}

	err = c.Burn(gas)

	if metrics.Enabled {
		metricPath := fmt.Sprintf("%s/%s/%s/%s", "fheos", "{{.Name}}", fheos.UtypeToString({{.OperationTypeName}}), "success/total")
		if err != nil {
			metricPath = fmt.Sprintf("%s/%s/%s/%s", "fheos", "{{.Name}}", fheos.UtypeToString({{.OperationTypeName}}), "error/fhe_gas_failure")
		}

		metrics.GetOrRegisterCounter(metricPath, nil).Inc(1)
	}

	return ret, err
}
`

	tmpl, err := template.New("functionTemplate").Parse(templateText)
	if err != nil {
		fmt.Printf("Error parsing template: %s\n", err)
		return nil
	}

	return tmpl
}

func CapitalizeFirstLetter(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func GenerateHandlerFunction(name string) (string, string) {
	capitalized := CapitalizeFirstLetter(name)
	templateText := fmt.Sprintf(`
func %sHandler(w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, precompiles.%s)
}
`, capitalized, capitalized)
	templateText2 := fmt.Sprintf(`{"/%s", %sHandler},`, name, capitalized)

	return templateText, templateText2
}
