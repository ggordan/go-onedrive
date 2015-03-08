package onedrive

import "net/http"

func buildRequestHeaders() *http.Request {

	return nil
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
