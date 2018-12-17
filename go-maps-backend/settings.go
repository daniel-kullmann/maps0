package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"log"
	"net/http"
	"os"
)

var db *sql.DB

const (
	token = "oh-we-totally-use-csrf"
)

func initDatabase() {
	_, err := os.Stat(DataBasePath)
	if os.IsNotExist(err) {
		log.Print("Need to create database")
		var sqlStmt = `CREATE TABLE "setting" ("name" varchar(255) NOT NULL PRIMARY KEY, "value" varchar(255) NOT NULL)`
		db, err := sql.Open("sqlite3", DataBasePath)
		if err != nil {
			log.Fatal(err)
		}
		_, err = db.Query("select name, value from setting limit 1")
		if err != nil {
			// Create table if it does not exist
			log.Print(err)
			_, err = db.Exec(sqlStmt)
			if err != nil {
				db.Close()
				log.Fatal(err)
			}
		}
		db.Close()
	}
}

func GetDb() (*sql.DB, error) {
	if db == nil {
		var err error
		db, err = sql.Open("sqlite3", DataBasePath)
		if err != nil {
			return nil, err
		}
	}
	return db, nil
}

func GetSettings(w http.ResponseWriter, r *http.Request) {
	db, err := GetDb()
	if err != nil {
		log.Fatal(err)
	}
	rows, err := db.Query("select name, value from setting")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	result := make(map[string]string)
	for rows.Next() {
		var name string
		var value string
		err = rows.Scan(&name, &value)
		if err != nil {
			log.Fatal(err)
		}
		result[name] = value
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	jsonResponse, err := json.Marshal(result)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Write([]byte(jsonResponse))
}

func GetToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte("{\"token\": \"" + token + "\"}"))
}

func SetSettings(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Access-Control-Allow-Origin", "*")

	if len(r.Header["X-Csrftoken"]) == 0 || r.Header["X-Csrftoken"][0] != token {
		w.WriteHeader(403)
		return
	}

	buffer := make([]uint8, 10*1024)
	count, err := r.Body.Read(buffer)
	if err != nil && err != io.EOF {
		log.Fatal(err)
	}

	dict := make(map[string]interface{})
	err = json.Unmarshal(buffer[0:count], &dict)
	if err != nil {
		log.Fatal(err)
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare("REPLACE INTO setting (name, value) VALUES (?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	for name, value := range dict {
		var stringValue string
		switch val := value.(type) {
		case string:
			stringValue = val
		case float64:
			stringValue = fmt.Sprintf("%f", val)
		default:
			log.Fatal("did not implement for this type")
		}
		_, err = stmt.Exec(name, stringValue)
		if err != nil {
			log.Fatal(err)
		}
	}
	tx.Commit()
	w.Write([]byte("{}"))
}
