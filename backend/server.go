package main

import (
	"net/http"
)

func tmp(resp http.ResponseWriter, req *http.Request) {}

func main() {
	//info := extractSshInfo()
	//executeSshCommand(info)
	//downloadViaCat(info, "./data/results.json")
	db := connect(dbInfo{"user", "1234"})
	jsonReader("./data/results.json", db)

	//mux := http.NewServeMux()
	//mux.HandleFunc("GET /dates", tmp)
	//mux.HandleFunc("GET /servers", tmp)
	//mux.HandleFunc("POST /servers", tmp)
}
