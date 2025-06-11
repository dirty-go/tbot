package tbot

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type responseParameters struct {
	MigrateToChatID int `json:"migrate_to_chat_id"`
	RetryAfter      int `json:"retry_after"`
}

type apiResponse struct {
	OK          bool                `json:"ok"`
	Result      json.RawMessage     `json:"result"`
	Description string              `json:"description"`
	ErrorCode   int                 `json:"error_code"`
	Parameter   *responseParameters `json:"parameter,omitempty"`
}

var netTransport = &http.Transport{
	TLSHandshakeTimeout:   10 * time.Second,
	MaxIdleConns:          10,
	IdleConnTimeout:       30 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
}

func (c *Client) sendRequest(method string, request url.Values, response any) error {
	var err error
	var req *http.Request
	var resp *http.Response
	endPoint := fmt.Sprintf(c.url, method)
	if request == nil {
		req, err = http.NewRequest(http.MethodPost, endPoint, nil)
	} else {
		req, err = http.NewRequest(http.MethodPost, endPoint, strings.NewReader(request.Encode()))

	}
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "application/json")

	netClient := &http.Client{
		Timeout:   120 * time.Second,
		Transport: netTransport,
	}

	resp, err = netClient.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	// resp.StatusCode is between 200 and 300.
	// This is because an HTTP status code with the form 2XX signifies a successful HTTP POST request
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var apiResp apiResponse
		err = json.NewDecoder(resp.Body).Decode(&apiResp)
		if err != nil {
			return fmt.Errorf("unable to decode response: %v", err)
		}
		err = resp.Body.Close()
		if err != nil {
			return fmt.Errorf("unable to close response body: %v", err)
		}

		if !apiResp.OK {
			return fmt.Errorf("%d : %s", apiResp.ErrorCode, apiResp.Description)
		}

		return json.Unmarshal(apiResp.Result, response)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return fmt.Errorf("%s : %s", resp.Status, string(b))
}

func (c *Client) sendRequestWithFiles(method string, request url.Values, response any, files ...inputFile) error {
	var err error
	var req *http.Request
	var resp *http.Response
	endPoint := fmt.Sprintf(c.url, method)

	r, w := io.Pipe()

	done := make(chan struct{})

	mw := multipart.NewWriter(w)
	netClient := &http.Client{
		Timeout:   120 * time.Second,
		Transport: netTransport,
	}
	go func() {
		defer close(done)
		req, err = http.NewRequest(http.MethodPost, endPoint, r)
		if err != nil {
			c.logger.Error(err)
			return
		}
		req.Header.Set("Content-Type", mw.FormDataContentType())
		resp, err = netClient.Do(req)
	}()

	for k := range request {
		_ = mw.WriteField(k, request.Get(k))
	}

	for _, file := range files {
		f, err := os.Open(file.name)
		if err != nil {
			return err
		}
		fileWriter, err := mw.CreateFormFile(file.field, file.name)
		if err != nil {
			return err
		}

		_, _ = io.Copy(fileWriter, f)
		_ = f.Close()
	}

	_ = mw.Close()
	_ = w.Close()

	<-done // post request is done
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	// resp.StatusCode is between 200 and 300.
	// This is because an HTTP status code with the form 2XX signifies a successful HTTP POST request
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var apiResp apiResponse
		err = json.NewDecoder(resp.Body).Decode(&apiResp)
		if err != nil {
			return fmt.Errorf("unable to decode response: %v", err)
		}
		err = resp.Body.Close()
		if err != nil {
			return fmt.Errorf("unable to close response body: %v", err)
		}

		if !apiResp.OK {
			return fmt.Errorf(apiResp.Description)
		}

		return json.Unmarshal(apiResp.Result, response)
	}

	return fmt.Errorf(resp.Status)
}
