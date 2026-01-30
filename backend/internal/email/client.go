package email

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	apiKey string
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
	}
}

func (c *Client) Send(to string, subject string, body string) error {

	type sender struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	type toEmail struct {
		Email string `json:"email"`
	}

	type sendRequest struct {
		Sender  sender    `json:"sender"`
		To      []toEmail `json:"to"`
		Subject string    `json:"subject"`
		Body    string    `json:"htmlContent"`
	}

	payload := sendRequest{
		Sender: sender{
			Name:  "Sentinel Test",
			Email: "kafkaesquesayhy@gmail.com",
		},
		To: []toEmail{
			{
				Email: to,
			},
		},
		Subject: subject,
		Body:    body,
	}
	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("Error marshalling data: %s", err)
	}

	Buffer := bytes.NewBuffer(jsonBytes)
	url := "https://api.brevo.com/v3/smtp/email"
	req, err := http.NewRequest(http.MethodPost, url, Buffer)
	req.Header.Set("content-type", "application/json")
	req.Header.Set("api-key", c.apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("Error sending request: %s", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 201 {
		errResp, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("Error Reading response body: %s", err)
		}
		return fmt.Errorf("http request rejected with StatusCode : %d and Error: %s", resp.StatusCode, string(errResp))
	}

	return nil
}
