package newsplugin

import (
	"strings"
	"testing"
	"time"
)

func TestCoordSectionsForTaskPost(t *testing.T) {
	t.Parallel()

	post := Post{
		Bundle: Bundle{
			InfoHash: "task-1",
			Message: Message{
				Extensions: map[string]any{
					"coord.type":                "task",
					"task.id":                   "sharing-task-001",
					"task.status":               "in_progress",
					"task.priority":             "high",
					"task.parent":               "sharing-epic-website",
					"task.required_assets":      []any{"theme", "content-plugin"},
					"task.expected_result_type": "pages",
				},
			},
		},
	}

	sections := CoordSectionsForPost(post)
	if len(sections) != 2 {
		t.Fatalf("len = %d, want 2", len(sections))
	}
	if sections[0].Title != "Execution state" {
		t.Fatalf("first title = %q", sections[0].Title)
	}
	if sections[0].Items[0].Label != "Task ID" {
		t.Fatalf("first label = %q", sections[0].Items[0].Label)
	}
	if sections[1].Items[1].Value != "theme, content-plugin" {
		t.Fatalf("required assets = %q", sections[1].Items[1].Value)
	}
}

func TestCoordActionsForAgentPost(t *testing.T) {
	t.Parallel()

	post := Post{
		Bundle: Bundle{
			InfoHash:  "agent-1",
			ArchiveMD: "/tmp/archive/post-agent-1.md",
			Message: Message{
				Extensions: map[string]any{
					"coord.type":         "agent",
					"agent.availability": "internal",
				},
			},
		},
		Topics:        []string{"website"},
		SourceName:    "writer-key",
		HasSourcePage: true,
	}

	actions := CoordActionsForPost(post)
	if len(actions) < 2 {
		t.Fatalf("len = %d, want at least 2", len(actions))
	}
	var hasFiltered bool
	for _, action := range actions {
		if strings.Contains(action.URL, "/agents?meta_key=agent.availability") {
			hasFiltered = true
		}
	}
	if !hasFiltered {
		t.Fatal("expected filtered agent availability action")
	}
}

func TestCoordActionsForIdeaPost(t *testing.T) {
	t.Parallel()

	post := Post{
		Bundle: Bundle{
			InfoHash: "idea-1",
			Message: Message{
				Extensions: map[string]any{
					"coord.type":   "idea",
					"coord.domain": "website",
				},
			},
		},
		Topics: []string{"website"},
	}

	actions := CoordActionsForPost(post)
	var hasTasks, hasTopicTasks bool
	for _, action := range actions {
		if action.URL == "/tasks" {
			hasTasks = true
		}
		if action.URL == "/tasks?topic=website" {
			hasTopicTasks = true
		}
	}
	if !hasTasks {
		t.Fatal("expected task board action")
	}
	if !hasTopicTasks {
		t.Fatal("expected topic task board action")
	}
}

func TestCoordContextActions(t *testing.T) {
	t.Parallel()

	post := Post{
		Bundle: Bundle{
			InfoHash:  "agent-1",
			ArchiveMD: "/tmp/archive/post-agent-1.md",
			Message: Message{
				Extensions: map[string]any{
					"coord.type": "agent",
				},
			},
		},
		Topics:        []string{"all", "website"},
		SourceName:    "writer-key",
		HasSourcePage: true,
	}

	actions := CoordContextActions(post)
	if len(actions) != 4 {
		t.Fatalf("len = %d, want 4", len(actions))
	}
	if actions[0].URL != "/topics/website" {
		t.Fatalf("first url = %q", actions[0].URL)
	}
	if actions[1].URL != "/sources/writer-key" {
		t.Fatalf("second url = %q", actions[1].URL)
	}
	if actions[3].URL != "/archive/messages/agent-1" {
		t.Fatalf("last url = %q", actions[3].URL)
	}
}

func TestCoordPathForIdeaPost(t *testing.T) {
	t.Parallel()

	post := Post{
		Bundle: Bundle{
			InfoHash: "idea-1",
			Message: Message{
				Extensions: map[string]any{
					"coord.type": "idea",
				},
			},
		},
	}

	if got := CoordPath(post); got != "/ideas/idea-1" {
		t.Fatalf("CoordPath() = %q, want /ideas/idea-1", got)
	}
	if got := CoordAPIPath(post); got != "/api/ideas/idea-1" {
		t.Fatalf("CoordAPIPath() = %q, want /api/ideas/idea-1", got)
	}
	if got := CoordCanonicalPath(post); got != "/ideas/idea-1" {
		t.Fatalf("CoordCanonicalPath() = %q, want /ideas/idea-1", got)
	}
	if got := CoordCanonicalAPIPath(post); got != "/api/ideas/idea-1" {
		t.Fatalf("CoordCanonicalAPIPath() = %q, want /api/ideas/idea-1", got)
	}
}

func TestCoordCanonicalPathsStayEmptyForUntypedPost(t *testing.T) {
	t.Parallel()

	post := Post{
		Bundle: Bundle{
			InfoHash: "post-1",
			Message: Message{
				Title: "Legacy post",
			},
		},
	}

	if got := CoordCanonicalPath(post); got != "" {
		t.Fatalf("CoordCanonicalPath() = %q, want empty", got)
	}
	if got := CoordCanonicalAPIPath(post); got != "" {
		t.Fatalf("CoordCanonicalAPIPath() = %q, want empty", got)
	}
}

func TestCoordDetailFocusForKnowledgePost(t *testing.T) {
	t.Parallel()

	post := Post{
		Bundle: Bundle{
			InfoHash: "doc-1",
			Message: Message{
				Title: "Planning note",
				Extensions: map[string]any{
					"coord.type": "markdown",
				},
			},
		},
		Topics: []string{"website"},
	}
	index := Index{
		Posts: []Post{
			post,
			{
				Bundle: Bundle{
					InfoHash: "task-1",
					Message: Message{
						Title: "Ship homepage",
						Extensions: map[string]any{
							"coord.type": "task",
						},
					},
				},
				Topics: []string{"website"},
			},
			{
				Bundle: Bundle{
					InfoHash: "code-1",
					Message: Message{
						Title: "Theme entry point",
						Extensions: map[string]any{
							"coord.type": "code",
						},
					},
				},
				Topics: []string{"website"},
			},
		},
	}

	focus := CoordDetailFocusForPost(index, post, 1)
	if focus == nil {
		t.Fatal("CoordDetailFocusForPost() = nil")
	}
	if focus.Title != "Turn this document into execution context" {
		t.Fatalf("title = %q", focus.Title)
	}
	if len(focus.Checklist) != 3 {
		t.Fatalf("checklist len = %d, want 3", len(focus.Checklist))
	}
	var hasKnowledgeBase, hasTaskLink, hasCodeLink bool
	for _, action := range focus.Actions {
		switch action.URL {
		case "/knowledge":
			hasKnowledgeBase = true
		case "/tasks/task-1":
			hasTaskLink = true
		case "/code/code-1":
			hasCodeLink = true
		}
	}
	if !hasKnowledgeBase || !hasTaskLink || !hasCodeLink {
		t.Fatalf("actions = %+v", focus.Actions)
	}
}

func TestCoordDetailFocusForIdeaPost(t *testing.T) {
	t.Parallel()

	post := Post{
		Bundle: Bundle{
			InfoHash: "idea-1",
			Message: Message{
				Title: "Improve workspace navigation",
				Extensions: map[string]any{
					"coord.type": "idea",
				},
			},
		},
		Topics: []string{"website"},
	}
	index := Index{
		Posts: []Post{
			post,
			{
				Bundle: Bundle{
					InfoHash: "task-1",
					Message: Message{
						Title: "Ship registry-linked homepage",
						Extensions: map[string]any{
							"coord.type": "task",
						},
					},
				},
				Topics: []string{"website"},
			},
		},
	}

	focus := CoordDetailFocusForPost(index, post, 1)
	if focus == nil {
		t.Fatal("CoordDetailFocusForPost() = nil")
	}
	if focus.Title != "Promote this idea into active delivery" {
		t.Fatalf("title = %q", focus.Title)
	}
	var hasTaskBoard, hasTopicTasks, hasRelatedTask bool
	for _, action := range focus.Actions {
		switch action.URL {
		case "/tasks":
			hasTaskBoard = true
		case "/tasks?topic=website":
			hasTopicTasks = true
		case "/tasks/task-1":
			hasRelatedTask = true
		}
	}
	if !hasTaskBoard || !hasTopicTasks || !hasRelatedTask {
		t.Fatalf("actions = %+v", focus.Actions)
	}
}

