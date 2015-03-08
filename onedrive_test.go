package onedrive

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	mux      *http.ServeMux
	server   *httptest.Server
	oneDrive *OneDrive
)

func setup() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
	oneDrive = NewOneDrive(http.DefaultClient, true)
	oneDrive.BaseURL = server.URL
}

func teardown() {
	server.Close()
}

func TestNewRequest(t *testing.T) {
	setup()
	defer teardown()

	tt := []struct {
		method string
		uri    string
		debug  bool
		body   interface{}
	}{
		{"GET", "/foo", false, ""},
		{"POST", "/foo/two", true, ""},
		{"PUT", "/abc", true, ""},
	}

	for i, tst := range tt {
		oneDrive.Debug = tst.debug
		req, err := oneDrive.newRequest(tst.method, tst.uri, tst.body)

		if tst.debug {
			if got, want := req.Header.Get("Accept"), "application/json;format=pretty"; got != want {
				t.Fatalf("[%d] Got %q Expected %q", i, got, want)
			}
		} else {
			if got, want := req.Header.Get("Accept"), "application/json"; got != want {
				t.Fatalf("[%d] Got %q Expected %q", i, got, want)
			}
		}

		if err != nil {
			t.Fatalf("[%d] %q", i, err.Error())
		}
		if got, want := req.Method, tst.method; got != want {
			t.Fatalf("[%d] Got %q Expected %q", i, got, want)
		}
		if got, want := req.URL.String(), oneDrive.BaseURL+tst.uri; got != want {
			t.Fatalf("[%d] Got %q Expected %q", i, got, want)
		}
	}
}

func TestDoValid(t *testing.T) {
	setup()
	defer teardown()

	type test1s struct{ A, B string }
	t1 := test1s{
		A: "hello",
		B: "world",
	}

	mux.HandleFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
		b, _ := json.Marshal(&t1)
		w.Write(b)
	})

	tt := []struct {
		method string
		uri    string
		into   interface{}
	}{
		{"GET", "/foo", new(test1s)},
	}

	for i, tst := range tt {
		req, err := oneDrive.newRequest(tst.method, tst.uri, nil)
		if err != nil {
			t.Fatal(err)
		}

		_, err = oneDrive.do(req, tst.into)
		if err != nil {
			t.Fatal(err)
		}

		if got, want := t1.A, tst.into.(*test1s).A; got != want {
			t.Fatalf("[%d] Got %q Expected %q", i, got, want)
		}

		if got, want := t1.B, tst.into.(*test1s).B; got != want {
			t.Fatalf("[%d] Got %q Expected %q", i, got, want)
		}

	}

}
