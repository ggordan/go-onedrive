package onedrive

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

func buildRequestHeaders() *http.Request {

	return nil
}

func createRequestBody(body interface{}) (io.ReadWriter, error) {
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}
	return buf, nil
}

func (od *OneDrive) newRequest(method, uri string, requestHeaders map[string]string, body interface{}) (*http.Request, error) {
	requestBody, err := createRequestBody(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, od.BaseURL+uri, requestBody)
	if err != nil {
		return nil, err
	}

	acceptHeader := "application/json"
	if od.Debug {
		acceptHeader += ";format=pretty"
	}

	req.Header.Add("Accept", acceptHeader)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", userAgent)

	if requestHeaders != nil {
		for header, value := range requestHeaders {
			req.Header.Set(header, value)
		}
	}

	return req, nil
}

func (od *OneDrive) do(req *http.Request, decodeInto interface{}) (*http.Response, error) {
	resp, err := od.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 && resp.StatusCode <= 507 {
		newErr := new(Error)
		if err := json.NewDecoder(resp.Body).Decode(newErr); err != nil {
			return resp, err
		}
		return resp, newErr
	}

	if decodeInto != nil {
		if err := json.NewDecoder(resp.Body).Decode(decodeInto); err != nil {
			return resp, err
		}
	}

	return resp, err
}
