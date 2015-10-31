package onedrive

import (
	"net/http"
	"time"
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
	// Private
	throttle time.Time
}

// NewOneDrive returns a new OneDrive client to enable you to communicate with
// the API
func NewOneDrive(c *http.Client, debug bool) *OneDrive {
	drive := OneDrive{
		Client:   c,
		BaseURL:  baseURL,
		Debug:    debug,
		throttle: time.Now(),
	}
	drive.Drives = &DriveService{&drive}
	drive.Items = &ItemService{&drive}
	return &drive
}

func (od *OneDrive) throttleRequest(time time.Time) {
	od.throttle = time
}
