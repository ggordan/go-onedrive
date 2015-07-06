package onedrive

import (
	"fmt"
	"net/http"
	"time"
)

// ItemService manages the communication with Item related API endpoints
type ItemService struct {
	*OneDrive
}

// The Thumbnail resource type represents a thumbnail for an image, video,
// document, or any file or folder on OneDrive that has a graphical representation.
// See: http://onedrive.github.io/resources/thumbnail.htm
type Thumbnail struct {
	Width  int    `json:"width"`
	Height int    `json:"height"`
	URL    string `json:"url"`
	// Relationships
	Content []byte `json:"content"`
}

// The ThumbnailSet type is a keyed collection of Thumbnail objects. It is used
// to represent a set of thumbnails associated with a single file on OneDrive.
// See: http://onedrive.github.io/resources/thumbnailSet.htm
type ThumbnailSet struct {
	ID     string     `json:"id"`
	Small  *Thumbnail `json:"small"`
	Medium *Thumbnail `json:"medium"`
	Large  *Thumbnail `json:"large"`
}

// Items represents a collection of Items
type Items struct {
	Collection []*Item `json:"value"`
}

// The ItemReference type groups data needed to reference a OneDrive item across
// the service into a single structure.
// See: http://onedrive.github.io/resources/itemReference.htm
type ItemReference struct {
	DriveID string `json:"driveId"`
	ID      string `json:"id"`
	Path    string `json:"path"`
}

// The Item resource type represents metadata for an item in OneDrive. All
// top-level filesystem objects in OneDrive are Item resources. If an item is
// a Folder or File facet, the Item resource will contain a value for either
// the folder or file property, respectively.
// See: http://onedrive.github.io/resources/item.htm
type Item struct {
	ID                   string         `json:"id"`
	Name                 string         `json:"name"`
	ETag                 string         `json:"eTag"`
	CTag                 string         `json:"cTag"`
	CreatedBy            *IdentitySet   `json:"createdBy"`
	LastModifiedBy       *IdentitySet   `json:"lastModifiedBy"`
	CreatedDateTime      time.Time      `json:"createdDateTime"`
	LastModifiedDateTime time.Time      `json:"lastModifiedDateTime"`
	Size                 int64          `json:"size"`
	ParentReference      *ItemReference `json:"parentReference"`
	WebURL               string         `json:"webUrl"`
	File                 *FileFacet     `json:"file"`
	Folder               *FolderFacet   `json:"folder"`
	Image                *ImageFacet    `json:"image"`
	Photo                *PhotoFacet    `json:"photo"`
	Audio                *AudioFacet    `json:"audio"`
	Video                *VideoFacet    `json:"video"`
	Location             *LocationFacet `json:"location"`
	Deleted              *DeletedFacet  `json:"deleted"`
	// Instance attributes
	ConflictBehaviour string `json:"@name.conflictBehavior"`
	DownloadURL       string `json:"@content.downloadUrl"`
	SourceURL         string `json:"@content.sourceUrl"`
	// Relationships
	Content    []byte        `json:"content"`
	Children   []*Item       `json:"children"`
	Thumbnails *ThumbnailSet `json:"thumbnails"`
}

// driveURIFromID returns a valid request URI based on the ID of the drive.
// Mostly exists to simplify special cases such as the default and root drives.
func itemURIFromID(itemID string) string {
	switch itemID {
	case "", "root":
		return "/drive/root"
	default:
		return fmt.Sprintf("/drive/items/%s", itemID)
	}
}

// Get returns an item with the specified ID.
func (is *ItemService) Get(itemID string) (*Item, *http.Response, error) {
	req, err := is.newRequest("GET", itemURIFromID(itemID), nil, nil)
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

// GetDefaultDriveRootFolder is a convenience function to return the root folder
// of the users default Drive
func (is *ItemService) GetDefaultDriveRootFolder() (*Item, *http.Response, error) {
	return is.Get("root")
}

// ListChildren returns a collection of all the Items under an Item
func (is *ItemService) ListChildren(itemID string) (*Items, *http.Response, error) {
	reqURI := fmt.Sprintf("/drive/items/%s/children", itemID)
	req, err := is.newRequest("GET", reqURI, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	items := new(Items)
	resp, err := is.do(req, items)
	if err != nil {
		return nil, resp, err
	}

	return items, resp, nil
}

type newFolder struct {
	Name   string       `json:"name"`
	Folder *FolderFacet `json:"folder"`
}

// CreateFolder creates a new folder within the parent.
func (is *ItemService) CreateFolder(parentID, folderName string) (*Item, *http.Response, error) {
	folder := newFolder{
		Name:   folderName,
		Folder: new(FolderFacet),
	}

	path := fmt.Sprintf("/drive/items/%s/children/%s", parentID, folderName)
	req, err := is.newRequest("PUT", path, nil, folder)
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

type newWebUpload struct {
	SourceURL string     `json:"@content.sourceUrl"`
	Name      string     `json:"name"`
	File      *FileFacet `json:"file"`
}

// Move changes the parent folder for a OneDrive Item resource.
// See: http://onedrive.github.io/items/move.htm
func (is ItemService) Update(item *Item, ifMatch bool) (*Item, *http.Response, error) {
	requestHeaders := make(map[string]string)
	if ifMatch {
		requestHeaders["if-match"] = item.ETag
	}

	path := fmt.Sprintf("/drive/items/%s", item.ID)
	req, err := is.newRequest("PATCH", path, requestHeaders, item)
	if err != nil {
		return nil, nil, err
	}

	resp, err := is.do(req, item)
	if err != nil {
		return nil, resp, err
	}

	return item, resp, nil
}

// Delete removed a OneDrive item by using its ID. Note that deleting items
// using this method will move the items to the Recycle Bin, instead of
// permanently deleting them.
// See: http://onedrive.github.io/items/delete.htm
func (is *ItemService) Delete(itemID, eTag string) (bool, *http.Response, error) {
	requestHeaders := make(map[string]string)
	if eTag != "" {
		requestHeaders["if-match"] = eTag
	}

	path := fmt.Sprintf("/drive/items/%s", itemID)
	req, err := is.newRequest("DELETE", path, requestHeaders, nil)
	if err != nil {
		return false, nil, err
	}

	resp, err := is.do(req, nil)
	if err != nil {
		return false, resp, err
	}

	return (resp.StatusCode == statusNoContent), resp, err
}

// Move changes the parent folder for a OneDrive Item resource.
// See: http://onedrive.github.io/items/move.htm
func (is ItemService) Move(itemID, parentReference ItemReference) (*Item, *http.Response, error) {
	path := fmt.Sprintf("/drive/items/%s", itemID)
	req, err := is.newRequest("PATCH", path, nil, parentReference)
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

// Move changes the parent folder for a OneDrive Item resource.
// See: http://onedrive.github.io/items/move.htm
func (is ItemService) Copy(itemID, name string, parentReference ItemReference) (*Item, *http.Response, error) {
	copyAction := struct {
		ParentReference *ItemReference `json:"parentReference"`
		Name            string         `json:"name,omitempty"`
	}{&parentReference, name}

	// The copy action requires a Prefer: respond-async header
	headers := map[string]string{"Prefer": "respond-async"}

	path := fmt.Sprintf("/drive/items/%s/action.copy", itemID)
	req, err := is.newRequest("POST", path, headers, copyAction)
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
