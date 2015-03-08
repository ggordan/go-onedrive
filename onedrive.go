package onedrive

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

const (
	version   = "0.1"
	baseURL   = "https://api.onedrive.com/v1.0"
	userAgent = "github.com/ggordan/go-onedrive; version " + version
)

// OneDrive is the entry point for the client. It manages the communication with
// Microsoft OneDrive API
type OneDrive struct {
	Client *http.Client
	// When debug is set to true, the JSON response is formatted for better readability
	Debug   bool
	BaseURL string
	// Services
	Drives *DriveService
	Items  *ItemService
}

// NewOneDrive returns a new OneDrive client to enable you to communicate with
// the API
func NewOneDrive(c *http.Client, debug bool) *OneDrive {
	drive := OneDrive{
		Client:  c,
		BaseURL: baseURL,
		Debug:   debug,
	}
	drive.Drives = &DriveService{&drive}
	drive.Items = &ItemService{&drive}
	return &drive
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
