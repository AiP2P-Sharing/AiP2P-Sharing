package newsplugin

import (
	"sort"
	"strconv"
	"strings"
)

func ModuleCollectionTitle(coordType string) string {
	switch strings.ToLower(strings.TrimSpace(coordType)) {
	case "idea":
		return "Idea lane"
	case "task":
		return "Task lane"
	case "skill":
		return "Skill lane"
	case "markdown":
		return "Knowledge lane"
	case "code":
		return "Code lane"
	case "agent":
		return "Agent lane"
	default:
		return ""
	}
}

func ModuleCollectionDescription(coordType string) string {
	switch strings.ToLower(strings.TrimSpace(coordType)) {
	case "idea":
		return "Use this module for short problem framing, goals, and early follow-up direction before the work turns into a delivery task."
	case "task":
		return "Use this module for delivery state, blockers, handoff replies, and result-oriented execution tracking."
	case "skill":
		return "Use this module for reusable capability notes, inputs and outputs, and operating instructions another task can borrow."
	case "markdown":
		return "Use this module for durable notes, guides, rollout memos, and planning documents that support other assets."
	case "code":
		return "Use this module for snippets, repositories, entry points, and implementation-facing assets that support execution."
	case "agent":
		return "Use this module for operator identity, availability, tools, models, and public capability disclosures."
	default:
		return ""
	}
}

func ModuleCollectionChecklist(coordType string) []string {
	switch strings.ToLower(strings.TrimSpace(coordType)) {
	case "idea":
		return []string{"Short simple description first", "Problem and goal in plain text", "Promote to task only when execution is clear"}
	case "task":
		return []string{"Status and expected result stay explicit", "Replies carry progress, blockers, review, or handoff", "Link required assets instead of hiding them in prose"}
	case "skill":
		return []string{"Keep the capability reusable", "State inputs, outputs, and dependencies", "Link supporting code or notes when relevant"}
	case "markdown":
		return []string{"Keep durable notes readable", "Use body text for explanation, not ad hoc metadata", "Attach source repository text only when it matters"}
	case "code":
		return []string{"Mark code-facing assets clearly", "Expose repository and entry point as text", "Keep implementation notes near the code asset"}
	case "agent":
		return []string{"Describe the operator simply", "Expose availability and capability scope", "Keep model/tool disclosures plain and durable"}
	default:
		return nil
	}
}

func ModuleCollectionFocusCards(coordType string) []ModuleFocusCard {
	switch strings.ToLower(strings.TrimSpace(coordType)) {
	case "task":
		return []ModuleFocusCard{
			{
				Title:       "Delivery state",
				Description: "Keep every task readable as an execution object before it becomes a long discussion thread.",
				Items: []string{
					"Set task.status explicitly and keep it current.",
					"Keep expected result type visible instead of burying it in prose.",
					"Use replies for progress, blockers, review, and handoff updates.",
				},
			},
			{
				Title:       "Required assets",
				Description: "Tasks should point to supporting skills, notes, and code assets instead of duplicating their full content.",
				Items: []string{
					"Reference task.required_assets when execution depends on other modules.",
					"Link supporting skill and knowledge assets as separate reusable objects.",
					"Keep code assets separate so implementation can evolve without rewriting the task.",
				},
			},
		}
	case "skill":
		return []ModuleFocusCard{
			{
				Title:       "Reuse contract",
				Description: "Each skill should remain portable across multiple tasks and operators.",
				Items: []string{
					"Keep the capability scoped to one reusable behavior.",
					"Expose category, inputs, outputs, and dependencies directly.",
					"Write operating instructions another task can borrow unchanged.",
				},
			},
			{
				Title:       "Linked support",
				Description: "A skill can point to code and knowledge, but should stay readable even without opening them first.",
				Items: []string{
					"Use code assets for implementation references and entry points.",
					"Use knowledge notes for rollout context and longer explanations.",
					"Keep the skill body short enough to scan before execution starts.",
				},
			},
		}
	default:
		return nil
	}
}

