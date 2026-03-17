package newsplugin

import (
	"strings"
	"time"
)

func BuildTypedCollectionPageData(app *App, index Index, spec CoordCollectionSpec, opts FeedOptions) TypedCollectionPageData {
	allPosts := FilterPostsByCoordType(index.FilterPosts(opts), spec.CoordType)
	posts, pagination := PaginatePosts(allPosts, opts, spec.Path)
	fullSet := FilterPostsByCoordType(index.FilterPosts(FeedOptions{Now: opts.Now}), spec.CoordType)
	return TypedCollectionPageData{
		Project:          app.ProjectName(),
		Version:          app.VersionString(),
		Title:            spec.Title,
		CoordType:        spec.CoordType,
		MetadataSchema:   CoordMetadataSchema(spec.CoordType),
		PublishGuide:     CoordPublishGuideFor(app.ProjectID(), spec.CoordType, CoordPublishContext{Topics: typedGuideTopics(opts), SourceName: strings.TrimSpace(opts.Source)}),
		Description:      spec.Description,
		Path:             spec.Path,
		APIPath:          spec.APIPath,
		Now:              opts.Now,
		Posts:            posts,
		Options:          opts,
		PageNav:          app.PageNav(spec.Path),
		SortOptions:      BuildSortOptions(opts, spec.Path),
		WindowOptions:    BuildWindowOptions(opts, spec.Path),
		PageSizeOptions:  BuildPageSizeOptions(opts, spec.Path),
		TopicFacetLabel:  spec.TopicFacetLabel,
		TopicFacets:      BuildFacetLinks(TopicStatsForPosts(fullSet), opts, spec.Path, "topic"),
		SourceFacetLabel: spec.SourceFacetLabel,
		SourceFacets:     BuildFacetLinks(SourceStatsForPosts(fullSet), opts, spec.Path, "source"),
		MetaFacetLabel:   spec.MetaFacetLabel,
		MetaFacets:       BuildMetadataFacetLinks(ExtensionFacetStatsForPosts(fullSet, spec.MetaExtensionKey), opts, spec.Path, spec.MetaExtensionKey),
		ActiveFilters:    BuildActiveFilters(opts, spec.Path),
		SummaryStats:     BuildSummaryStats(allPosts),
		TotalPostCount:   len(fullSet),
		Pagination:       pagination,
		EmptyTitle:       spec.EmptyTitle,
		EmptyCopy:        spec.EmptyCopy,
		NodeStatus:       app.NodeStatus(index),
	}
}

func BuildTypedCollectionAPIResponse(app *App, index Index, spec CoordCollectionSpec, opts FeedOptions, scope string) map[string]any {
	allPosts := FilterPostsByCoordType(index.FilterPosts(opts), spec.CoordType)
	posts, pagination := PaginatePosts(allPosts, opts, spec.APIPath)
	fullSet := FilterPostsByCoordType(index.FilterPosts(FeedOptions{Now: opts.Now}), spec.CoordType)
	return map[string]any{
		"project":         app.ProjectID(),
		"scope":           scope,
		"title":           spec.Title,
		"coord_type":      spec.CoordType,
		"metadata_schema": APICoordMetadataSchema(CoordMetadataSchema(spec.CoordType)),
		"publish_guide":   APICoordPublishGuide(CoordPublishGuideFor(app.ProjectID(), spec.CoordType, CoordPublishContext{Topics: typedGuideTopics(opts), SourceName: strings.TrimSpace(opts.Source)})),
		"options":         APIOptions(opts),
		"summary":         BuildSummaryStats(allPosts),
		"pagination":      pagination,
		"posts":           APIPosts(posts),
		"facets": map[string]any{
			"topics":   TopicStatsForPosts(fullSet),
			"sources":  SourceStatsForPosts(fullSet),
			"metadata": ExtensionFacetStatsForPosts(fullSet, spec.MetaExtensionKey),
		},
	}
}

