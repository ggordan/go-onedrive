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

// Items represents a list of Items
type Items struct {
	List []*Item `json:"value"`
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

// GetByID returns an item with the specified ID.
func (is *ItemService) GetByID(itemID string) (*Item, *http.Response, error) {
	path := fmt.Sprintf("/drive/items/%s", itemID)

	req, err := is.newRequest("GET", path, nil)
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

// GetChildrenByID returns a collection of all the Items under an Item
func (is *ItemService) GetChildrenByID(itemID string) (*Items, *http.Response, error) {
	path := fmt.Sprintf("/drive/items/%s/children", itemID)

	req, err := is.newRequest("GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	items := new(Items)
	resp, err := is.do(req, nil)
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
	req, err := is.newRequest("PUT", path, folder)
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
