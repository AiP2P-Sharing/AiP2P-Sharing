package newsdemocontent

import (
	"context"
	_ "embed"
	"strings"

	"aip2p.org/internal/apphost"
	newsplugin "aip2p.org/internal/plugins/newsdemo"
)

type Plugin struct{}

//go:embed aip2p.plugin.json
var pluginManifestJSON []byte

func (Plugin) Manifest() apphost.PluginManifest {
	return apphost.MustLoadPluginManifestJSON(pluginManifestJSON)
}

func (Plugin) Build(_ context.Context, cfg apphost.Config, theme apphost.WebTheme) (*apphost.Site, error) {
	cfg = newsplugin.ApplyDefaultConfig(cfg)
	app, err := newsplugin.NewWithThemeAndOptions(
		cfg.StoreRoot,
		cfg.Project,
		cfg.Version,
		cfg.ArchiveRoot,
		cfg.RulesPath,
		cfg.WriterPolicyPath,
		cfg.NetPath,
		theme,
		newsplugin.ContentOnlyAppOptions(),
	)
	if err != nil {
		return nil, err
	}
	staticFS, err := theme.StaticFS()
	if err != nil {
		return nil, err
	}
	disabled := disabledCoordTypes(cfg.PluginConfig)
	return &apphost.Site{
		Manifest: Plugin{}.Manifest(),
		Theme:    theme.Manifest(),
		Handler:  newHandler(app, staticFS, disabled),
	}, nil
}

func disabledCoordTypes(cfg map[string]any) map[string]struct{} {
	disabled := map[string]struct{}{}
	if len(cfg) == 0 {
		return disabled
	}
	value, ok := cfg["disabled_coord_types"]
	if !ok {
		return disabled
	}
	for _, item := range sliceStrings(value) {
		item = strings.ToLower(strings.TrimSpace(item))
		if item == "" {
			continue
		}
		disabled[item] = struct{}{}
	}
	return disabled
}

func sliceStrings(value any) []string {
	switch typed := value.(type) {
	case []string:
		return typed
	case []any:
		out := make([]string, 0, len(typed))
		for _, item := range typed {
			text, ok := item.(string)
			if ok {
				out = append(out, text)
			}
		}
		return out
	default:
		return nil
	}
}
