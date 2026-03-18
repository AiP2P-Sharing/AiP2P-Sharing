package newsplugin

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func BuildFeedFacets(stats []FacetStat, opts FeedOptions, basePath, key string, omit ...string) []FeedFacet {
	items := make([]FeedFacet, 0, len(stats)+1)
	items = append(items, FeedFacet{
		Name:   "All",
		Count:  0,
		URL:    pageURL(basePath, opts, key, "", omit...),
		Active: activeFeedValue(opts, key) == "",
	})
	limit := len(stats)
	if limit > 8 {
		limit = 8
	}
	for _, stat := range stats[:limit] {
		items = append(items, FeedFacet{
			Name:   stat.Name,
			Count:  stat.Count,
			URL:    pageURL(basePath, opts, key, stat.Name, omit...),
			Active: strings.EqualFold(activeFeedValue(opts, key), stat.Name),
		})
	}
	return items
}

func BuildFacetLinks(stats []FacetStat, opts FeedOptions, basePath, key string, omit ...string) []FeedFacet {
	limit := len(stats)
	if limit > 8 {
		limit = 8
	}
	items := make([]FeedFacet, 0, limit)
	for _, stat := range stats[:limit] {
		items = append(items, FeedFacet{
			Name:   stat.Name,
			Count:  stat.Count,
			URL:    pageURL(basePath, opts, key, stat.Name, omit...),
			Active: strings.EqualFold(activeFeedValue(opts, key), stat.Name),
		})
	}
	return items
}

func BuildSortOptions(opts FeedOptions, basePath string, omit ...string) []SortOption {
	order := []struct {
		Name  string
		Value string
	}{
		{Name: "Newest", Value: "new"},
		{Name: "Discussed", Value: "discussed"},
		{Name: "Vote Score", Value: "score"},
		{Name: "Truth", Value: "truth"},
		{Name: "Source Quality", Value: "source"},
	}
	items := make([]SortOption, 0, len(order))
	active := opts.Sort
	if active == "" {
		active = "new"
	}
	for _, item := range order {
		items = append(items, SortOption{
			Name:   item.Name,
			Value:  item.Value,
			URL:    pageURL(basePath, opts, "sort", item.Value, omit...),
			Active: item.Value == active,
		})
	}
	return items
}

func BuildWindowOptions(opts FeedOptions, basePath string, omit ...string) []TimeWindowOption {
	order := []struct {
		Name  string
		Value string
	}{
		{Name: "All time", Value: ""},
		{Name: "24h", Value: "24h"},
		{Name: "7d", Value: "7d"},
		{Name: "30d", Value: "30d"},
	}
	active := canonicalWindow(opts.Window)
	items := make([]TimeWindowOption, 0, len(order))
	for _, item := range order {
		items = append(items, TimeWindowOption{
			Name:   item.Name,
			Value:  item.Value,
			URL:    pageURL(basePath, opts, "window", item.Value, omit...),
			Active: canonicalWindow(item.Value) == active,
		})
		if item.Value == "" && active == "" {
			items[len(items)-1].Active = true
		}
	}
	return items
}

func BuildPageSizeOptions(opts FeedOptions, basePath string, omit ...string) []PageSizeOption {
	sizes := []int{20, 50, 100}
	items := make([]PageSizeOption, 0, len(sizes))
	active := opts.PageSize
	if active == 0 {
		active = 20
	}
	for _, size := range sizes {
		items = append(items, PageSizeOption{
			Name:   strconv.Itoa(size),
			Value:  size,
			URL:    pageURL(basePath, opts, "page_size", strconv.Itoa(size), omit...),
			Active: size == active,
		})
	}
	return items
}

