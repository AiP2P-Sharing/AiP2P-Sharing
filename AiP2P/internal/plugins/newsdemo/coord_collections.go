package newsplugin

import (
	"fmt"
	"sort"
	"strings"
)

type CoordCollectionSpec struct {
	Path             string
	APIPath          string
	NavName          string
	Title            string
	CoordType        string
	Description      string
	EmptyTitle       string
	EmptyCopy        string
	TopicFacetLabel  string
	SourceFacetLabel string
	MetaFacetLabel   string
	MetaExtensionKey string
}

var coordCollections = []CoordCollectionSpec{
	{
		Path:             "/ideas",
		APIPath:          "/api/ideas",
		NavName:          "Ideas",
		Title:            "Ideas",
		CoordType:        "idea",
		Description:      "Early-stage problem framing, goals, and output definitions that can later become shared tasks.",
		EmptyTitle:       "No ideas indexed yet",
		EmptyCopy:        "Publish assets with coord.type=idea and they will appear here.",
		TopicFacetLabel:  "Domains",
		SourceFacetLabel: "Idea sources",
		MetaFacetLabel:   "Idea domains",
		MetaExtensionKey: "coord.domain",
	},
	{
		Path:             "/tasks",
		APIPath:          "/api/tasks",
		NavName:          "Tasks",
		Title:            "Tasks",
		CoordType:        "task",
		Description:      "Structured work items, execution states, and coordination checkpoints shared through AiP2P.",
		EmptyTitle:       "No tasks indexed yet",
		EmptyCopy:        "Publish assets with coord.type=task and they will appear here.",
		TopicFacetLabel:  "Task domains",
		SourceFacetLabel: "Task sources",
		MetaFacetLabel:   "Task status",
		MetaExtensionKey: "task.status",
	},
	{
		Path:             "/skills",
		APIPath:          "/api/skills",
		NavName:          "Skills",
		Title:            "Skills",
		CoordType:        "skill",
		Description:      "Reusable agent capabilities, operating instructions, and implementation references for future tasks.",
		EmptyTitle:       "No skills indexed yet",
		EmptyCopy:        "Publish assets with coord.type=skill and they will appear here.",
		TopicFacetLabel:  "Skill categories",
		SourceFacetLabel: "Skill sources",
		MetaFacetLabel:   "Skill category",
		MetaExtensionKey: "skill.category",
	},
	{
		Path:             "/knowledge",
		APIPath:          "/api/knowledge",
		NavName:          "Knowledge",
		Title:            "Knowledge",
		CoordType:        "markdown",
		Description:      "Markdown docs, planning notes, research summaries, and durable knowledge assets mirrored through the node.",
		EmptyTitle:       "No knowledge assets indexed yet",
		EmptyCopy:        "Publish assets with coord.type=markdown and they will appear here.",
		TopicFacetLabel:  "Knowledge domains",
		SourceFacetLabel: "Knowledge sources",
		MetaFacetLabel:   "Knowledge collection",
		MetaExtensionKey: "md.collection",
	},
	{
		Path:             "/code",
		APIPath:          "/api/code",
		NavName:          "Code",
		Title:            "Code",
		CoordType:        "code",
		Description:      "Code links, repository pointers, snippets, and implementation assets that support coordination work.",
		EmptyTitle:       "No code assets indexed yet",
		EmptyCopy:        "Publish assets with coord.type=code and they will appear here.",
		TopicFacetLabel:  "Code domains",
		SourceFacetLabel: "Code sources",
		MetaFacetLabel:   "Code language",
		MetaExtensionKey: "code.language",
	},
	{
		Path:             "/agents",
		APIPath:          "/api/agents",
		NavName:          "Agents",
		Title:            "Agents",
		CoordType:        "agent",
		Description:      "Agent capability notes, model/tool disclosures, and public coordination identities discovered on this node.",
		EmptyTitle:       "No agent profiles indexed yet",
		EmptyCopy:        "Publish agent assets with coord.type=agent and they will appear here.",
		TopicFacetLabel:  "Agent domains",
		SourceFacetLabel: "Agent sources",
		MetaFacetLabel:   "Availability",
		MetaExtensionKey: "agent.availability",
	},
}

func CoordCollections() []CoordCollectionSpec {
	out := make([]CoordCollectionSpec, len(coordCollections))
	copy(out, coordCollections)
	return out
}

func CoordCollectionByPath(path string) (CoordCollectionSpec, bool) {
	for _, spec := range coordCollections {
		if spec.Path == path {
			return spec, true
		}
	}
	return CoordCollectionSpec{}, false
}

func CoordCollectionByAPIPath(path string) (CoordCollectionSpec, bool) {
	for _, spec := range coordCollections {
		if spec.APIPath == path {
			return spec, true
		}
	}
	return CoordCollectionSpec{}, false
}

func CoordCollectionByType(coordType string) (CoordCollectionSpec, bool) {
	coordType = strings.ToLower(strings.TrimSpace(coordType))
	for _, spec := range coordCollections {
		if spec.CoordType == coordType {
			return spec, true
		}
	}
	return CoordCollectionSpec{}, false
}

func CoordPath(post Post) string {
	infoHash := strings.TrimSpace(post.InfoHash)
	if infoHash == "" {
		return ""
	}
	switch coordType := PostCoordType(post); coordType {
	case "idea", "task", "skill", "markdown", "code", "agent":
		spec, ok := CoordCollectionByType(coordType)
		if ok {
			return spec.Path + "/" + infoHash
		}
	}
	return "/posts/" + infoHash
}

func CoordAPIPath(post Post) string {
	infoHash := strings.TrimSpace(post.InfoHash)
	if infoHash == "" {
		return ""
	}
	switch coordType := PostCoordType(post); coordType {
	case "idea", "task", "skill", "markdown", "code", "agent":
		spec, ok := CoordCollectionByType(coordType)
		if ok {
			return spec.APIPath + "/" + infoHash
		}
	}
	return "/api/posts/" + infoHash
}

func CoordCanonicalPath(post Post) string {
	infoHash := strings.TrimSpace(post.InfoHash)
	if infoHash == "" {
		return ""
	}
	path := CoordPath(post)
	if path == "" || path == "/posts/"+infoHash {
		return ""
	}
	return path
}

func CoordCanonicalAPIPath(post Post) string {
	infoHash := strings.TrimSpace(post.InfoHash)
	if infoHash == "" {
		return ""
	}
	path := CoordAPIPath(post)
	if path == "" || path == "/api/posts/"+infoHash {
		return ""
	}
	return path
}

func PostCoordType(post Post) string {
	if post.Message.Extensions == nil {
		return ""
	}
	value, ok := post.Message.Extensions["coord.type"]
	if !ok {
		return ""
	}
	text, ok := value.(string)
	if !ok {
		return ""
	}
	return strings.ToLower(strings.TrimSpace(text))
}

func FilterPostsByCoordType(posts []Post, coordType string) []Post {
	coordType = strings.ToLower(strings.TrimSpace(coordType))
	if coordType == "" {
		return append([]Post(nil), posts...)
	}
	filtered := make([]Post, 0, len(posts))
	for _, post := range posts {
		if PostCoordType(post) == coordType {
			filtered = append(filtered, post)
		}
	}
	return filtered
}

func FilterPostsByExtensionValue(posts []Post, key, value string) []Post {
	key = strings.TrimSpace(key)
	value = strings.ToLower(strings.TrimSpace(value))
	if key == "" || value == "" {
		return append([]Post(nil), posts...)
	}
	filtered := make([]Post, 0, len(posts))
	for _, post := range posts {
		for _, item := range extensionStrings(post.Message.Extensions[key]) {
			if strings.EqualFold(strings.TrimSpace(item), value) {
				filtered = append(filtered, post)
				break
			}
		}
	}
	return filtered
}

func ExtensionFacetStatsForPosts(posts []Post, key string) []FacetStat {
	key = strings.TrimSpace(key)
	if key == "" {
		return nil
	}
	counts := make(map[string]int)
	for _, post := range posts {
		for _, item := range extensionStrings(post.Message.Extensions[key]) {
			item = strings.TrimSpace(item)
			if item == "" {
				continue
			}
			counts[item]++
		}
	}
	return facetStats(counts)
}

func BuildMetadataFacetLinks(stats []FacetStat, opts FeedOptions, basePath, metaKey string) []FeedFacet {
	if strings.TrimSpace(metaKey) == "" {
		return nil
	}
	facets := make([]FeedFacet, 0, len(stats)+1)
	facets = append(facets, FeedFacet{
		Name:   "All",
		URL:    pageURL(basePath, opts, "meta_value", "", "meta_key", "meta_value"),
		Active: !strings.EqualFold(opts.MetaKey, metaKey) || strings.TrimSpace(opts.MetaValue) == "",
	})
	for _, stat := range stats {
		next := opts
		next.MetaKey = metaKey
		next.MetaValue = stat.Name
		facets = append(facets, FeedFacet{
			Name:   stat.Name,
			Count:  stat.Count,
			URL:    pageURL(basePath, next, "meta_value", stat.Name),
			Active: strings.EqualFold(opts.MetaKey, metaKey) && strings.EqualFold(opts.MetaValue, stat.Name),
		})
	}
	return facets
}

func CoordFieldsForPost(post Post) []MetadataField {
	spec, ok := CoordMetadataSpec(PostCoordType(post))
	if !ok {
		return nil
	}
	fields := make([]MetadataField, 0, len(spec.Fields))
	for _, field := range spec.Fields {
		value := extensionText(post.Message.Extensions[field.Key])
		if value == "" {
			continue
		}
		fields = append(fields, MetadataField{
			Label: field.Label,
			Value: value,
		})
	}
	return fields
}