func TestCoordDetailFocusForTaskPost(t *testing.T) {
	t.Parallel()

	post := Post{
		Bundle: Bundle{
			InfoHash: "task-1",
			Message: Message{
				Title: "Ship typed routes",
				Extensions: map[string]any{
					"coord.type":  "task",
					"task.status": "blocked",
				},
			},
		},
		Topics: []string{"website"},
	}
	index := Index{
		Posts: []Post{
			post,
			{
				Bundle: Bundle{
					InfoHash: "skill-1",
					Message: Message{
						Title: "Route authoring skill",
						Extensions: map[string]any{
							"coord.type": "skill",
						},
					},
				},
				Topics: []string{"website"},
			},
			{
				Bundle: Bundle{
					InfoHash: "agent-1",
					Message: Message{
						Title: "Delivery operator",
						Extensions: map[string]any{
							"coord.type": "agent",
						},
					},
				},
				Topics: []string{"website"},
			},
		},
	}

	focus := CoordDetailFocusForPost(index, post, 1)
	if focus == nil {
		t.Fatal("CoordDetailFocusForPost() = nil")
	}
	if focus.Title != "Clear blockers on this task" {
		t.Fatalf("title = %q", focus.Title)
	}
	var hasTaskBoard, hasBlockedFilter, hasSkillLink, hasAgentLink bool
	for _, action := range focus.Actions {
		switch action.URL {
		case "/tasks":
			hasTaskBoard = true
		case "/tasks?meta_key=task.status&meta_value=blocked":
			hasBlockedFilter = true
		case "/skills/skill-1":
			hasSkillLink = true
		case "/agents/agent-1":
			hasAgentLink = true
		}
	}
	if !hasTaskBoard || !hasBlockedFilter || !hasSkillLink || !hasAgentLink {
		t.Fatalf("actions = %+v", focus.Actions)
	}
}

func TestCoordDetailFocusForSkillPost(t *testing.T) {
	t.Parallel()

	post := Post{
		Bundle: Bundle{
			InfoHash: "skill-1",
			Message: Message{
				Title: "Seed runtime",
				Extensions: map[string]any{
					"coord.type":       "skill",
					"skill.category":   "ops",
					"skill.depends_on": []any{"identity"},
				},
			},
		},
		Topics: []string{"ops"},
	}
	index := Index{
		Posts: []Post{
			post,
			{
				Bundle: Bundle{
					InfoHash: "task-1",
					Message: Message{
						Title: "Run seed flow",
						Extensions: map[string]any{
							"coord.type": "task",
						},
					},
				},
				Topics: []string{"ops"},
			},
			{
				Bundle: Bundle{
					InfoHash: "agent-1",
					Message: Message{
						Title: "Ops operator",
						Extensions: map[string]any{
							"coord.type": "agent",
						},
					},
				},
				Topics: []string{"ops"},
			},
		},
	}

	focus := CoordDetailFocusForPost(index, post, 1)
	if focus == nil {
		t.Fatal("CoordDetailFocusForPost() = nil")
	}
	if focus.Title != "Attach this skill to active delivery" {
		t.Fatalf("title = %q", focus.Title)
	}
	var hasSkills, hasCategoryFilter, hasTaskLink, hasAgentLink bool
	for _, action := range focus.Actions {
		switch action.URL {
		case "/skills":
			hasSkills = true
		case "/skills?meta_key=skill.category&meta_value=ops":
			hasCategoryFilter = true
		case "/tasks/task-1":
			hasTaskLink = true
		case "/agents/agent-1":
			hasAgentLink = true
		}
	}
	if !hasSkills || !hasCategoryFilter || !hasTaskLink || !hasAgentLink {
		t.Fatalf("actions = %+v", focus.Actions)
	}
}

func TestCoordDetailFocusForAgentPost(t *testing.T) {
	t.Parallel()

	post := Post{
		Bundle: Bundle{
			InfoHash: "agent-1",
			Message: Message{
				Title: "Coordination operator",
				Extensions: map[string]any{
					"coord.type":         "agent",
					"agent.availability": "internal",
				},
			},
		},
		Topics: []string{"website"},
	}
	index := Index{
		Posts: []Post{
			post,
			{
				Bundle: Bundle{
					InfoHash: "task-1",
					Message: Message{
						Title: "Ship homepage",
						Extensions: map[string]any{
							"coord.type": "task",
						},
					},
				},
				Topics: []string{"website"},
			},
			{
				Bundle: Bundle{
					InfoHash: "skill-1",
					Message: Message{
						Title: "Theme rollout skill",
						Extensions: map[string]any{
							"coord.type": "skill",
						},
					},
				},
				Topics: []string{"website"},
			},
		},
	}

	focus := CoordDetailFocusForPost(index, post, 1)
	if focus == nil {
		t.Fatal("CoordDetailFocusForPost() = nil")
	}
	if focus.Title != "Route this agent into the right work" {
		t.Fatalf("title = %q", focus.Title)
	}
	var hasAgents, hasAvailabilityFilter, hasTaskLink, hasSkillLink bool
	for _, action := range focus.Actions {
		switch action.URL {
		case "/agents":
			hasAgents = true
		case "/agents?meta_key=agent.availability&meta_value=internal":
			hasAvailabilityFilter = true
		case "/tasks/task-1":
			hasTaskLink = true
		case "/skills/skill-1":
			hasSkillLink = true
		}
	}
	if !hasAgents || !hasAvailabilityFilter || !hasTaskLink || !hasSkillLink {
		t.Fatalf("actions = %+v", focus.Actions)
	}
}

func TestHomeWorkbenchLanesPreferActiveTasks(t *testing.T) {
	t.Parallel()

	posts := []Post{
		{
			Bundle: Bundle{
				InfoHash:  "task-active",
				CreatedAt: time.Date(2026, 3, 17, 10, 0, 0, 0, time.UTC),
				Message: Message{
					Title: "Active task",
					Extensions: map[string]any{
						"coord.type":  "task",
						"task.status": "in_progress",
					},
				},
			},
			Summary: "Active task summary",
		},
		{
			Bundle: Bundle{
				InfoHash:  "task-done",
				CreatedAt: time.Date(2026, 3, 17, 9, 0, 0, 0, time.UTC),
				Message: Message{
					Title: "Completed task",
					Extensions: map[string]any{
						"coord.type":  "task",
						"task.status": "done",
					},
				},
			},
		},
		{
			Bundle: Bundle{
				InfoHash:  "skill-1",
				CreatedAt: time.Date(2026, 3, 17, 8, 0, 0, 0, time.UTC),
				Message: Message{
					Title: "Skill asset",
					Extensions: map[string]any{
						"coord.type": "skill",
					},
				},
			},
		},
		{
			Bundle: Bundle{
				InfoHash:  "legacy-1",
				CreatedAt: time.Date(2026, 3, 17, 7, 0, 0, 0, time.UTC),
				Message: Message{
					Title: "Legacy note",
				},
			},
		},
	}

	lanes := HomeWorkbenchLanes(posts)
	if len(lanes) != 6 {
		t.Fatalf("len = %d, want 6", len(lanes))
	}
	if lanes[0].Title != "Active tasks" {
		t.Fatalf("first lane = %q", lanes[0].Title)
	}
	if len(lanes[0].Posts) != 1 || lanes[0].Posts[0].InfoHash != "task-active" {
		t.Fatalf("active tasks = %+v", lanes[0].Posts)
	}
	if len(lanes[1].Posts) != 1 || lanes[1].Posts[0].InfoHash != "skill-1" {
		t.Fatalf("skills lane = %+v", lanes[1].Posts)
	}
}

func TestScopedCoordActionsForTopic(t *testing.T) {
	t.Parallel()

	posts := []Post{
		{
			Bundle: Bundle{
				Message: Message{
					Extensions: map[string]any{"coord.type": "task"},
				},
			},
		},
		{
			Bundle: Bundle{
				Message: Message{
					Extensions: map[string]any{"coord.type": "task"},
				},
			},
		},
		{
			Bundle: Bundle{
				Message: Message{
					Extensions: map[string]any{"coord.type": "skill"},
				},
			},
		},
	}

	actions := ScopedCoordActions("topic", "website", posts)
	if len(actions) != 2 {
		t.Fatalf("len = %d, want 2", len(actions))
	}
	if actions[0].Label != "Tasks (2)" {
		t.Fatalf("first label = %q", actions[0].Label)
	}
	if actions[0].URL != "/tasks?topic=website" {
		t.Fatalf("first url = %q", actions[0].URL)
	}
	if actions[1].URL != "/skills?topic=website" {
		t.Fatalf("second url = %q", actions[1].URL)
	}
}

