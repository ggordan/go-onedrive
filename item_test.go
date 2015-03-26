package onedrive

import (
	"fmt"
	"net/http"
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

func newBaseItem(name, id string, size int64) *Item {
	return &Item{
		Name:                 name,
		ID:                   id,
		Size:                 size,
		CTag:                 "ctag",
		ETag:                 "etag",
		WebURL:               "https://onedrive.live.com/redir?page=self&resid=" + id,
		CreatedDateTime:      parseTime("2015-03-08T03:26:46.443Z"),
		LastModifiedDateTime: parseTime("2015-03-09T12:05:17.333Z"),
		DownloadURL:          "https://download-url.com/someid",
		CreatedBy: &IdentitySet{
			User: userIdentity,
		},
		LastModifiedBy: &IdentitySet{
			User:        userIdentity,
			Application: appIdentity,
		},
		ParentReference: &ItemReference{
			DriveID: "0123456789abc",
			ID:      "0123456789abc!104",
			Path:    "/drive/root:/Test folder 1",
		},
	}
}

func newAudioItem(name, id string, size int64, audio *AudioFacet, file *FileFacet) *Item {
	item := newBaseItem(name, id, size)
	item.Audio = audio
	item.File = file
	return item
}

func newFolderItem(name, id string, size int64, folder *FolderFacet) *Item {
	item := newBaseItem(name, id, size)
	item.DownloadURL = ""
	item.ParentReference = nil
	item.Folder = folder
	return item
}

func newImageItem(name, id string, size int64, image *ImageFacet, file *FileFacet) *Item {
	item := newBaseItem(name, id, size)
	item.Image = image
	item.File = file
	return item
}

func newPhotoItem(name, id string, size int64, image *ImageFacet, file *FileFacet, photo *PhotoFacet) *Item {
	item := newBaseItem(name, id, size)
	item.Image = image
	item.File = file
	item.Photo = photo
	return item
}

func newVideoItem(name, id string, size int64, file *FileFacet, photo *PhotoFacet, location *LocationFacet, video *VideoFacet) *Item {
	item := newBaseItem(name, id, size)
	item.Location = location
	item.Photo = photo
	item.File = file
	item.Video = video
	return item
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

	audioItem := newAudioItem("01 Perth.mp3", "0123456789abc!121", 7904129,
		newAudioFacet("Bon Iver", "Bon Iver", "Bon Iver", 238, "Justin Vernon", "", 1, 1, 262138, "\u00a7", false, true, "Perth", 1, 0, 2011),
		newFileFacet("audio/mpeg", newHashesFacet("61AA4245BAB442EB18920B293C3E24B44457E665", "CEF984EA")),
	)

	folderItem := newFolderItem("root", "0123456789abc!101", 10655823, newFolderFacet(3))
	folderItem.ParentReference = &ItemReference{
		DriveID: "0123456789abc",
		ID:      "0123456789abc!104",
		Path:    "/drive/root:/Test folder 1",
	}

	imageItem := newImageItem("sydney_opera_house_2011-1920x1080.jpg", "0123456789abc!110", 666657,
		newImageFacet(1080, 1920),
		newFileFacet("image/jpeg", newHashesFacet("6968B0F0934762EC44ADBC90959FAC6F03FBE211", "FEBB5160")),
	)

	photoItem := newPhotoItem("IMG_2538.JPG", "0123456789abc!119", 403305,
		newImageFacet(480, 720),
		newFileFacet("image/jpeg", newHashesFacet("D528F485B3A594A36F00ED7633DC2AE1C442A93D", "4DD1C268")),
		newPhotoFacet(parseTime("2013-11-28T11:57:27Z"), "Canon", "Canon EOS 600D", 9.0, 200.0, 1.0, 18.0, 0),
	)

	videoItem := newVideoItem("Video 10-03-2015 20 34 37.mov", "0123456789abc!123", 4114667,
		newFileFacet("video/mp4", newHashesFacet("990944543C492C90A703A31BFFEED09BBFCB65BC", "CBBE2450")),
		newPhotoFacet(parseTime("2015-03-10T13:34:35Z"), "Apple", "iPhone 5", 0.0, 0.0, 0.0, 0.0, 0),
		newLocationFacet(7.824, 51.5074, -0.2377),
		newVideoFacet(16382248, 1833, 1920, 1080),
	)

	tt := []struct {
		itemID         string
		expectedStatus int
		expectedOut    *Item
	}{
		{"audio", 200, audioItem},
		{"folder", 200, folderItem},
		{"image", 200, imageItem},
		{"photo", 200, photoItem},
		{"video", 200, videoItem},
	}
	for i, tst := range tt {
		mux.HandleFunc(itemURIFromID(tst.itemID), fileWrapperHandler(validFixtureFromItemID(tst.itemID), http.StatusOK))
		item, _, err := oneDrive.Items.Get(tst.itemID)
		if err != nil {
			t.Fatalf("Problem fetching the default drive: %s", err.Error())
		}
		if !reflect.DeepEqual(item, tst.expectedOut) {
			t.Errorf("[%d] Got \n%v Expected \n%v", i, *item, *tst.expectedOut)
		}
	}
}

func TestGetItemInvalid(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(itemURIFromID("missing"), fileWrapperHandler("fixtures/request.invalid.notFound.json", http.StatusNotFound))
	missingDrive, resp, err := oneDrive.Items.Get("missing")
	if missingDrive != nil {
		t.Fatalf("A drive was returned when an error was expected: %v", resp)
	}

	expectedErr := &Error{
		innerError{
			Code:    "itemNotFound",
			Message: "Item Does Not Exist",
			InnerError: &innerError{
				Code: "itemDoesNotExist",
				InnerError: &innerError{
					Code: "folderDoesNotExist",
				},
			},
		},
	}

	if !reflect.DeepEqual(expectedErr, err) {
		t.Errorf("Got %v Expected %v", err, expectedErr)
	}
}

func TestGetDefaultDriveRootFolder(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/drive/root", fileWrapperHandler("fixtures/drive.root.valid.json", http.StatusOK))
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

func TestListItemChildren(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/drive/items/some-id/children", fileWrapperHandler("fixtures/item.children.valid.json", http.StatusOK))
	items, _, err := oneDrive.Items.ListChildren("some-id")
	if err != nil {
		t.Fatal(err)
	}

	if got, want := len(items.Collection), 3; got != want {
		t.Fatalf("Got %d Expected %d", got, want)
	}

	if got, want := items.Collection[0].Folder.ChildCount, int64(10); got != want {
		t.Fatalf("Got %d Expected %d folder child items", got, want)
	}

	if got, want := items.Collection[1].Name, "Test folder 2"; got != want {
		t.Fatalf("Got %q Expected %q", got, want)
	}
}

func TestListItemChildrenInvalid(t *testing.T) {
	setup()
	defer teardown()
	mux.HandleFunc("/drive/items/some-id/children", fileWrapperHandler("fixtures/item.children.invalid.json", http.StatusOK))
	_, resp, err := oneDrive.Items.ListChildren("some-id")
	if err == nil {
		t.Fatalf("Expected error, got : %v", resp)
	}
}

func TestCreateFolder(t *testing.T) {
	setup()
	defer teardown()

	folderItem := newFolderItem("root", "0123456789abc!101", 10655823, newFolderFacet(3))
	mux.HandleFunc("/drive/items/0123456789abc!104/children/root", fileWrapperHandler("fixtures/item.folder.valid.json", http.StatusOK))
	item, _, err := oneDrive.Items.CreateFolder("0123456789abc!104", "root")
	if err != nil {
		t.Fatalf("An error occured while attempting to create a folder: %s", err)
	}

	if got, want := item, folderItem; reflect.DeepEqual(got, want) {
		t.Fatalf("Got %v Expected %v", *got, *want)
	}
}

func TestDeleteItem(t *testing.T) {
	setup()
	defer teardown()
}
