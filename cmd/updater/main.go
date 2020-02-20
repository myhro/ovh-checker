package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/myhro/ovh-checker/database"
	"github.com/nleof/goyesql"
)

func fetchAPI() (string, error) {
	url := "https://www.ovh.com/engine/api/dedicated/server/availabilities?country=pt"

	log.Print("Fetching servers availability")
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	log.Print("Done")

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func sleep() {
	time.Sleep(60 * time.Second)
}

func main() {
	db, err := database.New()
	if err != nil {
		log.Fatal(err)
	}

	queries, err := goyesql.ParseFile("sql/offer.sql")
	if err != nil {
		log.Fatal(err)
	}

	for {
		body, err := fetchAPI()
		if err != nil {
			log.Print(err)
			sleep()
			continue
		}

		log.Print("Updating servers offers")
		_, err = db.Exec(queries["import-json"], body)
		if err != nil {
			log.Print(err)
			sleep()
			continue
		}
		log.Print("Done")

		sleep()
	}
}