func TestAPIScopedCoordActionsForSource(t *testing.T) {
	t.Parallel()

	posts := []Post{
		{
			Bundle: Bundle{
				Message: Message{
					Extensions: map[string]any{"coord.type": "markdown"},
				},
			},
		},
	}

	actions := APIScopedCoordActions("source", "origin-1", posts)
	if len(actions) != 1 {
		t.Fatalf("len = %d, want 1", len(actions))
	}
	if actions[0].Label != "Knowledge (1)" {
		t.Fatalf("label = %q", actions[0].Label)
	}
	if actions[0].URL != "/api/knowledge?source=origin-1" {
		t.Fatalf("url = %q", actions[0].URL)
	}
}

func TestHomeTopicWorkbench(t *testing.T) {
	t.Parallel()

	index := Index{
		Posts: []Post{
			{
				Bundle: Bundle{
					Message: Message{
						Title: "Task one",
						Extensions: map[string]any{
							"coord.type": "task",
						},
					},
				},
				Topics: []string{"website"},
			},
			{
				Bundle: Bundle{
					Message: Message{
						Title: "Agent one",
						Extensions: map[string]any{
							"coord.type": "agent",
						},
					},
				},
				Topics: []string{"website"},
			},
			{
				Bundle: Bundle{
					Message: Message{
						Title: "Legacy note",
					},
				},
				Topics: []string{"website"},
			},
		},
	}

	cards := HomeTopicWorkbench(index, index.Posts, 3)
	if len(cards) != 1 {
		t.Fatalf("len = %d, want 1", len(cards))
	}
	if cards[0].Name != "website" {
		t.Fatalf("name = %q", cards[0].Name)
	}
	if cards[0].URL != "/topics/website" {
		t.Fatalf("url = %q", cards[0].URL)
	}
	if len(cards[0].ModuleActions) != 2 {
		t.Fatalf("module actions = %d, want 2", len(cards[0].ModuleActions))
	}
	if cards[0].JumpPanel == nil {
		t.Fatal("expected jump panel")
	}
	if cards[0].Summary[0].Label != "Visible assets" {
		t.Fatalf("summary label = %q", cards[0].Summary[0].Label)
	}
}

func TestHomePriorityCards(t *testing.T) {
	t.Parallel()

	posts := []Post{
		{
			Bundle: Bundle{
				Message: Message{
					Extensions: map[string]any{
						"coord.type":  "task",
						"task.status": "in_progress",
					},
				},
			},
			Topics: []string{"website"},
		},
		{
			Bundle: Bundle{
				Message: Message{
					Extensions: map[string]any{
						"coord.type": "idea",
					},
				},
			},
			Topics: []string{"website"},
		},
		{
			Bundle: Bundle{
				Message: Message{
					Extensions: map[string]any{
						"coord.type": "skill",
					},
				},
			},
			SourceName:    "source-a",
			HasSourcePage: true,
		},
	}

	cards := HomePriorityCards(posts)
	if len(cards) != 3 {
		t.Fatalf("len = %d, want 3", len(cards))
	}
	if cards[0].Title != "Execute active work" {
		t.Fatalf("first title = %q", cards[0].Title)
	}
	if cards[0].Stats[0].Value != "1" {
		t.Fatalf("active tasks = %q", cards[0].Stats[0].Value)
	}
	if cards[1].Stats[0].Label != "Live workstreams" {
		t.Fatalf("second first label = %q", cards[1].Stats[0].Label)
	}
	if cards[2].Actions[0].URL != "/skills" {
		t.Fatalf("third first action = %q", cards[2].Actions[0].URL)
	}
}

func TestHomeSourceWorkbench(t *testing.T) {
	t.Parallel()

	index := Index{
		Posts: []Post{
			{
				Bundle: Bundle{
					Message: Message{
						Title: "Knowledge asset",
						Extensions: map[string]any{
							"coord.type": "markdown",
						},
					},
				},
				SourceName:    "source-a",
				HasSourcePage: true,
				SourceURL:     "https://example.com/a",
			},
		},
	}

	cards := HomeSourceWorkbench(index, index.Posts, 3)
	if len(cards) != 1 {
		t.Fatalf("len = %d, want 1", len(cards))
	}
	if cards[0].URL != "/sources/source-a" {
		t.Fatalf("url = %q", cards[0].URL)
	}
	if cards[0].ExternalURL != "https://example.com/a" {
		t.Fatalf("external url = %q", cards[0].ExternalURL)
	}
	if len(cards[0].ModuleActions) != 1 {
		t.Fatalf("module actions = %d, want 1", len(cards[0].ModuleActions))
	}
	if cards[0].ModuleActions[0].URL != "/knowledge?source=source-a" {
		t.Fatalf("module action url = %q", cards[0].ModuleActions[0].URL)
	}
	if cards[0].JumpPanel == nil {
		t.Fatal("expected jump panel")
	}
}

func TestLegacyPostsOnlyReturnsUntypedPosts(t *testing.T) {
	t.Parallel()

	posts := []Post{
		{
			Bundle: Bundle{
				InfoHash: "typed-1",
				Message: Message{
					Extensions: map[string]any{
						"coord.type": "idea",
					},
				},
			},
		},
		{
			Bundle: Bundle{
				InfoHash: "legacy-1",
				Message:  Message{},
			},
		},
	}

	got := LegacyPosts(posts, 5)
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
	if got[0].InfoHash != "legacy-1" {
		t.Fatalf("infohash = %q", got[0].InfoHash)
	}
}

func TestCoordRelatedGroupsForTask(t *testing.T) {
	t.Parallel()

	index := Index{
		Posts: []Post{
			{
				Bundle: Bundle{
					InfoHash: "task-1",
					Message: Message{
						Extensions: map[string]any{
							"coord.type": "task",
						},
					},
				},
				Topics: []string{"website", "tasks"},
			},
			{
				Bundle: Bundle{
					InfoHash: "skill-1",
					Message: Message{
						Title: "Skill asset",
						Extensions: map[string]any{
							"coord.type": "skill",
						},
					},
				},
				Topics:  []string{"website", "skills"},
				Summary: "skill summary",
			},
			{
				Bundle: Bundle{
					InfoHash: "md-1",
					Message: Message{
						Title: "Knowledge asset",
						Extensions: map[string]any{
							"coord.type": "markdown",
						},
					},
				},
				Topics: []string{"website", "planning"},
			},
			{
				Bundle: Bundle{
					InfoHash: "code-1",
					Message: Message{
						Title: "Code asset",
						Extensions: map[string]any{
							"coord.type": "code",
						},
					},
				},
				Topics: []string{"website", "code"},
			},
			{
				Bundle: Bundle{
					InfoHash: "agent-1",
					Message: Message{
						Title: "Agent asset",
						Extensions: map[string]any{
							"coord.type": "agent",
						},
					},
				},
				Topics: []string{"website", "agents"},
			},
		},
	}

	groups := CoordRelatedGroups(index, index.Posts[0], 2)
	if len(groups) != 4 {
		t.Fatalf("len = %d, want 4", len(groups))
	}
	if groups[0].Kind != "skill" {
		t.Fatalf("first kind = %q", groups[0].Kind)
	}
	if groups[1].Kind != "markdown" {
		t.Fatalf("second kind = %q", groups[1].Kind)
	}
	if groups[2].Kind != "code" {
		t.Fatalf("third kind = %q", groups[2].Kind)
	}
	if groups[3].Kind != "agent" {
		t.Fatalf("fourth kind = %q", groups[3].Kind)
	}
}

func TestCoordExtendedContextGroupsForTaskAreDedupedAgainstWorkbench(t *testing.T) {
	t.Parallel()

	index := Index{
		Posts: []Post{
			{
				Bundle: Bundle{
					InfoHash: "task-1",
					Message: Message{
						Extensions: map[string]any{
							"coord.type": "task",
						},
					},
				},
				Topics: []string{"website"},
			},
			{
				Bundle: Bundle{
					InfoHash: "skill-1",
					Message: Message{
						Extensions: map[string]any{
							"coord.type": "skill",
						},
					},
				},
				Topics: []string{"website"},
			},
			{
				Bundle: Bundle{
					InfoHash: "md-1",
					Message: Message{
						Extensions: map[string]any{
							"coord.type": "markdown",
						},
					},
				},
				Topics: []string{"website"},
			},
			{
				Bundle: Bundle{
					InfoHash: "code-1",
					Message: Message{
						Extensions: map[string]any{
							"coord.type": "code",
						},
					},
				},
				Topics: []string{"website"},
			},
			{
				Bundle: Bundle{
					InfoHash: "agent-1",
					Message: Message{
						Extensions: map[string]any{
							"coord.type": "agent",
						},
					},
				},
				Topics: []string{"website"},
			},
		},
	}

	groups := CoordExtendedContextGroupsForPost(index, index.Posts[0], 2)
	if len(groups) != 0 {
		t.Fatalf("len = %d, want 0", len(groups))
	}
}

