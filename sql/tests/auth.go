package sqlsuite

import (
	"log"

	"github.com/icrowley/fake"
)

func addRandomUser() int {
	return addUser(fake.EmailAddress())
}

func addUser(email string) int {
	db := newDB()
	queries := newQueries("auth")

	_, err := db.Exec(queries["add-user"], email, fake.SimplePassword())
	if err != nil {
		log.Fatal("error when adding user: ", err)
	}

	var id int
	err = db.Get(&id, queries["user-exists"], email)
	if err != nil {
		log.Fatal("error when fetching user id: ", err)
	}

	return id
}