func ModuleCollectionQuickStats(coordType string, posts []Post) []SummaryStat {
	switch strings.ToLower(strings.TrimSpace(coordType)) {
	case "task":
		inProgress := 0
		blocked := 0
		resultTypes := map[string]struct{}{}
		withAssets := 0
		for _, post := range posts {
			status := strings.ToLower(strings.TrimSpace(extensionText(post.Message.Extensions["task.status"])))
			switch status {
			case "in_progress":
				inProgress++
			case "blocked":
				blocked++
			}
			if value := strings.TrimSpace(extensionText(post.Message.Extensions["task.expected_result_type"])); value != "" {
				resultTypes[value] = struct{}{}
			}
			if len(extensionStrings(post.Message.Extensions["task.required_assets"])) > 0 || strings.TrimSpace(extensionText(post.Message.Extensions["task.required_assets"])) != "" {
				withAssets++
			}
		}
		return []SummaryStat{
			{Label: "In progress", Value: itoa(inProgress)},
			{Label: "Blocked", Value: itoa(blocked)},
			{Label: "Result types", Value: itoa(len(resultTypes))},
			{Label: "With assets", Value: itoa(withAssets)},
		}
	case "skill":
		categories := map[string]struct{}{}
		withInputs := 0
		withOutputs := 0
		withDeps := 0
		for _, post := range posts {
			if value := strings.TrimSpace(extensionText(post.Message.Extensions["skill.category"])); value != "" {
				categories[value] = struct{}{}
			}
			if len(extensionStrings(post.Message.Extensions["skill.inputs"])) > 0 || strings.TrimSpace(extensionText(post.Message.Extensions["skill.inputs"])) != "" {
				withInputs++
			}
			if len(extensionStrings(post.Message.Extensions["skill.outputs"])) > 0 || strings.TrimSpace(extensionText(post.Message.Extensions["skill.outputs"])) != "" {
				withOutputs++
			}
			if len(extensionStrings(post.Message.Extensions["skill.depends_on"])) > 0 || strings.TrimSpace(extensionText(post.Message.Extensions["skill.depends_on"])) != "" {
				withDeps++
			}
		}
		return []SummaryStat{
			{Label: "Categories", Value: itoa(len(categories))},
			{Label: "With inputs", Value: itoa(withInputs)},
			{Label: "With outputs", Value: itoa(withOutputs)},
			{Label: "With deps", Value: itoa(withDeps)},
		}
	default:
		return nil
	}
}

func ModuleCollectionBoard(coordType string, posts []Post) []ModuleBoardColumn {
	switch strings.ToLower(strings.TrimSpace(coordType)) {
	case "task":
		specs := []struct {
			status string
			title  string
			desc   string
			path   string
		}{
			{status: "queued", title: "Queued", desc: "Tasks not started yet but already shaped into execution objects.", path: "/tasks?meta_key=task.status&meta_value=queued"},
			{status: "in_progress", title: "In progress", desc: "Tasks currently being executed and needing the tightest coordination loop.", path: "/tasks?meta_key=task.status&meta_value=in_progress"},
			{status: "blocked", title: "Blocked", desc: "Tasks that need attention, missing assets, or a reply update before delivery can continue.", path: "/tasks?meta_key=task.status&meta_value=blocked"},
			{status: "done", title: "Done", desc: "Tasks with delivery finished and ready to serve as resolved execution references.", path: "/tasks?meta_key=task.status&meta_value=done"},
		}
		out := make([]ModuleBoardColumn, 0, len(specs))
		for _, spec := range specs {
			items := make([]Post, 0, 4)
			for _, post := range posts {
				if strings.EqualFold(strings.TrimSpace(extensionText(post.Message.Extensions["task.status"])), spec.status) {
					items = append(items, post)
					if len(items) >= 3 {
						break
					}
				}
			}
			count := 0
			for _, post := range posts {
				if strings.EqualFold(strings.TrimSpace(extensionText(post.Message.Extensions["task.status"])), spec.status) {
					count++
				}
			}
			out = append(out, ModuleBoardColumn{
				Title:       spec.title,
				Description: spec.desc,
				Path:        spec.path,
				Count:       count,
				Posts:       items,
			})
		}
		return out
	case "skill":
		byCategory := map[string][]Post{}
		for _, post := range posts {
			category := strings.TrimSpace(extensionText(post.Message.Extensions["skill.category"]))
			if category == "" {
				category = "uncategorized"
			}
			byCategory[category] = append(byCategory[category], post)
		}
		type pair struct {
			name  string
			posts []Post
		}
		pairs := make([]pair, 0, len(byCategory))
		for name, items := range byCategory {
			pairs = append(pairs, pair{name: name, posts: items})
		}
		sort.Slice(pairs, func(i, j int) bool {
			if len(pairs[i].posts) == len(pairs[j].posts) {
				return pairs[i].name < pairs[j].name
			}
			return len(pairs[i].posts) > len(pairs[j].posts)
		})
		if len(pairs) > 4 {
			pairs = pairs[:4]
		}
		out := make([]ModuleBoardColumn, 0, len(pairs))
		for _, item := range pairs {
			posts := item.posts
			if len(posts) > 3 {
				posts = posts[:3]
			}
			out = append(out, ModuleBoardColumn{
				Title:       item.name,
				Description: "Reusable skill assets in this category.",
				Path:        "/skills?meta_key=skill.category&meta_value=" + item.name,
				Count:       len(item.posts),
				Posts:       posts,
			})
		}
		return out
	default:
		return nil
	}
}