func BuildHomePageData(app *App, index Index, rules SubscriptionRules, opts FeedOptions, agentView, showNetworkWarn bool) HomePageData {
	allPosts := index.FilterPosts(opts)
	posts, pagination := PaginatePosts(allPosts, opts, "/")
	workbenchPosts := index.FilterPosts(FeedOptions{Now: opts.Now})
	return HomePageData{
		Project:           app.ProjectName(),
		Version:           app.VersionString(),
		Posts:             posts,
		LegacyPosts:       LegacyPosts(workbenchPosts, 4),
		Now:               time.Now(),
		ListenAddr:        app.HTTPListenAddr(),
		AgentView:         agentView,
		ShowNetworkWarn:   showNetworkWarn,
		Options:           opts,
		PageNav:           app.PageNav("/"),
		TopicFacets:       BuildFeedFacets(index.TopicStats, opts, "/", "topic"),
		SourceFacets:      BuildFeedFacets(index.SourceStats, opts, "/", "source"),
		SortOptions:       BuildSortOptions(opts, "/"),
		WindowOptions:     BuildWindowOptions(opts, "/"),
		PageSizeOptions:   BuildPageSizeOptions(opts, "/"),
		ActiveFilters:     BuildActiveFilters(opts, "/"),
		CoordGroupFacets:  HomeCoordGroupFacets(workbenchPosts, 6),
		SummaryStats:      BuildSummaryStats(allPosts),
		PriorityCards:     HomePriorityCards(workbenchPosts),
		WorkspaceStats:    HomeWorkbenchStats(workbenchPosts),
		WorkspaceActions:  HomeWorkbenchActions(workbenchPosts),
		WorkspaceLanes:    HomeWorkbenchLanes(workbenchPosts),
		WorkspaceClusters: HomeWorkbenchClusters(index, workbenchPosts, 3, 2),
		AgentWorkspace:    HomeAgentWorkbench(index, workbenchPosts, 2, 2),
		TopicWorkspace:    HomeTopicWorkbench(index, workbenchPosts, 3),
		SourceWorkspace:   HomeSourceWorkbench(index, workbenchPosts, 3),
		TotalPostCount:    len(index.Posts),
		Pagination:        pagination,
		Subscriptions:     rules,
		NodeStatus:        app.NodeStatus(index),
	}
}

func BuildHomeAPIResponse(app *App, index Index, opts FeedOptions) map[string]any {
	allPosts := index.FilterPosts(opts)
	posts, pagination := PaginatePosts(allPosts, opts, "/api/feed")
	workbenchPosts := index.FilterPosts(FeedOptions{Now: opts.Now})
	return map[string]any{
		"project": app.ProjectID(),
		"scope":   "feed",
		"options": APIOptions(opts),
		"summary": BuildSummaryStats(allPosts),
		"priority_surface": APIHomePriorityCards(
			HomePriorityCards(workbenchPosts),
		),
		"pagination": pagination,
		"posts":      APIPosts(posts),
		"coordination_groups": APIFacets(
			HomeCoordGroupFacets(workbenchPosts, 6),
		),
		"workspace_clusters": APIWorkbenchClusters(
			HomeWorkbenchClusters(index, workbenchPosts, 3, 2),
		),
		"agent_workspace": APIAgentWorkbench(
			HomeAgentWorkbench(index, workbenchPosts, 2, 2),
		),
		"topic_workstreams": APIRegistryWorkbench(
			HomeTopicWorkbench(index, workbenchPosts, 3),
		),
		"source_workstreams": APIRegistryWorkbench(
			HomeSourceWorkbench(index, workbenchPosts, 3),
		),
		"facets": map[string]any{
			"channels": index.ChannelStats,
			"topics":   index.TopicStats,
			"sources":  index.SourceStats,
		},
	}
}