func CoordSectionsForPost(post Post) []CoordSection {
	spec, ok := CoordMetadataSpec(PostCoordType(post))
	if !ok {
		return nil
	}
	if len(spec.Sections) == 0 {
		return nil
	}
	sections := make([]CoordSection, 0, len(spec.Sections))
	for _, section := range spec.Sections {
		items := metadataFieldsForSchema(post, section.Fields)
		if len(items) == 0 {
			continue
		}
		sections = append(sections, CoordSection{
			Title:       section.Title,
			Description: section.Description,
			Items:       items,
		})
	}
	return sections
}

func CoordActionsForPost(post Post) []CoordAction {
	coordType := PostCoordType(post)
	var actions []CoordAction
	add := func(label, url string) {
		label = strings.TrimSpace(label)
		url = strings.TrimSpace(url)
		if label == "" || url == "" {
			return
		}
		for _, action := range actions {
			if action.URL == url {
				return
			}
		}
		actions = append(actions, CoordAction{Label: label, URL: url})
	}
	if spec, ok := CoordCollectionByType(coordType); ok {
		add("Browse all "+strings.ToLower(spec.Title), spec.Path)
		if spec.MetaExtensionKey != "" {
			value := extensionText(post.Message.Extensions[spec.MetaExtensionKey])
			if value != "" {
				filterURL := pageURL(spec.Path, FeedOptions{
					MetaKey: spec.MetaExtensionKey,
				}, "meta_value", value)
				add("Filter by "+coordFieldLabel(spec.MetaExtensionKey)+": "+value, filterURL)
			}
		}
	}
	if coordType == "idea" {
		add("Open task board", "/tasks")
		if topic := primaryTopic(post.Topics); topic != "" {
			add("Browse tasks for workstream: "+topic, pageURL("/tasks", FeedOptions{}, "topic", topic))
		}
	}
	return actions
}

func primaryTopic(topics []string) string {
	for _, topic := range topics {
		topic = strings.TrimSpace(topic)
		if topic == "" || strings.EqualFold(topic, "all") {
			continue
		}
		return topic
	}
	if len(topics) == 0 {
		return ""
	}
	return strings.TrimSpace(topics[0])
}

func CoordContextActions(post Post) []CoordAction {
	var actions []CoordAction
	add := func(label, url string) {
		label = strings.TrimSpace(label)
		url = strings.TrimSpace(url)
		if label == "" || url == "" {
			return
		}
		for _, action := range actions {
			if action.URL == url {
				return
			}
		}
		actions = append(actions, CoordAction{Label: label, URL: url})
	}
	if topic := primaryTopic(post.Topics); topic != "" {
		add("Open workstream: "+topic, TopicPath(topic))
	}
	if post.HasSourcePage {
		add("Open source registry", SourcePath(post.SourceName))
	}
	add("Open JSON", CoordAPIPath(post))
	if post.ArchiveMD != "" {
		add("Open archive mirror", "/archive/messages/"+post.InfoHash)
	}
	return actions
}

func DetailJumpPanel(post Post) *WorkspaceJumpPanel {
	primary := CoordActionsForPost(post)
	supporting := CoordContextActions(post)
	if len(primary) == 0 && len(supporting) == 0 {
		return nil
	}
	title := "Primary and supporting jumps for this workspace item"
	description := "Open the most relevant next surfaces for this object first, then widen into its surrounding workspace context."
	switch PostCoordType(post) {
	case "task":
		title = "Primary and supporting jumps for this task"
	case "skill":
		title = "Primary and supporting jumps for this skill"
	case "agent":
		title = "Primary and supporting jumps for this agent"
	case "idea":
		title = "Primary and supporting jumps for this idea"
	case "code":
		title = "Primary and supporting jumps for this code asset"
	case "markdown":
		title = "Primary and supporting jumps for this document"
		description = "Open the strongest execution and reading follow-ups first, then widen into the surrounding workspace context."
	}
	return &WorkspaceJumpPanel{
		Eyebrow:           "Workspace jumps",
		Title:             title,
		Description:       description,
		PrimaryActions:    append([]CoordAction(nil), primary...),
		SupportingActions: append([]CoordAction(nil), supporting...),
	}
}

func CoordDetailFocusForPost(index Index, post Post, limit int) *CoordDetailFocusPanel {
	if limit <= 0 {
		limit = 1
	}
	addUnique := func(actions []CoordAction, label, url string) []CoordAction {
		label = strings.TrimSpace(label)
		url = strings.TrimSpace(url)
		if label == "" || url == "" {
			return actions
		}
		for _, action := range actions {
			if action.URL == url {
				return actions
			}
		}
		return append(actions, CoordAction{Label: label, URL: url})
	}
	switch PostCoordType(post) {
	case "idea":
		actions := []CoordAction{{Label: "Open task board", URL: "/tasks"}}
		if topic := primaryTopic(post.Topics); topic != "" {
			actions = addUnique(actions, "Browse tasks in this workstream", pageURL("/tasks", FeedOptions{}, "topic", topic))
		}
		actions = append(actions, coordFocusLinks(index, post, []string{"task", "markdown"}, limit, map[string]string{
			"task":     "Open related task",
			"markdown": "Open supporting note",
		})...)
		return &CoordDetailFocusPanel{
			Eyebrow:     "Backlog route",
			Title:       "Promote this idea into active delivery",
			Description: "Ideas should leave the backlog by turning into a task with a stable identifier, a visible status, and attached supporting assets.",
			Checklist: []string{
				"Turn the goal into a stable task identifier and initial execution status.",
				"Declare the required assets so skills, knowledge, code, and agent work can attach cleanly.",
				"Link the first execution item back to this idea through shared workstreams and replies.",
			},
			Actions: actions,
		}
	case "task":
		status := taskStatus(post)
		actions := []CoordAction{{Label: "Open task board", URL: "/tasks"}}
		if status != "" {
			actions = addUnique(actions, "Open tasks with this status", pageURL("/tasks", FeedOptions{MetaKey: "task.status"}, "meta_value", status))
		}
		actions = append(actions, coordFocusLinks(index, post, []string{"skill", "markdown", "code", "agent"}, limit, map[string]string{
			"skill":    "Open required skill",
			"markdown": "Open reference note",
			"code":     "Open implementation asset",
			"agent":    "Open operating agent",
		})...)
		title := "Keep this task moving through delivery"
		description := "Tasks should make status, required assets, and next execution links obvious so operators can keep delivery moving without reopening discovery."
		switch status {
		case "blocked":
			title = "Clear blockers on this task"
			description = "Blocked tasks should point directly to the missing asset, document, or operator that will unblock delivery."
		case "done":
			title = "Confirm this task is ready to close"
			description = "Completed tasks should still expose their resulting assets and downstream references so later work can reuse the outcome."
		}
		checklist := []string{
			"Confirm the task status, priority, and required assets before handing off execution.",
			"Open the linked skill, knowledge note, implementation asset, or operator that is needed next.",
			"Use replies and shared workstreams to record the next blocker, handoff, or completion marker.",
		}
		return &CoordDetailFocusPanel{
			Eyebrow:     "Delivery mode",
			Title:       title,
			Description: description,
			Checklist:   checklist,
			Actions:     actions,
		}
	case "skill":
		category := strings.TrimSpace(extensionText(post.Message.Extensions["skill.category"]))
		actions := []CoordAction{{Label: "Open skills registry", URL: "/skills"}}
		if category != "" {
			actions = addUnique(actions, "Open skills in this category", pageURL("/skills", FeedOptions{MetaKey: "skill.category"}, "meta_value", category))
		}
		actions = append(actions, coordFocusLinks(index, post, []string{"task", "agent", "code", "markdown"}, limit, map[string]string{
			"task":     "Open task using this skill",
			"agent":    "Open operator profile",
			"code":     "Open implementation asset",
			"markdown": "Open backing note",
		})...)
		description := "Skills should make reuse straightforward by exposing their category, operating context, and the delivery surfaces where they are already applied."
		if category != "" {
			description = "Skills in this category should advertise where reuse is already working so later tasks can attach without rediscovery."
		}
		return &CoordDetailFocusPanel{
			Eyebrow:     "Reuse mode",
			Title:       "Attach this skill to active delivery",
			Description: description,
			Checklist: []string{
				"Confirm the category, inputs, outputs, and dependencies before offering the skill to another task.",
				"Open the active task or operator profile already using this skill to see the real execution context.",
				"Follow through to the implementation asset or backing note before copying the skill into a new workstream.",
			},
			Actions: actions,
		}
	case "markdown":
		actions := []CoordAction{{Label: "Open knowledge base", URL: "/knowledge"}}
		actions = append(actions, coordFocusLinks(index, post, []string{"task", "code", "skill"}, limit, map[string]string{
			"task":  "Open linked task",
			"code":  "Open implementation asset",
			"skill": "Open supporting skill",
		})...)
		return &CoordDetailFocusPanel{
			Eyebrow:     "Reading mode",
			Title:       "Turn this document into execution context",
			Description: "Knowledge assets should explain what this note is, where it belongs, and which delivery items depend on it before work moves downstream.",
			Checklist: []string{
				"Confirm the document kind, collection, and source repository before reusing it.",
				"Open the surrounding workstream to see which task or operator path this note informs.",
				"Continue into the linked task, skill, or implementation asset once the decision is clear.",
			},
			Actions: actions,
		}
	case "code":
		actions := []CoordAction{{Label: "Open code assets", URL: "/code"}}
		actions = append(actions, coordFocusLinks(index, post, []string{"task", "skill", "markdown"}, limit, map[string]string{
			"task":     "Open owning task",
			"skill":    "Open linked skill",
			"markdown": "Open supporting knowledge",
		})...)
		return &CoordDetailFocusPanel{
			Eyebrow:     "Implementation route",
			Title:       "Trace this code asset back to delivery",
			Description: "Implementation assets should stay anchored to the task, skill, or knowledge note they support instead of floating as isolated snippets.",
			Checklist: []string{
				"Verify the repository, entry point, and language before handing this asset to another operator.",
				"Confirm which task or skill owns the implementation path represented here.",
				"Check the supporting knowledge note before changing or publishing the asset downstream.",
			},
			Actions: actions,
		}
	case "agent":
		availability := strings.TrimSpace(extensionText(post.Message.Extensions["agent.availability"]))
		actions := []CoordAction{{Label: "Open agent profiles", URL: "/agents"}}
		if availability != "" {
			actions = addUnique(actions, "Open agents with this availability", pageURL("/agents", FeedOptions{MetaKey: "agent.availability"}, "meta_value", availability))
		}
		actions = append(actions, coordFocusLinks(index, post, []string{"task", "skill", "markdown", "code"}, limit, map[string]string{
			"task":     "Open linked task",
			"skill":    "Open supported skill",
			"markdown": "Open shared note",
			"code":     "Open implementation asset",
		})...)
		description := "Agent profiles should reveal where this operator is active, what it can safely do, and which workstream needs attention next."
		if availability == "internal" {
			description = "Internal operators should point directly to the tasks and reusable skills they can pick up next inside the current coordination workspace."
		}
		return &CoordDetailFocusPanel{
			Eyebrow:     "Operator mode",
			Title:       "Route this agent into the right work",
			Description: description,
			Checklist: []string{
				"Confirm the availability, models, tools, and declared skills before assigning new work.",
				"Open the linked task or supported skill that best matches the current delivery need.",
				"Check the shared note or implementation asset before asking the agent to operate beyond its current context.",
			},
			Actions: actions,
		}
	default:
		return nil
	}
}

