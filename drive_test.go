package onedrive

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

func TestGetDefaultDriveValid(t *testing.T) {
	setup()
	defer teardown()

	testFile := "fixtures/drive.valid.json"

	fb, err := ioutil.ReadFile(testFile)
	if err != nil {
		t.Fatal(err)
	}

	var expectedDrive Drive
	if err := json.Unmarshal(fb, &expectedDrive); err != nil {
		t.Fatal(err)
	}

	mux.HandleFunc("/drive", fileWrapperHandler(testFile))

	drive, resp, err := oneDrive.Drives.GetDefaultDrive()
	if err != nil {
		t.Fatalf("Problem fetching the default drive: %s", err.Error())
	}

	if resp.StatusCode != 200 {
		t.Fatalf("Expected response to be 200 got %d", resp.StatusCode)
	}

	if drive.ID != expectedDrive.ID {
		t.Fatalf("Expected ID to be %q, but got %q", expectedDrive.ID, drive.ID)
	}

}