func BuildDirectoryPageData(app *App, index Index, kind string) (DirectoryPageData, bool) {
	switch strings.ToLower(strings.TrimSpace(kind)) {
	case "sources", "source":
		return DirectoryPageData{
			Project:      app.ProjectName(),
			Version:      app.VersionString(),
			Kind:         "Sources",
			Path:         "/sources",
			APIPath:      "/api/sources",
			Now:          time.Now(),
			PageNav:      app.PageNav("/sources"),
			Items:        BuildSourceDirectory(index),
			SummaryStats: BuildDirectorySummaryStatsForKind("sources", index.SourceStats, index.Posts),
			NodeStatus:   app.NodeStatus(index),
		}, true
	case "topics", "topic":
		return DirectoryPageData{
			Project:      app.ProjectName(),
			Version:      app.VersionString(),
			Kind:         "Topics",
			Path:         "/topics",
			APIPath:      "/api/topics",
			Now:          time.Now(),
			PageNav:      app.PageNav("/topics"),
			Items:        BuildTopicDirectory(index),
			SummaryStats: BuildDirectorySummaryStatsForKind("topics", index.TopicStats, index.Posts),
			NodeStatus:   app.NodeStatus(index),
		}, true
	default:
		return DirectoryPageData{}, false
	}
}

func BuildDirectoryAPIResponse(app *App, index Index, kind string) (map[string]any, bool) {
	switch strings.ToLower(strings.TrimSpace(kind)) {
	case "sources", "source":
		return map[string]any{
			"project": app.ProjectID(),
			"scope":   "sources",
			"summary": BuildDirectorySummaryStatsForKind("sources", index.SourceStats, index.Posts),
			"items":   BuildSourceDirectory(index),
		}, true
	case "topics", "topic":
		return map[string]any{
			"project": app.ProjectID(),
			"scope":   "topics",
			"summary": BuildDirectorySummaryStatsForKind("topics", index.TopicStats, index.Posts),
			"items":   BuildTopicDirectory(index),
		}, true
	default:
		return nil, false
	}
}

type scopedCollectionConfig struct {
	kindLabel   string
	path        string
	apiPath     string
	directory   string
	navPath     string
	sideLabel   string
	facetKind   string
	omitField   string
	externalURL bool
}

func scopedCollectionConfigFor(scopeKind, name string) (scopedCollectionConfig, bool) {
	switch scopeKind {
	case "source":
		return scopedCollectionConfig{
			kindLabel:   "Source",
			path:        SourcePath(name),
			apiPath:     "/api" + SourcePath(name),
			directory:   "/sources",
			navPath:     "/sources",
			sideLabel:   "Topics from this source",
			facetKind:   "topic",
			omitField:   "source",
			externalURL: true,
		}, true
	case "topic":
		return scopedCollectionConfig{
			kindLabel: "Topic",
			path:      TopicPath(name),
			apiPath:   "/api" + TopicPath(name),
			directory: "/topics",
			navPath:   "/topics",
			sideLabel: "Sources covering this topic",
			facetKind: "source",
			omitField: "topic",
		}, true
	default:
		return scopedCollectionConfig{}, false
	}
}

func scopedCollectionOptions(scopeKind, name string, opts FeedOptions) FeedOptions {
	scoped := opts
	switch scopeKind {
	case "source":
		scoped.Source = name
	case "topic":
		scoped.Topic = name
	}
	return scoped
}

func scopedCollectionFullSet(index Index, scopeKind, name string, now time.Time) []Post {
	switch scopeKind {
	case "source":
		return index.FilterPosts(FeedOptions{Source: name, Now: now})
	case "topic":
		return index.FilterPosts(FeedOptions{Topic: name, Now: now})
	default:
		return nil
	}
}

func hasScopedCollection(index Index, scopeKind, name string) bool {
	switch scopeKind {
	case "source":
		return HasSource(index, name)
	case "topic":
		return HasTopic(index, name)
	default:
		return false
	}
}

func scopedCollectionSideFacets(scopeKind string, posts []Post, opts FeedOptions, path string) []FeedFacet {
	switch scopeKind {
	case "source":
		return BuildFacetLinks(TopicStatsForPosts(posts), opts, path, "topic", "source")
	case "topic":
		return BuildFacetLinks(SourceStatsForPosts(posts), opts, path, "source", "topic")
	default:
		return nil
	}
}

func scopedCollectionAPIFacets(scopeKind string, posts []Post) map[string]any {
	switch scopeKind {
	case "source":
		return map[string]any{
			"channels": ChannelStatsForPosts(posts),
			"topics":   TopicStatsForPosts(posts),
		}
	case "topic":
		return map[string]any{
			"channels": ChannelStatsForPosts(posts),
			"sources":  SourceStatsForPosts(posts),
		}
	default:
		return map[string]any{}
	}
}

