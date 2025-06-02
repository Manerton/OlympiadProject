package another_service

import "net/http"

type AnotherService struct {
	notifyService string
}

func New(notifyService string) *AnotherService {
	return &AnotherService{
		notifyService: notifyService,
	}
}

func SendNotifyAcceptAccount(email string) {
	url := "https://172.16.0.94/olymp_notification/web/index.php?r=email%2Fsend-code"

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {

	}
	req.Header.Set("requestToken", "1234567890")
	req.Header.Set("email", email)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {

	}

	defer resp.Body.Close()
}
