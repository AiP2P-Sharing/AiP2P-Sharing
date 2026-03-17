package newsdemocode

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

func TestPluginBuildServesCodePage(t *testing.T) {
	t.Parallel()
	site := buildCodeSite(t)
	req := httptest.NewRequest(http.MethodGet, "/code", nil)
	rec := httptest.NewRecorder()
	site.Handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "Code") {
		t.Fatalf("expected code page content, got %q", rec.Body.String())
	}
}

func TestPluginBuildServesCodeAPI(t *testing.T) {
	t.Parallel()
	site := buildCodeSite(t)
	req := httptest.NewRequest(http.MethodGet, "/api/code", nil)
	rec := httptest.NewRecorder()
	site.Handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}
	body := rec.Body.String()
	for _, want := range []string{`"scope": "code"`, `"coord_type": "code"`} {
		if !strings.Contains(body, want) {
			t.Fatalf("expected code API payload to contain %q, got %q", want, body)
		}
	}
}

func buildCodeSite(t *testing.T) *apphost.Site {
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
