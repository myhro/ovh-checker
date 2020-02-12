package main

import (
	"log"
	"time"

	"github.com/myhro/ovh-checker/notification"
	"github.com/myhro/ovh-checker/postgres"
	"github.com/nleof/goyesql"
)

func sleep() {
	time.Sleep(60 * time.Second)
}

func main() {
	db, err := postgres.New()
	if err != nil {
		log.Fatal(err)
	}

	queries, err := goyesql.ParseFile("sql/notification.sql")
	if err != nil {
		log.Fatal(err)
	}

	for {
		log.Print("Loading pending notifications")
		res := []notification.PendingNotification{}
		err = db.Select(&res, queries["pending-notifications"])
		if err != nil {
			log.Print(err)
			sleep()
			continue
		}
		log.Print("Done")

		for _, n := range res {
			log.Print("Sending email to user ", n.ID)
			err := sendEmail(n)
			if err != nil {
				log.Print(err)
				continue
			}

			_, err = db.Exec(queries["mark-as-sent"], time.Now(), n.ID)
			if err != nil {
				log.Print(err)
				continue
			}
			log.Print("Marked as sent to ", n.ID)
		}

		sleep()
	}
}