func TestCoordExtendedContextGroupsForIdeaRemainVisible(t *testing.T) {
	t.Parallel()

	index := Index{
		Posts: []Post{
			{
				Bundle: Bundle{
					InfoHash: "idea-1",
					Message: Message{
						Extensions: map[string]any{
							"coord.type": "idea",
						},
					},
				},
				Topics: []string{"website"},
			},
			{
				Bundle: Bundle{
					InfoHash: "task-1",
					Message: Message{
						Extensions: map[string]any{
							"coord.type": "task",
						},
					},
				},
				Topics: []string{"website"},
			},
			{
				Bundle: Bundle{
					InfoHash: "md-1",
					Message: Message{
						Extensions: map[string]any{
							"coord.type": "markdown",
						},
					},
				},
				Topics: []string{"website"},
			},
		},
	}

	groups := CoordExtendedContextGroupsForPost(index, index.Posts[0], 2)
	if len(groups) != 2 {
		t.Fatalf("len = %d, want 2", len(groups))
	}
	if groups[0].Kind != "task" {
		t.Fatalf("first kind = %q", groups[0].Kind)
	}
	if groups[1].Kind != "markdown" {
		t.Fatalf("second kind = %q", groups[1].Kind)
	}
}

func TestCoordRelationGroupsForTask(t *testing.T) {
	t.Parallel()

	index := Index{
		Posts: []Post{
			{
				Bundle: Bundle{
					InfoHash: "task-1",
					Message: Message{
						Extensions: map[string]any{
							"coord.type":                "task",
							"task.required_assets":      []any{"theme", "content-plugin"},
							"task.expected_result_type": "pages",
						},
					},
				},
				Topics: []string{"website"},
			},
			{
				Bundle: Bundle{
					InfoHash: "skill-1",
					Message: Message{
						Title: "Skill asset",
						Extensions: map[string]any{
							"coord.type":     "skill",
							"skill.category": "ops",
						},
					},
				},
				Topics: []string{"website"},
			},
			{
				Bundle: Bundle{
					InfoHash: "md-1",
					Message: Message{
						Title: "Knowledge asset",
						Extensions: map[string]any{
							"coord.type": "markdown",
						},
					},
				},
				Topics: []string{"website"},
			},
			{
				Bundle: Bundle{
					InfoHash: "code-1",
					Message: Message{
						Title: "Code asset",
						Extensions: map[string]any{
							"coord.type": "code",
						},
					},
				},
				Topics: []string{"website"},
			},
			{
				Bundle: Bundle{
					InfoHash: "agent-1",
					Message: Message{
						Title: "Agent asset",
						Extensions: map[string]any{
							"coord.type": "agent",
						},
					},
				},
				Topics: []string{"website"},
			},
		},
	}

	groups := CoordRelationGroupsForPost(index, index.Posts[0], 2)
	if len(groups) != 4 {
		t.Fatalf("len = %d, want 4", len(groups))
	}
	if groups[0].Key != "required_skills" || groups[0].Title != "Capability dependencies" {
		t.Fatalf("first group = %+v", groups[0])
	}
	if groups[1].Key != "supporting_knowledge" || groups[1].Kind != "markdown" {
		t.Fatalf("second group = %+v", groups[1])
	}
	if groups[2].Key != "implementation_assets" || groups[2].Kind != "code" {
		t.Fatalf("third group = %+v", groups[2])
	}
	if groups[3].Key != "operating_agents" || groups[3].Kind != "agent" {
		t.Fatalf("fourth group = %+v", groups[3])
	}
	if len(groups[0].Evidence) < 2 {
		t.Fatalf("required_skills evidence = %+v", groups[0].Evidence)
	}
	if groups[0].Evidence[0] != "Shared workstream: website" {
		t.Fatalf("first evidence = %q", groups[0].Evidence[0])
	}
	if groups[0].Evidence[1] != "Required assets: theme, content-plugin" {
		t.Fatalf("second evidence = %q", groups[0].Evidence[1])
	}
	if len(groups[0].Items) != 1 {
		t.Fatalf("required_skills items = %d, want 1", len(groups[0].Items))
	}
	if groups[0].Items[0].Post.InfoHash != "skill-1" {
		t.Fatalf("required_skills item = %q", groups[0].Items[0].Post.InfoHash)
	}
	if len(groups[0].Items[0].Evidence) < 2 {
		t.Fatalf("required_skills item evidence = %+v", groups[0].Items[0].Evidence)
	}
	if groups[0].Items[0].Evidence[1] != "Target category: ops" {
		t.Fatalf("expected target category evidence, got %+v", groups[0].Items[0].Evidence)
	}
	if groups[0].DominantReason != "Top match: Target category: ops · Required assets: theme, content-plugin" {
		t.Fatalf("dominant reason = %q", groups[0].DominantReason)
	}
}

func TestCoordRelationGroupsForIdea(t *testing.T) {
	t.Parallel()

	index := Index{
		Posts: []Post{
			{
				Bundle: Bundle{
					InfoHash: "idea-1",
					Message: Message{
						Extensions: map[string]any{
							"coord.type":            "idea",
							"coord.goal":            "Make AiP2P Sharing the default surface.",
							"coord.expected_output": "Typed pages",
						},
					},
				},
				Topics: []string{"website"},
			},
			{
				Bundle: Bundle{
					InfoHash: "task-1",
					Message: Message{
						Title: "Task asset",
						Extensions: map[string]any{
							"coord.type":                "task",
							"task.status":               "in_progress",
							"task.expected_result_type": "pages",
						},
					},
				},
				Topics: []string{"website"},
			},
			{
				Bundle: Bundle{
					InfoHash: "md-1",
					Message: Message{
						Title: "Knowledge asset",
						Extensions: map[string]any{
							"coord.type": "markdown",
						},
					},
				},
				Topics: []string{"website"},
			},
		},
	}

	groups := CoordRelationGroupsForPost(index, index.Posts[0], 2)
	if len(groups) != 2 {
		t.Fatalf("len = %d, want 2", len(groups))
	}
	if groups[0].Key != "promoted_tasks" || groups[0].Title != "Promoted tasks" {
		t.Fatalf("first group = %+v", groups[0])
	}
	if groups[1].Key != "shaping_notes" || groups[1].Kind != "markdown" {
		t.Fatalf("second group = %+v", groups[1])
	}
	if len(groups[0].Evidence) < 3 {
		t.Fatalf("promoted_tasks evidence = %+v", groups[0].Evidence)
	}
	if groups[0].Evidence[1] != "Goal: Make AiP2P Sharing the default surface." {
		t.Fatalf("goal evidence = %q", groups[0].Evidence[1])
	}
	if len(groups[0].Items) != 1 {
		t.Fatalf("promoted_tasks items = %d, want 1", len(groups[0].Items))
	}
	if len(groups[0].Items[0].Evidence) < 3 {
		t.Fatalf("promoted_tasks item evidence = %+v", groups[0].Items[0].Evidence)
	}
	if groups[0].Items[0].Evidence[1] != "Target status: in_progress" {
		t.Fatalf("target status evidence = %q", groups[0].Items[0].Evidence[1])
	}
	if groups[0].DominantReason != "Top match: Target status: in_progress · Goal: Make AiP2P Sharing the default surface." {
		t.Fatalf("dominant reason = %q", groups[0].DominantReason)
	}
}

