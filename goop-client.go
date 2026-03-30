package goop

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/html/charset"
)

// Init a new HTTP client for use when the client doesn't want to use their own.
var (
	defaultClient = &http.Client{}

	// Headers contains all HTTP headers to send
	Headers = make(map[string]string)

	// Cookies contains all HTTP cookies to send
	Cookies = make(map[string]string)
)

// Header sets a new HTTP header
func Header(n string, v string) {
	Headers[n] = v
}

// Cookie sets a cookie for http requests
func Cookie(n string, v string) {
	Cookies[n] = v
}

// GetWithClient returns the HTML returned by the url using a provided HTTP client
func GetWithClient(url string, client *http.Client) (string, error) {
	timer := startTimer("HTTP GET: "+url, DebugVerbose)
	defer timer.finish()

	logHTTPRequest("GET", url, Headers)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		if debug {
			panic("Couldn't create GET request to " + url)
		}
		return "", newError(ErrCreatingGetRequest, "error creating get request to "+url)
	}

	setHeadersAndCookies(req)

	// Perform request
	resp, err := client.Do(req)
	if err != nil {
		if debug {
			panic("Couldn't perform GET request to " + url)
		}
		return "", newError(ErrInGetRequest, "couldn't perform GET request to "+url)
	}
	defer resp.Body.Close()

	logHTTPResponse(resp.StatusCode, int(resp.ContentLength))

	utf8Body, err := charset.NewReader(resp.Body, resp.Header.Get("Content-Type"))
	if err != nil {
		return "", err
	}
	bytes, err := ioutil.ReadAll(utf8Body)
	if err != nil {
		if debug {
			panic("Unable to read the response body")
		}
		return "", newError(ErrReadingResponse, "unable to read the response body")
	}
	return string(bytes), nil
}

// setHeadersAndCookies helps build a request
func setHeadersAndCookies(req *http.Request) {
	// Set headers
	for hName, hValue := range Headers {
		req.Header.Set(hName, hValue)
	}
	// Set cookies
	for cName, cValue := range Cookies {
		req.AddCookie(&http.Cookie{
			Name:  cName,
			Value: cValue,
		})
	}
}

// getBodyReader serializes the body for a network request. See the test file for examples
func getBodyReader(rawBody interface{}) (io.Reader, error) {
	var bodyReader io.Reader

	if rawBody != nil {
		switch body := rawBody.(type) {
		case map[string]string:
			jsonBody, err := json.Marshal(body)
			if err != nil {
				if debug {
					panic("Unable to read the response body")
				}
				return nil, newError(ErrMarshallingPostRequest, "couldn't serialize map of strings to JSON.")
			}
			bodyReader = bytes.NewBuffer(jsonBody)
		case url.Values:
			bodyReader = strings.NewReader(body.Encode())
		case []byte: //expects JSON format
			bodyReader = bytes.NewBuffer(body)
		case string: //expects JSON format
			bodyReader = strings.NewReader(body)
		default:
			return nil, newError(ErrMarshallingPostRequest, fmt.Sprintf("Cannot handle body type %T", rawBody))
		}
	}

	return bodyReader, nil
}

// PostWithClient returns the HTML returned by the url using a provided HTTP client
// The type of the body must conform to one of the types listed in func getBodyReader()
func PostWithClient(url string, bodyType string, body interface{}, client *http.Client) (string, error) {
	timer := startTimer("HTTP POST: "+url, DebugVerbose)
	defer timer.finish()

	logHTTPRequest("POST", url, Headers)

	bodyReader, err := getBodyReader(body)
	if err != nil {
		return "todo:", err
	}

	req, _ := http.NewRequest("POST", url, bodyReader)
	Header("Content-Type", bodyType)
	setHeadersAndCookies(req)

	if debug {
		// Save a copy of this request for debugging.
		requestDump, err := httputil.DumpRequest(req, true)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(requestDump))
	}

	// Perform request
	resp, err := client.Do(req)

	if err != nil {
		if debug {
			panic("Couldn't perform POST request to " + url)
		}
		return "", newError(ErrCreatingPostRequest, "couldn't perform POST request to "+url)
	}
	defer resp.Body.Close()

	logHTTPResponse(resp.StatusCode, int(resp.ContentLength))

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		if debug {
			panic("Unable to read the response body")
		}
		return "", newError(ErrReadingResponse, "unable to read the response body")
	}
	return string(bytes), nil
}

// Get returns the HTML returned by the url as a string using the default HTTP client
func Get(url string) (string, error) {
	client := &http.Client{
		Timeout: DefaultTimeout,
	}
	return GetWithClient(url, client)
}

// GetWithTimeout returns the HTML returned by the url with a custom timeout
func GetWithTimeout(url string, timeout time.Duration) (string, error) {
	client := &http.Client{
		Timeout: timeout,
	}
	return GetWithClient(url, client)
}

// PostWithTimeout returns the HTML returned by the url with a custom timeout
func PostWithTimeout(url string, bodyType string, body interface{}, timeout time.Duration) (string, error) {
	client := &http.Client{
		Timeout: timeout,
	}
	return PostWithClient(url, bodyType, body, client)
}

// Post returns the HTML returned by the url as a string using the default HTTP client
func Post(url string, bodyType string, body interface{}) (string, error) {
	client := &http.Client{
		Timeout: DefaultTimeout,
	}
	return PostWithClient(url, bodyType, body, client)
}

// PostForm is a convenience method for POST requests that
func PostForm(url string, data url.Values) (string, error) {
	client := &http.Client{
		Timeout: DefaultTimeout,
	}
	return PostWithClient(url, "application/x-www-form-urlencoded", data, client)
}
