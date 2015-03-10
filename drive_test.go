package onedrive

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

func TestGetDefaultDrive(t *testing.T) {
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

	defaultDrive, resp, err := oneDrive.Drives.GetDefaultDrive()
	if err != nil {
		t.Fatalf("Problem fetching the default drive: %s", err.Error())
	}

	if resp.StatusCode != 200 {
		t.Fatalf("Expected response to be 200 got %d", resp.StatusCode)
	}

	if got, want := defaultDrive.ID, expectedDrive.ID; got != want {
		t.Fatalf("Got %q Expected %q", got, want)
	}

	if got, want := defaultDrive.Owner.User.DisplayName, expectedDrive.Owner.User.DisplayName; got != want {
		t.Fatalf("Got %q Expected %q", got, want)
	}

}

func TestGetRootDrive(t *testing.T) {
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

	mux.HandleFunc("/drive/root", fileWrapperHandler(testFile))

	rootDrive, resp, err := oneDrive.Drives.GetRootDrive()
	if err != nil {
		t.Fatalf("Problem fetching the root drive: %s", err.Error())
	}

	if resp.StatusCode != 200 {
		t.Fatalf("Expected response to be 200 got %d", resp.StatusCode)
	}

	if got, want := rootDrive.ID, expectedDrive.ID; got != want {
		t.Fatalf("Got %q Expected %q", got, want)
	}

	if got, want := rootDrive.Owner.User.DisplayName, expectedDrive.Owner.User.DisplayName; got != want {
		t.Fatalf("Got %q Expected %q", got, want)
	}

}
