package handlers

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/a-h/templ"
)

type failComponent struct{}

func (failComponent) Render(ctx context.Context, w io.Writer) error {
	_, _ = w.Write([]byte("partial"))
	return errors.New("render boom")
}

func TestRender_SuccessWritesStatusAndBody(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	Render(rr, req, templ.Raw("<p>ok</p>"), http.StatusOK)
	if rr.Code != http.StatusOK {
		t.Fatalf("status %d", rr.Code)
	}
	if ct := rr.Header().Get("Content-Type"); !strings.HasPrefix(ct, "text/html") {
		t.Fatalf("Content-Type %q", ct)
	}
	if body := rr.Body.String(); body != "<p>ok</p>" {
		t.Fatalf("body %q", body)
	}
}

func TestRender_FailurePreservesIntendedStatus(t *testing.T) {
	// Locked behavior: intended status is committed before render; on failure
	// http.Error runs but cannot replace an already-written status.
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	Render(rr, req, failComponent{}, http.StatusAccepted)

	if rr.Code != http.StatusAccepted {
		t.Fatalf("intended status must be preserved, got %d", rr.Code)
	}
	body := rr.Body.String()
	if !strings.Contains(body, "render boom") {
		t.Fatalf("want error helper text in body, got %q", body)
	}
	// partial write from the component may precede the error helper message
	if !strings.Contains(body, "partial") {
		t.Fatalf("want partial component output before error path, got %q", body)
	}
}
