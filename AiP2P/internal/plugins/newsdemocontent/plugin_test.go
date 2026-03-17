package newsdemocontent

import (
	"context"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"aip2p.org/internal/apphost"
	newsplugin "aip2p.org/internal/plugins/newsdemo"
	"aip2p.org/internal/themes/newsdemo"
)

func TestPluginBuildServesHomePage(t *testing.T) {
	t.Parallel()

	site := buildContentSite(t)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	site.Handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "AiP2P Sharing") {
		t.Fatalf("expected home page content, got %q", rec.Body.String())
	}
	body := rec.Body.String()
	for _, want := range []string{
		`href="/network"`,
		`href="/writer-policy"`,
		`href="/archive"`,
		`>Overall<`,
		`>Network<`,
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("expected home page to contain %q, got %q", want, body)
		}
	}
	for _, unwanted := range []string{
		"Bundle store",
		"Torrent refs",
		"Sync daemon",
	} {
		if strings.Contains(body, unwanted) {
			t.Fatalf("expected sidebar to hide %q, got %q", unwanted, body)
		}
	}
}

func TestPluginBuildServesFeedAPI(t *testing.T) {
	t.Parallel()

	site := buildContentSite(t)
	req := httptest.NewRequest(http.MethodGet, "/api/feed", nil)
	rec := httptest.NewRecorder()
	site.Handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "\"scope\": \"feed\"") {
		t.Fatalf("expected feed payload, got %q", rec.Body.String())
	}
}

func TestPluginBuildServesIdeasPage(t *testing.T) {
	t.Parallel()

	site := buildContentSite(t)
	req := httptest.NewRequest(http.MethodGet, "/ideas", nil)
	rec := httptest.NewRecorder()
	site.Handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "Ideas") {
		t.Fatalf("expected ideas page content, got %q", rec.Body.String())
	}
}

func TestPluginBuildServesIdeasAPI(t *testing.T) {
	t.Parallel()

	site := buildContentSite(t)
	req := httptest.NewRequest(http.MethodGet, "/api/ideas", nil)
	rec := httptest.NewRecorder()
	site.Handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}
	body := rec.Body.String()
	for _, want := range []string{`"scope": "ideas"`, `"coord_type": "idea"`} {
		if !strings.Contains(body, want) {
			t.Fatalf("expected ideas API payload to contain %q, got %q", want, body)
		}
	}
}

func buildContentSite(t *testing.T) *apphost.Site {
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

func TestPluginBuildCanDisableTypedRoutes(t *testing.T) {
	t.Parallel()

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
		PluginConfig: map[string]any{
			"disabled_coord_types": []any{"task"},
		},
	}
	site, err := Plugin{}.Build(context.Background(), cfg, newsdemo.Theme{})
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}
	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	rec := httptest.NewRecorder()
	site.Handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404; body = %s", rec.Code, rec.Body.String())
	}
}

func TestWithRawQuery(t *testing.T) {
	t.Parallel()

	if got := newsplugin.AppendRawQuery("/ideas/idea-1", "agent=1"); got != "/ideas/idea-1?agent=1" {
		t.Fatalf("withRawQuery() = %q", got)
	}
	if got := newsplugin.AppendRawQuery("/ideas/idea-1?view=compact", "agent=1"); got != "/ideas/idea-1?view=compact&agent=1" {
		t.Fatalf("withRawQuery() = %q", got)
	}
	if got := newsplugin.AppendRawQuery("/ideas/idea-1", ""); got != "/ideas/idea-1" {
		t.Fatalf("withRawQuery() = %q", got)
	}
}