func coordFocusLinks(index Index, post Post, targets []string, limit int, labels map[string]string) []CoordAction {
	actions := make([]CoordAction, 0, len(targets))
	add := func(label, url string) {
		label = strings.TrimSpace(label)
		url = strings.TrimSpace(url)
		if label == "" || url == "" {
			return
		}
		for _, action := range actions {
			if action.URL == url {
				return
			}
		}
		actions = append(actions, CoordAction{Label: label, URL: url})
	}
	for _, target := range targets {
		matches := relatedPostsByType(index, post, target, limit)
		if len(matches) == 0 {
			continue
		}
		label := strings.TrimSpace(labels[target])
		if label == "" {
			label = "Open"
		}
		add(label+": "+matches[0].Message.Title, CoordPath(matches[0]))
	}
	return actions
}

func CoordRelatedGroups(index Index, post Post, limit int) []CoordRelatedGroup {
	coordType := PostCoordType(post)
	if coordType == "" || limit <= 0 {
		return nil
	}
	targetsByType := map[string][]string{
		"task":     {"skill", "markdown", "code", "agent"},
		"skill":    {"task", "markdown", "code"},
		"markdown": {"task", "code", "skill"},
		"code":     {"task", "skill", "markdown"},
		"agent":    {"task", "skill"},
		"idea":     {"task", "markdown"},
	}
	targets := targetsByType[coordType]
	if len(targets) == 0 {
		return nil
	}
	groups := make([]CoordRelatedGroup, 0, len(targets))
	for _, target := range targets {
		matches := relatedPostsByType(index, post, target, limit)
		if len(matches) == 0 {
			continue
		}
		spec, ok := CoordCollectionByType(target)
		if !ok {
			continue
		}
		groups = append(groups, CoordRelatedGroup{
			Title: spec.Title,
			Kind:  target,
			Path:  spec.Path,
			Posts: matches,
		})
	}
	return groups
}

func CoordExtendedContextGroupsForPost(index Index, post Post, limit int) []CoordRelatedGroup {
	groups := CoordRelatedGroups(index, post, limit)
	if len(groups) == 0 {
		return nil
	}
	panel := CoordWorkbenchPanelForPost(index, post, limit)
	if panel == nil || len(panel.Groups) == 0 {
		return groups
	}
	seen := make(map[string]struct{}, len(panel.Groups))
	for _, group := range panel.Groups {
		kind := strings.TrimSpace(group.Kind)
		if kind == "" {
			continue
		}
		seen[kind] = struct{}{}
	}
	filtered := make([]CoordRelatedGroup, 0, len(groups))
	for _, group := range groups {
		if _, ok := seen[strings.TrimSpace(group.Kind)]; ok {
			continue
		}
		filtered = append(filtered, group)
	}
	if len(filtered) == 0 {
		return nil
	}
	return filtered
}

func CoordRelationGroupsForPost(index Index, post Post, limit int) []CoordRelationGroup {
	coordType := PostCoordType(post)
	if coordType == "" || limit <= 0 {
		return nil
	}
	type relationSpec struct {
		Key         string
		Kind        string
		Title       string
		Description string
		Priority    int
	}
	specsByType := map[string][]relationSpec{
		"idea": {
			{Key: "promoted_tasks", Kind: "task", Title: "Promoted tasks", Description: "Execution items that can carry this idea out of backlog mode.", Priority: 2},
			{Key: "shaping_notes", Kind: "markdown", Title: "Shaping notes", Description: "Planning or research notes that sharpen the problem framing.", Priority: 1},
		},
		"task": {
			{Key: "required_skills", Kind: "skill", Title: "Capability dependencies", Description: "Reusable capabilities this task likely depends on for execution.", Priority: 4},
			{Key: "supporting_knowledge", Kind: "markdown", Title: "Delivery notes", Description: "Notes and documents that explain decisions, rollout steps, or operating context.", Priority: 3},
			{Key: "implementation_assets", Kind: "code", Title: "Implementation surface", Description: "Code-facing assets that likely produce the task output.", Priority: 2},
			{Key: "operating_agents", Kind: "agent", Title: "Operator ownership", Description: "Agent profiles that can pick up or continue this delivery item.", Priority: 1},
		},
		"skill": {
			{Key: "applied_tasks", Kind: "task", Title: "Applied tasks", Description: "Delivery items currently reusing or calling for this capability.", Priority: 4},
			{Key: "backing_knowledge", Kind: "markdown", Title: "Backing knowledge", Description: "Guides or notes that define how the skill should be run.", Priority: 3},
			{Key: "implementation_assets", Kind: "code", Title: "Implementation assets", Description: "Code or scripts that embody the capability in practice.", Priority: 2},
			{Key: "operating_agents", Kind: "agent", Title: "Operating agents", Description: "Agent profiles that can execute or maintain this skill.", Priority: 1},
		},
		"markdown": {
			{Key: "informed_tasks", Kind: "task", Title: "Informed tasks", Description: "Delivery items that are likely to depend on this note before moving forward.", Priority: 3},
			{Key: "linked_code", Kind: "code", Title: "Linked code", Description: "Implementation assets that this document helps explain or justify.", Priority: 2},
			{Key: "linked_skills", Kind: "skill", Title: "Linked skills", Description: "Reusable capabilities that draw operating guidance from this note.", Priority: 1},
		},
		"code": {
			{Key: "owning_tasks", Kind: "task", Title: "Owning tasks", Description: "Delivery items that give this implementation asset its current purpose.", Priority: 3},
			{Key: "implemented_skills", Kind: "skill", Title: "Implemented skills", Description: "Capabilities this code asset helps make executable.", Priority: 2},
			{Key: "explaining_notes", Kind: "markdown", Title: "Explaining notes", Description: "Knowledge assets that describe why this implementation exists or how to operate it.", Priority: 1},
		},
		"agent": {
			{Key: "owned_tasks", Kind: "task", Title: "Delivery ownership", Description: "Delivery items this operator is positioned to execute or unblock.", Priority: 4},
			{Key: "supported_skills", Kind: "skill", Title: "Capability coverage", Description: "Capabilities this agent can safely run or maintain.", Priority: 3},
			{Key: "shared_knowledge", Kind: "markdown", Title: "Shared knowledge", Description: "Notes this operator is likely to read before acting.", Priority: 2},
			{Key: "implementation_surface", Kind: "code", Title: "Implementation touchpoints", Description: "Code-facing assets this operator is likely to touch while working.", Priority: 1},
		},
	}
	specs := specsByType[coordType]
	if len(specs) == 0 {
		return nil
	}
	groups := make([]CoordRelationGroup, 0, len(specs))
	for _, spec := range specs {
		matches := relatedPostsByType(index, post, spec.Kind, limit)
		if len(matches) == 0 {
			continue
		}
		collection, ok := CoordCollectionByType(spec.Kind)
		if !ok {
			continue
		}
		items := make([]CoordRelationItem, 0, len(matches))
		for _, match := range matches {
			items = append(items, CoordRelationItem{
				Post:     match,
				Evidence: coordRelationItemEvidence(post, match, spec.Key),
				Score:    coordRelationScore(post, match, spec.Key),
			})
		}
		sort.SliceStable(items, func(i, j int) bool {
			if items[i].Score != items[j].Score {
				return items[i].Score > items[j].Score
			}
			left := items[i].Post.CreatedAt
			right := items[j].Post.CreatedAt
			if !left.Equal(right) {
				return left.After(right)
			}
			return items[i].Post.InfoHash < items[j].Post.InfoHash
		})
		groups = append(groups, CoordRelationGroup{
			Key:            spec.Key,
			Title:          spec.Title,
			Description:    spec.Description,
			DominantReason: coordRelationDominantReason(coordRelationEvidence(post, spec.Key), items),
			Evidence:       coordRelationEvidence(post, spec.Key),
			Kind:           spec.Kind,
			Path:           collection.Path,
			Priority:       spec.Priority,
			Score:          coordRelationGroupScore(spec.Priority, items),
			Items:          items,
		})
	}
	sort.SliceStable(groups, func(i, j int) bool {
		if groups[i].Score != groups[j].Score {
			return groups[i].Score > groups[j].Score
		}
		if groups[i].Priority != groups[j].Priority {
			return groups[i].Priority > groups[j].Priority
		}
		if len(groups[i].Items) != len(groups[j].Items) {
			return len(groups[i].Items) > len(groups[j].Items)
		}
		return groups[i].Key < groups[j].Key
	})
	return groups
}

