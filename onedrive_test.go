package onedrive

import (
	"net/http"
	"net/http/httptest"
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

type testStruct1 struct {
	as, bs, cs string
}
type testStruct2 struct {
	as string
	af float64
	ai int
}

// func TestDoValid(t *testing.T) {
// 	setup()
// 	defer teardown()

// 	type test1s struct{ A, B string }
// 	t1 := test1s{
// 		A: "hello",
// 		B: "world",
// 	}

// 	mux.HandleFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
// 		var b []byte
// 		json.NewDecoder(r.Body).Decode(b)
// 		w.Write(b)
// 	})

// 	tt := []struct {
// 		method     string
// 		uri        string
// 		shouldPass bool
// 		into       interface{}
// 	}{
// 		{"GET", "/foo", true, new(test1s)},
// 		{"GET", "/foo", false, nil},
// 	}

// 	for i, tst := range tt {
// 		req, _ := oneDrive.newRequest(tst.method, tst.uri, tst.into)

// 		_, err := oneDrive.do(req, tst.into)
// 		if err != nil && tst.shouldPass {
// 			t.Fatalf("[%d] %s", i, err.Error())
// 		}

// 		if got, want := t1.A, tst.into.(*test1s).A; got != want && tst.shouldPass {
// 			t.Fatalf("[%d] Got %q Expected %q", i, got, want)
// 		}

// 		if got, want := t1.B, tst.into.(*test1s).B; got != want && tst.shouldPass {
// 			t.Fatalf("[%d] Got %q Expected %q", i, got, want)
// 		}

// 	}

// }
