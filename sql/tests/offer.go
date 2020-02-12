package sqlsuite

import (
	"log"
)

func loadOffers(file string) {
	db := newDB()
	queries := newQueries("offer")
	_, err := db.Exec(queries["import-json"], readFile(file))
	if err != nil {
		log.Fatal("LoadOffers: ", err)
	}
}
