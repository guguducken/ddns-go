package ipgetter

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/rs/zerolog/log"
)

type Httpbin struct {
	SimpleChecker
}

type HttpBinGetResponse struct {
	Args    map[string]string `json:"args"`
	Headers map[string]string `json:"headers"`
	Origin  string            `json:"origin"`
	URL     string            `json:"url"`
}

func NewHttpbinGetter(url string, token string) *Httpbin {
	return &Httpbin{
		SimpleChecker{
			Type:  HttpbinGetter,
			URL:   url,
			Token: token,
		},
	}
}

func (h *Httpbin) GetIP() (ip string, err error) {
	return h.GetIPWithContext(context.Background())
}

func (h *Httpbin) GetIPWithContext(ctx context.Context) (ip string, err error) {
	log.Debug().Str("ip_getter", HttpbinGetter).Msg("start get ip")
	url := h.GetURL()
	token := h.GetToken()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Accept", "application/json")

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

	httpbinResp := HttpBinGetResponse{}
	err = json.Unmarshal(reply, &httpbinResp)
	if err != nil {
		return "", err
	}

	if net.ParseIP(httpbinResp.Origin) == nil {
		return "", errors.Join(ErrInvalidResponseIP, errors.New(fmt.Sprintf("invalid ip is %s", httpbinResp.Origin)))
	}

	return httpbinResp.Origin, nil
}