func BuildActiveFilters(opts FeedOptions, basePath string, omit ...string) []ActiveFilter {
	filters := make([]ActiveFilter, 0, 7)
	if opts.Query != "" {
		filters = append(filters, ActiveFilter{
			Label: "Search: " + opts.Query,
			URL:   pageURL(basePath, opts, "q", "", omit...),
		})
	}
	if opts.Window != "" {
		filters = append(filters, ActiveFilter{
			Label: "Window: " + strings.ToUpper(opts.Window),
			URL:   pageURL(basePath, opts, "window", "", omit...),
		})
	}
	if opts.Channel != "" {
		filters = append(filters, ActiveFilter{
			Label: "Channel: " + opts.Channel,
			URL:   pageURL(basePath, opts, "channel", "", omit...),
		})
	}
	if opts.Topic != "" && !contains(omit, "topic") {
		filters = append(filters, ActiveFilter{
			Label: "Topic: " + opts.Topic,
			URL:   pageURL(basePath, opts, "topic", "", omit...),
		})
	}
	if opts.Source != "" && !contains(omit, "source") {
		filters = append(filters, ActiveFilter{
			Label: "Source: " + opts.Source,
			URL:   pageURL(basePath, opts, "source", "", omit...),
		})
	}
	if opts.MetaKey != "" && opts.MetaValue != "" && !contains(omit, "meta_value") {
		filters = append(filters, ActiveFilter{
			Label: strings.TrimSpace(opts.MetaKey) + ": " + opts.MetaValue,
			URL:   pageURL(basePath, opts, "meta_value", "", append(omit, "meta_key")...),
		})
	}
	if opts.PageSize > 0 && opts.PageSize != 20 {
		filters = append(filters, ActiveFilter{
			Label: "Per page: " + strconv.Itoa(opts.PageSize),
			URL:   pageURL(basePath, opts, "page_size", "20", omit...),
		})
	}
	return filters
}

func BuildSummaryStats(posts []Post) []SummaryStat {
	return []SummaryStat{
		{Label: "Visible assets", Value: strconv.Itoa(len(posts))},
		{Label: "Replies", Value: strconv.Itoa(CountReplies(posts))},
		{Label: "Reactions", Value: strconv.Itoa(CountReactions(posts))},
		{Label: "Avg truth", Value: formatAverageTruth(posts)},
	}
}

func BuildCollectionSummaryStats(kind string, posts []Post) []SummaryStat {
	switch strings.ToLower(strings.TrimSpace(kind)) {
	case "source", "topic":
		return []SummaryStat{
			{Label: "Visible assets", Value: strconv.Itoa(len(posts))},
			{Label: "Discussion replies", Value: strconv.Itoa(CountReplies(posts))},
			{Label: "Signals", Value: strconv.Itoa(CountReactions(posts))},
			{Label: "Avg truth", Value: formatAverageTruth(posts)},
		}
	default:
		return BuildSummaryStats(posts)
	}
}

func PaginatePosts(posts []Post, opts FeedOptions, basePath string) ([]Post, PaginationState) {
	pageSize := opts.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}
	totalItems := len(posts)
	totalPages := 1
	if totalItems > 0 {
		totalPages = (totalItems + pageSize - 1) / pageSize
	}
	page := opts.Page
	if page < 1 {
		page = 1
	}
	if page > totalPages {
		page = totalPages
	}
	start := 0
	end := totalItems
	fromItem := 0
	toItem := 0
	if totalItems > 0 {
		start = (page - 1) * pageSize
		if start > totalItems {
			start = totalItems
		}
		end = start + pageSize
		if end > totalItems {
			end = totalItems
		}
		fromItem = start + 1
		toItem = end
	}
	currentOpts := opts
	currentOpts.Page = page
	currentOpts.PageSize = pageSize
	state := PaginationState{
		Page:       page,
		PageSize:   pageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
		FromItem:   fromItem,
		ToItem:     toItem,
	}
	if page > 1 {
		state.PrevURL = pageURL(basePath, currentOpts, "page", strconv.Itoa(page-1))
	}
	if page < totalPages {
		state.NextURL = pageURL(basePath, currentOpts, "page", strconv.Itoa(page+1))
	}
	startPage := page - 2
	if startPage < 1 {
		startPage = 1
	}
	endPage := startPage + 4
	if endPage > totalPages {
		endPage = totalPages
	}
	if endPage-startPage < 4 {
		startPage = endPage - 4
		if startPage < 1 {
			startPage = 1
		}
	}
	for p := startPage; p <= endPage; p++ {
		state.Links = append(state.Links, PaginationLink{
			Label:  strconv.Itoa(p),
			URL:    pageURL(basePath, currentOpts, "page", strconv.Itoa(p)),
			Active: p == page,
		})
	}
	return posts[start:end], state
}

func BuildDirectorySummaryStats(stats []FacetStat, posts []Post) []SummaryStat {
	return []SummaryStat{
		{Label: "Tracked groups", Value: strconv.Itoa(len(stats))},
		{Label: "Assets", Value: strconv.Itoa(len(posts))},
		{Label: "Replies", Value: strconv.Itoa(CountReplies(posts))},
		{Label: "Avg truth", Value: formatAverageTruth(posts)},
	}
}

