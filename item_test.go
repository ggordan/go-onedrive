package onedrive

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func parseTime(t string) time.Time {
	pt, err := time.Parse(time.RFC3339Nano, t)
	if err != nil {
		panic(err)
	}
	return pt
}

func validFixtureFromItemID(itemID string) string {
	switch itemID {
	case "", "root":
		return "fixtures/drive.valid.default.json"
	default:
		return fmt.Sprintf("fixtures/item.%s.valid.json", itemID)
	}
}

var userIdentity = &Identity{
	DisplayName: "Gordan Grasarevic",
	ID:          "0123456789abc",
}
var deviceIdentity = &Identity{
	DisplayName: "test app",
	ID:          "test-id",
}
var appIdentity = &Identity{
	DisplayName: "OneDrive website",
	ID:          "44048800",
}

func newItem(name, id string, size int64) *Item {
	return &Item{
		CreatedDateTime:      parseTime("2015-03-08T03:26:46.443Z"),
		CTag:                 "ctag",
		ETag:                 "etag",
		ID:                   id,
		LastModifiedDateTime: parseTime("2015-03-09T12:05:17.333Z"),
		Name:                 name,
		Size:                 size,
		WebURL:               "https://onedrive.live.com/redir?page=self&resid=" + id,
	}
}

func TestItemURIFromID(t *testing.T) {
	tt := []struct {
		in, out string
	}{
		{"", "/drive/root"},
		{"root", "/drive/root"},
		{"123", "/drive/items/123"},
	}
	for i, tst := range tt {
		if got, want := itemURIFromID(tst.in), tst.out; got != want {
			t.Errorf("[%d] Got %q Expected %q", i, got, want)
		}
	}

}

func TestGetItem(t *testing.T) {
	setup()
	defer teardown()

	tt := []struct {
		itemID         string
		expectedStatus int
		expectedOut    *Drive
	}{
		{"folder", 200, expectedDefaultDrive},
		{"image", 200, expectedDefaultDrive},
		{"photo", 200, expectedDefaultDrive},
		{"video", 200, expectedDefaultDrive},
	}
	for i, tst := range tt {
		mux.HandleFunc(itemURIFromID(tst.itemID), fileWrapperHandler(validFixtureFromItemID(tst.itemID), 200))
		drive, _, err := oneDrive.Items.Get(tst.itemID)
		if err != nil {
			t.Fatalf("Problem fetching the default drive: %s", err.Error())
		}
		if !reflect.DeepEqual(drive, tst.expectedOut) {
			t.Errorf("[%d] Got %v Expected %v", i, drive, tst.expectedOut)
		}
	}
}

func TestGetDefaultDriveRootFolder(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/drive/root", fileWrapperHandler("fixtures/drive.root.valid.json", 200))
	root, _, err := oneDrive.Items.GetDefaultDriveRootFolder()
	if err != nil {
		t.Fatalf("Problem fetching the root drive: %s", err.Error())
	}

	expectedItem := newItem("root", "EBCEC5405197F0B!101", 17546845)
	expectedItem.CreatedBy = &IdentitySet{
		User: userIdentity,
	}
	expectedItem.LastModifiedBy = &IdentitySet{
		User:        userIdentity,
		Application: appIdentity,
	}
	expectedItem.Folder = &FolderFacet{
		ChildCount: 3,
	}

	if got, want := root, expectedItem; !reflect.DeepEqual(got, want) {
		t.Errorf("Got %v Expected %v", *got, *want)
	}
}
