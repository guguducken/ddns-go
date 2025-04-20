package requestutils

import (
	"time"

	"github.com/valyala/fasthttp"
)

type RequestOption func(*RequestConfig)

func WithRetryTimes(times int) RequestOption {
	return func(config *RequestConfig) {
		config.retryTimes = times
	}
}

func WithRetryDelay(delay time.Duration) RequestOption {
	return func(config *RequestConfig) {
		config.retryDelay = delay
	}
}

func WithPreRetryConfig(preConfig func(req *fasthttp.Request)) RequestOption {
	return func(config *RequestConfig) {
		config.preRetryConfig = preConfig
	}
}

func WithTimeOut(timeout time.Duration) RequestOption {
	return func(config *RequestConfig) {
		config.timeOut = timeout
	}
}

func WithMaxRedirect(maxRedirect int) RequestOption {
	return func(config *RequestConfig) {
		config.maxRedirect = maxRedirect
	}
}

func WithDefaultValue() RequestOption {
	return func(config *RequestConfig) {
		if config.retryTimes == 0 {
			config.retryTimes = 3
		}
		if config.retryDelay == 0 {
			config.retryDelay = 1 * time.Second
		}
		if config.timeOut == 0 {
			config.timeOut = 10 * time.Second
		}
		if config.maxRedirect == 0 {
			config.maxRedirect = 10
		}
	}
}
