package onedrive

import "testing"

func TestNewRequest(t *testing.T) {
	setup()
	defer teardown()

	tt := []struct {
		method         string
		uri            string
		requestHeaders map[string]string
		debug          bool
		body           interface{}
	}{
		{"GET", "/foo", map[string]string{"Content-Type": "text/plain"}, false, nil},
		{"POST", "/foo/two", nil, true, ""},
		{"DELETE", "/foo/two", nil, true, testStruct1{"a", "b", "c"}},
	}

	for i, tst := range tt {
		oneDrive.Debug = tst.debug

		req, err := oneDrive.newRequest(tst.method, tst.uri, tst.requestHeaders, tst.body)
		if err != nil {
			t.Fatal(err)
		}

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

		for k, v := range tst.requestHeaders {
			if got, want := req.Header.Get(k), v; got != want {
				t.Fatalf("[%d] Got %q Expected %q", i, got, want)
			}
		}
	}
}
