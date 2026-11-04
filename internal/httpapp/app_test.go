package httpapp

import (
	"compress/gzip"
	"io"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"
)

func testFS() (staticFS, mediaFS fs.FS) {
	staticFS = fstest.MapFS{
		"hello.txt": &fstest.MapFile{Data: []byte("static-hello")},
		"app.js":    &fstest.MapFile{Data: []byte("console.log(1)")},
	}
	mediaFS = fstest.MapFS{
		"pic.txt": &fstest.MapFile{Data: []byte("media-pic")},
	}
	return staticFS, mediaFS
}

func newTestHandler(origins ...string) http.Handler {
	s, m := testFS()
	var o []string
	if len(origins) > 0 {
		o = origins
	}
	return New(s, m, o)
}

func TestSeamCallable(t *testing.T) {
	h := newTestHandler()
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/nope", nil))
	if rr.Code != http.StatusNotFound {
		t.Fatalf("status %d", rr.Code)
	}
}

func TestNew_RejectsWildcardOrigins(t *testing.T) {
	s, m := testFS()
	for _, origins := range [][]string{
		{"*"},
		{"https://app.example.com", "*"},
		{"https://*.example.com"},
		{""},
		{"not-a-url"},
	} {
		func() {
			defer func() {
				if recover() == nil {
					t.Fatalf("New(%v): want panic", origins)
				}
			}()
			_ = New(s, m, origins)
		}()
	}
}

func TestHomeOK(t *testing.T) {
	h := newTestHandler()
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/", nil))
	if rr.Code != http.StatusOK {
		t.Fatalf("status %d body %q", rr.Code, rr.Body.String())
	}
	body := rr.Body.String()
	if !strings.Contains(body, "progressive web app") {
		t.Fatalf("home body missing marker: %q", body)
	}
	if !strings.Contains(body, "Welcome to the Site") {
		t.Fatalf("home body missing title: %q", body)
	}
	if ct := rr.Header().Get("Content-Type"); !strings.HasPrefix(ct, "text/html") {
		t.Fatalf("Content-Type %q", ct)
	}
}

func TestStaticAndMedia(t *testing.T) {
	h := newTestHandler()

	t.Run("static hit", func(t *testing.T) {
		rr := httptest.NewRecorder()
		h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/static/hello.txt", nil))
		if rr.Code != http.StatusOK || rr.Body.String() != "static-hello" {
			t.Fatalf("status %d body %q", rr.Code, rr.Body.String())
		}
	})

	t.Run("static js content-type", func(t *testing.T) {
		rr := httptest.NewRecorder()
		h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/static/app.js", nil))
		if rr.Code != http.StatusOK {
			t.Fatalf("status %d", rr.Code)
		}
		ct := rr.Header().Get("Content-Type")
		if ct != "application/javascript" {
			t.Fatalf("Content-Type %q", ct)
		}
	})

	t.Run("static miss", func(t *testing.T) {
		rr := httptest.NewRecorder()
		h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/static/missing.txt", nil))
		if rr.Code != http.StatusNotFound {
			t.Fatalf("status %d", rr.Code)
		}
	})

	t.Run("media hit", func(t *testing.T) {
		rr := httptest.NewRecorder()
		h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/media/pic.txt", nil))
		if rr.Code != http.StatusOK || rr.Body.String() != "media-pic" {
			t.Fatalf("status %d body %q", rr.Code, rr.Body.String())
		}
	})

	t.Run("media miss", func(t *testing.T) {
		rr := httptest.NewRecorder()
		h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/media/nope.txt", nil))
		if rr.Code != http.StatusNotFound {
			t.Fatalf("status %d", rr.Code)
		}
	})

	t.Run("prefix isolation", func(t *testing.T) {
		rr := httptest.NewRecorder()
		h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/static/pic.txt", nil))
		if rr.Code != http.StatusNotFound {
			t.Fatalf("static must not serve media file, status %d", rr.Code)
		}
		rr = httptest.NewRecorder()
		h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/media/hello.txt", nil))
		if rr.Code != http.StatusNotFound {
			t.Fatalf("media must not serve static file, status %d", rr.Code)
		}
	})

	t.Run("empty media MapFS miss", func(t *testing.T) {
		s, _ := testFS()
		h := New(s, fstest.MapFS{}, nil)
		rr := httptest.NewRecorder()
		h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/media/any.txt", nil))
		if rr.Code != http.StatusNotFound {
			t.Fatalf("status %d", rr.Code)
		}
	})

	t.Run("empty host media dir miss", func(t *testing.T) {
		s, _ := testFS()
		dir := t.TempDir()
		h := New(s, os.DirFS(dir), nil)
		rr := httptest.NewRecorder()
		h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/media/missing.bin", nil))
		if rr.Code != http.StatusNotFound {
			t.Fatalf("status %d", rr.Code)
		}
	})

	t.Run("missing media root path miss", func(t *testing.T) {
		s, _ := testFS()
		// DirFS of a path that does not exist: Open fails → FileServer 404 class.
		missing := filepath.Join(t.TempDir(), "no-such-media-root")
		h := New(s, os.DirFS(missing), nil)
		rr := httptest.NewRecorder()
		h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/media/x.txt", nil))
		if rr.Code != http.StatusNotFound {
			t.Fatalf("status %d", rr.Code)
		}
	})
}