func TestCoordRelationItemsAreSortedByScore(t *testing.T) {
	t.Parallel()

	index := Index{
		Posts: []Post{
			{
				Bundle: Bundle{
					InfoHash: "task-1",
					Message: Message{
						Title: "Ship typed theme routes",
						Extensions: map[string]any{
							"coord.type":           "task",
							"task.required_assets": []any{"theme", "routes"},
						},
					},
				},
				Summary: "Theme route rollout",
				Topics:  []string{"website"},
			},
			{
				Bundle: Bundle{
					InfoHash: "skill-strong",
					Message: Message{
						Title: "Theme route skill",
						Extensions: map[string]any{
							"coord.type":     "skill",
							"skill.category": "theme",
							"skill.outputs":  "routes",
						},
					},
				},
				Summary: "Theme route implementation",
				Topics:  []string{"website"},
			},
			{
				Bundle: Bundle{
					InfoHash: "skill-weak",
					Message: Message{
						Title: "Ops seed skill",
						Extensions: map[string]any{
							"coord.type":     "skill",
							"skill.category": "ops",
						},
					},
				},
				Summary: "Runtime bootstrap",
				Topics:  []string{"website"},
			},
		},
	}

	groups := CoordRelationGroupsForPost(index, index.Posts[0], 5)
	if len(groups) == 0 {
		t.Fatal("len = 0, want at least 1 group")
	}
	if len(groups[0].Items) != 2 {
		t.Fatalf("items len = %d, want 2", len(groups[0].Items))
	}
	if groups[0].Items[0].Post.InfoHash != "skill-strong" {
		t.Fatalf("first item = %q, want skill-strong", groups[0].Items[0].Post.InfoHash)
	}
	if groups[0].Items[0].Score <= groups[0].Items[1].Score {
		t.Fatalf("scores = %d, %d, want descending", groups[0].Items[0].Score, groups[0].Items[1].Score)
	}
}

func TestCoordRelationGroupsAreSortedByPriorityAndScore(t *testing.T) {
	t.Parallel()

	index := Index{
		Posts: []Post{
			{
				Bundle: Bundle{
					InfoHash: "task-1",
					Message: Message{
						Title: "Ship theme routes",
						Extensions: map[string]any{
							"coord.type":                "task",
							"task.required_assets":      []any{"theme", "routes"},
							"task.expected_result_type": "routes",
						},
					},
				},
				Summary: "Implement theme routes for the sharing workspace",
				Topics:  []string{"website"},
			},
			{
				Bundle: Bundle{
					InfoHash: "md-1",
					Message: Message{
						Title: "Planning note",
						Extensions: map[string]any{
							"coord.type":    "markdown",
							"md.collection": "planning",
							"md.kind":       "guide",
						},
					},
				},
				Summary: "General planning context",
				Topics:  []string{"website"},
			},
			{
				Bundle: Bundle{
					InfoHash: "code-1",
					Message: Message{
						Title: "Theme routes entry",
						Extensions: map[string]any{
							"coord.type":    "code",
							"code.kind":     "entrypoint",
							"code.entry":    "theme/routes",
							"code.language": "go",
						},
					},
				},
				Summary: "Implement theme routes in go",
				Topics:  []string{"website"},
			},
		},
	}

	groups := CoordRelationGroupsForPost(index, index.Posts[0], 5)
	if len(groups) != 2 {
		t.Fatalf("len = %d, want 2", len(groups))
	}
	if groups[0].Key != "implementation_assets" {
		t.Fatalf("first group = %q, want implementation_assets", groups[0].Key)
	}
	if groups[1].Key != "supporting_knowledge" {
		t.Fatalf("second group = %q, want supporting_knowledge", groups[1].Key)
	}
	if groups[0].Priority >= groups[1].Priority {
		t.Fatalf("priorities = %d, %d, want implementation_assets to win despite lower base priority", groups[0].Priority, groups[1].Priority)
	}
	if groups[0].Score <= groups[1].Score {
		t.Fatalf("scores = %d, %d, want implementation_assets to rank first", groups[0].Score, groups[1].Score)
	}
}

func TestWorkstreamRelationPaths(t *testing.T) {
	t.Parallel()

	posts := []Post{
		{
			Bundle: Bundle{
				InfoHash: "task-1",
				Message: Message{
					Title: "Ship theme routes",
					Extensions: map[string]any{
						"coord.type":                "task",
						"task.required_assets":      []any{"theme", "routes"},
						"task.expected_result_type": "routes",
					},
				},
			},
			Summary: "Implement theme routes for the workspace",
			Topics:  []string{"website"},
		},
		{
			Bundle: Bundle{
				InfoHash: "code-1",
				Message: Message{
					Title: "Theme route entry",
					Extensions: map[string]any{
						"coord.type":    "code",
						"code.entry":    "theme/routes",
						"code.language": "go",
					},
				},
			},
			Summary: "Go entry for theme routes",
			Topics:  []string{"website"},
		},
		{
			Bundle: Bundle{
				InfoHash: "skill-1",
				Message: Message{
					Title: "Theme skill",
					Extensions: map[string]any{
						"coord.type":     "skill",
						"skill.category": "theme",
					},
				},
			},
			Summary: "Skill for theme changes",
			Topics:  []string{"website"},
		},
	}

	paths := WorkstreamRelationPaths(posts, 4)
	if len(paths) == 0 {
		t.Fatal("len = 0, want at least 1 path")
	}
	if paths[0].From.InfoHash != "task-1" {
		t.Fatalf("first path source = %q, want task-1", paths[0].From.InfoHash)
	}
	if paths[0].Reason == "" {
		t.Fatal("reason = empty, want dominant reason")
	}
	if len(paths) != 3 {
		t.Fatalf("len = %d, want 3", len(paths))
	}
}

func TestWorkstreamRelationPanelForPosts(t *testing.T) {
	t.Parallel()

	posts := []Post{
		{
			Bundle: Bundle{
				InfoHash: "task-1",
				Message: Message{
					Title: "Task one",
					Extensions: map[string]any{
						"coord.type":                "task",
						"task.expected_result_type": "pages",
						"task.required_assets":      []any{"theme"},
					},
				},
			},
			Summary: "Primary task",
			Topics:  []string{"website"},
		},
		{
			Bundle: Bundle{
				InfoHash: "idea-1",
				Message: Message{
					Title: "Idea one",
					Extensions: map[string]any{
						"coord.type":            "idea",
						"coord.goal":            "Refresh the site",
						"coord.expected_output": "pages",
					},
				},
			},
			Summary: "Primary idea",
			Topics:  []string{"website"},
		},
		{
			Bundle: Bundle{
				InfoHash: "agent-1",
				Message: Message{
					Title: "Agent one",
					Extensions: map[string]any{
						"coord.type":         "agent",
						"agent.availability": "internal",
						"agent.skills":       "theme",
					},
				},
			},
			Summary: "Supporting operator",
			Topics:  []string{"website"},
		},
		{
			Bundle: Bundle{
				InfoHash: "code-1",
				Message: Message{
					Title: "Code one",
					Extensions: map[string]any{
						"coord.type":    "code",
						"code.entry":    "web/theme",
						"code.language": "go",
					},
				},
			},
			Summary: "Supporting implementation",
			Topics:  []string{"website"},
		},
	}

	panel := WorkstreamRelationPanelForPosts(posts, 4)
	if panel == nil {
		t.Fatal("panel = nil, want non-nil")
	}
	if len(panel.PrimaryPaths) != 2 {
		t.Fatalf("primary len = %d, want 2", len(panel.PrimaryPaths))
	}
	if len(panel.SupportingPaths) == 0 {
		t.Fatal("supporting len = 0, want at least 1")
	}
	if panel.PrimaryPaths[0].From.InfoHash != "task-1" {
		t.Fatalf("first primary path source = %q, want task-1", panel.PrimaryPaths[0].From.InfoHash)
	}
}

func TestHomeWorkbenchClustersForActiveTask(t *testing.T) {
	t.Parallel()

	index := Index{
		Posts: []Post{
			{
				Bundle: Bundle{
					InfoHash: "task-1",
					Message: Message{
						Title: "Active task",
						Extensions: map[string]any{
							"coord.type":  "task",
							"task.status": "in_progress",
						},
					},
				},
				Topics:        []string{"website", "tasks"},
				SourceName:    "source-a",
				HasSourcePage: true,
			},
			{
				Bundle: Bundle{
					InfoHash: "skill-1",
					Message: Message{
						Title: "Skill asset",
						Extensions: map[string]any{
							"coord.type": "skill",
						},
					},
				},
				Topics: []string{"website", "skills"},
			},
			{
				Bundle: Bundle{
					InfoHash: "md-1",
					Message: Message{
						Title: "Knowledge asset",
						Extensions: map[string]any{
							"coord.type": "markdown",
						},
					},
				},
				Topics: []string{"website", "planning"},
			},
		},
	}

	clusters := HomeWorkbenchClusters(index, index.Posts, 3, 2)
	if len(clusters) != 1 {
		t.Fatalf("len = %d, want 1", len(clusters))
	}
	if clusters[0].Task.InfoHash != "task-1" {
		t.Fatalf("task infohash = %q", clusters[0].Task.InfoHash)
	}
	if len(clusters[0].Related) != 2 {
		t.Fatalf("related len = %d, want 2", len(clusters[0].Related))
	}
	if clusters[0].Related[0].Kind != "skill" {
		t.Fatalf("first related kind = %q", clusters[0].Related[0].Kind)
	}
	if len(clusters[0].RegistryActions) != 3 {
		t.Fatalf("registry actions len = %d, want 3", len(clusters[0].RegistryActions))
	}
	if clusters[0].RegistryActions[0].URL != "/topics/website" {
		t.Fatalf("first registry action = %q", clusters[0].RegistryActions[0].URL)
	}
	if clusters[0].RegistryActions[2].URL != "/sources/source-a" {
		t.Fatalf("source registry action = %q", clusters[0].RegistryActions[2].URL)
	}
}

