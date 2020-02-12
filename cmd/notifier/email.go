package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/myhro/ovh-checker/notification"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func name(email string) string {
	replacer := strings.NewReplacer(".", " ", "_", " ")
	email = replacer.Replace(email)

	email = strings.Split(email, "@")[0]
	email = strings.Split(email, "+")[0]

	return strings.Title(email)
}

func sendEmail(notif notification.PendingNotification) error {
	from := mail.NewEmail("OVH Checker", os.Getenv("OVH_CHECKER_EMAIL"))
	to := mail.NewEmail(name(notif.Email), notif.Email)
	subject := "OVH server is available"
	content := fmt.Sprintf(
		"OVH server %v (%v %vc/%vt %vGB %v) is available in %v.",
		notif.Server,
		notif.Processor,
		notif.Cores,
		notif.Threads,
		notif.Memory,
		notif.Storage,
		notif.Country,
	)
	message := mail.NewSingleEmail(from, subject, to, content, content)

	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	resp, err := client.Send(message)
	if err != nil {
		return err
	} else if resp.StatusCode == http.StatusAccepted {
		return nil
	}

	return fmt.Errorf("%v %v", resp.StatusCode, resp.Body)
}
