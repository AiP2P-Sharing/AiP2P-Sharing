package newsdemoknowledge

import (
	"io/fs"
	"net/http"

	newsplugin "aip2p.org/internal/plugins/newsdemo"
)

func newHandler(app *newsplugin.App, staticFS fs.FS) http.Handler {
	return newsplugin.NewTypedModuleHandler(app, staticFS, "/knowledge", "/knowledge/", "/api/knowledge", "/api/knowledge/", "markdown", "knowledge")
}
