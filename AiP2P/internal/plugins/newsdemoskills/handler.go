package newsdemoskills

import (
	"io/fs"
	"net/http"

	newsplugin "aip2p.org/internal/plugins/newsdemo"
)

func newHandler(app *newsplugin.App, staticFS fs.FS) http.Handler {
	return newsplugin.NewTypedModuleHandler(app, staticFS, "/skills", "/skills/", "/api/skills", "/api/skills/", "skill", "skill")
}
