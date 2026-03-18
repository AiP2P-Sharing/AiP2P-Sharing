package newsdemocontent

import (
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	newsplugin "aip2p.org/internal/plugins/newsdemo"
)

func newHandler(app *newsplugin.App, staticFS fs.FS, disabledCoordTypes map[string]struct{}) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handleHome(app, w, r)
	})
	mux.HandleFunc("/posts/", func(w http.ResponseWriter, r *http.Request) {
		handlePost(app, w, r)
	})
	mux.HandleFunc("/sources", func(w http.ResponseWriter, r *http.Request) {
		handleSources(app, w, r)
	})
	mux.HandleFunc("/sources/", func(w http.ResponseWriter, r *http.Request) {
		handleSource(app, w, r)
	})
	mux.HandleFunc("/topics", func(w http.ResponseWriter, r *http.Request) {
		handleTopics(app, w, r)
	})
	mux.HandleFunc("/topics/", func(w http.ResponseWriter, r *http.Request) {
		handleTopic(app, w, r)
	})
	for _, spec := range newsplugin.CoordCollections() {
		if _, disabled := disabledCoordTypes[strings.ToLower(strings.TrimSpace(spec.CoordType))]; disabled {
			continue
		}
		spec := spec
		mux.HandleFunc(spec.Path, func(w http.ResponseWriter, r *http.Request) {
			handleCoordCollection(app, w, r, spec)
		})
		mux.HandleFunc(spec.APIPath, func(w http.ResponseWriter, r *http.Request) {
			handleAPICoordCollection(app, w, r, spec)
		})
	}
	mux.HandleFunc("/api/feed", func(w http.ResponseWriter, r *http.Request) {
		handleAPIFeed(app, w, r)
	})
	mux.HandleFunc("/api/posts/", func(w http.ResponseWriter, r *http.Request) {
		handleAPIPost(app, w, r)
	})
	mux.HandleFunc("/api/torrents/", func(w http.ResponseWriter, r *http.Request) {
		handleAPITorrent(app, w, r)
	})
	mux.HandleFunc("/api/sources", func(w http.ResponseWriter, r *http.Request) {
		handleAPISources(app, w, r)
	})
	mux.HandleFunc("/api/sources/", func(w http.ResponseWriter, r *http.Request) {
		handleAPISource(app, w, r)
	})
	mux.HandleFunc("/api/topics", func(w http.ResponseWriter, r *http.Request) {
		handleAPITopics(app, w, r)
	})
	mux.HandleFunc("/api/topics/", func(w http.ResponseWriter, r *http.Request) {
		handleAPITopic(app, w, r)
	})
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))
	return mux
}

func handleHome(app *newsplugin.App, w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	index, err := app.Index()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rules, err := app.SubscriptionRules()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	opts := readFeedOptions(r)
	showNetworkWarn := newsplugin.ShouldShowNetworkWarning(r)
	if showNetworkWarn {
		http.SetCookie(w, &http.Cookie{
			Name:     "aip2p_news_network_warning_seen",
			Value:    "1",
			Path:     "/",
			MaxAge:   180 * 24 * 60 * 60,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		})
	}
	data := newsplugin.BuildHomePageData(app, index, rules, opts, newsplugin.IsAgentViewer(r), showNetworkWarn)
	if err := app.Templates().ExecuteTemplate(w, "home.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handlePost(app *newsplugin.App, w http.ResponseWriter, r *http.Request) {
	infoHash := newsplugin.PathValue("/posts/", r.URL.Path)
	if infoHash == "" {
		http.NotFound(w, r)
		return
	}
	index, ok := loadIndexOrError(app, w)
	if !ok {
		return
	}
	post, ok := index.PostByInfoHash[strings.ToLower(infoHash)]
	if !ok {
		http.NotFound(w, r)
		return
	}
	if canonical := newsplugin.CoordCanonicalPath(post); canonical != "" {
		http.Redirect(w, r, newsplugin.AppendRawQuery(canonical, r.URL.RawQuery), http.StatusTemporaryRedirect)
		return
	}
	data := newsplugin.BuildPostPageData(app, index, "/", post, "")
	data.AgentView = newsplugin.IsAgentViewer(r)
	renderTemplate(w, app, "post.html", data)
}

func handleSources(app *newsplugin.App, w http.ResponseWriter, r *http.Request) {
	handleDirectoryPage(app, w, r, "/sources", "sources")
}

func handleSource(app *newsplugin.App, w http.ResponseWriter, r *http.Request) {
	handleScopedCollectionPage(app, w, r, "/sources/", "source")
}

func handleTopics(app *newsplugin.App, w http.ResponseWriter, r *http.Request) {
	handleDirectoryPage(app, w, r, "/topics", "topics")
}

func handleTopic(app *newsplugin.App, w http.ResponseWriter, r *http.Request) {
	handleScopedCollectionPage(app, w, r, "/topics/", "topic")
}

func handleCoordCollection(app *newsplugin.App, w http.ResponseWriter, r *http.Request, spec newsplugin.CoordCollectionSpec) {
	if r.URL.Path != spec.Path {
		http.NotFound(w, r)
		return
	}
	index, ok := loadIndexOrError(app, w)
	if !ok {
		return
	}
	opts := readFeedOptions(r)
	data := newsplugin.BuildTypedCollectionPageData(app, index, spec, opts)
	data.AgentView = newsplugin.IsAgentViewer(r)
	renderTemplate(w, app, "typed_collection.html", data)
}

func handleAPIFeed(app *newsplugin.App, w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/feed" {
		http.NotFound(w, r)
		return
	}
	index, ok := loadIndexOrError(app, w)
	if !ok {
		return
	}
	opts := readFeedOptions(r)
	newsplugin.WriteJSON(w, http.StatusOK, newsplugin.BuildHomeAPIResponse(app, index, opts))
}