func coordRelationGroupScore(priority int, items []CoordRelationItem) int {
	if len(items) == 0 {
		return priority * 8
	}
	score := priority * 8
	score += items[0].Score
	if len(items) > 1 {
		score += minInt(12, items[1].Score/3)
	}
	score += minInt(9, len(items)*3)
	return score
}

func coordRelationDominantReason(groupEvidence []string, items []CoordRelationItem) string {
	parts := make([]string, 0, 2)
	if len(items) > 0 {
		if line := preferredRelationEvidenceLine(items[0].Evidence, "Matched workstream"); line != "" {
			parts = append(parts, line)
		}
	}
	if line := preferredRelationEvidenceLine(groupEvidence, "Shared workstream"); line != "" {
		duplicate := false
		for _, part := range parts {
			if part == line {
				duplicate = true
				break
			}
		}
		if !duplicate {
			parts = append(parts, line)
		}
	}
	if len(parts) == 0 {
		return ""
	}
	if len(parts) > 2 {
		parts = parts[:2]
	}
	return "Top match: " + strings.Join(parts, " · ")
}

func preferredRelationEvidenceLine(lines []string, deprioritizedPrefix string) string {
	deprioritizedPrefix = strings.TrimSpace(deprioritizedPrefix)
	fallback := ""
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if fallback == "" {
			fallback = line
		}
		if deprioritizedPrefix != "" && strings.HasPrefix(line, deprioritizedPrefix+":") {
			continue
		}
		return line
	}
	return fallback
}

func coordRelationScore(base Post, target Post, key string) int {
	score := 0
	shared := sharedTopics(base.Topics, target.Topics)
	if len(shared) > 0 {
		score += 30
		if len(shared) > 1 {
			score += minInt(10, (len(shared)-1)*4)
		}
	}
	if strings.EqualFold(strings.TrimSpace(base.SourceName), strings.TrimSpace(target.SourceName)) && strings.TrimSpace(base.SourceName) != "" {
		score += 6
	}
	overlap := tokenOverlapCount(relationBaseTexts(base, key), relationTargetTexts(target, key))
	score += minInt(32, overlap*8)
	presence := 0
	for _, text := range relationTargetTexts(target, key) {
		if strings.TrimSpace(text) != "" {
			presence++
		}
	}
	score += minInt(12, presence*2)
	return score
}

func coordRelationEvidence(post Post, key string) []string {
	add := func(lines []string, label, value string) []string {
		label = strings.TrimSpace(label)
		value = strings.TrimSpace(value)
		if label == "" || value == "" {
			return lines
		}
		line := label + ": " + value
		for _, existing := range lines {
			if existing == line {
				return lines
			}
		}
		return append(lines, line)
	}
	var lines []string
	if topic := primaryTopic(post.Topics); topic != "" {
		lines = add(lines, "Shared workstream", topic)
	}
	switch PostCoordType(post) {
	case "idea":
		switch key {
		case "promoted_tasks":
			lines = add(lines, "Goal", extensionText(post.Message.Extensions["coord.goal"]))
			lines = add(lines, "Expected output", extensionText(post.Message.Extensions["coord.expected_output"]))
		case "shaping_notes":
			lines = add(lines, "Problem", extensionText(post.Message.Extensions["coord.problem"]))
			lines = add(lines, "Domain", extensionText(post.Message.Extensions["coord.domain"]))
		}
	case "task":
		switch key {
		case "required_skills":
			lines = add(lines, "Required assets", extensionText(post.Message.Extensions["task.required_assets"]))
			lines = add(lines, "Status", extensionText(post.Message.Extensions["task.status"]))
		case "supporting_knowledge":
			lines = add(lines, "Parent", extensionText(post.Message.Extensions["task.parent"]))
			lines = add(lines, "Expected result", extensionText(post.Message.Extensions["task.expected_result_type"]))
		case "implementation_assets":
			lines = add(lines, "Expected result", extensionText(post.Message.Extensions["task.expected_result_type"]))
			lines = add(lines, "Required assets", extensionText(post.Message.Extensions["task.required_assets"]))
		case "operating_agents":
			lines = add(lines, "Status", extensionText(post.Message.Extensions["task.status"]))
			lines = add(lines, "Priority", extensionText(post.Message.Extensions["task.priority"]))
		}
	case "skill":
		switch key {
		case "applied_tasks":
			lines = add(lines, "Category", extensionText(post.Message.Extensions["skill.category"]))
			lines = add(lines, "Depends on", extensionText(post.Message.Extensions["skill.depends_on"]))
		case "backing_knowledge":
			lines = add(lines, "Depends on", extensionText(post.Message.Extensions["skill.depends_on"]))
			lines = add(lines, "Repository", extensionText(post.Message.Extensions["skill.repo"]))
		case "implementation_assets":
			lines = add(lines, "Outputs", extensionText(post.Message.Extensions["skill.outputs"]))
			lines = add(lines, "Repository", extensionText(post.Message.Extensions["skill.repo"]))
		case "operating_agents":
			lines = add(lines, "Category", extensionText(post.Message.Extensions["skill.category"]))
			lines = add(lines, "Inputs", extensionText(post.Message.Extensions["skill.inputs"]))
		}
	case "markdown":
		switch key {
		case "informed_tasks":
			lines = add(lines, "Collection", extensionText(post.Message.Extensions["md.collection"]))
			lines = add(lines, "Document kind", extensionText(post.Message.Extensions["md.kind"]))
		case "linked_code":
			lines = add(lines, "Collection", extensionText(post.Message.Extensions["md.collection"]))
			lines = add(lines, "Source repository", extensionText(post.Message.Extensions["md.source_repo"]))
		case "linked_skills":
			lines = add(lines, "Document kind", extensionText(post.Message.Extensions["md.kind"]))
			lines = add(lines, "Collection", extensionText(post.Message.Extensions["md.collection"]))
		}
	case "code":
		switch key {
		case "owning_tasks":
			lines = add(lines, "Code kind", extensionText(post.Message.Extensions["code.kind"]))
			lines = add(lines, "Entry point", extensionText(post.Message.Extensions["code.entry"]))
		case "implemented_skills":
			lines = add(lines, "Language", extensionText(post.Message.Extensions["code.language"]))
			lines = add(lines, "Entry point", extensionText(post.Message.Extensions["code.entry"]))
		case "explaining_notes":
			lines = add(lines, "Repository", extensionText(post.Message.Extensions["code.repo"]))
			lines = add(lines, "Entry point", extensionText(post.Message.Extensions["code.entry"]))
		}
	case "agent":
		switch key {
		case "owned_tasks":
			lines = add(lines, "Availability", extensionText(post.Message.Extensions["agent.availability"]))
			lines = add(lines, "Skills", extensionText(post.Message.Extensions["agent.skills"]))
		case "supported_skills":
			lines = add(lines, "Skills", extensionText(post.Message.Extensions["agent.skills"]))
			lines = add(lines, "Models", extensionText(post.Message.Extensions["agent.models"]))
		case "shared_knowledge":
			lines = add(lines, "Tools", extensionText(post.Message.Extensions["agent.tools"]))
			lines = add(lines, "Skills", extensionText(post.Message.Extensions["agent.skills"]))
		case "implementation_surface":
			lines = add(lines, "Tools", extensionText(post.Message.Extensions["agent.tools"]))
			lines = add(lines, "Availability", extensionText(post.Message.Extensions["agent.availability"]))
		}
	}
	if len(lines) > 3 {
		return append([]string(nil), lines[:3]...)
	}
	return lines
}

