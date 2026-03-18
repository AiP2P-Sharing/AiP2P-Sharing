package newsplugin

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPageURLEncodesMetadataFilters(t *testing.T) {
	t.Parallel()

	url := pageURL("/tasks", FeedOptions{
		MetaKey:   "task.status",
		MetaValue: "in_progress",
		PageSize:  50,
	}, "sort", "score")
	if !strings.Contains(url, "meta_key=task.status") {
		t.Fatalf("url = %q, missing meta_key", url)
	}
	if !strings.Contains(url, "meta_value=in_progress") {
		t.Fatalf("url = %q, missing meta_value", url)
	}
	if !strings.Contains(url, "sort=score") {
		t.Fatalf("url = %q, missing sort", url)
	}
}

func TestBuildActiveFiltersIncludesMetadataFilter(t *testing.T) {
	t.Parallel()

	filters := BuildActiveFilters(FeedOptions{
		MetaKey:   "task.status",
		MetaValue: "in_progress",
	}, "/tasks")
	if len(filters) != 1 {
		t.Fatalf("len = %d, want 1", len(filters))
	}
	if filters[0].Label != "task.status: in_progress" {
		t.Fatalf("label = %q", filters[0].Label)
	}
	if !strings.Contains(filters[0].URL, "/tasks") {
		t.Fatalf("url = %q, want /tasks path", filters[0].URL)
	}
}

func TestBuildCollectionSummaryStatsForSource(t *testing.T) {
	t.Parallel()

	posts := []Post{
		{},
		{},
	}

	stats := BuildCollectionSummaryStats("source", posts)
	if len(stats) != 4 {
		t.Fatalf("len = %d, want 4", len(stats))
	}
	if stats[0].Label != "Visible assets" {
		t.Fatalf("first label = %q", stats[0].Label)
	}
	if stats[1].Label != "Discussion replies" {
		t.Fatalf("second label = %q", stats[1].Label)
	}
	if stats[2].Label != "Signals" {
		t.Fatalf("third label = %q", stats[2].Label)
	}
}

func TestBuildDirectorySummaryStatsForTopics(t *testing.T) {
	t.Parallel()

	stats := BuildDirectorySummaryStatsForKind("topics", []FacetStat{{Name: "website", Count: 3}}, []Post{{}, {}})
	if len(stats) != 4 {
		t.Fatalf("len = %d, want 4", len(stats))
	}
	if stats[0].Label != "Tracked workstreams" {
		t.Fatalf("first label = %q", stats[0].Label)
	}
	if stats[1].Label != "Indexed assets" {
		t.Fatalf("second label = %q", stats[1].Label)
	}
}

func TestAPIWorkspaceNavigationIncludesLegacyAliases(t *testing.T) {
	t.Parallel()

	panel := &WorkspaceJumpPanel{
		Eyebrow:           "Workspace jumps",
		Title:             "Primary and supporting jumps",
		Description:       "Test panel",
		PrimaryActions:    []CoordAction{{Label: "Open task", URL: "/tasks/task-1"}},
		SupportingActions: []CoordAction{{Label: "Open source", URL: "/sources/source-1"}},
	}
	payload := APIWorkspaceNavigation(panel, map[string][]CoordAction{
		"registry_actions": {{Label: "Topic workstream", URL: "/topics/website"}},
	})
	if payload["canonical_field"] != "jump_panel" {
		t.Fatalf("canonical_field = %v", payload["canonical_field"])
	}
	if _, ok := payload["jump_panel"]; ok {
		t.Fatalf("payload = %+v, unexpected nested jump_panel duplicate", payload)
	}
	legacy, ok := payload["legacy_fields"].(map[string]any)
	if !ok {
		t.Fatalf("legacy_fields type = %T", payload["legacy_fields"])
	}
	if _, ok := legacy["registry_actions"]; !ok {
		t.Fatalf("legacy_fields = %+v", legacy)
	}
}

