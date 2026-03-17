package builtin

import (
	_ "embed"
	"fmt"
	"strings"

	"aip2p.org/internal/apphost"
	newsdemoagents "aip2p.org/internal/plugins/newsdemoagents"
	newsdemoarchive "aip2p.org/internal/plugins/newsdemoarchive"
	newsdemocode "aip2p.org/internal/plugins/newsdemocode"
	newsdemocontent "aip2p.org/internal/plugins/newsdemocontent"
	newsdemogovernance "aip2p.org/internal/plugins/newsdemogovernance"
	newsdemoideas "aip2p.org/internal/plugins/newsdemoideas"
	newsdemoknowledge "aip2p.org/internal/plugins/newsdemoknowledge"
	newsdemoops "aip2p.org/internal/plugins/newsdemoops"
	newsdemoskills "aip2p.org/internal/plugins/newsdemoskills"
	newsdemotasks "aip2p.org/internal/plugins/newsdemotasks"
	"aip2p.org/internal/themes/newsdemo"
)

//go:embed news-demo.app.json
var newsDemoAppJSON []byte

func DefaultRegistry() *apphost.Registry {
	registry := apphost.NewRegistry()
	registry.MustRegisterTheme(newsdemo.Theme{})
	registry.MustRegisterPlugin(newsdemocontent.Plugin{})
	registry.MustRegisterPlugin(newsdemoideas.Plugin{})
	registry.MustRegisterPlugin(newsdemoknowledge.Plugin{})
	registry.MustRegisterPlugin(newsdemocode.Plugin{})
	registry.MustRegisterPlugin(newsdemoagents.Plugin{})
	registry.MustRegisterPlugin(newsdemoskills.Plugin{})
	registry.MustRegisterPlugin(newsdemotasks.Plugin{})
	registry.MustRegisterPlugin(newsdemoarchive.Plugin{})
	registry.MustRegisterPlugin(newsdemogovernance.Plugin{})
	registry.MustRegisterPlugin(newsdemoops.Plugin{})
	return registry
}

func DefaultApps() []apphost.AppManifest {
	return []apphost.AppManifest{
		apphost.MustLoadAppManifestJSON(newsDemoAppJSON),
	}
}

func ResolveApp(id string) (apphost.AppManifest, error) {
	id = strings.ToLower(strings.TrimSpace(id))
	for _, app := range DefaultApps() {
		if strings.ToLower(strings.TrimSpace(app.ID)) == id {
			return app, nil
		}
	}
	return apphost.AppManifest{}, fmt.Errorf("app %q not found", id)
}