func TestHomeCoordGroupFacetsUsesTypedTopics(t *testing.T) {
	t.Parallel()

	posts := []Post{
		{
			Bundle: Bundle{
				InfoHash: "idea-1",
				Message: Message{
					Extensions: map[string]any{
						"coord.type": "idea",
					},
				},
			},
			Topics: []string{"all", "website", "coordination"},
		},
		{
			Bundle: Bundle{
				InfoHash: "task-1",
				Message: Message{
					Extensions: map[string]any{
						"coord.type": "task",
					},
				},
			},
			Topics: []string{"website", "tasks"},
		},
		{
			Bundle: Bundle{
				InfoHash: "legacy-1",
				Message:  Message{},
			},
			Topics: []string{"legacy"},
		},
	}

	facets := HomeCoordGroupFacets(posts, 4)
	if len(facets) != 3 {
		t.Fatalf("len = %d, want 3", len(facets))
	}
	if facets[0].Name != "website" {
		t.Fatalf("first facet = %q", facets[0].Name)
	}
	for _, facet := range facets {
		if facet.Name == "all" || facet.Name == "legacy" {
			t.Fatalf("unexpected facet %q", facet.Name)
		}
	}
}

func TestCoordWorkbenchPanelForTask(t *testing.T) {
	t.Parallel()

	index := Index{
		Posts: []Post{
			{
				Bundle: Bundle{
					InfoHash: "task-1",
					Message: Message{
						Extensions: map[string]any{
							"coord.type": "task",
						},
					},
				},
				Topics: []string{"website"},
			},
			{
				Bundle: Bundle{
					InfoHash: "skill-1",
					Message: Message{
						Extensions: map[string]any{
							"coord.type": "skill",
						},
					},
				},
				Topics: []string{"website"},
			},
			{
				Bundle: Bundle{
					InfoHash: "code-1",
					Message: Message{
						Extensions: map[string]any{
							"coord.type": "code",
						},
					},
				},
				Topics: []string{"website"},
			},
		},
	}

	panel := CoordWorkbenchPanelForPost(index, index.Posts[0], 2)
	if panel == nil {
		t.Fatal("expected panel")
	}
	if panel.Title != "Task delivery workspace" {
		t.Fatalf("title = %q", panel.Title)
	}
	if len(panel.Groups) != 2 {
		t.Fatalf("groups len = %d, want 2", len(panel.Groups))
	}
	if panel.Groups[0].Title != "Required skills" {
		t.Fatalf("first group title = %q", panel.Groups[0].Title)
	}
}

func TestCoordWorkbenchPanelForAgent(t *testing.T) {
	t.Parallel()

	index := Index{
		Posts: []Post{
			{
				Bundle: Bundle{
					InfoHash: "agent-1",
					Message: Message{
						Extensions: map[string]any{
							"coord.type": "agent",
						},
					},
				},
				Topics: []string{"website"},
			},
			{
				Bundle: Bundle{
					InfoHash: "task-1",
					Message: Message{
						Extensions: map[string]any{
							"coord.type": "task",
						},
					},
				},
				Topics: []string{"website"},
			},
			{
				Bundle: Bundle{
					InfoHash: "md-1",
					Message: Message{
						Extensions: map[string]any{
							"coord.type": "markdown",
						},
					},
				},
				Topics: []string{"website"},
			},
		},
	}

	panel := CoordWorkbenchPanelForPost(index, index.Posts[0], 2)
	if panel == nil {
		t.Fatal("expected panel")
	}
	if panel.Title != "Agent coordination workspace" {
		t.Fatalf("title = %q", panel.Title)
	}
	if len(panel.Groups) != 2 {
		t.Fatalf("groups len = %d, want 2", len(panel.Groups))
	}
}

func TestHomeAgentWorkbench(t *testing.T) {
	t.Parallel()

	index := Index{
		Posts: []Post{
			{
				Bundle: Bundle{
					InfoHash: "agent-1",
					Message: Message{
						Title: "Agent profile",
						Extensions: map[string]any{
							"coord.type": "agent",
						},
					},
				},
				Topics: []string{"website"},
			},
			{
				Bundle: Bundle{
					InfoHash: "task-1",
					Message: Message{
						Extensions: map[string]any{
							"coord.type": "task",
						},
					},
				},
				Topics: []string{"website"},
			},
		},
	}

	cards := HomeAgentWorkbench(index, index.Posts, 2, 2)
	if len(cards) != 1 {
		t.Fatalf("len = %d, want 1", len(cards))
	}
	if cards[0].Agent.InfoHash != "agent-1" {
		t.Fatalf("agent infohash = %q", cards[0].Agent.InfoHash)
	}
	if cards[0].Panel == nil || cards[0].Panel.Title != "Agent coordination workspace" {
		t.Fatalf("panel = %+v", cards[0].Panel)
	}
}

func TestWorkstreamJumpPanelForTopic(t *testing.T) {
	t.Parallel()

	posts := []Post{
		{
			Bundle: Bundle{
				InfoHash: "task-1",
				Message: Message{
					Extensions: map[string]any{
						"coord.type":  "task",
						"task.status": "in_progress",
					},
				},
			},
			Topics: []string{"website"},
		},
		{
			Bundle: Bundle{
				InfoHash: "idea-1",
				Message: Message{
					Extensions: map[string]any{
						"coord.type": "idea",
					},
				},
			},
			Topics: []string{"website"},
		},
		{
			Bundle: Bundle{
				InfoHash: "skill-1",
				Message: Message{
					Extensions: map[string]any{
						"coord.type": "skill",
					},
				},
			},
			Topics: []string{"website"},
		},
		{
			Bundle: Bundle{
				InfoHash: "md-1",
				Message: Message{
					Extensions: map[string]any{
						"coord.type": "markdown",
					},
				},
			},
			Topics: []string{"website"},
		},
	}

	panel := WorkstreamJumpPanel("topic", "website", posts)
	if panel == nil {
		t.Fatal("expected panel")
	}
	if panel.Title != "Primary and supporting jumps for this workstream" {
		t.Fatalf("title = %q", panel.Title)
	}
	if len(panel.PrimaryActions) != 3 {
		t.Fatalf("primary len = %d, want 3", len(panel.PrimaryActions))
	}
	if panel.PrimaryActions[0].URL != "/tasks?meta_key=task.status&meta_value=in_progress&topic=website" {
		t.Fatalf("first primary = %q", panel.PrimaryActions[0].URL)
	}
	if len(panel.SupportingActions) == 0 || panel.SupportingActions[len(panel.SupportingActions)-1].URL != "/sources" {
		t.Fatalf("supporting = %+v", panel.SupportingActions)
	}
}

func TestClusterJumpPanel(t *testing.T) {
	t.Parallel()

	task := Post{
		Bundle: Bundle{
			InfoHash: "task-1",
			Message: Message{
				Extensions: map[string]any{
					"coord.type": "task",
				},
			},
		},
		Topics:        []string{"website"},
		SourceName:    "writer-a",
		HasSourcePage: true,
	}

	panel := ClusterJumpPanel(task)
	if panel == nil {
		t.Fatal("expected panel")
	}
	if panel.Title != "Primary and supporting jumps for this task" {
		t.Fatalf("title = %q", panel.Title)
	}
	if len(panel.PrimaryActions) < 3 {
		t.Fatalf("primary len = %d", len(panel.PrimaryActions))
	}
	if panel.PrimaryActions[0].URL != "/tasks/task-1" {
		t.Fatalf("first primary = %q", panel.PrimaryActions[0].URL)
	}
	if panel.PrimaryActions[1].URL != "/topics/website" {
		t.Fatalf("second primary = %q", panel.PrimaryActions[1].URL)
	}
	if panel.PrimaryActions[2].URL != "/sources/writer-a" {
		t.Fatalf("third primary = %q", panel.PrimaryActions[2].URL)
	}
	if len(panel.SupportingActions) == 0 || panel.SupportingActions[0].URL != "/tasks?meta_key=task.status&meta_value=in_progress" {
		t.Fatalf("supporting = %+v", panel.SupportingActions)
	}
}