func TestCORS(t *testing.T) {
	const permitted = "https://app.site.com"
	const other = "https://evil.example.com"

	t.Run("empty origins no CORS", func(t *testing.T) {
		h := newTestHandler()
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Origin", permitted)
		h.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Fatalf("status %d", rr.Code)
		}
		if acao := rr.Header().Get("Access-Control-Allow-Origin"); acao != "" {
			t.Fatalf("unexpected ACAO %q", acao)
		}
	})

	t.Run("permitted origin on app", func(t *testing.T) {
		h := newTestHandler(permitted, "https://hda.site.com")
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Origin", permitted)
		h.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Fatalf("status %d", rr.Code)
		}
		if acao := rr.Header().Get("Access-Control-Allow-Origin"); acao != permitted {
			t.Fatalf("ACAO %q", acao)
		}
	})

	t.Run("non-permitted origin on app", func(t *testing.T) {
		h := newTestHandler(permitted)
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Origin", other)
		h.ServeHTTP(rr, req)
		if acao := rr.Header().Get("Access-Control-Allow-Origin"); acao != "" {
			t.Fatalf("unexpected ACAO %q", acao)
		}
	})

	t.Run("no Origin header still OK", func(t *testing.T) {
		h := newTestHandler(permitted)
		rr := httptest.NewRecorder()
		h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/", nil))
		if rr.Code != http.StatusOK {
			t.Fatalf("status %d", rr.Code)
		}
	})

	t.Run("static and media unaffected by CORS", func(t *testing.T) {
		h := newTestHandler(permitted)
		for _, path := range []string{"/static/hello.txt", "/media/pic.txt"} {
			rr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, path, nil)
			req.Header.Set("Origin", permitted)
			h.ServeHTTP(rr, req)
			if rr.Code != http.StatusOK {
				t.Fatalf("%s status %d", path, rr.Code)
			}
			if acao := rr.Header().Get("Access-Control-Allow-Origin"); acao != "" {
				t.Fatalf("%s unexpected ACAO %q", path, acao)
			}
		}
	})

	t.Run("OPTIONS preflight permitted", func(t *testing.T) {
		h := newTestHandler(permitted)
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodOptions, "/", nil)
		req.Header.Set("Origin", permitted)
		req.Header.Set("Access-Control-Request-Method", "GET")
		req.Header.Set("Access-Control-Request-Headers", "hx-request,content-type")
		h.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Fatalf("preflight status %d", rr.Code)
		}
		if acao := rr.Header().Get("Access-Control-Allow-Origin"); acao != permitted {
			t.Fatalf("ACAO %q", acao)
		}
		allowMeth := rr.Header().Get("Access-Control-Allow-Methods")
		if !strings.Contains(strings.ToUpper(allowMeth), "GET") {
			t.Fatalf("Allow-Methods %q", allowMeth)
		}
		allowHdr := strings.ToLower(rr.Header().Get("Access-Control-Allow-Headers"))
		if !strings.Contains(allowHdr, "hx-request") || !strings.Contains(allowHdr, "content-type") {
			t.Fatalf("Allow-Headers %q", allowHdr)
		}
	})

	t.Run("OPTIONS preflight rejected origin", func(t *testing.T) {
		h := newTestHandler(permitted)
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodOptions, "/", nil)
		req.Header.Set("Origin", other)
		req.Header.Set("Access-Control-Request-Method", "GET")
		h.ServeHTTP(rr, req)
		if acao := rr.Header().Get("Access-Control-Allow-Origin"); acao != "" {
			t.Fatalf("unexpected ACAO %q", acao)
		}
		if am := rr.Header().Get("Access-Control-Allow-Methods"); am != "" {
			t.Fatalf("unexpected Allow-Methods %q", am)
		}
	})

	t.Run("static and media preflight have no CORS", func(t *testing.T) {
		h := newTestHandler(permitted)
		for _, path := range []string{"/static/hello.txt", "/media/pic.txt"} {
			rr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodOptions, path, nil)
			req.Header.Set("Origin", permitted)
			req.Header.Set("Access-Control-Request-Method", "GET")
			h.ServeHTTP(rr, req)
			if acao := rr.Header().Get("Access-Control-Allow-Origin"); acao != "" {
				t.Fatalf("%s unexpected ACAO %q", path, acao)
			}
			if am := rr.Header().Get("Access-Control-Allow-Methods"); am != "" {
				t.Fatalf("%s unexpected Allow-Methods %q", path, am)
			}
		}
	})
}

