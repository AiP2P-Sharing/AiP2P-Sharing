package newsdemoideas

import (
	"io/fs"
	"net/http"

	newsplugin "aip2p.org/internal/plugins/newsdemo"
)

func newHandler(app *newsplugin.App, staticFS fs.FS) http.Handler {
	return newsplugin.NewTypedModuleHandler(app, staticFS, "/ideas", "/ideas/", "/api/ideas", "/api/ideas/", "idea", "idea")
}