func TestAPIRegistryWorkbenchExposesModuleActionsAndLegacyAlias(t *testing.T) {
	t.Parallel()

	payload := APIRegistryWorkbench([]RegistryWorkbenchCard{{
		Kind:          "topic",
		Name:          "website",
		URL:           "/topics/website",
		ModuleActions: []CoordAction{{Label: "Open tasks", URL: "/tasks?topic=website"}},
		JumpPanel: &WorkspaceJumpPanel{
			Title:          "Workspace jumps",
			PrimaryActions: []CoordAction{{Label: "Open topic", URL: "/topics/website"}},
		},
	}})
	if len(payload) != 1 {
		t.Fatalf("len = %d, want 1", len(payload))
	}
	if _, ok := payload[0]["module_actions"]; !ok {
		t.Fatalf("payload = %+v, missing module_actions", payload[0])
	}
	if _, ok := payload[0]["typed_actions"]; ok {
		t.Fatalf("payload = %+v, unexpected typed_actions top-level alias", payload[0])
	}
	navigation, ok := payload[0]["navigation"].(map[string]any)
	if !ok {
		t.Fatalf("navigation type = %T", payload[0]["navigation"])
	}
	if _, ok := navigation["jump_panel"]; ok {
		t.Fatalf("navigation = %+v, unexpected nested jump_panel duplicate", navigation)
	}
	legacy, ok := navigation["legacy_fields"].(map[string]any)
	if !ok {
		t.Fatalf("legacy_fields type = %T", navigation["legacy_fields"])
	}
	if _, ok := legacy["typed_actions"]; !ok {
		t.Fatalf("legacy_fields = %+v", legacy)
	}
}

func TestBuildScopedCollectionAPIResponseExposesModuleActionsAndLegacyAlias(t *testing.T) {
	t.Parallel()

	app := &App{project: "aip2p.sharing"}
	index := Index{
		Posts: []Post{
			{
				Bundle: Bundle{
					InfoHash: "task-1",
					Message: Message{
						Channel: "aip2p.sharing/tasks",
						Title:   "Task 1",
						Extensions: map[string]any{
							"coord.type":  "task",
							"task.status": "in_progress",
						},
					},
				},
				Topics: []string{"website"},
				SourceName: "source-a",
			},
		},
		SourceStats: []FacetStat{{Name: "source-a", Count: 1}},
		TopicStats:  []FacetStat{{Name: "website", Count: 1}},
	}
	payload, ok := BuildScopedCollectionAPIResponse(app, index, "topic", "website", FeedOptions{})
	if !ok {
		t.Fatal("ok = false, want true")
	}
	if _, ok := payload["module_actions"]; !ok {
		t.Fatalf("payload = %+v, missing module_actions", payload)
	}
	if _, ok := payload["typed_entry_points"]; ok {
		t.Fatalf("payload = %+v, unexpected typed_entry_points top-level alias", payload)
	}
	navigation, ok := payload["navigation"].(map[string]any)
	if !ok {
		t.Fatalf("navigation type = %T", payload["navigation"])
	}
	if _, ok := navigation["jump_panel"]; ok {
		t.Fatalf("navigation = %+v, unexpected nested jump_panel duplicate", navigation)
	}
	legacy, ok := navigation["legacy_fields"].(map[string]any)
	if !ok {
		t.Fatalf("legacy_fields type = %T", navigation["legacy_fields"])
	}
	if _, ok := legacy["typed_entry_points"]; !ok {
		t.Fatalf("legacy_fields = %+v", legacy)
	}
}

