package main

import (
	"github.com/fhenixprotocol/fheos/precompiles"
	"net/http"
)

type HandlerDef struct {
	Name    string
	Handler func(w http.ResponseWriter, r *http.Request)
}

func AddHandler(w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, precompiles.Add)
}

func AndHandler(w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, precompiles.And)
}

func DivHandler(w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, precompiles.Div)
}

func EqHandler(w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, precompiles.Eq)
}

func GtHandler(w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, precompiles.Gt)
}

func GteHandler(w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, precompiles.Gte)
}

func LtHandler(w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, precompiles.Lt)
}

func LteHandler(w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, precompiles.Lte)
}

func MaxHandler(w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, precompiles.Max)
}

func MinHandler(w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, precompiles.Min)
}

func MulHandler(w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, precompiles.Mul)
}

func NeHandler(w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, precompiles.Ne)
}

func OrHandler(w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, precompiles.Or)
}

func RemHandler(w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, precompiles.Rem)
}

func RolHandler(w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, precompiles.Rol)
}

func RorHandler(w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, precompiles.Ror)
}

func ShlHandler(w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, precompiles.Shl)
}

func ShrHandler(w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, precompiles.Shr)
}

func SubHandler(w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, precompiles.Sub)
}

func XorHandler(w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, precompiles.Xor)
}
func getHandlers() []HandlerDef {
	return []HandlerDef{
		{"/Add", AddHandler}, {"/And", AndHandler}, {"/Div", DivHandler}, {"/Eq", EqHandler}, {"/Gt", GtHandler}, {"/Gte", GteHandler}, {"/Lt", LtHandler}, {"/Lte", LteHandler}, {"/Max", MaxHandler}, {"/Min", MinHandler}, {"/Mul", MulHandler}, {"/Ne", NeHandler}, {"/Or", OrHandler}, {"/Rem", RemHandler}, {"/Rol", RolHandler}, {"/Ror", RorHandler}, {"/Shl", ShlHandler}, {"/Shr", ShrHandler}, {"/Sub", SubHandler}, {"/Xor", XorHandler}}
}