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

func TestGet(t *testing.T) {
	setup()
	defer teardown()

	tt := []struct {
		driveID        string
		expectedPath   string
		expectedStatus int
		expectedOut    *Drive
	}{
		{"", "/drive", 200, expectedDefaultDrive},
		{"root", "/drive/root", 200, expectedDefaultDrive},
		{"test-id", "/drives/test-id", 200, &Drive{
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
		if got, want := drivePathFromID(tst.driveID), tst.expectedPath; got != want {
			t.Errorf("Got %q Expected %q", got, want)
		}

		mux.HandleFunc(tst.expectedPath, fileWrapperHandler(validFixtureFromDriveID(tst.driveID)))
		drive, _, err := oneDrive.Drives.Get(tst.driveID)
		if err != nil {
			t.Fatalf("Problem fetching the default drive: %s", err.Error())
		}

		if !reflect.DeepEqual(drive, tst.expectedOut) {
			t.Errorf("[%d] Got %v Expected %v", i, drive, tst.expectedOut)
		}
	}
}

func TestGetDefaultDrive(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/drive", fileWrapperHandler("fixtures/drive.valid.default.json"))
	defaultDrive, resp, err := oneDrive.Drives.GetDefaultDrive()
	if err != nil {
		t.Fatalf("Problem fetching the default drive: %s", err.Error())
	}

	if resp.StatusCode != 200 {
		t.Fatalf("Expected response to be 200 got %d", resp.StatusCode)
	}

	if !reflect.DeepEqual(defaultDrive, expectedDefaultDrive) {
		t.Errorf("Got %v Expected %v", defaultDrive, expectedDefaultDrive)
	}
}

func TestGetRootDrive(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/drive/root", fileWrapperHandler("fixtures/drive.valid.default.json"))
	rootDrive, resp, err := oneDrive.Drives.GetRootDrive()
	if err != nil {
		t.Fatalf("Problem fetching the root drive: %s", err.Error())
	}

	if resp.StatusCode != 200 {
		t.Fatalf("Expected response to be 200 got %d", resp.StatusCode)
	}

	if !reflect.DeepEqual(rootDrive, expectedDefaultDrive) {
		t.Errorf("Got %v Expected %v", rootDrive, expectedDefaultDrive)
	}
}
