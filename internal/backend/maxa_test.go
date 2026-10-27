package backend

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/ZenNotes/tui/internal/remote"
)

func TestMaxaIncrementalWritesRetryLatest(t *testing.T) {
	for _, op := range []string{"append", "prepend"} {
		t.Run(op, func(t *testing.T) {
			body := "# Original\n"
			writes := 0
			var mu sync.Mutex
			version := func() string { return fmt.Sprintf("%x", sha256.Sum256([]byte(body))) }
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				mu.Lock()
				defer mu.Unlock()
				w.Header().Set("Content-Type", "application/json")
				switch r.URL.Path {
				case "/api/notes/read":
					_ = json.NewEncoder(w).Encode(map[string]any{"path": "note.md", "body": body, "version": version()})
				case "/api/notes/write":
					var req struct {
						Body, BaseVersion string
						Force             bool
					}
					if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
						t.Error(err)
					}
					writes++
					if writes == 1 {
						body += "Agent external change\n"
						w.WriteHeader(409)
						_ = json.NewEncoder(w).Encode(map[string]string{"code": "VERSION_CONFLICT", "version": version(), "body": body})
						return
					}
					if req.BaseVersion != version() || req.Force {
						t.Error("write did not use latest version", req)
					}
					body = req.Body
					_ = json.NewEncoder(w).Encode(map[string]string{"path": "note.md", "version": version()})
				default:
					http.NotFound(w, r)
				}
			}))
			defer server.Close()
			r := &Remote{client: remote.NewClient(server.URL, "")}
			var err error
			if op == "append" {
				_, err = r.AppendToNote(context.Background(), "note.md", "CLI change")
			} else {
				_, err = r.PrependToNote(context.Background(), "note.md", "CLI change")
			}
			if err != nil {
				t.Fatal(err)
			}
			if writes != 2 || !strings.Contains(body, "Agent external change") || !strings.Contains(body, "CLI change") {
				t.Fatal(writes, body)
			}
		})
	}
}

func TestMaxaPersistentConflictFailsWithoutForce(t *testing.T) {
	writes := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			_, _ = w.Write([]byte(`{"path":"note.md","body":"latest","version":"version"}`))
			return
		}
		writes++
		w.WriteHeader(409)
	}))
	defer server.Close()
	r := &Remote{client: remote.NewClient(server.URL, "")}
	_, err := r.AppendToNote(context.Background(), "note.md", "change")
	if err == nil || writes != 5 {
		t.Fatal(writes, err)
	}
}
