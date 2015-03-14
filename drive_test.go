package onedrive

import (
	"fmt"
	"reflect"
	"testing"
)

var expectedDefaultDrive = &Drive{
	ID:        "0123456789abc",
	DriveType: "personal",
	Owner: &IdentitySet{
		User: &Identity{
			DisplayName: "Gordan Grasarevic",
			ID:          "0123456789abc",
		},
	},
	Quota: &Quota{
		Deleted:   0,
		Remaining: 16095471537,
		State:     "normal",
		Total:     16106127360,
		Used:      10655823,
	},
}

func validFixtureFromDriveID(driveID string) string {
	switch driveID {
	case "", "root":
		return "fixtures/drive.valid.default.json"
	default:
		return fmt.Sprintf("fixtures/drive.valid.%s.json", driveID)
	}
}

func TestDriveURIFromID(t *testing.T) {
	tt := []struct {
		in, out string
	}{
		{"", "/drive"},
		{"123", "/drives/123"},
	}
	for i, tst := range tt {
		if got, want := driveURIFromID(tst.in), tst.out; got != want {
			t.Errorf("[%d] Got %q Expected %q", i, got, want)
		}
	}
}

func TestDriveChildrenURIFromID(t *testing.T) {
	tt := []struct {
		in, out string
	}{
		{"", "/drive/root/children"},
		{"root", "/drive/root/children"},
		{"default", "/drive/root/children"},
		{"test-drive", "/drives/test-drive/root/children"},
	}
	for i, tst := range tt {
		if got, want := driveChildrenURIFromID(tst.in), tst.out; got != want {
			t.Errorf("[%d] Got %q Expected %q", i, got, want)
		}
	}
}

func TestGet(t *testing.T) {
	setup()
	defer teardown()

	tt := []struct {
		driveID        string
		expectedStatus int
		expectedOut    *Drive
	}{
		{"", 200, expectedDefaultDrive},
		{"root", 200, expectedDefaultDrive},
		{"test-id", 200, &Drive{
			ID:        "test-id",
			DriveType: "consumer",
			Owner: &IdentitySet{
				Device: &Identity{
					DisplayName: "test app",
					ID:          "test-id",
				},
			},
			Quota: &Quota{
				Deleted: 1,
			},
		},
		},
	}
	for i, tst := range tt {
		mux.HandleFunc(driveURIFromID(tst.driveID), fileWrapperHandler(validFixtureFromDriveID(tst.driveID), 200))
		drive, _, err := oneDrive.Drives.Get(tst.driveID)
		if err != nil {
			t.Fatalf("Problem fetching the default drive: %s", err.Error())
		}
		if !reflect.DeepEqual(drive, tst.expectedOut) {
			t.Errorf("[%d] Got %v Expected %v", i, drive, tst.expectedOut)
		}
	}
}

func TestGetMissing(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/drives/missing-drive", fileWrapperHandler("fixtures/request.invalid.notFound.json", 404))
	missingDrive, resp, err := oneDrive.Drives.Get("missing-drive")
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

func TestGetMalformed(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/drives/malformed-drive", fileWrapperHandler("fixtures/drive.invalid.malformed.json", 200))
	_, _, err := oneDrive.Drives.Get("malformed-drive")
	if err == nil {
		t.Fatalf("Expected error, got: %v", err)
	}
}

func TestGetDefault(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/drive", fileWrapperHandler("fixtures/drive.valid.default.json", 200))
	defaultDrive, _, err := oneDrive.Drives.GetDefault()
	if err != nil {
		t.Fatalf("Problem fetching the default drive: %s", err.Error())
	}
	if !reflect.DeepEqual(defaultDrive, expectedDefaultDrive) {
		t.Errorf("Got %v Expected %v", defaultDrive, expectedDefaultDrive)
	}
}

func TestListAllDrives(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/drives", fileWrapperHandler("fixtures/drive.collection.valid.json", 200))
	drives, _, err := oneDrive.Drives.ListAll()
	if err != nil {
		t.Fatalf("Problem fetching the drive list: %s", err.Error())
	}

	expectedDrives := &Drives{
		Collection: []*Drive{
			expectedDefaultDrive,
		},
	}
	if !reflect.DeepEqual(drives, expectedDrives) {
		t.Errorf("Got %v Expected %v", *drives, *expectedDrives)
	}
}

func TestListAllDrivesInvalid(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/drives", fileWrapperHandler("fixtures/request.invalid.badArgument.json", 400))
	_, resp, err := oneDrive.Drives.ListAll()
	if err == nil {
		t.Fatalf("Expected err, got : %v", resp)
	}

	expectedErr := &Error{
		innerError{
			Code:    "invalidArgument",
			Message: "Bad Argument",
			InnerError: &innerError{
				Code: "badArgument",
			},
		},
	}

	if !reflect.DeepEqual(err, expectedErr) {
		t.Errorf("Got %v Expected %v", err, expectedErr)
	}
}

func TestListDriveChildren(t *testing.T) {
	setup()
	defer teardown()

	testFolder1 := newItem("Test folder 1", "0123456789abc!104", 9989166)
	testFolder1.Folder = &FolderFacet{
		ChildCount: 10,
	}
	testFolder1.ParentReference = &ItemReference{
		DriveID: "0123456789abc",
		ID:      "0123456789abc!101",
		Path:    "/drive/root:",
	}
	testFolder1.CreatedBy = &IdentitySet{
		Application: appIdentity,
		User:        userIdentity,
	}
	testFolder1.LastModifiedBy = &IdentitySet{
		Application: appIdentity,
		User:        userIdentity,
	}

	mux.HandleFunc(driveChildrenURIFromID(""), fileWrapperHandler("fixtures/drive.children.valid.json", 200))
	items, _, err := oneDrive.Drives.ListChildren("")
	if err != nil {
		t.Fatalf("Problem fetching the drive children: %s", err.Error())
	}

	if got, want := len(items.Collection), 3; got != want {
		t.Fatalf("Got %d, Expected %d children", got, want)
	}

	if got, want := items.Collection[0], testFolder1; !reflect.DeepEqual(got, want) {
		t.Errorf("Got %v Expected %v", *got, *want)
	}
}

func TestListDriveChildrenInvalid(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(driveChildrenURIFromID(""), fileWrapperHandler("fixtures/drive.children.invalid.json", 200))
	_, resp, err := oneDrive.Drives.ListChildren("")
	if err == nil {
		t.Fatalf("Expected error, got : %v", resp)
	}
}