func coordRelationItemEvidence(base Post, target Post, key string) []string {
	add := func(lines []string, label, value string) []string {
		label = strings.TrimSpace(label)
		value = strings.TrimSpace(value)
		if label == "" || value == "" {
			return lines
		}
		line := label + ": " + value
		for _, existing := range lines {
			if existing == line {
				return lines
			}
		}
		return append(lines, line)
	}
	var lines []string
	if topic := primarySharedTopic(base.Topics, target.Topics); topic != "" {
		lines = add(lines, "Matched workstream", topic)
	}
	switch key {
	case "promoted_tasks":
		lines = add(lines, "Target status", extensionText(target.Message.Extensions["task.status"]))
		lines = add(lines, "Target result", extensionText(target.Message.Extensions["task.expected_result_type"]))
	case "shaping_notes":
		lines = add(lines, "Target collection", extensionText(target.Message.Extensions["md.collection"]))
		lines = add(lines, "Target kind", extensionText(target.Message.Extensions["md.kind"]))
	case "required_skills":
		lines = add(lines, "Target category", extensionText(target.Message.Extensions["skill.category"]))
		lines = add(lines, "Target outputs", extensionText(target.Message.Extensions["skill.outputs"]))
	case "supporting_knowledge":
		lines = add(lines, "Target collection", extensionText(target.Message.Extensions["md.collection"]))
		lines = add(lines, "Target kind", extensionText(target.Message.Extensions["md.kind"]))
	case "implementation_assets":
		lines = add(lines, "Target language", extensionText(target.Message.Extensions["code.language"]))
		lines = add(lines, "Target entry", extensionText(target.Message.Extensions["code.entry"]))
	case "operating_agents":
		lines = add(lines, "Target availability", extensionText(target.Message.Extensions["agent.availability"]))
		lines = add(lines, "Target skills", extensionText(target.Message.Extensions["agent.skills"]))
	case "applied_tasks":
		lines = add(lines, "Target status", extensionText(target.Message.Extensions["task.status"]))
		lines = add(lines, "Target priority", extensionText(target.Message.Extensions["task.priority"]))
	case "backing_knowledge":
		lines = add(lines, "Target collection", extensionText(target.Message.Extensions["md.collection"]))
		lines = add(lines, "Target source", extensionText(target.Message.Extensions["md.source_repo"]))
	case "informed_tasks", "owning_tasks", "owned_tasks", "delivery_ownership":
		lines = add(lines, "Target status", extensionText(target.Message.Extensions["task.status"]))
		lines = add(lines, "Target result", extensionText(target.Message.Extensions["task.expected_result_type"]))
	case "implemented_skills", "supported_skills":
		lines = add(lines, "Target category", extensionText(target.Message.Extensions["skill.category"]))
		lines = add(lines, "Target inputs", extensionText(target.Message.Extensions["skill.inputs"]))
	case "linked_code", "implementation_surface":
		lines = add(lines, "Target language", extensionText(target.Message.Extensions["code.language"]))
		lines = add(lines, "Target entry", extensionText(target.Message.Extensions["code.entry"]))
	case "linked_skills", "shared_knowledge", "explaining_notes":
		lines = add(lines, "Target collection", extensionText(target.Message.Extensions["md.collection"]))
		lines = add(lines, "Target kind", extensionText(target.Message.Extensions["md.kind"]))
	}
	switch PostCoordType(target) {
	case "task":
		lines = add(lines, "Target status", extensionText(target.Message.Extensions["task.status"]))
	case "skill":
		lines = add(lines, "Target category", extensionText(target.Message.Extensions["skill.category"]))
	case "markdown":
		lines = add(lines, "Target collection", extensionText(target.Message.Extensions["md.collection"]))
	case "code":
		lines = add(lines, "Target language", extensionText(target.Message.Extensions["code.language"]))
	case "agent":
		lines = add(lines, "Target availability", extensionText(target.Message.Extensions["agent.availability"]))
	}
	if len(lines) > 3 {
		return append([]string(nil), lines[:3]...)
	}
	return lines
}

func primarySharedTopic(baseTopics, targetTopics []string) string {
	shared := sharedTopics(baseTopics, targetTopics)
	if len(shared) == 0 {
		return ""
	}
	return shared[0]
}

func sharedTopics(baseTopics, targetTopics []string) []string {
	allowed := normalizedTopics(baseTopics)
	if len(allowed) == 0 {
		return nil
	}
	out := make([]string, 0, len(targetTopics))
	seen := make(map[string]struct{}, len(targetTopics))
	for _, topic := range targetTopics {
		name := strings.ToLower(strings.TrimSpace(topic))
		if name == "" || name == "all" {
			continue
		}
		if _, ok := allowed[name]; !ok {
			continue
		}
		if _, ok := seen[name]; ok {
			continue
		}
		seen[name] = struct{}{}
		out = append(out, name)
	}
	return out
}

func relationBaseTexts(post Post, key string) []string {
	switch PostCoordType(post) {
	case "idea":
		return []string{
			post.Message.Title,
			post.Summary,
			extensionText(post.Message.Extensions["coord.problem"]),
			extensionText(post.Message.Extensions["coord.goal"]),
			extensionText(post.Message.Extensions["coord.expected_output"]),
			extensionText(post.Message.Extensions["coord.domain"]),
		}
	case "task":
		return []string{
			post.Message.Title,
			post.Summary,
			extensionText(post.Message.Extensions["task.required_assets"]),
			extensionText(post.Message.Extensions["task.expected_result_type"]),
			extensionText(post.Message.Extensions["task.status"]),
			extensionText(post.Message.Extensions["task.priority"]),
			extensionText(post.Message.Extensions["task.parent"]),
		}
	case "skill":
		return []string{
			post.Message.Title,
			post.Summary,
			extensionText(post.Message.Extensions["skill.category"]),
			extensionText(post.Message.Extensions["skill.depends_on"]),
			extensionText(post.Message.Extensions["skill.inputs"]),
			extensionText(post.Message.Extensions["skill.outputs"]),
		}
	case "markdown":
		return []string{
			post.Message.Title,
			post.Summary,
			extensionText(post.Message.Extensions["md.kind"]),
			extensionText(post.Message.Extensions["md.collection"]),
			extensionText(post.Message.Extensions["md.source_repo"]),
		}
	case "code":
		return []string{
			post.Message.Title,
			post.Summary,
			extensionText(post.Message.Extensions["code.kind"]),
			extensionText(post.Message.Extensions["code.entry"]),
			extensionText(post.Message.Extensions["code.language"]),
			extensionText(post.Message.Extensions["code.repo"]),
		}
	case "agent":
		return []string{
			post.Message.Title,
			post.Summary,
			extensionText(post.Message.Extensions["agent.skills"]),
			extensionText(post.Message.Extensions["agent.availability"]),
			extensionText(post.Message.Extensions["agent.models"]),
			extensionText(post.Message.Extensions["agent.tools"]),
		}
	default:
		return []string{post.Message.Title, post.Summary}
	}
}

func relationTargetTexts(post Post, key string) []string {
	_ = key
	switch PostCoordType(post) {
	case "task":
		return []string{
			post.Message.Title,
			post.Summary,
			extensionText(post.Message.Extensions["task.status"]),
			extensionText(post.Message.Extensions["task.expected_result_type"]),
			extensionText(post.Message.Extensions["task.priority"]),
		}
	case "skill":
		return []string{
			post.Message.Title,
			post.Summary,
			extensionText(post.Message.Extensions["skill.category"]),
			extensionText(post.Message.Extensions["skill.inputs"]),
			extensionText(post.Message.Extensions["skill.outputs"]),
		}
	case "markdown":
		return []string{
			post.Message.Title,
			post.Summary,
			extensionText(post.Message.Extensions["md.kind"]),
			extensionText(post.Message.Extensions["md.collection"]),
			extensionText(post.Message.Extensions["md.source_repo"]),
		}
	case "code":
		return []string{
			post.Message.Title,
			post.Summary,
			extensionText(post.Message.Extensions["code.kind"]),
			extensionText(post.Message.Extensions["code.entry"]),
			extensionText(post.Message.Extensions["code.language"]),
		}
	case "agent":
		return []string{
			post.Message.Title,
			post.Summary,
			extensionText(post.Message.Extensions["agent.skills"]),
			extensionText(post.Message.Extensions["agent.availability"]),
			extensionText(post.Message.Extensions["agent.tools"]),
		}
	default:
		return []string{post.Message.Title, post.Summary}
	}
}

func tokenOverlapCount(left []string, right []string) int {
	lhs := tokenSet(left...)
	rhs := tokenSet(right...)
	if len(lhs) == 0 || len(rhs) == 0 {
		return 0
	}
	count := 0
	for token := range lhs {
		if _, ok := rhs[token]; ok {
			count++
		}
	}
	return count
}

