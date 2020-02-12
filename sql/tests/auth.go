package sqlsuite

import (
	"log"

	"github.com/icrowley/fake"
)

func addRandomUser() string {
	email := fake.EmailAddress()
	addUser(email)
	return email
}

func addUser(email string) {
	db := newDB()
	queries := newQueries("auth")
	_, err := db.Exec(queries["add-user"], email, fake.SimplePassword())
	if err != nil {
		log.Fatal("AddUser: ", err)
	}
}