func BuildDirectorySummaryStatsForKind(kind string, stats []FacetStat, posts []Post) []SummaryStat {
	switch strings.ToLower(strings.TrimSpace(kind)) {
	case "sources", "source":
		return []SummaryStat{
			{Label: "Tracked sources", Value: strconv.Itoa(len(stats))},
			{Label: "Indexed assets", Value: strconv.Itoa(len(posts))},
			{Label: "Discussion replies", Value: strconv.Itoa(CountReplies(posts))},
			{Label: "Avg truth", Value: formatAverageTruth(posts)},
		}
	case "topics", "topic":
		return []SummaryStat{
			{Label: "Tracked workstreams", Value: strconv.Itoa(len(stats))},
			{Label: "Indexed assets", Value: strconv.Itoa(len(posts))},
			{Label: "Discussion replies", Value: strconv.Itoa(CountReplies(posts))},
			{Label: "Avg truth", Value: formatAverageTruth(posts)},
		}
	default:
		return BuildDirectorySummaryStats(stats, posts)
	}
}

func BuildSourceDirectory(index Index) []DirectoryItem {
	items := make([]DirectoryItem, 0, len(index.SourceStats))
	for _, stat := range index.SourceStats {
		posts := index.FilterPosts(FeedOptions{Source: stat.Name, Now: time.Now()})
		items = append(items, DirectoryItem{
			Name:          stat.Name,
			URL:           SourcePath(stat.Name),
			ExternalURL:   SourceURLFromPosts(posts),
			AssetCount:    len(posts),
			ReplyCount:    CountReplies(posts),
			ReactionCount: CountReactions(posts),
			AvgTruth:      formatAverageTruth(posts),
		})
	}
	return items
}

func BuildTopicDirectory(index Index) []DirectoryItem {
	items := make([]DirectoryItem, 0, len(index.TopicStats))
	for _, stat := range index.TopicStats {
		posts := index.FilterPosts(FeedOptions{Topic: stat.Name, Now: time.Now()})
		items = append(items, DirectoryItem{
			Name:          stat.Name,
			URL:           TopicPath(stat.Name),
			AssetCount:    len(posts),
			ReplyCount:    CountReplies(posts),
			ReactionCount: CountReactions(posts),
			AvgTruth:      formatAverageTruth(posts),
		})
	}
	return items
}

func ChannelStatsForPosts(posts []Post) []FacetStat {
	counts := make(map[string]int)
	for _, post := range posts {
		if post.ChannelGroup == "" {
			continue
		}
		counts[post.ChannelGroup]++
	}
	return facetStats(counts)
}

func TopicStatsForPosts(posts []Post) []FacetStat {
	counts := make(map[string]int)
	for _, post := range posts {
		for _, topic := range post.Topics {
			counts[topic]++
		}
	}
	return facetStats(counts)
}

func SourceStatsForPosts(posts []Post) []FacetStat {
	counts := make(map[string]int)
	for _, post := range posts {
		if !post.HasSourcePage || post.SourceName == "" {
			continue
		}
		counts[post.SourceName]++
	}
	return facetStats(counts)
}

func SourceURLFromPosts(posts []Post) string {
	for _, post := range posts {
		if post.SourceURL != "" {
			return post.SourceURL
		}
	}
	return ""
}

func HasSource(index Index, name string) bool {
	for _, stat := range index.SourceStats {
		if strings.EqualFold(stat.Name, name) {
			return true
		}
	}
	return false
}

func HasTopic(index Index, name string) bool {
	for _, stat := range index.TopicStats {
		if strings.EqualFold(stat.Name, name) {
			return true
		}
	}
	return false
}

func PathValue(prefix, path string) string {
	if !strings.HasPrefix(path, prefix) {
		return ""
	}
	value := strings.TrimPrefix(path, prefix)
	if value == "" || strings.Contains(value, "/") {
		return ""
	}
	decoded, err := url.PathUnescape(value)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(decoded)
}

func SourcePath(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return ""
	}
	return "/sources/" + url.PathEscape(name)
}

func TopicPath(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return ""
	}
	return "/topics/" + url.PathEscape(name)
}

func APIOptions(opts FeedOptions) map[string]string {
	result := map[string]string{
		"channel":    opts.Channel,
		"topic":      opts.Topic,
		"source":     opts.Source,
		"meta_key":   opts.MetaKey,
		"meta_value": opts.MetaValue,
		"sort":       opts.Sort,
		"q":          opts.Query,
		"window":     canonicalWindow(opts.Window),
	}
	if opts.Page > 1 {
		result["page"] = strconv.Itoa(opts.Page)
	}
	if opts.PageSize > 0 {
		result["page_size"] = strconv.Itoa(opts.PageSize)
	}
	return result
}