func tokenSet(values ...string) map[string]struct{} {
	out := make(map[string]struct{})
	for _, value := range values {
		value = strings.ToLower(strings.TrimSpace(value))
		if value == "" {
			continue
		}
		fields := strings.FieldsFunc(value, func(r rune) bool {
			return !(r >= 'a' && r <= 'z') && !(r >= '0' && r <= '9')
		})
		for _, field := range fields {
			field = strings.TrimSpace(field)
			if len(field) < 3 {
				continue
			}
			out[field] = struct{}{}
		}
	}
	return out
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func CoordWorkbenchPanelForPost(index Index, post Post, limit int) *CoordWorkbenchPanel {
	coordType := PostCoordType(post)
	if coordType == "" || limit <= 0 {
		return nil
	}
	switch coordType {
	case "task":
		groups := coordWorkbenchGroups(index, post, []string{"skill", "markdown", "code", "agent"}, limit, map[string]string{
			"skill":    "Required skills",
			"markdown": "Reference docs",
			"code":     "Implementation assets",
			"agent":    "Operating agents",
		})
		if len(groups) == 0 {
			return nil
		}
		return &CoordWorkbenchPanel{
			Eyebrow:     "Execution map",
			Title:       "Task delivery workspace",
			Description: "Assets and operators currently attached to this task by shared coordination context.",
			Groups:      groups,
		}
	case "agent":
		groups := coordWorkbenchGroups(index, post, []string{"task", "skill", "markdown", "code"}, limit, map[string]string{
			"task":     "Linked tasks",
			"skill":    "Supported skills",
			"markdown": "Shared docs",
			"code":     "Implementation surface",
		})
		if len(groups) == 0 {
			return nil
		}
		return &CoordWorkbenchPanel{
			Eyebrow:     "Contribution map",
			Title:       "Agent coordination workspace",
			Description: "Nearby work and reusable assets currently associated with this agent profile.",
			Groups:      groups,
		}
	case "skill":
		groups := coordWorkbenchGroups(index, post, []string{"task", "markdown", "code", "agent"}, limit, map[string]string{
			"task":     "Tasks using this skill",
			"markdown": "Docs backing this skill",
			"code":     "Code implementing this skill",
			"agent":    "Agents operating this skill",
		})
		if len(groups) == 0 {
			return nil
		}
		return &CoordWorkbenchPanel{
			Eyebrow:     "Usage map",
			Title:       "Skill reuse workspace",
			Description: "Where this skill is being applied and which assets or operators surround it.",
			Groups:      groups,
		}
	default:
		return nil
	}
}

func WorkstreamRelationPaths(posts []Post, limit int) []WorkstreamRelationPath {
	if limit <= 0 || len(posts) == 0 {
		return nil
	}
	index := Index{Posts: posts}
	out := make([]WorkstreamRelationPath, 0, limit)
	seen := make(map[string]struct{}, len(posts))
	for _, post := range posts {
		coordType := PostCoordType(post)
		if coordType == "" {
			continue
		}
		groups := CoordRelationGroupsForPost(index, post, 2)
		if len(groups) == 0 || len(groups[0].Items) == 0 {
			continue
		}
		top := groups[0]
		target := top.Items[0]
		key := post.InfoHash + "|" + top.Key + "|" + target.Post.InfoHash
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, WorkstreamRelationPath{
			From:          post,
			RelationKey:   top.Key,
			RelationTitle: top.Title,
			Kind:          top.Kind,
			Path:          top.Path,
			Target:        target.Post,
			Reason:        top.DominantReason,
			Score:         top.Score,
		})
	}
	sort.SliceStable(out, func(i, j int) bool {
		leftPriority := workstreamPathSourcePriority(out[i].From)
		rightPriority := workstreamPathSourcePriority(out[j].From)
		if leftPriority != rightPriority {
			return leftPriority > rightPriority
		}
		if out[i].Score != out[j].Score {
			return out[i].Score > out[j].Score
		}
		left := out[i].From.CreatedAt
		right := out[j].From.CreatedAt
		if !left.Equal(right) {
			return left.After(right)
		}
		if out[i].RelationTitle != out[j].RelationTitle {
			return out[i].RelationTitle < out[j].RelationTitle
		}
		return out[i].From.InfoHash < out[j].From.InfoHash
	})
	if len(out) > limit {
		return append([]WorkstreamRelationPath(nil), out[:limit]...)
	}
	return out
}

func WorkstreamRelationPanelForPosts(posts []Post, limit int) *WorkstreamRelationPanel {
	paths := WorkstreamRelationPaths(posts, limit)
	if len(paths) == 0 {
		return nil
	}
	primaryCount := 2
	if len(paths) < primaryCount {
		primaryCount = len(paths)
	}
	panel := &WorkstreamRelationPanel{
		Title:        "Coordination paths",
		Description:  "Primary typed links currently anchoring this scoped workspace, followed by lower-priority supporting paths.",
		PrimaryPaths: append([]WorkstreamRelationPath(nil), paths[:primaryCount]...),
	}
	if len(paths) > primaryCount {
		panel.SupportingPaths = append([]WorkstreamRelationPath(nil), paths[primaryCount:]...)
	}
	return panel
}

func workstreamPathSourcePriority(post Post) int {
	switch PostCoordType(post) {
	case "task":
		return 6
	case "idea":
		return 5
	case "agent":
		return 4
	case "skill":
		return 3
	case "markdown":
		return 2
	case "code":
		return 1
	default:
		return 0
	}
}

func coordWorkbenchGroups(index Index, post Post, targets []string, limit int, titles map[string]string) []CoordRelatedGroup {
	groups := make([]CoordRelatedGroup, 0, len(targets))
	for _, target := range targets {
		matches := relatedPostsByType(index, post, target, limit)
		if len(matches) == 0 {
			continue
		}
		spec, ok := CoordCollectionByType(target)
		if !ok {
			continue
		}
		title := spec.Title
		if custom := strings.TrimSpace(titles[target]); custom != "" {
			title = custom
		}
		groups = append(groups, CoordRelatedGroup{
			Title: title,
			Kind:  target,
			Path:  spec.Path,
			Posts: matches,
		})
	}
	return groups
}

func HomeWorkbenchStats(posts []Post) []SummaryStat {
	typedCount := 0
	activeTasks := 0
	skillCount := 0
	knowledgeCount := 0
	agentCount := 0
	legacyCount := 0
	for _, post := range posts {
		switch PostCoordType(post) {
		case "":
			legacyCount++
		default:
			typedCount++
		}
		switch PostCoordType(post) {
		case "task":
			if isActiveTaskStatus(taskStatus(post)) {
				activeTasks++
			}
		case "skill":
			skillCount++
		case "markdown":
			knowledgeCount++
		case "agent":
			agentCount++
		}
	}
	return []SummaryStat{
		{Label: "Typed assets", Value: fmt.Sprintf("%d", typedCount)},
		{Label: "Active tasks", Value: fmt.Sprintf("%d", activeTasks)},
		{Label: "Skills", Value: fmt.Sprintf("%d", skillCount)},
		{Label: "Knowledge docs", Value: fmt.Sprintf("%d", knowledgeCount)},
		{Label: "Agents", Value: fmt.Sprintf("%d", agentCount)},
		{Label: "Legacy carry-over", Value: fmt.Sprintf("%d", legacyCount)},
	}
}

func relatedPostsByType(index Index, base Post, targetType string, limit int) []Post {
	if limit <= 0 {
		return nil
	}
	topics := normalizedTopics(base.Topics)
	if len(topics) == 0 {
		return nil
	}
	matches := make([]Post, 0, limit)
	for _, post := range index.Posts {
		if post.InfoHash == base.InfoHash {
			continue
		}
		if PostCoordType(post) != targetType {
			continue
		}
		if !sharesTopic(post.Topics, topics) {
			continue
		}
		matches = append(matches, post)
		if len(matches) == limit {
			break
		}
	}
	return matches
}

func normalizedTopics(topics []string) map[string]struct{} {
	if len(topics) == 0 {
		return nil
	}
	out := make(map[string]struct{}, len(topics))
	for _, topic := range topics {
		topic = strings.ToLower(strings.TrimSpace(topic))
		if topic == "" {
			continue
		}
		out[topic] = struct{}{}
	}
	return out
}

func sharesTopic(topics []string, allowed map[string]struct{}) bool {
	if len(topics) == 0 || len(allowed) == 0 {
		return false
	}
	for _, topic := range topics {
		topic = strings.ToLower(strings.TrimSpace(topic))
		if _, ok := allowed[topic]; ok {
			return true
		}
	}
	return false
}

func HomeWorkbenchActions(posts []Post) []CoordAction {
	actions := []CoordAction{
		{Label: "Open tasks", URL: "/tasks"},
		{Label: "Active tasks", URL: "/tasks?meta_key=task.status&meta_value=in_progress"},
		{Label: "Skills registry", URL: "/skills"},
		{Label: "Knowledge base", URL: "/knowledge"},
		{Label: "Agent profiles", URL: "/agents"},
		{Label: "Archive", URL: "/archive"},
		{Label: "JSON feed", URL: "/api/feed"},
	}
	if len(FilterPostsByCoordType(posts, "idea")) > 0 {
		actions = append([]CoordAction{{Label: "Idea backlog", URL: "/ideas"}}, actions...)
	}
	return actions
}

func HomePriorityCards(posts []Post) []HomePriorityCard {
	typedPosts := filterTypedPosts(posts)
	liveTopics := 0
	for _, stat := range TopicStatsForPosts(typedPosts) {
		if name := strings.TrimSpace(stat.Name); name != "" && !strings.EqualFold(name, "all") {
			liveTopics++
		}
	}
	activeTasks := len(filterTaskPostsByStatus(posts, len(posts), []string{"in_progress", "blocked", "queued", "ready"}))
	taskBacklog := len(FilterPostsByCoordType(posts, "task"))
	ideaBacklog := len(FilterPostsByCoordType(posts, "idea"))
	skills := len(FilterPostsByCoordType(posts, "skill"))
	knowledge := len(FilterPostsByCoordType(posts, "markdown"))
	codeAssets := len(FilterPostsByCoordType(posts, "code"))
	agents := len(FilterPostsByCoordType(posts, "agent"))
	sources := len(SourceStatsForPosts(typedPosts))

	return []HomePriorityCard{
		{
			Eyebrow:     "Priority 01",
			Title:       "Execute active work",
			Description: "Start from the active task surface when the goal is delivery, blockers, or immediate next actions.",
			URL:         "/tasks?meta_key=task.status&meta_value=in_progress",
			Stats: []SummaryStat{
				{Label: "Active tasks", Value: fmt.Sprintf("%d", activeTasks)},
				{Label: "Task backlog", Value: fmt.Sprintf("%d", taskBacklog)},
				{Label: "Idea backlog", Value: fmt.Sprintf("%d", ideaBacklog)},
			},
			Actions: []CoordAction{
				{Label: "Open active tasks", URL: "/tasks?meta_key=task.status&meta_value=in_progress"},
				{Label: "Open all tasks", URL: "/tasks"},
				{Label: "Browse ideas", URL: "/ideas"},
			},
		},
		{
			Eyebrow:     "Priority 02",
			Title:       "Coordinate workstreams",
			Description: "Open the registry surfaces when you need to widen scope from one task into its live topics, sources, and operators.",
			URL:         "/topics",
			Stats: []SummaryStat{
				{Label: "Live workstreams", Value: fmt.Sprintf("%d", liveTopics)},
				{Label: "Tracked sources", Value: fmt.Sprintf("%d", sources)},
				{Label: "Agent profiles", Value: fmt.Sprintf("%d", agents)},
			},
			Actions: []CoordAction{
				{Label: "Topic registry", URL: "/topics"},
				{Label: "Source registry", URL: "/sources"},
				{Label: "Agent profiles", URL: "/agents"},
			},
		},
		{
			Eyebrow:     "Priority 03",
			Title:       "Reuse proven assets",
			Description: "Reach for reusable capabilities and durable references before creating new coordination debt.",
			URL:         "/skills",
			Stats: []SummaryStat{
				{Label: "Skills", Value: fmt.Sprintf("%d", skills)},
				{Label: "Knowledge docs", Value: fmt.Sprintf("%d", knowledge)},
				{Label: "Code assets", Value: fmt.Sprintf("%d", codeAssets)},
			},
			Actions: []CoordAction{
				{Label: "Skills registry", URL: "/skills"},
				{Label: "Knowledge base", URL: "/knowledge"},
				{Label: "Code assets", URL: "/code"},
			},
		},
	}
}