func TestBuildTypedCollectionAPIResponseIncludesMetadataSchema(t *testing.T) {
	t.Parallel()

	app := &App{project: "aip2p.sharing"}
	spec, ok := CoordCollectionByType("task")
	if !ok {
		t.Fatal("missing task collection spec")
	}
	index := Index{
		Posts: []Post{{
			Bundle: Bundle{
				InfoHash: "task-1",
				Message: Message{
					Channel: "aip2p.sharing/tasks",
					Title:   "Task 1",
					Extensions: map[string]any{
						"coord.type":  "task",
						"task.status": "in_progress",
					},
				},
			},
		}},
	}
	payload := BuildTypedCollectionAPIResponse(app, index, spec, FeedOptions{}, "tasks")
	schema, ok := payload["metadata_schema"].(map[string]any)
	if !ok {
		t.Fatalf("metadata_schema type = %T", payload["metadata_schema"])
	}
	if schema["coord_type"] != "task" {
		t.Fatalf("coord_type = %v, want task", schema["coord_type"])
	}
	fields, ok := schema["fields"].([]map[string]any)
	if !ok || len(fields) == 0 {
		t.Fatalf("fields = %#v", schema["fields"])
	}
	if fields[0]["key"] != "task.id" {
		t.Fatalf("first key = %v, want task.id", fields[0]["key"])
	}
}

func TestBuildTypedCollectionPageDataIncludesMetadataSchema(t *testing.T) {
	t.Parallel()

	app := &App{project: "aip2p.sharing", version: "0.1.0.6"}
	spec, ok := CoordCollectionByType("task")
	if !ok {
		t.Fatal("missing task collection spec")
	}
	index := Index{
		Posts: []Post{{
			Bundle: Bundle{
				InfoHash: "task-1",
				Message: Message{
					Channel: "aip2p.sharing/tasks",
					Title:   "Task 1",
					Extensions: map[string]any{
						"coord.type":  "task",
						"task.status": "in_progress",
					},
				},
			},
		}},
	}
	data := BuildTypedCollectionPageData(app, index, spec, FeedOptions{Topic: "website", Source: "writer-a"})
	if data.MetadataSchema == nil {
		t.Fatal("metadata schema = nil")
	}
	if data.MetadataSchema.CoordType != "task" {
		t.Fatalf("coord type = %q, want task", data.MetadataSchema.CoordType)
	}
	if data.PublishGuide == nil {
		t.Fatal("publish guide = nil")
	}
	if data.PublishGuide.Channel != "aip2p.sharing/tasks" {
		t.Fatalf("channel = %q, want aip2p.sharing/tasks", data.PublishGuide.Channel)
	}
	if !strings.Contains(data.PublishGuide.StarterCommand, "cmd/aip2p publish") {
		t.Fatalf("starter command = %q", data.PublishGuide.StarterCommand)
	}
	if !strings.Contains(data.PublishGuide.StarterExtensionsJSON, "\"topics\": [\n    \"website\"\n  ]") {
		t.Fatalf("starter json = %q", data.PublishGuide.StarterExtensionsJSON)
	}
	if len(data.PublishGuide.ContextNotes) == 0 {
		t.Fatal("context notes = empty")
	}
	if len(data.PublishGuide.Templates) != 9 {
		t.Fatalf("templates len = %d, want 9", len(data.PublishGuide.Templates))
	}
	if data.PublishGuide.Templates[0].Kind != "reply" {
		t.Fatalf("template kind = %q, want reply", data.PublishGuide.Templates[0].Kind)
	}
	if !strings.Contains(data.PublishGuide.Templates[2].StarterExtensionsJSON, "\"thread.role\": \"review-request\"") {
		t.Fatalf("template json = %q", data.PublishGuide.Templates[2].StarterExtensionsJSON)
	}
	if !strings.Contains(data.PublishGuide.Templates[3].StarterExtensionsJSON, "\"thread.role\": \"handoff\"") {
		t.Fatalf("template json = %q", data.PublishGuide.Templates[3].StarterExtensionsJSON)
	}
	if data.PublishGuide.Templates[4].Channel != "aip2p.sharing/skills" {
		t.Fatalf("template channel = %q, want aip2p.sharing/skills", data.PublishGuide.Templates[4].Channel)
	}
	if !strings.Contains(data.PublishGuide.Templates[5].StarterExtensionsJSON, "\"md.kind\": \"delivery-note\"") {
		t.Fatalf("template json = %q", data.PublishGuide.Templates[5].StarterExtensionsJSON)
	}
	if data.PublishGuide.Templates[7].Kind != "reaction" {
		t.Fatalf("template kind = %q, want reaction", data.PublishGuide.Templates[7].Kind)
	}
	if !strings.Contains(data.PublishGuide.Templates[7].StarterCommand, "--kind reaction") {
		t.Fatalf("template command = %q", data.PublishGuide.Templates[7].StarterCommand)
	}
	if len(data.PublishGuide.Workflow) != 3 {
		t.Fatalf("workflow len = %d, want 3", len(data.PublishGuide.Workflow))
	}
	if len(data.PublishGuide.Workflow[0].Templates) != 4 {
		t.Fatalf("workflow templates len = %d, want 4", len(data.PublishGuide.Workflow[0].Templates))
	}
}

