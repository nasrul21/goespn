package goespn

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"time"
)

// Client : client interface
type Client interface {
	NewRequest(method string, fullPath string, body io.Reader) (*http.Request, error)
	ExecuteRequest(req *http.Request, v interface{}) error
	Call(method, path string, body io.Reader, v interface{}) error
}

// Client struct
type client struct {
	BaseURL  string
	LogLevel int
	Logger   *log.Logger
}

// NewClient : this function will always be called when the library is in use
func NewClient() Client {
	return &client{
		BaseURL:  "http://site.api.espn.com/apis",
		LogLevel: 2,
		Logger:   log.New(os.Stderr, "", log.LstdFlags),
	}
}

// NewRequest : send new request
func (c *client) NewRequest(method string, fullPath string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, fullPath, nil)
	if err != nil {
		return nil, err
	}
	return req, nil
}

var defHTTPTimeout = 15 * time.Second
var httpClient = &http.Client{Timeout: defHTTPTimeout}

// ExecuteRequest : execute request
func (c *client) ExecuteRequest(req *http.Request, v interface{}) error {
	logLevel := c.LogLevel
	logger := c.Logger

	if logLevel > 1 {
		logger.Println("Request ", req.Method, ": ", req.URL.Host, req.URL.Path)
	}

	start := time.Now()

	res, err := httpClient.Do(req)
	if err != nil {
		if logLevel > 0 {
			logger.Println("Cannot send request: ", err)
		}
		return err
	}
	defer res.Body.Close()

	if logLevel > 2 {
		logger.Println("Completed in ", time.Since(start))
	}

	if err != nil {
		if logLevel > 0 {
			logger.Println("Request failed: ", err)
		}
		return err
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		if logLevel > 0 {
			logger.Println("Cannot read response body: ", err)
		}
		return err
	}

	if logLevel > 2 {
		logger.Println("Midtrans response: ", string(resBody))
	}

	if v != nil {
		if err = json.Unmarshal(resBody, v); err != nil {
			return err
		}

		// when return unexpected error, midtrans not return `status_message` but `message`, so this to catch it
		error := make(map[string]string)
		if res.StatusCode >= 500 {
			err := json.Unmarshal(resBody, &error)
			if err != nil {
				return err
			}
		}

		// we're safe to reflect status_code if response not return status code
		if reflect.ValueOf(v).Elem().Kind() == reflect.Struct {
			if reflect.ValueOf(v).Elem().FieldByName("StatusCode").Len() == 0 {
				reflect.ValueOf(v).Elem().FieldByName("StatusCode").SetString(strconv.Itoa(res.StatusCode))
				// response of snap transaction not return StatusMessage
				if req.URL.Path != "/snap/v1/transactions" {
					reflect.ValueOf(v).Elem().FieldByName("StatusMessage").SetString(error["message"])
				}
			}
		}
	}

	return nil
}

// Call the ESPN API at specific `path` using the specified HTTP `method`. The result will be
// given to `v` if there is no error. If any error occurred, the return of this function is the error
// itself, otherwise nil.
func (c *client) Call(method, path string, body io.Reader, v interface{}) error {

	path = c.BaseURL + path

	req, err := c.NewRequest(method, path, body)
	if err != nil {
		return err
	}

	return c.ExecuteRequest(req, v)
}
