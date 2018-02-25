// Copyright (c) 2018 Tim Heckman
//
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file at the root of this repository.

package pwnedpasswords

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDefaultURL(t *testing.T) {
	const want = "https://api.pwnedpasswords.com/range/"

	if DefaultURL != want {
		t.Fatalf("DefaultURL = %q, want %q", DefaultURL, want)
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		s string
		i string
		e bool
	}{
		{s: "DefaultURL", i: DefaultURL, e: false},
		{s: "Valid URL", i: "http://localhost/", e: false},
		{s: "Invalid URL", i: "~notAurl", e: true},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.s, func(t *testing.T) {
			var c Client
			var err error

			c, err = New(tt.i)

			if err != nil {
				if tt.e {
					return
				}

				t.Fatalf("unexpected error: %s", err)
			}

			if c.h != tt.i {
				t.Fatalf("Client.h = %q, want %q", c.h, tt.i)
			}

			if c.c == nil {
				t.Fatal("c.c = nil, want allocated *http.Client")
			}
		})
	}
}

func TestHashPassword(t *testing.T) {
	tests := []struct {
		s   string
		i   string
		pre string
		suf string
	}{
		{s: "EmptyString", i: "", pre: "DA39A", suf: "3EE5E6B4B0D3255BFEF95601890AFD80709"},
		{s: "password", i: "password", pre: "5BAA6", suf: "1E4C9B93F3F0682250B6CF8331B7EE68FD8"},
		{s: "test word", i: "test word", pre: "DBE16", suf: "4D0591EEAAC5A33064731DBD53F6A819DCE"},
		{s: "unicode", i: "世界", pre: "CF165", suf: "6101ED511A094D1E4E515BBF8D32B266090"},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.s, func(t *testing.T) {
			var prefix, suffix string

			prefix, suffix = HashPassword([]byte(tt.i))

			if prefix != tt.pre {
				t.Fatalf("prefix = %q, want %q", prefix, tt.pre)
			}

			if suffix != tt.suf {
				t.Fatalf("suffix = %q, want %q", suffix, tt.suf)
			}
		})
	}
}

func Test_readHashes(t *testing.T) {
	var h []hashCount
	var err error

	r := strings.NewReader(testAPIOutput)

	h, err = readHashes(r)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(h) != 3 {
		t.Fatalf("len(h) = %d, want %d", len(h), 3)
	}

	testHashes(t, h)
}

func Test_getHashes(t *testing.T) {
	testServer := setupTestServer()

	defer testServer.Close()

	client, _ := New(testServer.URL + "/range/")

	var h []hashCount
	var err error

	h, err = getHashes(client, "DA39A")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	testHashes(t, h)
}

func Test_ClientCheck(t *testing.T) {
	testServer := setupTestServer()
	defer testServer.Close()

	client, _ := New(testServer.URL + "/range/")

	tests := []struct {
		s string
		i string
		o int
	}{
		{s: "EmptyString", i: "", o: 0},
		{s: "password", i: "password", o: 3303003},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.s, func(t *testing.T) {
			var n int
			var err error

			n, err = client.Check([]byte(tt.i))
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if n != tt.o {
				t.Fatalf("client.Check([]byte(%s)) = %d, want %d", tt.i, n, tt.o)
			}
		})
	}
}

func testHashes(t *testing.T, h []hashCount) {
	res := []hashCount{
		{suffix: "1E4C9B93F3F0682250B6CF8331B7EE68FD8", count: 3303003},
		{suffix: "4D0591EEAAC5A33064731DBD53F6A819DCE", count: 0},
		{suffix: "6101ED511A094D1E4E515BBF8D32B266090", count: 42},
	}

	for i := range res {
		if h[i].suffix != res[i].suffix {
			t.Fatalf("h[%d].suffix = %q, want %q", i, h[i].suffix, res[i].suffix)
		}

		if h[i].count != res[i].count {
			t.Fatalf("h[%d].count = %d, want %d", i, h[i].count, res[i].count)
		}
	}
}

func setupTestServer() *httptest.Server {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("User-Agent"), "go-pwnedpasswords") {
			io.WriteString(w, "ERROR: User-Agent appears invalid (missing go-pwnedpasswords)")
		}
		io.WriteString(w, testAPIOutput)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", http.NotFound)
	mux.HandleFunc("/range/DA39A", fn)
	mux.HandleFunc("/range/5BAA6", fn)

	return httptest.NewServer(mux)
}

var testAPIOutput = `1E4C9B93F3F0682250B6CF8331B7EE68FD8:3303003
4D0591EEAAC5A33064731DBD53F6A819DCE:0
6101ED511A094D1E4E515BBF8D32B266090:42
`
