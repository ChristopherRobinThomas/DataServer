package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func listServers(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	servers := getServers()

	err := json.NewEncoder(w).Encode(servers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func listEntries(w http.ResponseWriter, req *http.Request) {
	name := req.PathValue("name")
	id, err := getServerByName(name)
	if err != nil {
		log.Fatalf("%v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	entries := getEntriesByServer(id)
	err = json.NewEncoder(w).Encode(entries)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