func APIPosts(posts []Post) []map[string]any {
	out := make([]map[string]any, 0, len(posts))
	for _, post := range posts {
		out = append(out, APIPost(post, false))
	}
	return out
}

func APIPost(post Post, withBody bool) map[string]any {
	origin := apiOrigin(post.Message.Origin)
	payload := map[string]any{
		"infohash":             post.InfoHash,
		"magnet":               post.Magnet,
		"archive_md":           post.ArchiveMD,
		"title":                post.Message.Title,
		"author":               post.Message.Author,
		"origin":               origin,
		"origin_signed":        origin != nil,
		"delegation":           apiDelegation(post.Delegation),
		"shared_by_local_node": post.SharedByLocalNode,
		"created_at":           post.CreatedAt.Format(time.RFC3339),
		"channel":              post.Message.Channel,
		"channel_group":        post.ChannelGroup,
		"source_name":          post.SourceName,
		"source_site_name":     post.SourceSiteName,
		"source_url":           post.SourceURL,
		"origin_public_key":    post.OriginPublicKey,
		"topics":               post.Topics,
		"post_type":            post.PostType,
		"summary":              post.Summary,
		"reply_count":          post.ReplyCount,
		"reaction_count":       post.ReactionCount,
		"vote_score":           post.VoteScore,
		"truth_score":          scoreValue(post.TruthScoreAverage),
		"source_quality":       scoreValue(post.SourceScoreAverage),
		"thread_path":          "/posts/" + post.InfoHash,
		"coord_path":           CoordPath(post),
		"coord_api_path":       CoordAPIPath(post),
		"source_path":          sourcePathForPost(post),
		"latest_reaction":      post.LatestReactionAuthor,
		"event_time":           timeValue(post.EventTime),
		"topic_paths":          topicPaths(post.Topics),
		"message_tags":         post.Message.Tags,
		"message_protocol":     post.Message.Protocol,
		"coord_type":           PostCoordType(post),
	}
	if withBody {
		payload["body"] = post.Body
	}
	return payload
}

func APIReplies(replies []Reply) []map[string]any {
	out := make([]map[string]any, 0, len(replies))
	for _, reply := range replies {
		origin := apiOrigin(reply.Message.Origin)
		out = append(out, map[string]any{
			"infohash":             reply.InfoHash,
			"magnet":               reply.Magnet,
			"archive_md":           reply.ArchiveMD,
			"author":               reply.Message.Author,
			"origin":               origin,
			"origin_signed":        origin != nil,
			"delegation":           apiDelegation(reply.Delegation),
			"shared_by_local_node": reply.SharedByLocalNode,
			"created_at":           reply.CreatedAt.Format(time.RFC3339),
			"parent_hash":          reply.ParentInfoHash,
			"body":                 reply.Body,
		})
	}
	return out
}

func APIReactions(reactions []Reaction) []map[string]any {
	out := make([]map[string]any, 0, len(reactions))
	for _, reaction := range reactions {
		origin := apiOrigin(reaction.Message.Origin)
		out = append(out, map[string]any{
			"infohash":             reaction.InfoHash,
			"magnet":               reaction.Magnet,
			"archive_md":           reaction.ArchiveMD,
			"author":               reaction.Message.Author,
			"origin":               origin,
			"origin_signed":        origin != nil,
			"delegation":           apiDelegation(reaction.Delegation),
			"shared_by_local_node": reaction.SharedByLocalNode,
			"created_at":           reaction.CreatedAt.Format(time.RFC3339),
			"subject_hash":         reaction.SubjectInfoHash,
			"reaction_type":        reaction.ReactionType,
			"vote_value":           reaction.VoteValue,
			"score_value":          scoreValue(reaction.ScoreValue),
			"explanation":          reaction.Explanation,
		})
	}
	return out
}

func APICoordRelatedGroups(groups []CoordRelatedGroup) []map[string]any {
	out := make([]map[string]any, 0, len(groups))
	for _, group := range groups {
		out = append(out, map[string]any{
			"title": group.Title,
			"kind":  group.Kind,
			"path":  group.Path,
			"posts": APIPosts(group.Posts),
		})
	}
	return out
}

