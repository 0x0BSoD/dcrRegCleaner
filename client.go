package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	Url        string
	http       *http.Client
	TagsToStay int
}

// Parameters for new client
type Parameters struct {
	Url string
}

// ===========================================
// Get - raw request, return []byte and error
func (c *Client) Get(endpoint string) ([]byte, http.Header, error) {

	resp, err := c.getResponse("GET", endpoint, nil)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("request: %s", endpoint)
		return nil, nil, fmt.Errorf("error during executing request, status: %d, error if exist: %s", resp.StatusCode, err)
	}

	bodyByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("error during read response body: %s", err)
	}

	return bodyByte, resp.Header, nil
}

// ===========================================
// Delete - raw request, return []byte and error
func (c *Client) Delete(endpoint string) ([]byte, http.Header, error) {

	resp, err := c.getResponse("DELETE", endpoint, nil)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("error during executing request: %s", err)
	}

	bodyByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("error during read response body: %s", err)
	}

	return bodyByte, resp.Header, nil
}

// ===========================================
// getResponse return http.Response and error [PRIVATE]
func (c *Client) getResponse(method, endpoint string, body []byte) (*http.Response, error) {
	urlStr := c.Url + endpoint

	// check full URL
	_url, err := url.Parse(urlStr)
	if err != nil {
		return nil, fmt.Errorf("error during parsing request URL: %s", err)
	}

	// read body if present
	var bodyReader io.Reader
	if body != nil {
		bodyReader = bytes.NewReader(body)
	}

	req, err := http.NewRequest(method, _url.String(), bodyReader)
	if err != nil {
		return nil, fmt.Errorf("error during creation of request: %s", err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Accept", "application/vnd.docker.distribution.manifest.list.v2+json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	// defer resp.Body.Close()

	return resp, nil
}

// ===========================================
// NewClient return client instance
func NewClient(c Config) *Client {
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	return &Client{
		Url:        fmt.Sprintf("%s/%s", c.Registry, c.ApiVersion),
		http:       &http.Client{Transport: tr},
		TagsToStay: c.TagsToStay,
	}
}
