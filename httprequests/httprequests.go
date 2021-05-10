package httprequests

import (
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Lord-Y/cypress-parallel-cli/logger"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/rs/zerolog/log"
)

func init() {
	logger.SetLoggerLogLevel()
}

// retryQuery permit to perform retries when HTTP requests fails
func retryQuery() (z *http.Client) {
	retryClient := retryablehttp.NewClient()

	if strings.TrimSpace(os.Getenv("HTTP_RETRY_MAX")) != "" {
		max, err := strconv.Atoi(os.Getenv("HTTP_RETRY_MAX"))
		if err != nil {
			log.Warn().Err(err).Msgf("Error occured while converting string to integer")
		}
		if max > 0 {
			retryClient.RetryMax = max
		}
	}

	if strings.TrimSpace(os.Getenv("HTTP_RETRY_WAIT_MIN")) != "" {
		waitMin, err := strconv.Atoi(os.Getenv("HTTP_RETRY_WAIT_MIN"))
		if err != nil {
			log.Warn().Err(err).Msgf("Error occured while converting string to integer")
		}
		if waitMin > 0 {
			retryClient.RetryWaitMin = time.Duration(waitMin) * time.Second
		}
	}

	if strings.TrimSpace(os.Getenv("HTTP_RETRY_WAIT_MAX")) != "" {
		waitMax, err := strconv.Atoi(os.Getenv("HTTP_RETRY_WAIT_MAX"))
		if err != nil {
			log.Warn().Err(err).Msgf("Error occured while converting string to integer")
		}
		if waitMax > 0 {
			retryClient.RetryWaitMax = time.Duration(waitMax) * time.Second
		}
	}

	retryClient.Logger = nil
	retryClient.RequestLogHook = func(_ retryablehttp.Logger, req *http.Request, attempt int) {
		retryClient.ResponseLogHook = func(_ retryablehttp.Logger, resp *http.Response) {
			wait := retryClient.Backoff(retryClient.RetryWaitMin, retryClient.RetryWaitMax, attempt, resp)
			if resp.StatusCode == 0 || (resp.StatusCode >= 500 && resp.StatusCode != 501) {
				if req.Header.Get("X-Request-Id") != "" {
					log.Warn().Str("requestId", req.Header.Get("X-Request-Id")).Msgf("Error occured while performing http request method %s %s with statusCode %d for attempt number %d, retrying in %s", req.Method, req.URL, resp.StatusCode, attempt, wait)
				} else {
					log.Warn().Msgf("Error occured while performing http request method %s %s with statusCode %d for attempt number %d, retrying in %s", req.Method, req.URL, resp.StatusCode, attempt, wait)
				}
			}
		}
	}
	return retryClient.StandardClient()
}

// PerformRequests permit to perform HTTP requests
func PerformRequests(headers map[string]string, method string, url string, payload string, retryProvider string) (body []byte, resp *http.Response, err error) {
	var (
		req *http.Request
	)

	switch retryProvider {
	default:
		standardClient := retryQuery()
		switch method {
		case "GET":
			req, err = http.NewRequest(method, url, nil)
		case "POST":
			req, err = http.NewRequest(method, url, strings.NewReader(payload))
		default:
			if payload == "" {
				req, err = http.NewRequest(method, url, nil)
			} else {
				req, err = http.NewRequest(method, url, strings.NewReader(payload))
			}
		}
		if err != nil {
			log.Error().Err(err).Msgf("Error occured while initializing http request")
			return nil, nil, err
		}
		if len(headers) > 0 {
			for k, v := range headers {
				req.Header.Set(k, v)
			}
		}
		if payload != "" {
			req.Header.Set("Content-Length", strconv.Itoa(len(payload)))
			log.Debug().Msgf("%s %s headers: %s body %v", req.Method, req.URL.Redacted(), req.Header, payload)
		} else {
			log.Debug().Msgf("%s %s headers: %s body %v", req.Method, req.URL.Redacted(), req.Header, nil)
		}
		resp, err = standardClient.Do(req)
		if err != nil {
			return nil, nil, err
		}
		defer resp.Body.Close()
		body, err = ioutil.ReadAll(resp.Body)
		log.Debug().Msgf("%s %s statusCode %d headers: %s body %v", req.Method, req.URL.Redacted(), resp.StatusCode, resp.Header, string(body))
	}
	return body, resp, err
}