func TestDetailJumpPanel(t *testing.T) {
	t.Parallel()

	post := Post{
		Bundle: Bundle{
			InfoHash:  "task-1",
			ArchiveMD: "/tmp/archive/task-1.md",
			Message: Message{
				Extensions: map[string]any{
					"coord.type":  "task",
					"task.id":     "sharing-task-001",
					"task.status": "in_progress",
				},
			},
		},
		Topics:        []string{"website"},
		SourceName:    "writer-a",
		HasSourcePage: true,
	}

	panel := DetailJumpPanel(post)
	if panel == nil {
		t.Fatal("expected panel")
	}
	if panel.Title != "Primary and supporting jumps for this task" {
		t.Fatalf("title = %q", panel.Title)
	}
	if len(panel.PrimaryActions) == 0 || panel.PrimaryActions[0].URL != "/tasks" {
		t.Fatalf("primary = %+v", panel.PrimaryActions)
	}
	if len(panel.SupportingActions) == 0 || panel.SupportingActions[0].URL != "/topics/website" {
		t.Fatalf("supporting = %+v", panel.SupportingActions)
	}
}

func TestBuildPostPageDataFiltersRelatedByCoordType(t *testing.T) {
	t.Parallel()

	app := &App{
		storeRoot: t.TempDir(),
		project:   "aip2p.sharing",
		version:   "test",
	}
	base := Post{
		Bundle: Bundle{
			InfoHash: "task-1",
			Message: Message{
				Title: "Ship workspace",
				Extensions: map[string]any{
					"coord.type": "task",
				},
			},
			CreatedAt: time.Unix(10, 0).UTC(),
		},
		Topics:       []string{"website"},
		SourceName:   "writer-a",
		ChannelGroup: "work",
	}
	taskRelated := Post{
		Bundle: Bundle{
			InfoHash: "task-2",
			Message: Message{
				Title: "Refine registry",
				Extensions: map[string]any{
					"coord.type": "task",
				},
			},
			CreatedAt: time.Unix(9, 0).UTC(),
		},
		Topics:       []string{"website"},
		SourceName:   "writer-a",
		ChannelGroup: "work",
	}
	skillRelated := Post{
		Bundle: Bundle{
			InfoHash: "skill-1",
			Message: Message{
				Title: "Theme ops",
				Extensions: map[string]any{
					"coord.type": "skill",
				},
			},
			CreatedAt: time.Unix(8, 0).UTC(),
		},
		Topics:       []string{"website"},
		SourceName:   "writer-a",
		ChannelGroup: "work",
	}
	index := Index{
		Posts: []Post{base, taskRelated, skillRelated},
		PostByInfoHash: map[string]Post{
			"task-1":  base,
			"task-2":  taskRelated,
			"skill-1": skillRelated,
		},
		RepliesByPost: map[string][]Reply{
			"task-1": {{Bundle: Bundle{InfoHash: "reply-1"}}},
		},
		ReactionsByPost: map[string][]Reaction{
			"task-1": {{Bundle: Bundle{InfoHash: "reaction-1"}}},
		},
	}

	data := BuildPostPageData(app, index, "/tasks", base, "task")
	if data.Project != "AiP2P Sharing" {
		t.Fatalf("project = %q", data.Project)
	}
	if len(data.Related) != 1 || data.Related[0].InfoHash != "task-2" {
		t.Fatalf("related = %+v", data.Related)
	}
	if len(data.Replies) != 1 || data.Replies[0].InfoHash != "reply-1" {
		t.Fatalf("replies = %+v", data.Replies)
	}
	if len(data.Reactions) != 1 || data.Reactions[0].InfoHash != "reaction-1" {
		t.Fatalf("reactions = %+v", data.Reactions)
	}
}

func TestBuildPostAPIResponseUsesScopeAndUnfilteredRelatedForGenericPost(t *testing.T) {
	t.Parallel()

	app := &App{
		storeRoot: t.TempDir(),
		project:   "aip2p.sharing",
		version:   "test",
	}
	base := Post{
		Bundle: Bundle{
			InfoHash: "legacy-1",
			Message: Message{
				Title: "Legacy asset",
			},
			CreatedAt: time.Unix(10, 0).UTC(),
		},
		Topics:       []string{"website"},
		SourceName:   "writer-a",
		ChannelGroup: "work",
	}
	taskRelated := Post{
		Bundle: Bundle{
			InfoHash: "task-2",
			Message: Message{
				Extensions: map[string]any{
					"coord.type": "task",
				},
			},
			CreatedAt: time.Unix(9, 0).UTC(),
		},
		Topics:       []string{"website"},
		SourceName:   "writer-a",
		ChannelGroup: "work",
	}
	skillRelated := Post{
		Bundle: Bundle{
			InfoHash: "skill-1",
			Message: Message{
				Extensions: map[string]any{
					"coord.type": "skill",
				},
			},
			CreatedAt: time.Unix(8, 0).UTC(),
		},
		Topics:       []string{"website"},
		SourceName:   "writer-a",
		ChannelGroup: "work",
	}
	index := Index{
		Posts: []Post{base, taskRelated, skillRelated},
		PostByInfoHash: map[string]Post{
			"legacy-1": base,
			"task-2":   taskRelated,
			"skill-1":  skillRelated,
		},
		RepliesByPost:   map[string][]Reply{},
		ReactionsByPost: map[string][]Reaction{},
	}

	payload := BuildPostAPIResponse(app, index, "post", base, "")
	if payload["scope"] != "post" {
		t.Fatalf("scope = %v", payload["scope"])
	}
	related, ok := payload["related"].([]map[string]any)
	if !ok {
		t.Fatalf("related type = %T", payload["related"])
	}
	if len(related) != 2 {
		t.Fatalf("related len = %d, want 2", len(related))
	}
}

func TestBuildPostAPIResponseIncludesNavigationEnvelope(t *testing.T) {
	t.Parallel()

	app := &App{
		storeRoot: t.TempDir(),
		project:   "aip2p.sharing",
		version:   "test",
	}
	post := Post{
		Bundle: Bundle{
			InfoHash:  "task-1",
			ArchiveMD: "/tmp/archive/task-1.md",
			Message: Message{
				Extensions: map[string]any{
					"coord.type":  "task",
					"task.id":     "sharing-task-001",
					"task.status": "in_progress",
				},
			},
		},
		Topics:        []string{"website"},
		SourceName:    "writer-a",
		HasSourcePage: true,
	}
	index := Index{
		Posts:           []Post{post},
		PostByInfoHash:  map[string]Post{"task-1": post},
		RepliesByPost:   map[string][]Reply{},
		ReactionsByPost: map[string][]Reaction{},
	}

	payload := BuildPostAPIResponse(app, index, "task", post, "task")
	navigation, ok := payload["navigation"].(map[string]any)
	if !ok {
		t.Fatalf("navigation type = %T", payload["navigation"])
	}
	if navigation["canonical_field"] != "jump_panel" {
		t.Fatalf("canonical_field = %v", navigation["canonical_field"])
	}
	if _, ok := navigation["jump_panel"]; ok {
		t.Fatalf("navigation = %+v, unexpected nested jump_panel duplicate", navigation)
	}
	legacy, ok := navigation["legacy_fields"].(map[string]any)
	if !ok {
		t.Fatalf("legacy_fields type = %T", navigation["legacy_fields"])
	}
	if _, ok := payload["coord_actions"]; ok {
		t.Fatalf("payload = %+v, unexpected coord_actions top-level alias", payload)
	}
	if _, ok := payload["context_links"]; ok {
		t.Fatalf("payload = %+v, unexpected context_links top-level alias", payload)
	}
	if _, ok := legacy["coord_actions"]; !ok {
		t.Fatalf("legacy_fields = %+v", legacy)
	}
	if _, ok := legacy["context_links"]; !ok {
		t.Fatalf("legacy_fields = %+v", legacy)
	}
	schema, ok := payload["metadata_schema"].(map[string]any)
	if !ok {
		t.Fatalf("metadata_schema type = %T", payload["metadata_schema"])
	}
	if schema["coord_type"] != "task" {
		t.Fatalf("coord_type = %v, want task", schema["coord_type"])
	}
}