func ScopedCoordActions(scopeKind, scopeName string, posts []Post) []CoordAction {
	scopeKind = strings.ToLower(strings.TrimSpace(scopeKind))
	scopeName = strings.TrimSpace(scopeName)
	if scopeName == "" {
		return nil
	}
	var key string
	switch scopeKind {
	case "source":
		key = "source"
	case "topic":
		key = "topic"
	default:
		return nil
	}
	opts := FeedOptions{}
	actions := make([]CoordAction, 0, len(coordCollections))
	for _, spec := range CoordCollections() {
		count := len(FilterPostsByCoordType(posts, spec.CoordType))
		if count == 0 {
			continue
		}
		label := fmt.Sprintf("%s (%d)", spec.Title, count)
		actions = append(actions, CoordAction{
			Label: label,
			URL:   pageURL(spec.Path, opts, key, scopeName),
		})
	}
	return actions
}

func APIScopedCoordActions(scopeKind, scopeName string, posts []Post) []CoordAction {
	scopeKind = strings.ToLower(strings.TrimSpace(scopeKind))
	scopeName = strings.TrimSpace(scopeName)
	if scopeName == "" {
		return nil
	}
	var key string
	switch scopeKind {
	case "source":
		key = "source"
	case "topic":
		key = "topic"
	default:
		return nil
	}
	opts := FeedOptions{}
	actions := make([]CoordAction, 0, len(coordCollections))
	for _, spec := range CoordCollections() {
		count := len(FilterPostsByCoordType(posts, spec.CoordType))
		if count == 0 {
			continue
		}
		label := fmt.Sprintf("%s (%d)", spec.Title, count)
		actions = append(actions, CoordAction{
			Label: label,
			URL:   pageURL(spec.APIPath, opts, key, scopeName),
		})
	}
	return actions
}

func WorkstreamJumpPanel(scopeKind, scopeName string, posts []Post) *WorkspaceJumpPanel {
	scopeKind = strings.ToLower(strings.TrimSpace(scopeKind))
	scopeName = strings.TrimSpace(scopeName)
	if scopeName == "" {
		return nil
	}
	var key string
	var title string
	var description string
	var supportingRoot CoordAction
	switch scopeKind {
	case "source":
		key = "source"
		title = "Primary and supporting jumps for this source"
		description = "Open the most likely execution and reuse surfaces for this source before widening into the full workspace."
		supportingRoot = CoordAction{Label: "Topic registry", URL: "/topics"}
	case "topic":
		key = "topic"
		title = "Primary and supporting jumps for this workstream"
		description = "Start with the strongest typed surfaces in this workstream, then widen into supporting modules and adjacent registries."
		supportingRoot = CoordAction{Label: "Source registry", URL: "/sources"}
	default:
		return nil
	}
	actions := make([]CoordAction, 0, len(coordCollections)+2)
	add := func(label, url string) {
		label = strings.TrimSpace(label)
		url = strings.TrimSpace(url)
		if label == "" || url == "" {
			return
		}
		for _, action := range actions {
			if action.URL == url {
				return
			}
		}
		actions = append(actions, CoordAction{Label: label, URL: url})
	}
	taskCount := len(FilterPostsByCoordType(posts, "task"))
	activeTaskCount := len(filterTaskPostsByStatus(posts, len(posts), []string{"in_progress"}))
	if activeTaskCount > 0 {
		opts := FeedOptions{MetaKey: "task.status", MetaValue: "in_progress"}
		add(fmt.Sprintf("Open active tasks (%d)", activeTaskCount), pageURL("/tasks", opts, key, scopeName))
	} else if taskCount > 0 {
		add(fmt.Sprintf("Open tasks (%d)", taskCount), pageURL("/tasks", FeedOptions{}, key, scopeName))
	}
	priorityOrder := []string{"idea", "skill", "markdown", "code", "agent"}
	for _, coordType := range priorityOrder {
		spec, ok := CoordCollectionByType(coordType)
		if !ok {
			continue
		}
		count := len(FilterPostsByCoordType(posts, coordType))
		if count == 0 {
			continue
		}
		add(fmt.Sprintf("Open %s (%d)", strings.ToLower(spec.Title), count), pageURL(spec.Path, FeedOptions{}, key, scopeName))
	}
	if supportingRoot.URL != "" {
		add(supportingRoot.Label, supportingRoot.URL)
	}
	if len(actions) == 0 {
		return nil
	}
	primary, supporting := splitCoordActions(actions, 3)
	return &WorkspaceJumpPanel{
		Eyebrow:           "Workspace jumps",
		Title:             title,
		Description:       description,
		PrimaryActions:    primary,
		SupportingActions: supporting,
	}
}

func ClusterJumpPanel(task Post) *WorkspaceJumpPanel {
	if strings.TrimSpace(task.InfoHash) == "" {
		return nil
	}
	var primary []CoordAction
	var supporting []CoordAction
	addPrimary := func(label, url string) {
		label = strings.TrimSpace(label)
		url = strings.TrimSpace(url)
		if label == "" || url == "" {
			return
		}
		for _, action := range primary {
			if action.URL == url {
				return
			}
		}
		primary = append(primary, CoordAction{Label: label, URL: url})
	}
	addSupporting := func(label, url string) {
		label = strings.TrimSpace(label)
		url = strings.TrimSpace(url)
		if label == "" || url == "" {
			return
		}
		for _, action := range primary {
			if action.URL == url {
				return
			}
		}
		for _, action := range supporting {
			if action.URL == url {
				return
			}
		}
		supporting = append(supporting, CoordAction{Label: label, URL: url})
	}
	addPrimary("Open task", CoordPath(task))
	for _, topic := range task.Topics {
		topic = strings.TrimSpace(topic)
		if topic == "" || strings.EqualFold(topic, "all") {
			continue
		}
		addPrimary("Open topic workstream: "+topic, TopicPath(topic))
		break
	}
	if task.HasSourcePage && strings.TrimSpace(task.SourceName) != "" {
		addPrimary("Open source workstream", SourcePath(task.SourceName))
	}
	addSupporting("Open active tasks", "/tasks?meta_key=task.status&meta_value=in_progress")
	addSupporting("Browse task module", "/tasks")
	if len(primary) == 0 && len(supporting) == 0 {
		return nil
	}
	return &WorkspaceJumpPanel{
		Eyebrow:           "Workspace jumps",
		Title:             "Primary and supporting jumps for this task",
		Description:       "Open the task itself first, then widen into the workstreams and registries that support its delivery path.",
		PrimaryActions:    primary,
		SupportingActions: supporting,
	}
}

func splitCoordActions(actions []CoordAction, primaryLimit int) ([]CoordAction, []CoordAction) {
	if len(actions) == 0 {
		return nil, nil
	}
	if primaryLimit <= 0 || primaryLimit >= len(actions) {
		return append([]CoordAction(nil), actions...), nil
	}
	primary := append([]CoordAction(nil), actions[:primaryLimit]...)
	supporting := append([]CoordAction(nil), actions[primaryLimit:]...)
	return primary, supporting
}

func HomeCoordGroupFacets(posts []Post, limit int) []FeedFacet {
	if limit <= 0 {
		return nil
	}
	stats := TopicStatsForPosts(filterTypedPosts(posts))
	facets := make([]FeedFacet, 0, limit)
	for _, stat := range stats {
		name := strings.TrimSpace(stat.Name)
		if name == "" || strings.EqualFold(name, "all") {
			continue
		}
		facets = append(facets, FeedFacet{
			Name:  name,
			Count: stat.Count,
			URL:   TopicPath(name),
		})
		if len(facets) == limit {
			break
		}
	}
	return facets
}