func ModulePrimaryFieldsForPost(post Post) []MetadataField {
	fields := make([]MetadataField, 0, 4)
	add := func(label, value string) {
		value = strings.TrimSpace(value)
		if value == "" {
			return
		}
		fields = append(fields, MetadataField{Label: label, Value: value})
	}

	ext := post.Message.Extensions
	switch strings.ToLower(strings.TrimSpace(PostCoordType(post))) {
	case "idea":
		add("Problem", extensionText(ext["coord.problem"]))
		add("Goal", extensionText(ext["coord.goal"]))
		add("Expected output", extensionText(ext["coord.expected_output"]))
	case "task":
		add("Task ID", extensionText(ext["task.id"]))
		add("Status", extensionText(ext["task.status"]))
		add("Required assets", extensionText(ext["task.required_assets"]))
		add("Expected result", extensionText(ext["task.expected_result_type"]))
	case "skill":
		add("Skill ID", extensionText(ext["skill.id"]))
		add("Category", extensionText(ext["skill.category"]))
		add("Inputs", extensionText(ext["skill.inputs"]))
		add("Outputs", extensionText(ext["skill.outputs"]))
		add("Depends on", extensionText(ext["skill.depends_on"]))
	case "markdown":
		add("Document kind", extensionText(ext["md.kind"]))
		add("Collection", extensionText(ext["md.collection"]))
		add("Repository", extensionText(ext["md.source_repo"]))
	case "code":
		add("Code kind", extensionText(ext["code.kind"]))
		add("Language", extensionText(ext["code.language"]))
		add("Repository", extensionText(ext["code.repo"]))
		add("Entry point", extensionText(ext["code.entry"]))
	case "agent":
		add("Agent ID", extensionText(ext["agent.id"]))
		add("Availability", extensionText(ext["agent.availability"]))
		add("Skills", extensionText(ext["agent.skills"]))
		add("Models", extensionText(ext["agent.models"]))
		add("Tools", extensionText(ext["agent.tools"]))
	}
	return fields
}

func ModuleDetailQuickStats(post Post) []SummaryStat {
	ext := post.Message.Extensions
	switch strings.ToLower(strings.TrimSpace(PostCoordType(post))) {
	case "task":
		return []SummaryStat{
			{Label: "Status", Value: fallbackValue(extensionText(ext["task.status"]), "unset")},
			{Label: "Result", Value: fallbackValue(extensionText(ext["task.expected_result_type"]), "unset")},
			{Label: "Required assets", Value: countField(ext["task.required_assets"])},
			{Label: "Topics", Value: itoa(len(post.Topics))},
		}
	case "skill":
		return []SummaryStat{
			{Label: "Category", Value: fallbackValue(extensionText(ext["skill.category"]), "unset")},
			{Label: "Inputs", Value: countField(ext["skill.inputs"])},
			{Label: "Outputs", Value: countField(ext["skill.outputs"])},
			{Label: "Depends on", Value: countField(ext["skill.depends_on"])},
		}
	default:
		return nil
	}
}

func ModuleDetailFocusCards(post Post) []ModuleFocusCard {
	ext := post.Message.Extensions
	switch strings.ToLower(strings.TrimSpace(PostCoordType(post))) {
	case "task":
		status := strings.TrimSpace(extensionText(ext["task.status"]))
		if status == "" {
			status = "unset"
		}
		result := strings.TrimSpace(extensionText(ext["task.expected_result_type"]))
		if result == "" {
			result = "not declared yet"
		}
		required := strings.TrimSpace(extensionText(ext["task.required_assets"]))
		if required == "" {
			required = "no explicit linked assets yet"
		}
		return []ModuleFocusCard{
			{
				Title:       "Execution snapshot",
				Description: "Task detail keeps the current delivery state explicit before anyone opens the full trail.",
				Items: []string{
					"Current status: " + status,
					"Expected result: " + result,
					"Required assets: " + required,
				},
			},
			{
				Title:       "Reply protocol",
				Description: "Use reply updates to keep the task timeline consistent for humans and AI agents.",
				Items: []string{
					"Progress update when work moves forward.",
					"Blocker update when delivery is blocked.",
					"Review or handoff note when ownership changes.",
				},
			},
		}
	case "skill":
		category := strings.TrimSpace(extensionText(ext["skill.category"]))
		if category == "" {
			category = "uncategorized"
		}
		inputs := strings.TrimSpace(extensionText(ext["skill.inputs"]))
		if inputs == "" {
			inputs = "not declared yet"
		}
		outputs := strings.TrimSpace(extensionText(ext["skill.outputs"]))
		if outputs == "" {
			outputs = "not declared yet"
		}
		deps := strings.TrimSpace(extensionText(ext["skill.depends_on"]))
		if deps == "" {
			deps = "no explicit dependencies"
		}
		return []ModuleFocusCard{
			{
				Title:       "Capability contract",
				Description: "Skill detail should tell another operator exactly when to reuse this capability.",
				Items: []string{
					"Category: " + category,
					"Inputs: " + inputs,
					"Outputs: " + outputs,
				},
			},
			{
				Title:       "Execution support",
				Description: "Dependencies and supporting assets stay visible without turning the skill into a task clone.",
				Items: []string{
					"Depends on: " + deps,
					"Keep linked code assets separate from the reusable skill note.",
					"Keep linked knowledge notes separate from the concise operator contract.",
				},
			},
		}
	default:
		return nil
	}
}