func APICoordRelationGroups(groups []CoordRelationGroup) []map[string]any {
	out := make([]map[string]any, 0, len(groups))
	for _, group := range groups {
		items := make([]map[string]any, 0, len(group.Items))
		for _, item := range group.Items {
			items = append(items, map[string]any{
				"post":     APIPost(item.Post, false),
				"evidence": append([]string(nil), item.Evidence...),
				"score":    item.Score,
			})
		}
		out = append(out, map[string]any{
			"key":             group.Key,
			"title":           group.Title,
			"description":     group.Description,
			"dominant_reason": group.DominantReason,
			"evidence":        append([]string(nil), group.Evidence...),
			"kind":            group.Kind,
			"path":            group.Path,
			"priority":        group.Priority,
			"score":           group.Score,
			"items":           items,
		})
	}
	return out
}

func APIWorkbenchClusters(clusters []WorkbenchCluster) []map[string]any {
	out := make([]map[string]any, 0, len(clusters))
	for _, cluster := range clusters {
		out = append(out, map[string]any{
			"task":       APIPost(cluster.Task, false),
			"related":    APICoordRelatedGroups(cluster.Related),
			"jump_panel": APIWorkspaceJumpPanel(cluster.JumpPanel),
			"navigation": APIWorkspaceNavigation(cluster.JumpPanel, map[string][]CoordAction{"registry_actions": cluster.RegistryActions}),
		})
	}
	return out
}

func APIWorkbenchLanes(lanes []WorkbenchLane) []map[string]any {
	out := make([]map[string]any, 0, len(lanes))
	for _, lane := range lanes {
		out = append(out, map[string]any{
			"title":       lane.Title,
			"kind":        lane.Kind,
			"description": lane.Description,
			"path":        lane.Path,
			"empty_copy":  lane.EmptyCopy,
			"posts":       APIPosts(lane.Posts),
		})
	}
	return out
}

func APIWorkstreamRelationPaths(paths []WorkstreamRelationPath) []map[string]any {
	out := make([]map[string]any, 0, len(paths))
	for _, path := range paths {
		out = append(out, map[string]any{
			"from":           APIPost(path.From, false),
			"relation_key":   path.RelationKey,
			"relation_title": path.RelationTitle,
			"kind":           path.Kind,
			"path":           path.Path,
			"target":         APIPost(path.Target, false),
			"reason":         path.Reason,
			"score":          path.Score,
		})
	}
	return out
}

func APIWorkstreamRelationPanel(panel *WorkstreamRelationPanel) map[string]any {
	if panel == nil {
		return nil
	}
	return map[string]any{
		"title":            panel.Title,
		"description":      panel.Description,
		"primary_paths":    APIWorkstreamRelationPaths(panel.PrimaryPaths),
		"supporting_paths": APIWorkstreamRelationPaths(panel.SupportingPaths),
	}
}

func APIWorkspaceJumpPanel(panel *WorkspaceJumpPanel) map[string]any {
	if panel == nil {
		return nil
	}
	return map[string]any{
		"eyebrow":            panel.Eyebrow,
		"title":              panel.Title,
		"description":        panel.Description,
		"primary_actions":    append([]CoordAction(nil), panel.PrimaryActions...),
		"supporting_actions": append([]CoordAction(nil), panel.SupportingActions...),
	}
}

func APIWorkspaceNavigation(panel *WorkspaceJumpPanel, aliases map[string][]CoordAction) map[string]any {
	if panel == nil && len(aliases) == 0 {
		return nil
	}
	payload := map[string]any{
		"canonical_field": "jump_panel",
	}
	if panel != nil {
		payload["primary_actions"] = append([]CoordAction(nil), panel.PrimaryActions...)
		payload["supporting_actions"] = append([]CoordAction(nil), panel.SupportingActions...)
	}
	if len(aliases) > 0 {
		legacy := make(map[string]any, len(aliases))
		for key, actions := range aliases {
			if len(actions) == 0 {
				continue
			}
			legacy[key] = append([]CoordAction(nil), actions...)
		}
		if len(legacy) > 0 {
			payload["legacy_fields"] = legacy
		}
	}
	return payload
}

