package newsdemoagents

import (
	"context"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"aip2p.org/internal/apphost"
	"aip2p.org/internal/themes/newsdemo"
)

func TestPluginBuildServesAgentsPage(t *testing.T) {
	t.Parallel()
	site := buildAgentsSite(t)
	req := httptest.NewRequest(http.MethodGet, "/agents", nil)
	rec := httptest.NewRecorder()
	site.Handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "Agents") {
		t.Fatalf("expected agents page content, got %q", rec.Body.String())
	}
}

func TestPluginBuildServesAgentsAPI(t *testing.T) {
	t.Parallel()
	site := buildAgentsSite(t)
	req := httptest.NewRequest(http.MethodGet, "/api/agents", nil)
	rec := httptest.NewRecorder()
	site.Handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}
	body := rec.Body.String()
	for _, want := range []string{`"scope": "agents"`, `"coord_type": "agent"`} {
		if !strings.Contains(body, want) {
			t.Fatalf("expected agents API payload to contain %q, got %q", want, body)
		}
	}
}

func buildAgentsSite(t *testing.T) *apphost.Site {
	t.Helper()
	root := t.TempDir()
	cfg := apphost.Config{
		RuntimeRoot:      filepath.Join(root, "runtime"),
		StoreRoot:        filepath.Join(root, "store"),
		ArchiveRoot:      filepath.Join(root, "archive"),
		RulesPath:        filepath.Join(root, "config", "subscriptions.json"),
		WriterPolicyPath: filepath.Join(root, "config", "writer_policy.json"),
		NetPath:          filepath.Join(root, "config", "aip2p_net.inf"),
		Project:          "aip2p.sharing",
		Version:          "test",
	}
	site, err := Plugin{}.Build(context.Background(), cfg, newsdemo.Theme{})
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}
	return site
}
