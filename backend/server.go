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
	mux.HandleFunc("GET /runs/{date}/{server}", listRuns)

	log.Println("Server listening on :9000...")
	log.Fatal(http.ListenAndServe(":9000", mux))
}