func HomeWorkbenchLanes(posts []Post) []WorkbenchLane {
	activeTasks := filterTaskPostsByStatus(posts, 3, []string{"in_progress", "blocked", "queued", "ready"})
	if len(activeTasks) == 0 {
		activeTasks = RecentCoordPosts(posts, "task", 3)
	}
	return []WorkbenchLane{
		{
			Title:       "Active tasks",
			Kind:        "task",
			Description: "Execution items that are currently moving or blocking delivery.",
			Path:        "/tasks?meta_key=task.status&meta_value=in_progress",
			EmptyCopy:   "No active tasks yet. Publish coord.type=task with task.status=in_progress.",
			Posts:       activeTasks,
		},
		{
			Title:       "Reusable skills",
			Kind:        "skill",
			Description: "Capabilities that can be reused across future coordination work.",
			Path:        "/skills",
			EmptyCopy:   "No reusable skills indexed yet.",
			Posts:       RecentCoordPosts(posts, "skill", 3),
		},
		{
			Title:       "Latest knowledge",
			Kind:        "markdown",
			Description: "Planning docs and durable notes that anchor the current rollout.",
			Path:        "/knowledge",
			EmptyCopy:   "No knowledge assets indexed yet.",
			Posts:       RecentCoordPosts(posts, "markdown", 3),
		},
		{
			Title:       "Code assets",
			Kind:        "code",
			Description: "Implementation notes, entry points, and code-facing coordination assets.",
			Path:        "/code",
			EmptyCopy:   "No code assets indexed yet.",
			Posts:       RecentCoordPosts(posts, "code", 3),
		},
		{
			Title:       "Agent profiles",
			Kind:        "agent",
			Description: "Named operators, tools, and declared capability surfaces on this node.",
			Path:        "/agents",
			EmptyCopy:   "No agent profiles indexed yet.",
			Posts:       RecentCoordPosts(posts, "agent", 3),
		},
		{
			Title:       "Idea backlog",
			Kind:        "idea",
			Description: "Problem framing and upstream goals before they turn into concrete tasks.",
			Path:        "/ideas",
			EmptyCopy:   "No idea backlog items indexed yet.",
			Posts:       RecentCoordPosts(posts, "idea", 3),
		},
	}
}

func HomeWorkbenchClusters(index Index, posts []Post, taskLimit, relatedLimit int) []WorkbenchCluster {
	if taskLimit <= 0 {
		return nil
	}
	tasks := filterTaskPostsByStatus(posts, taskLimit, []string{"in_progress", "blocked", "queued", "ready"})
	if len(tasks) == 0 {
		tasks = RecentCoordPosts(posts, "task", taskLimit)
	}
	clusters := make([]WorkbenchCluster, 0, len(tasks))
	for _, task := range tasks {
		clusters = append(clusters, WorkbenchCluster{
			Task:            task,
			Related:         CoordRelatedGroups(index, task, relatedLimit),
			RegistryActions: ClusterRegistryActions(task),
			JumpPanel:       ClusterJumpPanel(task),
		})
	}
	return clusters
}

func ClusterRegistryActions(task Post) []CoordAction {
	var actions []CoordAction
	add := func(label, url string) {
		label = strings.TrimSpace(label)
		url = strings.TrimSpace(url)
		if label == "" || url == "" {
			return
		}
		for _, action := range actions {
			if action.URL == url {
				return
			}
		}
		actions = append(actions, CoordAction{Label: label, URL: url})
	}
	for _, topic := range task.Topics {
		topic = strings.TrimSpace(topic)
		if topic == "" || strings.EqualFold(topic, "all") {
			continue
		}
		add("Topic workstream: "+topic, TopicPath(topic))
	}
	if task.HasSourcePage && strings.TrimSpace(task.SourceName) != "" {
		add("Source workstream", SourcePath(task.SourceName))
	}
	return actions
}

func HomeAgentWorkbench(index Index, posts []Post, limit, relatedLimit int) []AgentWorkbenchCard {
	agents := RecentCoordPosts(posts, "agent", limit)
	if len(agents) == 0 {
		return nil
	}
	cards := make([]AgentWorkbenchCard, 0, len(agents))
	for _, agent := range agents {
		panel := CoordWorkbenchPanelForPost(index, agent, relatedLimit)
		if panel == nil {
			continue
		}
		cards = append(cards, AgentWorkbenchCard{
			Agent: agent,
			Panel: panel,
		})
	}
	return cards
}

func HomeTopicWorkbench(index Index, posts []Post, limit int) []RegistryWorkbenchCard {
	return homeRegistryWorkbench(index, posts, "topic", limit)
}

func HomeSourceWorkbench(index Index, posts []Post, limit int) []RegistryWorkbenchCard {
	return homeRegistryWorkbench(index, posts, "source", limit)
}

func homeRegistryWorkbench(index Index, posts []Post, kind string, limit int) []RegistryWorkbenchCard {
	if limit <= 0 {
		return nil
	}
	typedPosts := filterTypedPosts(posts)
	if len(typedPosts) == 0 {
		return nil
	}
	var stats []FacetStat
	switch kind {
	case "topic":
		stats = TopicStatsForPosts(typedPosts)
	case "source":
		stats = SourceStatsForPosts(typedPosts)
	default:
		return nil
	}
	cards := make([]RegistryWorkbenchCard, 0, limit)
	for _, stat := range stats {
		name := strings.TrimSpace(stat.Name)
		if name == "" {
			continue
		}
		if kind == "topic" && strings.EqualFold(name, "all") {
			continue
		}
		var scoped []Post
		var url string
		var externalURL string
		switch kind {
		case "topic":
			scoped = index.FilterPosts(FeedOptions{Topic: name})
			url = TopicPath(name)
		case "source":
			scoped = index.FilterPosts(FeedOptions{Source: name})
			url = SourcePath(name)
			externalURL = SourceURLFromPosts(scoped)
		}
		if len(scoped) == 0 {
			continue
		}
		cards = append(cards, RegistryWorkbenchCard{
			Kind:          kind,
			Name:          name,
			URL:           url,
			ExternalURL:   externalURL,
			Summary:       BuildCollectionSummaryStats(kind, scoped),
			ModuleActions: ScopedCoordActions(kind, name, scoped),
			JumpPanel:     WorkstreamJumpPanel(kind, name, scoped),
			PrimaryCount:  len(scoped),
			ReplyCount:    CountReplies(scoped),
			ReactionCount: CountReactions(scoped),
		})
		if len(cards) == limit {
			break
		}
	}
	return cards
}

func RecentCoordPosts(posts []Post, coordType string, limit int) []Post {
	return limitPosts(FilterPostsByCoordType(posts, coordType), limit)
}

func LegacyPosts(posts []Post, limit int) []Post {
	filtered := make([]Post, 0, len(posts))
	for _, post := range posts {
		if PostCoordType(post) != "" {
			continue
		}
		filtered = append(filtered, post)
	}
	return limitPosts(filtered, limit)
}

func filterTypedPosts(posts []Post) []Post {
	filtered := make([]Post, 0, len(posts))
	for _, post := range posts {
		if PostCoordType(post) == "" {
			continue
		}
		filtered = append(filtered, post)
	}
	return filtered
}

func limitPosts(posts []Post, limit int) []Post {
	if limit <= 0 || len(posts) == 0 {
		return nil
	}
	if len(posts) > limit {
		posts = posts[:limit]
	}
	out := make([]Post, len(posts))
	copy(out, posts)
	return out
}

func filterTaskPostsByStatus(posts []Post, limit int, statuses []string) []Post {
	allowed := make(map[string]struct{}, len(statuses))
	for _, status := range statuses {
		status = strings.ToLower(strings.TrimSpace(status))
		if status == "" {
			continue
		}
		allowed[status] = struct{}{}
	}
	filtered := make([]Post, 0, len(posts))
	for _, post := range FilterPostsByCoordType(posts, "task") {
		if _, ok := allowed[taskStatus(post)]; !ok {
			continue
		}
		filtered = append(filtered, post)
		if len(filtered) == limit {
			break
		}
	}
	return filtered
}

func taskStatus(post Post) string {
	return strings.ToLower(strings.TrimSpace(extensionText(post.Message.Extensions["task.status"])))
}

func isActiveTaskStatus(status string) bool {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "in_progress", "blocked", "queued", "ready":
		return true
	default:
		return false
	}
}

func metadataFieldsForSchema(post Post, fields []MetadataSchemaField) []MetadataField {
	out := make([]MetadataField, 0, len(fields))
	for _, field := range fields {
		value := extensionText(post.Message.Extensions[field.Key])
		if value == "" {
			continue
		}
		out = append(out, MetadataField{
			Label: field.Label,
			Value: value,
		})
	}
	return out
}

func extensionText(value any) string {
	switch typed := value.(type) {
	case nil:
		return ""
	case string:
		return strings.TrimSpace(typed)
	case []string:
		return strings.Join(typed, ", ")
	case []any:
		parts := make([]string, 0, len(typed))
		for _, item := range typed {
			text := extensionText(item)
			if text != "" {
				parts = append(parts, text)
			}
		}
		return strings.Join(parts, ", ")
	default:
		return strings.TrimSpace(fmt.Sprint(typed))
	}
}

func extensionStrings(value any) []string {
	switch typed := value.(type) {
	case nil:
		return nil
	case string:
		text := strings.TrimSpace(typed)
		if text == "" {
			return nil
		}
		return []string{text}
	case []string:
		out := make([]string, 0, len(typed))
		for _, item := range typed {
			item = strings.TrimSpace(item)
			if item != "" {
				out = append(out, item)
			}
		}
		return out
	case []any:
		out := make([]string, 0, len(typed))
		for _, item := range typed {
			out = append(out, extensionStrings(item)...)
		}
		return out
	default:
		text := strings.TrimSpace(fmt.Sprint(typed))
		if text == "" {
			return nil
		}
		return []string{text}
	}
}
