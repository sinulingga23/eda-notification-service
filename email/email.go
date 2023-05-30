package email

import (
	"errors"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"sync"
)

var (
	senderEmail    string
	senderPassword string
	smtpHost       string
	smptPort       string
)

func init() {
	senderEmail = os.Getenv("SENDER_EMAIL")
	senderPassword = os.Getenv("SENDER_PASSWORD")
	smtpHost = os.Getenv("SMTP_SERVER_HOST")
	smptPort = os.Getenv("SMTP_SERVER_PORT")
}

type Receiver struct {
	Email   string
	Subject string
	Message string
}

func SendTo(receiver Receiver) error {
	auth := smtp.PlainAuth("", senderEmail, senderPassword, smtpHost)

	message := []byte(receiver.Message)

	err := smtp.SendMail(
		fmt.Sprintf("%s:%s", smtpHost, smptPort),
		auth,
		senderEmail,
		[]string{receiver.Email},
		message)
	if err != nil {
		log.Printf("[SendTo]: Error: %v", err)
		return err
	}

	return nil
}

func SendBatch(receivers []Receiver) error {
	if len(receivers) == 0 {
		return errors.New("Data empty.")
	}

	wg := new(sync.WaitGroup)

	for currentWorker := 0; currentWorker < len(receivers); currentWorker++ {
		wg.Add(1)
		go func(currentWorker int, wg *sync.WaitGroup) {
			defer wg.Done()
			receiver := receivers[currentWorker]
			if err := SendTo(receiver); err != nil {
				log.Printf("CurrentWoker: %v, %v", currentWorker, err)
			}
		}(currentWorker, wg)
	}

	wg.Wait()
	return nil
}
