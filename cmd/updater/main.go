package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/myhro/ovh-checker/postgres"
	"github.com/nleof/goyesql"
)

func main() {
	url := "https://www.ovh.com/engine/api/dedicated/server/availabilities?country=pt"

	db, err := postgres.New()
	if err != nil {
		log.Fatal(err)
	}

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