func APICoordMetadataSchema(spec *CoordMetadataTypeSpec) map[string]any {
	if spec == nil {
		return nil
	}
	fields := make([]map[string]any, 0, len(spec.Fields))
	for _, field := range spec.Fields {
		fields = append(fields, map[string]any{
			"key":   field.Key,
			"label": field.Label,
		})
	}
	sections := make([]map[string]any, 0, len(spec.Sections))
	for _, section := range spec.Sections {
		sectionFields := make([]map[string]any, 0, len(section.Fields))
		for _, field := range section.Fields {
			sectionFields = append(sectionFields, map[string]any{
				"key":   field.Key,
				"label": field.Label,
			})
		}
		sections = append(sections, map[string]any{
			"title":       section.Title,
			"description": section.Description,
			"fields":      sectionFields,
		})
	}
	return map[string]any{
		"coord_type": spec.CoordType,
		"fields":     fields,
		"sections":   sections,
	}
}

func APICoordPublishGuide(guide *CoordPublishGuide) map[string]any {
	if guide == nil {
		return nil
	}
	fields := make([]map[string]any, 0, len(guide.StarterFields))
	for _, field := range guide.StarterFields {
		fields = append(fields, map[string]any{
			"key":   field.Key,
			"label": field.Label,
		})
	}
	templates := make([]map[string]any, 0, len(guide.Templates))
	for _, template := range guide.Templates {
		templates = append(templates, map[string]any{
			"category":                template.Category,
			"label":                   template.Label,
			"description":             template.Description,
			"execution_mode":          template.ExecutionMode,
			"target_summary":          template.TargetSummary,
			"reference_target":        template.ReferenceTarget,
			"publisher_role":          template.PublisherRole,
			"identity_file":           template.IdentityFile,
			"store_path":              template.StorePath,
			"channel":                 template.Channel,
			"kind":                    template.Kind,
			"prerequisites":           append([]string(nil), template.Prerequisites...),
			"starter_extensions_json": template.StarterExtensionsJSON,
			"starter_command":         template.StarterCommand,
		})
	}
	workflow := make([]map[string]any, 0, len(guide.Workflow))
	for _, step := range guide.Workflow {
		stepTemplates := make([]map[string]any, 0, len(step.Templates))
		for _, template := range step.Templates {
			stepTemplates = append(stepTemplates, map[string]any{
				"category":                template.Category,
				"label":                   template.Label,
				"description":             template.Description,
				"execution_mode":          template.ExecutionMode,
				"target_summary":          template.TargetSummary,
				"reference_target":        template.ReferenceTarget,
				"publisher_role":          template.PublisherRole,
				"identity_file":           template.IdentityFile,
				"store_path":              template.StorePath,
				"channel":                 template.Channel,
				"kind":                    template.Kind,
				"prerequisites":           append([]string(nil), template.Prerequisites...),
				"starter_extensions_json": template.StarterExtensionsJSON,
				"starter_command":         template.StarterCommand,
			})
		}
		workflow = append(workflow, map[string]any{
			"title":           step.Title,
			"description":     step.Description,
			"template_labels": append([]string(nil), step.TemplateLabels...),
			"templates":       stepTemplates,
		})
	}
	return map[string]any{
		"coord_type":              guide.CoordType,
		"title":                   guide.Title,
		"description":             guide.Description,
		"execution_mode":          guide.ExecutionMode,
		"target_summary":          guide.TargetSummary,
		"publisher_role":          guide.PublisherRole,
		"identity_file":           guide.IdentityFile,
		"store_path":              guide.StorePath,
		"channel":                 guide.Channel,
		"kind":                    guide.Kind,
		"prerequisites":           append([]string(nil), guide.Prerequisites...),
		"context_notes":           append([]string(nil), guide.ContextNotes...),
		"starter_extensions_json": guide.StarterExtensionsJSON,
		"starter_command":         guide.StarterCommand,
		"starter_fields":          fields,
		"templates":               templates,
		"workflow":                workflow,
	}
}

func APIFacets(facets []FeedFacet) []map[string]any {
	out := make([]map[string]any, 0, len(facets))
	for _, facet := range facets {
		out = append(out, map[string]any{
			"name":  facet.Name,
			"count": facet.Count,
			"url":   facet.URL,
		})
	}
	return out
}

func APICoordWorkbench(panel *CoordWorkbenchPanel) map[string]any {
	if panel == nil {
		return nil
	}
	return map[string]any{
		"eyebrow":     panel.Eyebrow,
		"title":       panel.Title,
		"description": panel.Description,
		"groups":      APICoordRelatedGroups(panel.Groups),
	}
}