func ModuleReplyLanesForPost(post Post, replies []Reply) []ModuleReplyLane {
	if strings.ToLower(strings.TrimSpace(PostCoordType(post))) != "task" {
		return nil
	}
	specs := []struct {
		key         string
		title       string
		description string
	}{
		{key: "progress", title: "Progress", description: "Execution updates that show the task moving forward."},
		{key: "blocker", title: "Blockers", description: "Replies that explain why delivery is blocked or waiting."},
		{key: "review", title: "Review", description: "Review requests and review-facing coordination replies."},
		{key: "handoff", title: "Handoff", description: "Replies that move ownership or package work for another operator."},
		{key: "other", title: "Other updates", description: "Replies that still belong to the task trail but do not match the main execution roles."},
	}
	grouped := map[string][]Reply{
		"progress": {},
		"blocker":  {},
		"review":   {},
		"handoff":  {},
		"other":    {},
	}
	for _, reply := range replies {
		key := moduleReplyLaneKey(reply)
		grouped[key] = append(grouped[key], reply)
	}
	out := make([]ModuleReplyLane, 0, len(specs))
	for _, spec := range specs {
		items := grouped[spec.key]
		if len(items) == 0 {
			continue
		}
		out = append(out, ModuleReplyLane{
			Key:         spec.key,
			Title:       spec.title,
			Description: spec.description,
			Count:       len(items),
			Replies:     items,
		})
	}
	return out
}

func ModuleReplyHighlightsForPost(post Post, replies []Reply) []ModuleReplyHighlight {
	if strings.ToLower(strings.TrimSpace(PostCoordType(post))) != "task" {
		return nil
	}
	lanes := ModuleReplyLanesForPost(post, replies)
	targets := []string{"progress", "blocker", "review"}
	out := make([]ModuleReplyHighlight, 0, len(targets))
	for _, key := range targets {
		for _, lane := range lanes {
			if lane.Key != key || len(lane.Replies) == 0 {
				continue
			}
			reply := lane.Replies[len(lane.Replies)-1]
			replyCopy := reply
			out = append(out, ModuleReplyHighlight{
				Key:         lane.Key,
				Title:       "Latest " + strings.ToLower(lane.Title),
				Description: lane.Description,
				Reply:       &replyCopy,
			})
			break
		}
	}
	return out
}

func moduleReplyLaneKey(reply Reply) string {
	role := strings.ToLower(strings.TrimSpace(extensionText(reply.Message.Extensions["thread.role"])))
	replyType := strings.ToLower(strings.TrimSpace(extensionText(reply.Message.Extensions["reply_type"])))
	value := role
	if value == "" {
		value = replyType
	}
	switch value {
	case "progress", "progress-update", "progress_update":
		return "progress"
	case "blocker", "blocked", "blocker-update", "blocker_update":
		return "blocker"
	case "review", "review-request", "review_request":
		return "review"
	case "handoff", "handoff-note", "handoff_note":
		return "handoff"
	default:
		return "other"
	}
}

func countField(value any) string {
	items := extensionStrings(value)
	if len(items) > 0 {
		return itoa(len(items))
	}
	if strings.TrimSpace(extensionText(value)) != "" {
		return "1"
	}
	return "0"
}

func fallbackValue(value, fallback string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return fallback
	}
	return value
}

func itoa(v int) string {
	if v < 0 {
		v = 0
	}
	return strconv.Itoa(v)
}
