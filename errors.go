package onedrive

// Error defines the basic structure of errors that are returned from the OneDrive
// API.
// See: http://onedrive.github.io/misc/errors.htm

type InnerError struct {
	Code       string      `json:"code"`
	Message    string      `json:"message"`
	InnerError *InnerError `json:"innererror"`
}

type Error struct {
	InnerError `json:"error"`
}

func (e Error) Error() string {
	return e.Message
}
