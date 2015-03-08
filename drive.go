package onedrive

import (
	"fmt"
	"net/http"
)

type StorageState string

const (
	NormalStorageState   StorageState = "normal"
	NearingStorageState  StorageState = "nearing"
	CriticalStorageState StorageState = "critical"
	ExceededStorageState StorageState = "exceeded"
)

// DriveService manages the communication with Drive related API endpoints
type DriveService struct {
	*OneDrive
}

// The Identity type represents an identity of an actor. For example, and actor
// can be a user, device, or application.
// See: http://onedrive.github.io/resources/identity.htm
type Identity struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
}

// The IdentitySet type is a keyed collection of Identity objects. It is used
// to represent a set of identities associated with various events for an item,
// such as created by or last modified by.
// See: http://onedrive.github.io/resources/identitySet.htm
type IdentitySet struct {
	User        *Identity `json:"user"`
	Application *Identity `json:"application"`
	Device      *Identity `json:"device"`
}

// The Quota facet groups storage space quota-related information on OneDrive
// into a single structure.
// See: http://onedrive.github.io/facets/quotainfo_facet.htm
type Quota struct {
	Total     int64  `json:"total"`
	Used      int64  `json:"used"`
	Remaining int64  `json:"remaining"`
	Deleted   int64  `json:"deleted"`
	State     string `json:"state"`
}

// Drives represents a list of Drives
type Drives struct {
	List []*Drive `json:"value"`
}

// The Drive resource represents a drive in OneDrive. It provides information
// about the owner of the drive, total and available storage space, and exposes
// a collection of all the Items contained within the drive.
// See: http://onedrive.github.io/resources/drive.htm
type Drive struct {
	ID        string      `json:"id"`
	DriveType string      `json:"driveType"`
	Owner     IdentitySet `json:"owner"`
	Quota     Quota       `json:"quota"`
}

// Get returns a Drive for the authenticated user. If no driveID is provided
// the users default Drive is returned. A user will always have at least one
// Drive available -- the default Drive.
func (ds *DriveService) Get(driveID string) (*Drive, *http.Response, error) {
	var path string

	switch driveID {
	case "":
		path = "/drive"
		break
	case "root":
		path = "/drive/root"
		break
	default:
		path = fmt.Sprintf("/drives/%s", driveID)
	}

	req, err := ds.newRequest("GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	drive := new(Drive)
	resp, err := ds.do(req, drive)
	if err != nil {
		return nil, resp, err
	}

	return drive, resp, nil
}

// GetDefaultDrive is a convenience function to return the users default Drive
func (ds *DriveService) GetDefaultDrive() (*Drive, *http.Response, error) {
	return ds.Get("")
}

// GetRootDrive is a convenience function to return the users root folder of the
// users default Drive
func (ds *DriveService) GetRootDrive() (*Drive, *http.Response, error) {
	return ds.Get("root")
}

// List returns all the Drives available to the authenticated user
func (ds *DriveService) List() (*Drives, *http.Response, error) {
	req, err := ds.newRequest("GET", "/drives", nil)
	if err != nil {
		return nil, nil, err
	}

	drives := new(Drives)
	resp, err := ds.do(req, drives)
	if err != nil {
		return nil, resp, err
	}

	return drives, resp, nil
}

func (ds *DriveService) ListChildren(driveID string) (*Items, *http.Response, error) {
	var path string

	switch driveID {
	case "root":
		path = "/drive/root/children"
		break
	default:
		path = fmt.Sprintf("/drives/%s/root/children", driveID)
	}

	req, err := ds.newRequest("GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	items := new(Items)
	resp, err := ds.do(req, items)
	if err != nil {
		return nil, resp, err
	}

	return items, resp, nil
}
