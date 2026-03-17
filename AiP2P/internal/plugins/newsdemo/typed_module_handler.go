package newsplugin

import (
	"io/fs"
	"net/http"
	"strings"
)

func NewTypedModuleHandler(app *App, staticFS fs.FS, path, itemPrefix, apiPath, apiItemPrefix, coordType, scope string) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		handleTypedCollection(app, w, r, path)
	})
	mux.HandleFunc(itemPrefix, func(w http.ResponseWriter, r *http.Request) {
		handleTypedPost(app, w, r, itemPrefix, path, coordType)
	})
	mux.HandleFunc(apiPath, func(w http.ResponseWriter, r *http.Request) {
		handleAPITypedCollection(app, w, r, apiPath)
	})
	mux.HandleFunc(apiItemPrefix, func(w http.ResponseWriter, r *http.Request) {
		handleAPITypedPost(app, w, r, apiItemPrefix, coordType, scope)
	})
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))
	return mux
}

func handleTypedCollection(app *App, w http.ResponseWriter, r *http.Request, path string) {
	if r.URL.Path != path {
		http.NotFound(w, r)
		return
	}
	spec, ok := CoordCollectionByPath(path)
	if !ok {
		http.NotFound(w, r)
		return
	}
	index, err := app.Index()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	opts := FeedOptionsFromRequest(r)
	data := BuildTypedCollectionPageData(app, index, spec, opts)
	if err := app.Templates().ExecuteTemplate(w, "typed_collection.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleTypedPost(app *App, w http.ResponseWriter, r *http.Request, prefix, navPath, coordType string) {
	infoHash := PathValue(prefix, r.URL.Path)
	if infoHash == "" {
		http.NotFound(w, r)
		return
	}
	index, err := app.Index()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	post, ok := index.PostByInfoHash[strings.ToLower(infoHash)]
	if !ok || PostCoordType(post) != coordType {
		http.NotFound(w, r)
		return
	}
	data := BuildPostPageData(app, index, navPath, post, coordType)
	if err := app.Templates().ExecuteTemplate(w, "post.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleAPITypedCollection(app *App, w http.ResponseWriter, r *http.Request, apiPath string) {
	if r.URL.Path != apiPath {
		http.NotFound(w, r)
		return
	}
	spec, ok := CoordCollectionByAPIPath(apiPath)
	if !ok {
		http.NotFound(w, r)
		return
	}
	index, err := app.Index()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	opts := FeedOptionsFromRequest(r)
	WriteJSON(w, http.StatusOK, BuildTypedCollectionAPIResponse(app, index, spec, opts, strings.TrimPrefix(spec.Path, "/")))
}

func handleAPITypedPost(app *App, w http.ResponseWriter, r *http.Request, prefix, coordType, scope string) {
	infoHash := PathValue(prefix, r.URL.Path)
	if infoHash == "" {
		http.NotFound(w, r)
		return
	}
	index, err := app.Index()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	post, ok := index.PostByInfoHash[strings.ToLower(infoHash)]
	if !ok || PostCoordType(post) != coordType {
		http.NotFound(w, r)
		return
	}
	WriteJSON(w, http.StatusOK, BuildPostAPIResponse(app, index, scope, post, coordType))
}
