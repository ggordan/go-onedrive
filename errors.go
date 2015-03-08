package onedrive

// Error defines the basic structure of errors that are returned from the OneDrive
// API.
// See: http://onedrive.github.io/misc/errors.htm
type Error struct {
	Body struct {
		Code       string `json:"code"`
		Message    string `json:"message"`
		InnerError struct {
			Code string `json:"code"`
		} `json:"innererror"`
	} `json:"error"`
}

func (e Error) Error() string {
	return e.Body.Message
}
