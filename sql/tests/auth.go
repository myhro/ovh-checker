package sqlsuite

import (
	"log"

	"github.com/icrowley/fake"
)

func AddRandomUser() string {
	email := fake.EmailAddress()
	AddUser(email)
	return email
}

func AddUser(email string) {
	db := NewDB()
	queries := NewQueries("auth")
	_, err := db.Exec(queries["add-user"], email, fake.SimplePassword())
	if err != nil {
		log.Fatal("AddUser: ", err)
	}
}
