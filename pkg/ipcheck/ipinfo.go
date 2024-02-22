package ipcheck

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type IPInfo struct {
	SimpleChecker
}

func NewIPInfoGetter(url string, token string) IPInfo {
	return IPInfo{
		SimpleChecker{
			Type:  IpInfoIO,
			URL:   url,
			Token: token,
		},
	}
}

func (i IPInfo) GetIP() (ip string, err error) {
	return i.GetIPWithContext(context.Background())
}

func (i IPInfo) GetIPWithContext(ctx context.Context) (ip string, err error) {
	url := i.GetURL()
	token := i.GetToken()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}

	if token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	if resp.StatusCode >= 300 {
		return "", errors.New(fmt.Sprintf("httpbin server response status code is not 2xx, is %d", resp.StatusCode))
	}

	reply, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return "", err
	}

	return string(reply), nil
}