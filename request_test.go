package onedrive

import (
	"testing"
	"time"
)

func TestCalculateThrottle(t *testing.T) {

	time1, _ := time.Parse("02 Jan 06 15:04 MST", "01 Jan 15 15:04 GMT")
	time1e, _ := time.Parse("02 Jan 06 15:04 MST", "01 Jan 15 16:04 GMT")
	time2, _ := time.Parse("02 Jan 06 15:04 MST", "01 Jan 15 15:04 GMT")
	time2e, _ := time.Parse("02 Jan 06 15:04 MST", "01 Jan 15 15:05 GMT")

	tt := []struct {
		currentTime  time.Time
		expectedTime time.Time
		retryAfter   string
	}{
		{time1, time1e, "3600"},
		{time2, time2e, "60"},
	}

	for i, tst := range tt {
		tm, err := calculateThrottle(tst.currentTime, tst.retryAfter)
		if err != nil {
			t.Fatalf("[%d] Couldn't calculate retry after: %s", i, err.Error())
		}
		if got, want := tm, tst.expectedTime; !got.Equal(want) {
			t.Fatalf("[%d] Got %s Expected %s", i, got, want)
		}
	}
}

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

// func TestDoValid(t *testing.T) {
// 	setup()
// 	defer teardown()

// 	t1 := testStruct1{"a", "b", "c"}

// 	tt := []struct {
// 		method   string
// 		uri      string
// 		into     interface{}
// 		intoType interface{}
// 	}{
// 		{"GET", "/foo", t1, &testStruct1{}},
// 		{"GET", "/foo", nil, &testStruct1{}},
// 	}

// 	for i, tst := range tt {
// 		req, err := oneDrive.newRequest(tst.method, tst.uri, nil, tst.into)
// 		if err != nil {
// 			t.Fatalf("[%d] %s", i, err.Error())
// 		}

// 		resp, err := oneDrive.do(req, &tst.intoType)
// 		if err != nil {
// 			t.Fatalf("[%d] %s", i, resp, err.Error())
// 		}

// 		fmt.Println(tst.intoType)

// 	}

// }