func handleAPICoordCollection(app *newsplugin.App, w http.ResponseWriter, r *http.Request, spec newsplugin.CoordCollectionSpec) {
	if r.URL.Path != spec.APIPath {
		http.NotFound(w, r)
		return
	}
	index, ok := loadIndexOrError(app, w)
	if !ok {
		return
	}
	opts := readFeedOptions(r)
	newsplugin.WriteJSON(w, http.StatusOK, newsplugin.BuildTypedCollectionAPIResponse(app, index, spec, opts, strings.TrimPrefix(spec.Path, "/")))
}

func handleAPIPost(app *newsplugin.App, w http.ResponseWriter, r *http.Request) {
	infoHash := newsplugin.PathValue("/api/posts/", r.URL.Path)
	if infoHash == "" {
		http.NotFound(w, r)
		return
	}
	index, ok := loadIndexOrError(app, w)
	if !ok {
		return
	}
	post, ok := index.PostByInfoHash[strings.ToLower(infoHash)]
	if !ok {
		http.NotFound(w, r)
		return
	}
	if canonical := newsplugin.CoordCanonicalAPIPath(post); canonical != "" {
		http.Redirect(w, r, newsplugin.AppendRawQuery(canonical, r.URL.RawQuery), http.StatusTemporaryRedirect)
		return
	}
	newsplugin.WriteJSON(w, http.StatusOK, newsplugin.BuildPostAPIResponse(app, index, "post", post, ""))
}

func handleAPITorrent(app *newsplugin.App, w http.ResponseWriter, r *http.Request) {
	infoHash := newsplugin.PathValue("/api/torrents/", r.URL.Path)
	infoHash = strings.TrimSuffix(strings.ToLower(strings.TrimSpace(infoHash)), ".torrent")
	if infoHash == "" {
		http.NotFound(w, r)
		return
	}
	path := filepath.Join(app.StoreRoot(), "torrents", infoHash+".torrent")
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/x-bittorrent")
	http.ServeFile(w, r, path)
}

func handleAPISources(app *newsplugin.App, w http.ResponseWriter, r *http.Request) {
	handleDirectoryAPI(app, w, r, "/api/sources", "sources")
}

func handleAPISource(app *newsplugin.App, w http.ResponseWriter, r *http.Request) {
	handleScopedCollectionAPI(app, w, r, "/api/sources/", "source")
}

func handleAPITopics(app *newsplugin.App, w http.ResponseWriter, r *http.Request) {
	handleDirectoryAPI(app, w, r, "/api/topics", "topics")
}

func handleAPITopic(app *newsplugin.App, w http.ResponseWriter, r *http.Request) {
	handleScopedCollectionAPI(app, w, r, "/api/topics/", "topic")
}

func readFeedOptions(r *http.Request) newsplugin.FeedOptions {
	return newsplugin.FeedOptionsFromRequest(r)
}

func loadIndexOrError(app *newsplugin.App, w http.ResponseWriter) (newsplugin.Index, bool) {
	index, err := app.Index()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return newsplugin.Index{}, false
	}
	return index, true
}

func renderTemplate(w http.ResponseWriter, app *newsplugin.App, name string, data any) {
	if err := app.Templates().ExecuteTemplate(w, name, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleDirectoryPage(app *newsplugin.App, w http.ResponseWriter, r *http.Request, path, kind string) {
	if r.URL.Path != path {
		http.NotFound(w, r)
		return
	}
	index, ok := loadIndexOrError(app, w)
	if !ok {
		return
	}
	data, _ := newsplugin.BuildDirectoryPageData(app, index, kind)
	renderTemplate(w, app, "directory.html", data)
}

func handleScopedCollectionPage(app *newsplugin.App, w http.ResponseWriter, r *http.Request, prefix, scope string) {
	name := newsplugin.PathValue(prefix, r.URL.Path)
	if name == "" {
		http.NotFound(w, r)
		return
	}
	index, ok := loadIndexOrError(app, w)
	if !ok {
		return
	}
	data, found := newsplugin.BuildScopedCollectionPageData(app, index, scope, name, readFeedOptions(r))
	if !found {
		http.NotFound(w, r)
		return
	}
	renderTemplate(w, app, "collection.html", data)
}

func handleDirectoryAPI(app *newsplugin.App, w http.ResponseWriter, r *http.Request, path, kind string) {
	if r.URL.Path != path {
		http.NotFound(w, r)
		return
	}
	index, ok := loadIndexOrError(app, w)
	if !ok {
		return
	}
	payload, _ := newsplugin.BuildDirectoryAPIResponse(app, index, kind)
	newsplugin.WriteJSON(w, http.StatusOK, payload)
}

func handleScopedCollectionAPI(app *newsplugin.App, w http.ResponseWriter, r *http.Request, prefix, scope string) {
	name := newsplugin.PathValue(prefix, r.URL.Path)
	if name == "" {
		http.NotFound(w, r)
		return
	}
	index, ok := loadIndexOrError(app, w)
	if !ok {
		return
	}
	payload, found := newsplugin.BuildScopedCollectionAPIResponse(app, index, scope, name, readFeedOptions(r))
	if !found {
		http.NotFound(w, r)
		return
	}
	newsplugin.WriteJSON(w, http.StatusOK, payload)
}