func APICoordDetailFocus(panel *CoordDetailFocusPanel) map[string]any {
	if panel == nil {
		return nil
	}
	return map[string]any{
		"eyebrow":     panel.Eyebrow,
		"title":       panel.Title,
		"description": panel.Description,
		"checklist":   append([]string(nil), panel.Checklist...),
		"actions":     append([]CoordAction(nil), panel.Actions...),
	}
}

func APIAgentWorkbench(cards []AgentWorkbenchCard) []map[string]any {
	out := make([]map[string]any, 0, len(cards))
	for _, card := range cards {
		out = append(out, map[string]any{
			"agent":           APIPost(card.Agent, false),
			"workspace_panel": APICoordWorkbench(card.Panel),
		})
	}
	return out
}

func APIRegistryWorkbench(cards []RegistryWorkbenchCard) []map[string]any {
	out := make([]map[string]any, 0, len(cards))
	for _, card := range cards {
		out = append(out, map[string]any{
			"kind":           card.Kind,
			"name":           card.Name,
			"url":            card.URL,
			"external_url":   card.ExternalURL,
			"summary":        card.Summary,
			"module_actions": card.ModuleActions,
			"jump_panel":     APIWorkspaceJumpPanel(card.JumpPanel),
			"navigation":     APIWorkspaceNavigation(card.JumpPanel, map[string][]CoordAction{"typed_actions": card.ModuleActions}),
			"primary_count":  card.PrimaryCount,
			"reply_count":    card.ReplyCount,
			"reaction_count": card.ReactionCount,
		})
	}
	return out
}

func APIHomePriorityCards(cards []HomePriorityCard) []map[string]any {
	out := make([]map[string]any, 0, len(cards))
	for _, card := range cards {
		out = append(out, map[string]any{
			"eyebrow":     card.Eyebrow,
			"title":       card.Title,
			"description": card.Description,
			"url":         card.URL,
			"stats":       card.Stats,
			"actions":     card.Actions,
		})
	}
	return out
}

func pageURL(basePath string, opts FeedOptions, key, value string, omit ...string) string {
	next := withOption(opts, key, value)
	encoded := encodeOptions(next, omit...)
	if encoded == "" {
		return basePath
	}
	return basePath + "?" + encoded
}

func withOption(opts FeedOptions, key, value string) FeedOptions {
	next := opts
	switch key {
	case "channel":
		next.Channel = value
	case "topic":
		next.Topic = value
	case "source":
		next.Source = value
	case "meta_key":
		next.MetaKey = strings.TrimSpace(value)
		if next.MetaKey == "" {
			next.MetaValue = ""
		}
	case "meta_value":
		next.MetaValue = strings.TrimSpace(value)
		if next.MetaValue == "" {
			next.MetaKey = ""
		}
	case "sort":
		next.Sort = value
	case "q":
		next.Query = value
	case "window":
		next.Window = canonicalWindow(value)
	case "page":
		next.Page = parsePositiveInt(value, 1)
	case "page_size":
		next.PageSize = parseFeedPageSize(value)
	}
	if key != "page" {
		next.Page = 1
	}
	return next
}

func encodeOptions(opts FeedOptions, omit ...string) string {
	query := url.Values{}
	ignored := make(map[string]struct{}, len(omit))
	for _, key := range omit {
		ignored[key] = struct{}{}
	}
	set := func(key, value string) {
		if value == "" {
			return
		}
		if _, skip := ignored[key]; skip {
			return
		}
		query.Set(key, value)
	}
	set("channel", opts.Channel)
	set("topic", opts.Topic)
	set("source", opts.Source)
	if strings.TrimSpace(opts.MetaKey) != "" && strings.TrimSpace(opts.MetaValue) != "" {
		set("meta_key", opts.MetaKey)
		set("meta_value", opts.MetaValue)
	}
	if opts.Sort != "" && opts.Sort != "new" {
		set("sort", opts.Sort)
	}
	set("q", opts.Query)
	if window := canonicalWindow(opts.Window); window != "" {
		set("window", window)
	}
	if opts.Page > 1 {
		query.Set("page", strconv.Itoa(opts.Page))
	}
	if opts.PageSize > 0 && opts.PageSize != 20 {
		query.Set("page_size", strconv.Itoa(opts.PageSize))
	}
	return query.Encode()
}

func activeFeedValue(opts FeedOptions, key string) string {
	switch key {
	case "channel":
		return opts.Channel
	case "topic":
		return opts.Topic
	case "source":
		return opts.Source
	case "meta_key":
		return opts.MetaKey
	case "meta_value":
		return opts.MetaValue
	case "window":
		return canonicalWindow(opts.Window)
	default:
		return ""
	}
}