func BuildScopedCollectionPageData(app *App, index Index, scopeKind, name string, opts FeedOptions) (CollectionPageData, bool) {
	cfg, ok := scopedCollectionConfigFor(scopeKind, name)
	if !ok || !hasScopedCollection(index, scopeKind, name) {
		return CollectionPageData{}, false
	}
	scopedOpts := scopedCollectionOptions(scopeKind, name, opts)
	allPosts := index.FilterPosts(scopedOpts)
	posts, pagination := PaginatePosts(allPosts, scopedOpts, cfg.path)
	fullSet := scopedCollectionFullSet(index, scopeKind, name, scopedOpts.Now)
	data := CollectionPageData{
		Project:         app.ProjectName(),
		Version:         app.VersionString(),
		Kind:            cfg.kindLabel,
		Name:            name,
		Path:            cfg.path,
		DirectoryURL:    cfg.directory,
		APIPath:         cfg.apiPath,
		Now:             time.Now(),
		Posts:           posts,
		Options:         scopedOpts,
		PageNav:         app.PageNav(cfg.navPath),
		SortOptions:     BuildSortOptions(scopedOpts, cfg.path, cfg.omitField),
		WindowOptions:   BuildWindowOptions(scopedOpts, cfg.path, cfg.omitField),
		PageSizeOptions: BuildPageSizeOptions(scopedOpts, cfg.path, cfg.omitField),
		SideLabel:       cfg.sideLabel,
		SideFacets:      scopedCollectionSideFacets(scopeKind, fullSet, scopedOpts, cfg.path),
		ActiveFilters:   BuildActiveFilters(scopedOpts, cfg.path, cfg.omitField),
		SummaryStats:    BuildCollectionSummaryStats(scopeKind, allPosts),
		WorkspaceStats:  HomeWorkbenchStats(fullSet),
		WorkspaceLanes:  HomeWorkbenchLanes(fullSet),
		JumpPanel:       WorkstreamJumpPanel(scopeKind, name, fullSet),
		RelationPanel:   WorkstreamRelationPanelForPosts(fullSet, 4),
		ModuleActions:   ScopedCoordActions(scopeKind, name, fullSet),
		TotalPostCount:  len(fullSet),
		Pagination:      pagination,
		NodeStatus:      app.NodeStatus(index),
	}
	if cfg.externalURL {
		data.ExternalURL = SourceURLFromPosts(fullSet)
	}
	return data, true
}

func BuildScopedCollectionAPIResponse(app *App, index Index, scopeKind, name string, opts FeedOptions) (map[string]any, bool) {
	cfg, ok := scopedCollectionConfigFor(scopeKind, name)
	if !ok || !hasScopedCollection(index, scopeKind, name) {
		return nil, false
	}
	scopedOpts := scopedCollectionOptions(scopeKind, name, opts)
	posts := index.FilterPosts(scopedOpts)
	fullSet := scopedCollectionFullSet(index, scopeKind, name, scopedOpts.Now)
	jumpPanel := WorkstreamJumpPanel(scopeKind, name, fullSet)
	moduleActions := APIScopedCoordActions(scopeKind, name, fullSet)
	payload := map[string]any{
		"project":            app.ProjectID(),
		"scope":              scopeKind,
		"name":               name,
		"options":            APIOptions(scopedOpts),
		"summary":            BuildCollectionSummaryStats(scopeKind, posts),
		"posts":              APIPosts(posts),
		"workstream_summary": HomeWorkbenchStats(fullSet),
		"workstream_lanes":   APIWorkbenchLanes(HomeWorkbenchLanes(fullSet)),
		"coordination_panel": APIWorkstreamRelationPanel(WorkstreamRelationPanelForPosts(fullSet, 4)),
		"module_actions":     moduleActions,
		"jump_panel":         APIWorkspaceJumpPanel(jumpPanel),
		"navigation":         APIWorkspaceNavigation(jumpPanel, map[string][]CoordAction{"typed_entry_points": moduleActions}),
		"facets":             scopedCollectionAPIFacets(scopeKind, fullSet),
	}
	if cfg.externalURL {
		payload["source_url"] = SourceURLFromPosts(fullSet)
	}
	return payload, true
}

