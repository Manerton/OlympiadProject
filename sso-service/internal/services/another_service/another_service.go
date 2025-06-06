package another_service

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type AnotherService struct {
	notifyService string
}

func New(notifyService string) *AnotherService {
	return &AnotherService{
		notifyService: notifyService,
	}
}

func SendNotifyAcceptAccount(email string) error {
	url := "https://172.16.0.94/olymp_notification/web/index.php?r=email%2Fsend-code"

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("failed create request")
	}
	// Header set
	req.Header.Set("requestToken", "1234567890")
	req.Header.Set("email", email)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed get respose %w", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read body %w", err)
	}

	type NotifyResponse struct {
		Code string `json:"code"`
	}

	notifyRes := NotifyResponse{}

	err = json.Unmarshal(body, &notifyRes)
	if err != nil {
		return fmt.Errorf("failed to parse body %w", err)
	}

	return nil
}
