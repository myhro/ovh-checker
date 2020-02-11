package sqlsuite

import (
	"log"
)

func LoadOffers(file string) {
	db := NewDB()
	queries := NewQueries("offer")
	_, err := db.Exec(queries["import-json"], ReadFile(file))
	if err != nil {
		log.Fatal("LoadOffers: ", err)
	}
}