func BuildPostPageData(app *App, index Index, navPath string, post Post, relatedCoordType string) PostPageData {
	return PostPageData{
		Project:        app.ProjectName(),
		Version:        app.VersionString(),
		PageNav:        app.PageNav(navPath),
		Post:           post,
		CoordType:      PostCoordType(post),
		MetadataSchema: CoordMetadataSchema(PostCoordType(post)),
		PublishGuide:   CoordPublishGuideFor(app.ProjectID(), PostCoordType(post), CoordPublishContext{Topics: post.Topics, SourceName: post.SourceName, Extensions: post.Message.Extensions, InfoHash: post.InfoHash, Magnet: post.Magnet}),
		CoordFields:    CoordFieldsForPost(post),
		CoordSections:  CoordSectionsForPost(post),
		JumpPanel:      DetailJumpPanel(post),
		DetailFocus:    CoordDetailFocusForPost(index, post, 2),
		CoordRelations: CoordRelationGroupsForPost(index, post, 2),
		CoordWorkspace: CoordWorkbenchPanelForPost(index, post, 2),
		CoordRelated:   CoordExtendedContextGroupsForPost(index, post, 2),
		Replies:        index.RepliesByPost[strings.ToLower(post.InfoHash)],
		Reactions:      index.ReactionsByPost[strings.ToLower(post.InfoHash)],
		Related:        relatedPostsForScope(index, post.InfoHash, relatedCoordType, 4),
		NodeStatus:     app.NodeStatus(index),
	}
}

func BuildPostAPIResponse(app *App, index Index, scope string, post Post, relatedCoordType string) map[string]any {
	coordActions := CoordActionsForPost(post)
	contextActions := CoordContextActions(post)
	jumpPanel := DetailJumpPanel(post)
	return map[string]any{
		"project":         app.ProjectID(),
		"scope":           scope,
		"coord_type":      PostCoordType(post),
		"metadata_schema": APICoordMetadataSchema(CoordMetadataSchema(PostCoordType(post))),
		"publish_guide":   APICoordPublishGuide(CoordPublishGuideFor(app.ProjectID(), PostCoordType(post), CoordPublishContext{Topics: post.Topics, SourceName: post.SourceName, Extensions: post.Message.Extensions, InfoHash: post.InfoHash, Magnet: post.Magnet})),
		"post":            APIPost(post, true),
		"coord_fields":    CoordFieldsForPost(post),
		"coord_sections":  CoordSectionsForPost(post),
		"jump_panel":      APIWorkspaceJumpPanel(jumpPanel),
		"navigation":      APIWorkspaceNavigation(jumpPanel, map[string][]CoordAction{"coord_actions": coordActions, "context_links": contextActions}),
		"detail_focus":    APICoordDetailFocus(CoordDetailFocusForPost(index, post, 2)),
		"coord_relations": APICoordRelationGroups(CoordRelationGroupsForPost(index, post, 2)),
		"workspace_panel": APICoordWorkbench(CoordWorkbenchPanelForPost(index, post, 2)),
		"coord_related":   APICoordRelatedGroups(CoordExtendedContextGroupsForPost(index, post, 2)),
		"replies":         APIReplies(index.RepliesByPost[strings.ToLower(post.InfoHash)]),
		"reactions":       APIReactions(index.ReactionsByPost[strings.ToLower(post.InfoHash)]),
		"related":         APIPosts(relatedPostsForScope(index, post.InfoHash, relatedCoordType, 4)),
	}
}

func relatedPostsForScope(index Index, infoHash, relatedCoordType string, limit int) []Post {
	related := index.RelatedPosts(infoHash, limit)
	if strings.TrimSpace(relatedCoordType) == "" {
		return related
	}
	return FilterPostsByCoordType(related, relatedCoordType)
}

func typedGuideTopics(opts FeedOptions) []string {
	if strings.TrimSpace(opts.Topic) == "" {
		return nil
	}
	return []string{opts.Topic}
}
