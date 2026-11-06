package main

import (
	"net/http"
)

func tmp(resp http.ResponseWriter, req *http.Request) {}

func main() {
	info := extractSshInfo()
	executeSshCommand(info)
	downloadViaCat(info, "./data/results.json")
	jsonReader("./data/results.json")

	mux := http.NewServeMux()
	mux.HandleFunc("GET /dates", tmp)
	mux.HandleFunc("GET /servers", tmp)
	mux.HandleFunc("POST /servers", tmp)
}
