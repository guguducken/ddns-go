package requestutils

import (
	"net/http"
	"sync"
	"time"

	"github.com/valyala/fasthttp"

	"github.com/guguducken/ddns-go/pkg/errno"
	"github.com/guguducken/ddns-go/pkg/utils/logutil"
)

type RequestConfig struct {
	retryTimes  int
	retryDelay  time.Duration
	maxRedirect int

	timeOut        time.Duration
	preRetryConfig func(req *fasthttp.Request)
}

func Get(url string, headers http.Header, options ...RequestOption) (*fasthttp.Response, error) {
	return Request("GET", url, headers, nil, options...)
}

func Delete(url string, headers http.Header, options ...RequestOption) (*fasthttp.Response, error) {
	return Request("DELETE", url, headers, nil, options...)
}

func Head(url string, headers http.Header, options ...RequestOption) (*fasthttp.Response, error) {
	return Request("HEAD", url, headers, nil, options...)
}

func Post(url string, headers http.Header, body []byte, options ...RequestOption) (*fasthttp.Response, error) {
	return Request("POST", url, headers, body, options...)
}

func Put(url string, headers http.Header, body []byte, options ...RequestOption) (*fasthttp.Response, error) {
	return Request("PUT", url, headers, body, options...)
}

func Patch(url string, headers http.Header, body []byte, options ...RequestOption) (*fasthttp.Response, error) {
	return Request("PATCH", url, headers, body, options...)
}

func InitFastHttpRequest(method string, url string, body []byte) *fasthttp.Request {
	req := fasthttp.AcquireRequest()
	req.Header.SetMethod(method)
	req.SetRequestURI(url)
	if method != "GET" && method != "HEAD" && body != nil {
		req.SetBody(body)
	}
	return req
}

func ReleaseResponse(resp *fasthttp.Response) {
	fasthttp.ReleaseResponse(resp)
}

var (
	defaultHttpClient *fasthttp.Client
	clientInitOnce    sync.Once
)

func InjectHttpClient(client *fasthttp.Client) {
	if client == nil {
		panic("inject client is nil")
	}
	if defaultHttpClient != nil {
		panic("http client is already initialized, must be injected into the client when not called")
	}
	defaultHttpClient = client
}

func defaultFastHttpClient() *fasthttp.Client {
	clientInitOnce.Do(func() {
		if defaultHttpClient == nil {
			defaultHttpClient = &fasthttp.Client{
				MaxConnsPerHost:     2000,
				StreamResponseBody:  true,
				MaxIdleConnDuration: 5 * time.Second,
				MaxConnWaitTimeout:  5 * time.Second,
				DialTimeout:         fasthttp.DialDualStackTimeout,
			}
		}
	})

	return defaultHttpClient
}

func Request(method string, url string, headers http.Header, body []byte, options ...RequestOption) (*fasthttp.Response, error) {
	o := deployRequestOptions(options)

	req := InitFastHttpRequest(method, url, body)
	req.SetTimeout(o.timeOut)
	for k := range headers {
		req.Header.Add(k, headers.Get(k))
	}
	defer fasthttp.ReleaseRequest(req)

	for t := range o.retryTimes {
		resp, err := do(req, o.maxRedirect)
		if err != nil {
			logutil.Error(err, "do request failed, will retry", logutil.Int("try_time", t))
			if o.preRetryConfig != nil {
				o.preRetryConfig(req)
			}
			time.Sleep(o.retryDelay)
			continue
		}
		return resp, nil
	}
	return nil, errno.ErrRequestMaxTimes
}

// do http request
// if request failed, will return ErrHttpRequest
func do(req *fasthttp.Request, maxRedirectTimes int) (*fasthttp.Response, error) {
	// init fast http response
	resp := fasthttp.AcquireResponse()
	if err := defaultFastHttpClient().DoRedirects(req, resp, maxRedirectTimes); err != nil {
		fasthttp.ReleaseResponse(resp)
		return nil, errno.OverrideError(errno.ErrRequestFailed, errno.OverrideMessage(err.Error()))
	}
	return resp, nil
}

func deployRequestOptions(options []RequestOption) *RequestConfig {
	cfg := &RequestConfig{}
	for _, opt := range options {
		opt(cfg)
	}
	WithDefaultValue()(cfg)
	return cfg
}