func TestBuildScopedCollectionPageDataBuildsSourceSurface(t *testing.T) {
	t.Parallel()

	app := &App{project: "aip2p.sharing", version: "0.1.0.6"}
	index := Index{
		Posts: []Post{{
			Bundle: Bundle{
				InfoHash: "asset-1",
				Message: Message{
					Channel: "aip2p.sharing/general",
					Title:   "Asset 1",
					Extensions: map[string]any{
						"coord.type": "markdown",
					},
				},
			},
			SourceName: "source-a",
			Topics:     []string{"website"},
		}},
		SourceStats: []FacetStat{{Name: "source-a", Count: 1}},
		TopicStats:  []FacetStat{{Name: "website", Count: 1}},
	}
	data, ok := BuildScopedCollectionPageData(app, index, "source", "source-a", FeedOptions{})
	if !ok {
		t.Fatal("ok = false, want true")
	}
	if data.Kind != "Source" {
		t.Fatalf("kind = %q, want Source", data.Kind)
	}
	if len(data.ModuleActions) == 0 {
		t.Fatalf("module actions = %d, want > 0", len(data.ModuleActions))
	}
	if data.JumpPanel == nil {
		t.Fatal("jump panel = nil, want non-nil")
	}
	if data.ExternalURL != "" {
		t.Fatalf("external url = %q, want empty for local source", data.ExternalURL)
	}
}

func TestBuildDirectoryPageDataBuildsTopicsDirectory(t *testing.T) {
	t.Parallel()

	app := &App{project: "aip2p.sharing", version: "0.1.0.6"}
	index := Index{
		Posts:      []Post{{}},
		TopicStats: []FacetStat{{Name: "website", Count: 1}},
	}
	data, ok := BuildDirectoryPageData(app, index, "topics")
	if !ok {
		t.Fatal("ok = false, want true")
	}
	if data.Kind != "Topics" {
		t.Fatalf("kind = %q, want Topics", data.Kind)
	}
	if data.Path != "/topics" {
		t.Fatalf("path = %q, want /topics", data.Path)
	}
}

func TestBuildDirectoryAPIResponseBuildsSourcesDirectory(t *testing.T) {
	t.Parallel()

	app := &App{project: "aip2p.sharing"}
	index := Index{
		Posts:       []Post{{}},
		SourceStats: []FacetStat{{Name: "source-a", Count: 1}},
	}
	payload, ok := BuildDirectoryAPIResponse(app, index, "sources")
	if !ok {
		t.Fatal("ok = false, want true")
	}
	if payload["scope"] != "sources" {
		t.Fatalf("scope = %v, want sources", payload["scope"])
	}
	if _, ok := payload["items"]; !ok {
		t.Fatalf("payload = %+v, missing items", payload)
	}
}

