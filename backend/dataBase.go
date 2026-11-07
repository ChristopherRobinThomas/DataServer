package main

import (
	"database/sql"
	// "encoding/json"
	"errors"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type dbInfo struct {
	User     string
	Password string
}

type Entry struct {
	Server int    `json:"server"`
	Date   string `json:"date"`
}

type Server struct {
	Name string `json:"name"`
	Addr string `json:"addr"`
	Arch string `json:"arch"`
	Dir  string `json:"dir"`
}

type Run struct {
	Name    string  `json:"name"`
	Entry   int     `json:"entry"`
	Overall float32 `json:"overall"`
	Cpu     float32 `json:"cpu"`
	Memory  float32 `json:"memory"`
	Io      float32 `json:"io"`
	System  float32 `json:"system"`
	Mt      float32 `json:"mt"`
	St      float32 `json:"st"`
	Os      float32 `json:"os"`
	Sleep   float32 `json:"sleep"`
}

func connect(info dbInfo) *sql.DB {
	connStr := fmt.Sprintf("postgresql://%s:%s@localhost:5432/mydb?sslmode=disable", info.User, info.Password)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("open")
	}
	return db
}

func insertEntry(entry Entry) {
	rows, err := db.Query(fmt.Sprintf("INSERT INTO entries (server, date) VALUES (%d,'%s');", entry.Server, entry.Date))
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}
	rows.Close()
}

func insertServer(server Server) {
	rows, err := db.Query(fmt.Sprintf("INSERT INTO servers (name, addr, dir, arch) VALUES('%s', '%s', '%s', '%s');", server.Name, server.Addr, server.Dir, server.Arch))
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}
	rows.Close()
}

func insertRun(run Run) {
	rows, err := db.Query(fmt.Sprintf("INSERT INTO runs (name, entry, overall, cpu, memory, io, system, mt, st, os, sleep) VALUES ('%s',%d,%v,%v,%v,%v,%v,%v,%v,%v,%v);",
		run.Name,
		run.Entry,
		run.Overall,
		run.Cpu,
		run.Memory,
		run.Io,
		run.System,
		run.Mt,
		run.St,
		run.Os,
		run.Sleep,
	))
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}
	rows.Close()
}

func getRunsByEntry(entryId int) []Run {
	rows, err := db.Query(fmt.Sprintf("SELECT entry, name, overall, memory, io, system, mt, st, os, sleep FROM runs WHERE entry = %d;", entryId))
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}
	defer rows.Close()
	arr := make([]Run, 0)
	for rows.Next() {
		var run Run
		err := rows.Scan(
			&run.Entry,
			&run.Name,
			&run.Overall,
			&run.Memory,
			&run.Io,
			&run.System,
			&run.Mt,
			&run.St,
			&run.Os,
			&run.Sleep,
		)
		if err != nil {
			log.Fatalf("Scan failed: %v", err)
		}
		arr = append(arr, run)
	}
	return arr
}

func getEntriesByServer(serverId int) []Entry {
	rows, err := db.Query(fmt.Sprintf("SELECT server, date FROM entries WHERE server = %d;", serverId))
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}
	defer rows.Close()
	arr := make([]Entry, 0)
	for rows.Next() {
		var entry Entry
		err := rows.Scan(
			&entry.Server,
			&entry.Date,
		)
		if err != nil {
			log.Fatalf("Scan failed: %v", err)
		}
		arr = append(arr, entry)
	}
	return arr
}

func getServers() []Server {
	rows, err := db.Query("SELECT addr, arch, dir, name FROM servers;")
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}
	defer rows.Close()
	arr := make([]Server, 0)
	for rows.Next() {
		var server Server
		err := rows.Scan(
			&server.Addr,
			&server.Arch,
			&server.Dir,
			&server.Name,
		)
		if err != nil {
			log.Fatalf("Scan failed: %v", err)
		}
		arr = append(arr, server)
	}
	return arr
}

func getEntryIdByDate(date string, server int) (int, error) {
	rows, err := db.Query(fmt.Sprintf("SELECT id, date, server FROM entries WHERE date = '%s' AND server = %d;", date, server))
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}
	defer rows.Close()
	var id int = -1
	var d string
	var s int
	for rows.Next() {
		err := rows.Scan(&id, &d, &s)
		if err != nil {
			log.Fatalf("Scan failed: %v", err)
			return id, errors.New("Can not find entry")
		}
		return id, nil
	}
	return id, errors.New("Can not find entry")
}

func getServerByName(name string) (int, error) {
	rows, err := db.Query(fmt.Sprintf("SELECT id, name FROM servers WHERE name = '%s';", name))
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}
	defer rows.Close()
	var id int = -1
	var n string
	for rows.Next() {
		err := rows.Scan(&id, &n)
		if err != nil {
			log.Fatalf("Scan failed: %v", err)
			return id, errors.New("Can not find server")
		}
		return id, nil
	}
	return id, errors.New("Can not find server")
}
