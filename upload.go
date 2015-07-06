package onedrive

import (
	"fmt"
	"net/http"
	"os"
)

const oneHundredMB = 104857600

// UploadFromURL allows your app to upload an item to OneDrive by providing a URL.
// OneDrive will download the file directly from a remote server so your app
// doesn't have to upload the file's bytes.
// See: http://onedrive.github.io/items/upload_url.htm
func (is *ItemService) UploadFromURL(parentID, name, webURL string) (*Item, *http.Response, error) {
	requestHeaders := map[string]string{
		"Prefer": "respond-async",
	}

	newFile := newWebUpload{
		webURL, name, new(FileFacet),
	}

	path := fmt.Sprintf("/drive/items/%s/children", parentID)
	req, err := is.newRequest("POST", path, requestHeaders, newFile)
	if err != nil {
		return nil, nil, err
	}

	item := new(Item)
	resp, err := is.do(req, item)
	if err != nil {
		return nil, resp, err
	}

	return item, resp, nil
}

// SimpleUpload allows you to provide the contents of a new file or update the
// contents of an existing file in a single API call. This method only supports
// files up to 100MB in size. For larger files use ResumableUpload().
// See: https://dev.onedrive.com/items/upload_put.htm
func (is ItemService) SimpleUpload(folderID string, file *os.File) (*Item, *http.Response, error) {
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, nil, err
	}

	if fileInfo.Size() >= oneHundredMB {
		return nil, nil, ErrFileTooLarge
	}

	path := fmt.Sprintf("/drive/items/%s/children/%s/content", folderID, file.Name())
	req, err := is.newRequest("PUT", path, nil, file)

	if err != nil {
		return nil, nil, err
	}

	item := new(Item)
	resp, err := is.do(req, item)
	if err != nil {
		return nil, resp, err
	}

	return item, resp, nil
}
