package main

import (
	"database/sql"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"github.com/nleof/goyesql"
)

var db *sql.DB

func init() {
	conn, ok := os.LookupEnv("POSTGRES_CONN")
	if !ok {
		conn = "dbname=ovh sslmode=disable"
	}

	var err error
	db, err = sql.Open("postgres", conn)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	url := "https://www.ovh.com/engine/api/dedicated/server/availabilities?country=pt"

	log.Print("Fetching servers availability")
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Done")

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	queries, err := goyesql.ParseFile("sql/json.sql")
	if err != nil {
		log.Fatal(err)
	}

	log.Print("Updating servers offers")
	_, err = db.Exec(queries["import-json"], string(body))
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Done")
}