func TestBuildHomePageDataIncludesWorkbenchAndFilters(t *testing.T) {
	t.Parallel()

	app := &App{project: "aip2p.sharing", version: "0.1.0.6", listenAddr: "127.0.0.1:1818"}
	index := Index{
		Posts: []Post{{
			Bundle: Bundle{
				InfoHash: "idea-1",
				Message: Message{
					Channel: "aip2p.sharing/general",
					Title:   "Idea 1",
					Extensions: map[string]any{
						"coord.type": "idea",
					},
				},
			},
			Topics:     []string{"website"},
			SourceName: "source-a",
		}},
		TopicStats:  []FacetStat{{Name: "website", Count: 1}},
		SourceStats: []FacetStat{{Name: "source-a", Count: 1}},
	}
	data := BuildHomePageData(app, index, SubscriptionRules{Topics: []string{"all"}}, FeedOptions{}, true, true)
	if !data.AgentView {
		t.Fatal("agent view = false, want true")
	}
	if !data.ShowNetworkWarn {
		t.Fatal("show network warn = false, want true")
	}
	if len(data.WorkspaceStats) == 0 {
		t.Fatalf("workspace stats = %d, want > 0", len(data.WorkspaceStats))
	}
	if len(data.TopicFacets) == 0 {
		t.Fatalf("topic facets = %d, want > 0", len(data.TopicFacets))
	}
}

func TestBuildHomeAPIResponseIncludesWorkspacePayloads(t *testing.T) {
	t.Parallel()

	app := &App{project: "aip2p.sharing"}
	index := Index{
		Posts: []Post{{
			Bundle: Bundle{
				InfoHash: "task-1",
				Message: Message{
					Channel: "aip2p.sharing/tasks",
					Title:   "Task 1",
					Extensions: map[string]any{
						"coord.type":  "task",
						"task.status": "in_progress",
					},
				},
			},
			Topics:     []string{"website"},
			SourceName: "source-a",
		}},
		TopicStats:  []FacetStat{{Name: "website", Count: 1}},
		SourceStats: []FacetStat{{Name: "source-a", Count: 1}},
	}
	payload := BuildHomeAPIResponse(app, index, FeedOptions{})
	if payload["scope"] != "feed" {
		t.Fatalf("scope = %v, want feed", payload["scope"])
	}
	if _, ok := payload["workspace_clusters"]; !ok {
		t.Fatalf("payload = %+v, missing workspace_clusters", payload)
	}
	if _, ok := payload["agent_workspace"]; !ok {
		t.Fatalf("payload = %+v, missing agent_workspace", payload)
	}
	if _, ok := payload["topic_workstreams"]; !ok {
		t.Fatalf("payload = %+v, missing topic_workstreams", payload)
	}
}

func TestAppendRawQuery(t *testing.T) {
	t.Parallel()

	if got := AppendRawQuery("/ideas/idea-1", "agent=1"); got != "/ideas/idea-1?agent=1" {
		t.Fatalf("AppendRawQuery() = %q", got)
	}
	if got := AppendRawQuery("/ideas/idea-1?view=compact", "agent=1"); got != "/ideas/idea-1?view=compact&agent=1" {
		t.Fatalf("AppendRawQuery() = %q", got)
	}
	if got := AppendRawQuery("/ideas/idea-1", ""); got != "/ideas/idea-1" {
		t.Fatalf("AppendRawQuery() = %q", got)
	}
}

func TestShouldShowNetworkWarning(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	if !ShouldShowNetworkWarning(req) {
		t.Fatal("ShouldShowNetworkWarning() = false, want true without cookie")
	}
	req.AddCookie(&http.Cookie{Name: "aip2p_news_network_warning_seen", Value: "1"})
	if ShouldShowNetworkWarning(req) {
		t.Fatal("ShouldShowNetworkWarning() = true, want false with cookie")
	}
}

func TestIsAgentViewer(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/?agent=1", nil)
	if !IsAgentViewer(req) {
		t.Fatal("IsAgentViewer() = false, want true for agent query")
	}
	browser := httptest.NewRequest(http.MethodGet, "/", nil)
	browser.Header.Set("User-Agent", "Mozilla/5.0")
	if IsAgentViewer(browser) {
		t.Fatal("IsAgentViewer() = true, want false for browser ua")
	}
	client := httptest.NewRequest(http.MethodGet, "/", nil)
	client.Header.Set("User-Agent", "curl/8.6.0")
	if !IsAgentViewer(client) {
		t.Fatal("IsAgentViewer() = false, want true for cli ua")
	}
}
