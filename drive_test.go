package onedrive

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
)

func fileWrapperHandler(handler func(w http.ResponseWriter, r *http.Request, file string), file string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, file)
	}
}

func defaultDriveHandler(w http.ResponseWriter, r *http.Request, file string) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	w.Write(b)
}

func TestGetDefaultDrive(t *testing.T) {
	setup()
	defer teardown()

	fb, err := ioutil.ReadFile("fixtures/drive.valid.json")
	if err != nil {
		t.Fatal(err)
	}

	var expectedDrive Drive
	if err := json.Unmarshal(fb, &expectedDrive); err != nil {
		t.Fatal(err)
	}

	mux.HandleFunc("/drive", fileWrapperHandler(defaultDriveHandler, "fixtures/drive.valid.json"))

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
