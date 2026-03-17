package newsdemoknowledge

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

func TestPluginBuildServesKnowledgePage(t *testing.T) {
	t.Parallel()
	site := buildKnowledgeSite(t)
	req := httptest.NewRequest(http.MethodGet, "/knowledge", nil)
	rec := httptest.NewRecorder()
	site.Handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "Knowledge") {
		t.Fatalf("expected knowledge page content, got %q", rec.Body.String())
	}
}

func TestPluginBuildServesKnowledgeAPI(t *testing.T) {
	t.Parallel()
	site := buildKnowledgeSite(t)
	req := httptest.NewRequest(http.MethodGet, "/api/knowledge", nil)
	rec := httptest.NewRecorder()
	site.Handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}
	body := rec.Body.String()
	for _, want := range []string{`"scope": "knowledge"`, `"coord_type": "markdown"`} {
		if !strings.Contains(body, want) {
			t.Fatalf("expected knowledge API payload to contain %q, got %q", want, body)
		}
	}
}

func buildKnowledgeSite(t *testing.T) *apphost.Site {
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