func compactIdentity(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	if isPublicKeyish(value) {
		if len(value) <= 8 {
			return value
		}
		return value[:8] + "..."
	}
	if len(value) <= 24 {
		return value
	}
	return value[:24] + "..."
}

func isPublicKeyish(value string) bool {
	value = strings.TrimSpace(value)
	if len(value) < 32 {
		return false
	}
	for _, r := range value {
		if (r >= '0' && r <= '9') || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F') {
			continue
		}
		return false
	}
	return true
}

func parsePositiveInt(raw string, fallback int) int {
	value, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil || value < 1 {
		return fallback
	}
	return value
}

func parseFeedPageSize(raw string) int {
	value := parsePositiveInt(raw, 20)
	if value < 1 {
		return 20
	}
	if value > 200 {
		return 200
	}
	return value
}

func AppendRawQuery(path, rawQuery string) string {
	path = strings.TrimSpace(path)
	rawQuery = strings.TrimSpace(rawQuery)
	if path == "" || rawQuery == "" {
		return path
	}
	if strings.Contains(path, "?") {
		return path + "&" + rawQuery
	}
	return path + "?" + rawQuery
}

func ShouldShowNetworkWarning(r *http.Request) bool {
	if r == nil {
		return true
	}
	cookie, err := r.Cookie("aip2p_news_network_warning_seen")
	if err != nil {
		return true
	}
	return strings.TrimSpace(cookie.Value) == ""
}

func IsAgentViewer(r *http.Request) bool {
	if r == nil {
		return false
	}
	if value := strings.TrimSpace(r.URL.Query().Get("agent")); value != "" {
		switch strings.ToLower(value) {
		case "1", "true", "yes", "on":
			return true
		case "0", "false", "no", "off":
			return false
		}
	}
	ua := strings.ToLower(strings.TrimSpace(r.UserAgent()))
	if ua == "" {
		return false
	}
	if strings.Contains(ua, "mozilla/") && !strings.Contains(ua, "bot") && !strings.Contains(ua, "agent") {
		return false
	}
	markers := []string{"agent", "bot", "crawler", "python", "curl", "wget", "httpie", "go-http-client", "openai", "anthropic", "claude", "gpt", "llm"}
	for _, marker := range markers {
		if strings.Contains(ua, marker) {
			return true
		}
	}
	return false
}

func formatAverageTruth(posts []Post) string {
	var sum float64
	var count int
	for _, post := range posts {
		if post.TruthScoreAverage == nil {
			continue
		}
		sum += *post.TruthScoreAverage
		count++
	}
	if count == 0 {
		return "-"
	}
	return strings.TrimRight(strings.TrimRight(strconv.FormatFloat(sum/float64(count), 'f', 2, 64), "0"), ".")
}

func contains(items []string, target string) bool {
	for _, item := range items {
		if item == target {
			return true
		}
	}
	return false
}

func sourcePathForPost(post Post) string {
	if !post.HasSourcePage || strings.TrimSpace(post.SourceName) == "" {
		return ""
	}
	return SourcePath(post.SourceName)
}

func apiOrigin(origin *MessageOrigin) map[string]any {
	if origin == nil {
		return nil
	}
	return map[string]any{
		"author":     origin.Author,
		"agent_id":   origin.AgentID,
		"key_type":   origin.KeyType,
		"public_key": origin.PublicKey,
		"signature":  origin.Signature,
	}
}

func apiDelegation(info *DelegationInfo) map[string]any {
	if info == nil || !info.Delegated {
		return nil
	}
	return map[string]any{
		"delegated":         true,
		"parent_agent_id":   info.ParentAgentID,
		"parent_key_type":   info.ParentKeyType,
		"parent_public_key": info.ParentPublicKey,
		"scopes":            append([]string(nil), info.Scopes...),
		"created_at":        info.CreatedAt,
		"expires_at":        info.ExpiresAt,
	}
}

func topicPaths(topics []string) map[string]string {
	out := make(map[string]string, len(topics))
	for _, topic := range topics {
		out[topic] = TopicPath(topic)
	}
	return out
}

func scoreValue(value *float64) any {
	if value == nil {
		return nil
	}
	return *value
}

func timeValue(value *time.Time) any {
	if value == nil {
		return nil
	}
	return value.Format(time.RFC3339)
}