func TestBuildPostPageDataIncludesMetadataSchema(t *testing.T) {
	t.Parallel()

	app := &App{project: "aip2p.sharing", version: "0.1.0"}
	post := Post{
		Bundle: Bundle{
			InfoHash: "task-1",
			Message: Message{
				Title: "Task 1",
				Extensions: map[string]any{
					"coord.type":  "task",
					"task.id":     "sharing-task-001",
					"task.status": "in_progress",
				},
			},
		},
		Topics:     []string{"all", "website"},
		SourceName: "writer-a",
	}
	index := Index{}
	data := BuildPostPageData(app, index, "/tasks", post, "task")
	if data.MetadataSchema == nil {
		t.Fatal("metadata schema = nil")
	}
	if data.MetadataSchema.CoordType != "task" {
		t.Fatalf("coord type = %q, want task", data.MetadataSchema.CoordType)
	}
	if data.PublishGuide == nil {
		t.Fatal("publish guide = nil")
	}
	if !strings.Contains(data.PublishGuide.StarterExtensionsJSON, "\"coord.type\": \"task\"") {
		t.Fatalf("starter json = %q", data.PublishGuide.StarterExtensionsJSON)
	}
	if !strings.Contains(data.PublishGuide.StarterCommand, "aip2p.sharing/tasks") {
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
	if !strings.Contains(data.PublishGuide.Templates[0].StarterCommand, "--kind reply") {
		t.Fatalf("template command = %q", data.PublishGuide.Templates[0].StarterCommand)
	}
	if !strings.Contains(data.PublishGuide.Templates[0].StarterCommand, post.InfoHash) {
		t.Fatalf("template command = %q", data.PublishGuide.Templates[0].StarterCommand)
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
	if !strings.Contains(data.PublishGuide.Templates[4].StarterCommand, post.Message.Extensions["task.id"].(string)) {
		t.Fatalf("template command = %q", data.PublishGuide.Templates[4].StarterCommand)
	}
	if !strings.Contains(data.PublishGuide.Templates[6].StarterExtensionsJSON, "\"code.kind\": \"implementation\"") {
		t.Fatalf("template json = %q", data.PublishGuide.Templates[6].StarterExtensionsJSON)
	}
	if data.PublishGuide.Templates[7].Kind != "reaction" {
		t.Fatalf("template kind = %q, want reaction", data.PublishGuide.Templates[7].Kind)
	}
	if !strings.Contains(data.PublishGuide.Templates[7].StarterExtensionsJSON, "\"reaction_type\": \"vote\"") {
		t.Fatalf("template json = %q", data.PublishGuide.Templates[7].StarterExtensionsJSON)
	}
	if !strings.Contains(data.PublishGuide.Templates[8].StarterExtensionsJSON, "\"value\": -1") {
		t.Fatalf("template json = %q", data.PublishGuide.Templates[8].StarterExtensionsJSON)
	}
}

func TestCoordPublishGuideForIdeaIncludesKnowledgeCodeAndReactionTemplates(t *testing.T) {
	t.Parallel()

	guide := CoordPublishGuideFor("aip2p.sharing", "idea", CoordPublishContext{
		Topics:   []string{"website", "coordination"},
		InfoHash: "idea-1",
		Magnet:   "magnet:?xt=urn:btih:idea-1",
	})
	if guide == nil {
		t.Fatal("guide = nil")
	}
	if len(guide.Templates) != 4 {
		t.Fatalf("templates len = %d, want 4", len(guide.Templates))
	}
	if guide.Templates[0].Label != "Promote to task" {
		t.Fatalf("first label = %q", guide.Templates[0].Label)
	}
	if !strings.Contains(guide.Templates[1].StarterExtensionsJSON, "\"md.kind\": \"framing-note\"") {
		t.Fatalf("template json = %q", guide.Templates[1].StarterExtensionsJSON)
	}
	if !strings.Contains(guide.Templates[2].StarterExtensionsJSON, "\"code.kind\": \"prototype\"") {
		t.Fatalf("template json = %q", guide.Templates[2].StarterExtensionsJSON)
	}
	if guide.Templates[3].Kind != "reaction" {
		t.Fatalf("template kind = %q, want reaction", guide.Templates[3].Kind)
	}
	if !strings.Contains(guide.Templates[3].StarterExtensionsJSON, "\"reaction_type\": \"vote\"") {
		t.Fatalf("template json = %q", guide.Templates[3].StarterExtensionsJSON)
	}
	if !strings.Contains(guide.Templates[3].StarterExtensionsJSON, "\"infohash\": \"idea-1\"") {
		t.Fatalf("template json = %q", guide.Templates[3].StarterExtensionsJSON)
	}
	if len(guide.Workflow) != 2 {
		t.Fatalf("workflow len = %d, want 2", len(guide.Workflow))
	}
	if len(guide.Workflow[0].Templates) != 2 {
		t.Fatalf("workflow templates len = %d, want 2", len(guide.Workflow[0].Templates))
	}
}

func TestCoordPublishGuideForAgentIncludesSkillAndOperatorNoteTemplates(t *testing.T) {
	t.Parallel()

	guide := CoordPublishGuideFor("aip2p.sharing", "agent", CoordPublishContext{
		Topics: []string{"agents", "website"},
		Extensions: map[string]any{
			"agent.id": "codex-site-operator",
		},
	})
	if guide == nil {
		t.Fatal("guide = nil")
	}
	if len(guide.Templates) != 3 {
		t.Fatalf("templates len = %d, want 3", len(guide.Templates))
	}
	if guide.Templates[0].Label != "Availability update" {
		t.Fatalf("first label = %q", guide.Templates[0].Label)
	}
	if guide.Templates[1].Channel != "aip2p.sharing/skills" {
		t.Fatalf("template channel = %q, want aip2p.sharing/skills", guide.Templates[1].Channel)
	}
	if !strings.Contains(guide.Templates[1].StarterCommand, "codex-site-operator") {
		t.Fatalf("template command = %q", guide.Templates[1].StarterCommand)
	}
	if !strings.Contains(guide.Templates[2].StarterExtensionsJSON, "\"md.kind\": \"operator-note\"") {
		t.Fatalf("template json = %q", guide.Templates[2].StarterExtensionsJSON)
	}
	if guide.Templates[1].Category != "Linked asset" {
		t.Fatalf("template category = %q, want Linked asset", guide.Templates[1].Category)
	}
}

func TestCoordPublishGuideForKnowledgeIncludesImplementationAndTaskTemplates(t *testing.T) {
	t.Parallel()

	guide := CoordPublishGuideFor("aip2p.sharing", "markdown", CoordPublishContext{
		Topics: []string{"website", "knowledge"},
	})
	if guide == nil {
		t.Fatal("guide = nil")
	}
	if len(guide.Templates) != 2 {
		t.Fatalf("templates len = %d, want 2", len(guide.Templates))
	}
	if guide.Templates[0].Channel != "aip2p.sharing/code" {
		t.Fatalf("template channel = %q, want aip2p.sharing/code", guide.Templates[0].Channel)
	}
	if !strings.Contains(guide.Templates[0].StarterExtensionsJSON, "\"code.kind\": \"implementation\"") {
		t.Fatalf("template json = %q", guide.Templates[0].StarterExtensionsJSON)
	}
	if guide.Templates[1].Label != "Promote to task" {
		t.Fatalf("second label = %q", guide.Templates[1].Label)
	}
}

func TestCoordPublishGuideForCodeIncludesKnowledgeAndSkillTemplates(t *testing.T) {
	t.Parallel()

	guide := CoordPublishGuideFor("aip2p.sharing", "code", CoordPublishContext{
		Topics: []string{"website", "code"},
	})
	if guide == nil {
		t.Fatal("guide = nil")
	}
	if len(guide.Templates) != 2 {
		t.Fatalf("templates len = %d, want 2", len(guide.Templates))
	}
	if !strings.Contains(guide.Templates[0].StarterExtensionsJSON, "\"md.kind\": \"implementation-note\"") {
		t.Fatalf("template json = %q", guide.Templates[0].StarterExtensionsJSON)
	}
	if guide.Templates[1].Channel != "aip2p.sharing/skills" {
		t.Fatalf("template channel = %q, want aip2p.sharing/skills", guide.Templates[1].Channel)
	}
}
