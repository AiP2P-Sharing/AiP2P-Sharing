package newsdemotasks

import (
	"io/fs"
	"net/http"

	newsplugin "aip2p.org/internal/plugins/newsdemo"
)

func newHandler(app *newsplugin.App, staticFS fs.FS) http.Handler {
	return newsplugin.NewTypedModuleHandler(app, staticFS, "/tasks", "/tasks/", "/api/tasks", "/api/tasks/", "task", "task")
}