func TestRecoverer(t *testing.T) {
	// Same-package constructor yields the mux without requiring the public
	// New seam to expose *chi.Mux.
	s, m := testFS()
	mux := newApp(s, m, nil)
	mux.Get("/__panic", func(http.ResponseWriter, *http.Request) {
		panic("boom")
	})

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/__panic", nil))
	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("want recovered 500, got %d", rr.Code)
	}
}

func TestCompressionGzip(t *testing.T) {
	const permitted = "https://app.site.com"
	h := newTestHandler(permitted)

	// Static .js sets Content-Type before write, so Chi Compress can gzip it.
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/static/app.js", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status %d", rr.Code)
	}
	if ce := rr.Header().Get("Content-Encoding"); ce != "gzip" {
		t.Fatalf("Content-Encoding %q, want gzip", ce)
	}
	gr, err := gzip.NewReader(rr.Body)
	if err != nil {
		t.Fatal(err)
	}
	body, err := io.ReadAll(gr)
	if err != nil {
		t.Fatal(err)
	}
	if err := gr.Close(); err != nil {
		t.Fatal(err)
	}
	if string(body) != "console.log(1)" {
		t.Fatalf("body %q", body)
	}

	// Home HTML with permitted Origin + gzip: 200, exact ACAO, gzip body.
	rr = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Origin", permitted)
	req.Header.Set("Accept-Encoding", "gzip")
	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("home gzip status %d", rr.Code)
	}
	if acao := rr.Header().Get("Access-Control-Allow-Origin"); acao != permitted {
		t.Fatalf("home gzip ACAO %q", acao)
	}
	if ce := rr.Header().Get("Content-Encoding"); ce != "gzip" {
		t.Fatalf("home Content-Encoding %q, want gzip", ce)
	}
	gr, err = gzip.NewReader(rr.Body)
	if err != nil {
		t.Fatal(err)
	}
	homeBody, err := io.ReadAll(gr)
	if err != nil {
		t.Fatal(err)
	}
	if err := gr.Close(); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(homeBody), "progressive web app") {
		t.Fatalf("gzip home missing marker: %q", homeBody)
	}

	// Brotli-only client: must not require br success.
	rr = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept-Encoding", "br")
	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status %d", rr.Code)
	}
	if ce := rr.Header().Get("Content-Encoding"); ce == "br" {
		t.Fatalf("Brotli must not be configured, got Content-Encoding=br")
	}
}
