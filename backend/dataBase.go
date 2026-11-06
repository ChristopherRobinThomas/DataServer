package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type dbInfo struct {
	user     string
	password string
}

type Entry struct {
	server int
	date   string
}

type Server struct {
	name string
	addr string
	arch string
	dir  string
}

type Run struct {
	name   string
	entry  int
	all    float32
	cpu    float32
	memory float32
	io     float32
	system float32
	mt     float32
	st     float32
	os     float32
	sleep  float32
}

func connect(info dbInfo) *sql.DB {
	connStr := fmt.Sprintf("postgresql://%s:%s@localhost:5432/mydb?sslmode=disable", info.user, info.password)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("open")
	}
	return db
}

func insertEntry(db *sql.DB, entry Entry) {
	rows, err := db.Query(fmt.Sprintf("INSERT INTO entries (server, date) VALUES (%s,%s)", entry.server, entry.date))
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}
	rows.Close()
}

func insertServer(db *sql.DB, server Server) {
	rows, err := db.Query(fmt.Sprintf("INSERT INTO servers (name, addr, dir, arch) VALUES", server.name, server.addr, server.dir, server.arch))
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}
	rows.Close()
}

func insertRun(db *sql.DB, run Run) {
	rows, err := db.Query(fmt.Sprintf("INSERT INTO run (name, entry, all, cpu, memory, io, system, mt, st, os, sleep) VALUES (%s,%d,%v,%v,%v,%v,%v,%v,%v,%v,%v,)",
		run.name,
		run.entry,
		run.all,
		run.memory,
		run.io,
		run.system,
		run.mt,
		run.st,
		run.os,
		run.sleep,
	))
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}
	rows.Close()
}

func getRunsByEntry(db *sql.DB, entryId int) []Run {
	rows, err := db.Query("SELECT entry, name, all, memory, io, system, mt, st, os, sleep FROM runs WHERE entry == %d", entryId)
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}
	defer rows.Close()
	arr := make([]Run, 0)
	for rows.Next() {
		var run Run
		err := rows.Scan(
			&run.entry,
			&run.name,
			&run.all,
			&run.memory,
			&run.io,
			&run.system,
			&run.mt,
			&run.st,
			&run.os,
			&run.sleep,
		)
		if err != nil {
			log.Fatalf("Scan failed: %v", err)
		}
		arr = append(arr, run)
	}
	return arr
}

func getEntriesByServer(db *sql.DB, serverId int) []Entry {
	rows, err := db.Query("SELECT server, date FROM entries WHERE entry == %d", serverId)
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}
	defer rows.Close()
	arr := make([]Entry, 0)
	for rows.Next() {
		var entry Entry
		err := rows.Scan(
			&entry.server,
			&entry.date,
		)
		if err != nil {
			log.Fatalf("Scan failed: %v", err)
		}
		arr = append(arr, entry)
	}
	return arr
}

func getServers(db *sql.DB) []Server {
	rows, err := db.Query("SELECT addr, arch, dir, name FROM servers")
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}
	defer rows.Close()
	arr := make([]Server, 0)
	for rows.Next() {
		var server Server
		err := rows.Scan(
			&server.addr,
			&server.arch,
			&server.dir,
			&server.name,
		)
		if err != nil {
			log.Fatalf("Scan failed: %v", err)
		}
		arr = append(arr, server)
	}
	return arr
}
