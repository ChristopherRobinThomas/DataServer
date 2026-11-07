package main

import (
	"database/sql"
	"log"
	"net/http"
)

var db *sql.DB

func main() {
	// info := extractSshInfo()
	// executeSshCommand(info)
	// downloadViaCat(info, "./data/results.json")
	// jsonReader("./data/results.json")
	db = connect(dbInfo{"user", "1234"})

	mux := http.NewServeMux()
	mux.HandleFunc("GET /servers", listServers)
	mux.HandleFunc("GET /entries/{name}", listEntries)

	log.Println("Server listening on :8080...")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
