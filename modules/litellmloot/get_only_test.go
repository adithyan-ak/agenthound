package litellmloot

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/adithyan-ak/agenthound/sdk/action"
)

// TestLoot_GETOnly is a hard regression guard on the read-only contract
// in sdk/action.Looter. The looter MUST issue ONLY GET (and HEAD)
// requests — no POST, PUT, PATCH, DELETE, etc. — because mutating
// methods would leave evidence in the LiteLLM gateway's audit log
// and (worse) could change upstream provider state.
//
// This test installs an http.Handler that records every method and
// fails the test if anything other than GET appears.
func TestLoot_GETOnly(t *testing.T) {
	var (
		mu      sync.Mutex
		methods []string
	)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		methods = append(methods, r.Method)
		mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/model/info":
			_, _ = w.Write([]byte(happyPathModelInfo))
		case "/key/list":
			_, _ = w.Write([]byte(happyPathKeyList))
		default:
			w.WriteHeader(404)
		}
	}))
	defer srv.Close()

	l := &Looter{}
	_, err := l.Loot(context.Background(), action.Target{
		Address: strings.TrimPrefix(srv.URL, "http://"),
	}, action.LootOptions{
		Credentials: map[string]string{"master_key": fakeMasterKey},
	})
	if err != nil {
		t.Fatalf("Loot: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()
	if len(methods) == 0 {
		t.Fatal("looter issued zero requests")
	}
	for _, m := range methods {
		if m != "GET" && m != "HEAD" {
			t.Errorf("looter issued non-read-only method %q (Looter contract violation; see sdk/action/looter.go doc)", m)
		}
	}
}
