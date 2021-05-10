package httprequests

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPerformRequests_GET(t *testing.T) {
	assert := assert.New(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "hello")
	}))
	defer ts.Close()

	headers := make(map[string]string)
	headers["X-Request-Id"] = "TestPerformRequests_GET"
	body, resp, err := PerformRequests(headers, "GET", ts.URL, "", "")

	if err != nil {
		t.Fail()
	} else {
		if resp.StatusCode == 200 {
			assert.Equal(200, resp.StatusCode, "Failed to perform http GET request")
			assert.Contains(string(body), "hello")
		}
	}
}

func TestPerformRequests_POST(t *testing.T) {
	assert := assert.New(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "hello")
	}))
	defer ts.Close()

	headers := make(map[string]string)
	headers["X-Request-Id"] = "TestPerformRequests_POST"
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	body, resp, err := PerformRequests(headers, "POST", ts.URL, "a=b", "")

	if err != nil {
		t.Fail()
	} else {
		if resp.StatusCode == 200 {
			assert.Equal(200, resp.StatusCode, "Failed to perform http POST request")
			assert.Contains(string(body), "hello")
		}
	}
}

func TestPerformRequests_500(t *testing.T) {
	assert := assert.New(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Internal Server Error")
	}))
	defer ts.Close()

	headers := make(map[string]string)
	headers["X-Request-Id"] = "TestPerformRequests_500"
	_, _, err := PerformRequests(headers, "GET", ts.URL, "", "")
	assert.Error(err)
}

func TestPerformRequests_500_nothing(t *testing.T) {
	assert := assert.New(t)
	headers := make(map[string]string)
	headers["X-Request-Id"] = "TestPerformRequests_500"
	_, _, err := PerformRequests(headers, "POST", "", "", "")
	assert.Error(err)
}

func TestPerformRequests_503(t *testing.T) {
	assert := assert.New(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintln(w, "Internal Server Error")
	}))
	defer ts.Close()

	headers := make(map[string]string)
	_, _, err := PerformRequests(headers, "GET", ts.URL, "", "")
	assert.Error(err)
}

func TestPerformRequests_500_set_retry_max(t *testing.T) {
	assert := assert.New(t)
	os.Setenv("HTTP_RETRY_MAX", "1")
	defer os.Unsetenv("HTTP_RETRY_MAX")

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Internal Server Error")
	}))
	defer ts.Close()

	headers := make(map[string]string)
	headers["X-Request-Id"] = "TestPerformRequests_500"
	_, _, err := PerformRequests(headers, "GET", ts.URL, "", "")
	assert.Error(err)
}

func TestPerformRequests_500_set_retry_wait_min(t *testing.T) {
	assert := assert.New(t)
	os.Setenv("HTTP_RETRY_WAIT_MIN", "5")
	defer os.Unsetenv("HTTP_RETRY_WAIT_MIN")

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Internal Server Error")
	}))
	defer ts.Close()

	headers := make(map[string]string)
	headers["X-Request-Id"] = "TestPerformRequests_500"
	_, _, err := PerformRequests(headers, "GET", ts.URL, "", "")
	assert.Error(err)
}

func TestPerformRequests_500_set_retry_wait_max(t *testing.T) {
	assert := assert.New(t)
	os.Setenv("HTTP_RETRY_WAIT_MAX", "5")
	defer os.Unsetenv("HTTP_RETRY_WAIT_MAX")

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Internal Server Error")
	}))
	defer ts.Close()

	headers := make(map[string]string)
	headers["X-Request-Id"] = "TestPerformRequests_500"
	_, _, err := PerformRequests(headers, "GET", ts.URL, "", "")
	assert.Error(err)
}

func TestPerformRequests_500_put(t *testing.T) {
	assert := assert.New(t)
	os.Setenv("HTTP_RETRY_MAX", "1")
	defer os.Unsetenv("HTTP_RETRY_MAX")

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Internal Server Error")
	}))
	defer ts.Close()

	headers := make(map[string]string)
	headers["X-Request-Id"] = "TestPerformRequests_500"
	_, _, err := PerformRequests(headers, "PUT", ts.URL, "", "")
	assert.Error(err)
}

func TestPerformRequests_500_put_payload(t *testing.T) {
	assert := assert.New(t)
	os.Setenv("HTTP_RETRY_MAX", "1")
	defer os.Unsetenv("HTTP_RETRY_MAX")

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Internal Server Error")
	}))
	defer ts.Close()

	headers := make(map[string]string)
	headers["X-Request-Id"] = "TestPerformRequests_500"
	_, _, err := PerformRequests(headers, "PUT", ts.URL, "a=b", "")
	assert.Error(err)
}
