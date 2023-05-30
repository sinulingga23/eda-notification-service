package whatsapp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	fromPhoneNumberId   string
	accessTokenWhatsapp string
)

func init() {
	fromPhoneNumberId = os.Getenv("FROM_PHONE_NUMBER_ID")
	accessTokenWhatsapp = os.Getenv("ACCESS_TOKEN_WHATSAPP")

}

type (
	Receiver struct {
		ToPhoneNumber string
		Message       string
	}

	sendMessageRequest struct {
		MessagingProduct string `json:"messaging_product"`
		RecipientType    string `json:"individual"`
		To               string `json:"to"`
		Type             string `json:"type"`
		Text             struct {
			PreviewUrl bool   `json:"preview_url"`
			Body       string `json:"body"`
		}
	}
)

func SendMessage(receiver Receiver, isContainUrl ...bool) error {
	containUrl := false
	if len(isContainUrl) > 0 {
		containUrl = isContainUrl[0]
	}

	payload := sendMessageRequest{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               receiver.ToPhoneNumber,
		Type:             "text",
		Text: struct {
			PreviewUrl bool   `json:"preview_url"`
			Body       string `json:"body"`
		}{
			PreviewUrl: containUrl,
			Body:       receiver.Message,
		},
	}

	bytesPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	request, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf(`https://graph.facebook.com/v17.0/%s/messages`, fromPhoneNumberId),
		bytes.NewReader(bytesPayload),
	)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessTokenWhatsapp))

	client := http.Client{Timeout: 10 * time.Second}

	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	log.Println(response)

	return nil
}
