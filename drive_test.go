package onedrive

import (
	"reflect"
	"testing"
)

func TestGet(t *testing.T) {
}

func TestGetDefaultDrive(t *testing.T) {
	setup()
	defer teardown()

	expectedDrive := &Drive{
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

	mux.HandleFunc("/drive", fileWrapperHandler("fixtures/drive.valid.json"))
	defaultDrive, resp, err := oneDrive.Drives.GetDefaultDrive()
	if err != nil {
		t.Fatalf("Problem fetching the default drive: %s", err.Error())
	}

	if resp.StatusCode != 200 {
		t.Fatalf("Expected response to be 200 got %d", resp.StatusCode)
	}

	if !reflect.DeepEqual(defaultDrive, expectedDrive) {
		t.Errorf("Got %v Expected %v", defaultDrive, expectedDrive)
	}
}

func TestGetRootDrive(t *testing.T) {
	setup()
	defer teardown()

	rootDrive := &Drive{
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

	mux.HandleFunc("/drive/root", fileWrapperHandler("fixtures/drive.valid.json"))
	defaultDrive, resp, err := oneDrive.Drives.GetRootDrive()
	if err != nil {
		t.Fatalf("Problem fetching the root drive: %s", err.Error())
	}

	if resp.StatusCode != 200 {
		t.Fatalf("Expected response to be 200 got %d", resp.StatusCode)
	}

	if !reflect.DeepEqual(defaultDrive, rootDrive) {
		t.Errorf("Got %v Expected %v", defaultDrive, rootDrive)
	}
}
