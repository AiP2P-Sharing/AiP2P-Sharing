package newsplugin

import (
	"encoding/json"
	"fmt"
	"strings"
)

type CoordPublishGuide struct {
	CoordType             string
	Title                 string
	Description           string
	ExecutionMode         string
	TargetSummary         string
	PublisherRole         string
	IdentityFile          string
	StorePath             string
	Channel               string
	Kind                  string
	Prerequisites         []string
	ContextNotes          []string
	StarterExtensionsJSON string
	StarterCommand        string
	StarterFields         []MetadataSchemaField
	Templates             []CoordPublishTemplate
	Workflow              []CoordPublishWorkflowStep
}

type CoordPublishTemplate struct {
	Category              string
	Label                 string
	Description           string
	ExecutionMode         string
	TargetSummary         string
	ReferenceTarget       string
	PublisherRole         string
	IdentityFile          string
	StorePath             string
	Channel               string
	Kind                  string
	Prerequisites         []string
	StarterExtensionsJSON string
	StarterCommand        string
}

type CoordPublishWorkflowStep struct {
	Title          string
	Description    string
	TemplateLabels []string
	Templates      []CoordPublishTemplate
}

type CoordPublishContext struct {
	Topics     []string
	SourceName string
	Extensions map[string]any
	InfoHash   string
	Magnet     string
}

func CoordPublishGuideFor(projectID, coordType string, ctx CoordPublishContext) *CoordPublishGuide {
	spec, ok := CoordMetadataSpec(coordType)
	if !ok {
		return nil
	}
	channel := publishChannelForCoordType(projectID, spec.CoordType)
	fields := append([]MetadataSchemaField{{Key: "coord.type", Label: "Coord type"}}, spec.Fields...)
	contextTopics := normalizeGuideTopics(ctx.Topics)
	contextNotes := publishGuideContextNotes(contextTopics, ctx.SourceName)
	templates := publishTemplatesFor(projectID, spec.CoordType, contextTopics, ctx)
	workflow := publishWorkflowFor(spec.CoordType)
	return &CoordPublishGuide{
		CoordType:             spec.CoordType,
		Title:                 publishGuideTitle(spec.CoordType),
		Description:           publishGuideDescription(spec.CoordType),
		ExecutionMode:         "Create asset",
		TargetSummary:         publishGuideTargetSummary(spec.CoordType),
		PublisherRole:         publisherRoleForGuide(spec.CoordType),
		IdentityFile:          identityFileForGuide(spec.CoordType),
		StorePath:             publishStorePath(),
		Channel:               channel,
		Kind:                  "post",
		Prerequisites:         publishGuidePrerequisites(spec.CoordType),
		ContextNotes:          contextNotes,
		StarterExtensionsJSON: starterExtensionsJSON(projectID, spec.CoordType, spec.Fields, contextTopics),
		StarterCommand:        starterPublishCommand(projectID, spec.CoordType, channel, spec.Fields, contextTopics),
		StarterFields:         fields,
		Templates:             templates,
		Workflow:              bindWorkflowTemplates(workflow, templates),
	}
}

func publishGuidePrerequisites(coordType string) []string {
	items := []string{
		"Ensure the target identity file exists and matches the role that will publish this asset.",
		"Ensure the local store path is initialized on this node before running the command.",
	}
	switch coordType {
	case "task", "idea", "skill", "markdown", "code", "agent":
		items = append(items, "Adjust placeholder title, body, and metadata values before publishing.")
	}
	return items
}

func publishStorePath() string {
	return "$HOME/.aip2p-sharing/aip2p/.aip2p"
}

func publisherRoleForGuide(coordType string) string {
	switch coordType {
	case "task":
		return "Task publisher"
	case "idea":
		return "Idea publisher"
	case "skill":
		return "Skill publisher"
	case "markdown":
		return "Knowledge publisher"
	case "code":
		return "Code publisher"
	case "agent":
		return "Agent profile owner"
	default:
		return "Workspace publisher"
	}
}

func identityFileForGuide(coordType string) string {
	return fmt.Sprintf("$HOME/.aip2p-sharing/identities/%s-01.json", coordType)
}

func identityFileForThread() string {
	return "$HOME/.aip2p-sharing/identities/task-operator-01.json"
}

func identityFileForReaction(coordType string) string {
	return fmt.Sprintf("$HOME/.aip2p-sharing/identities/%s-operator-01.json", coordType)
}

func publishGuideTargetSummary(coordType string) string {
	switch coordType {
	case "idea":
		return "Creates a new idea asset in the backlog channel."
	case "task":
		return "Creates a new task asset in the execution channel."
	case "skill":
		return "Creates a new reusable skill asset in the capability channel."
	case "markdown":
		return "Creates a new knowledge asset in the reading and note channel."
	case "code":
		return "Creates a new code asset in the implementation channel."
	case "agent":
		return "Creates a new operator profile asset in the agent channel."
	default:
		return "Creates a new typed coordination asset."
	}
}

func publishChannelForCoordType(projectID, coordType string) string {
	projectID = strings.TrimSpace(projectID)
	if projectID == "" {
		projectID = "aip2p.sharing"
	}
	if spec, ok := CoordCollectionByType(coordType); ok {
		path := strings.TrimPrefix(spec.Path, "/")
		if path != "" {
			return projectID + "/" + path
		}
	}
	return projectID + "/general"
}

func publishGuideTitle(coordType string) string {
	switch coordType {
	case "idea":
		return "Publish a new idea asset"
	case "task":
		return "Publish a new task asset"
	case "skill":
		return "Publish a new skill asset"
	case "markdown":
		return "Publish a new knowledge asset"
	case "code":
		return "Publish a new code asset"
	case "agent":
		return "Publish a new agent profile"
	default:
		return "Publish a new coordination asset"
	}
}

func publishGuideDescription(coordType string) string {
	switch coordType {
	case "idea":
		return "Start with problem framing, goal, and expected output so the idea can later become an active task."
	case "task":
		return "Include stable delivery fields so operators can track status, dependencies, and expected results without reopening discovery."
	case "skill":
		return "Declare capability inputs, outputs, and dependencies so the skill can be reused across later tasks."
	case "markdown":
		return "Record durable document metadata so operators know where this note belongs before using it in delivery."
	case "code":
		return "Publish implementation metadata that points back to the repository, entry point, and language used in delivery."
	case "agent":
		return "Publish a stable operator profile with declared availability, models, tools, and reusable skill areas."
	default:
		return "Publish a typed coordination asset with stable metadata."
	}
}

func starterExtensionsJSON(projectID, coordType string, fields []MetadataSchemaField, topics []string) string {
	payload := map[string]any{
		"project":    strings.TrimSpace(projectID),
		"coord.type": coordType,
	}
	for _, field := range fields {
		payload[field.Key] = starterValueForField(field.Key)
	}
	if len(topics) > 0 {
		payload["topics"] = topics
	}
	data, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return fmt.Sprintf("{\"project\":\"%s\",\"coord.type\":\"%s\"}", projectID, coordType)
	}
	return string(data)
}

func starterExtensionsJSONCompact(projectID, coordType string, fields []MetadataSchemaField, topics []string) string {
	payload := map[string]any{
		"project":    strings.TrimSpace(projectID),
		"coord.type": coordType,
	}
	for _, field := range fields {
		payload[field.Key] = starterValueForField(field.Key)
	}
	if len(topics) > 0 {
		payload["topics"] = topics
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Sprintf("{\"project\":\"%s\",\"coord.type\":\"%s\"}", projectID, coordType)
	}
	return string(data)
}

func starterPublishCommand(projectID, coordType, channel string, fields []MetadataSchemaField, topics []string) string {
	identityFile := identityFileForGuide(coordType)
	extensions := starterExtensionsJSONCompact(projectID, coordType, fields, topics)
	return fmt.Sprintf(
		"go -C ./AiP2P run ./cmd/aip2p publish --store \"$HOME/.aip2p-sharing/aip2p/.aip2p\" --identity-file \"%s\" --kind post --channel \"%s\" --title \"<title>\" --body \"<body>\" --extensions-json '%s'",
		identityFile,
		channel,
		extensions,
	)
}

func starterValueForField(key string) any {
	switch key {
	case "task.required_assets", "skill.inputs", "skill.outputs", "skill.depends_on", "agent.models", "agent.tools", "agent.skills":
		return []string{"<value>"}
	default:
		return "<value>"
	}
}

func normalizeGuideTopics(topics []string) []string {
	out := make([]string, 0, len(topics))
	seen := make(map[string]struct{}, len(topics))
	for _, topic := range topics {
		topic = strings.TrimSpace(topic)
		if topic == "" || strings.EqualFold(topic, reservedTopicAll) {
			continue
		}
		folded := strings.ToLower(topic)
		if _, ok := seen[folded]; ok {
			continue
		}
		seen[folded] = struct{}{}
		out = append(out, topic)
	}
	return out
}

func publishGuideContextNotes(topics []string, sourceName string) []string {
	notes := make([]string, 0, 2)
	if len(topics) > 0 {
		notes = append(notes, "Starter topics inherited from the current workspace context: "+strings.Join(topics, ", "))
	}
	sourceName = strings.TrimSpace(sourceName)
	if sourceName != "" {
		notes = append(notes, "Current source scope is "+sourceName+". Source grouping comes from the signing identity or source page, not from extensions-json.")
	}
	return notes
}

func publishTemplatesFor(projectID, coordType string, topics []string, ctx CoordPublishContext) []CoordPublishTemplate {
	switch coordType {
	case "task":
		return []CoordPublishTemplate{
			buildThreadTemplate(projectID, "Progress update", "Publish a delivery update as a reply so the task timeline stays attached to the current execution item.", "update", "progress-note", topics, ctx),
			buildThreadTemplate(projectID, "Blocker update", "Publish a blocker note as a reply so the next operator can see what is currently preventing delivery.", "blocker", "blocker-note", topics, ctx),
			buildThreadTemplate(projectID, "Request review", "Publish a review request as a reply so another operator can review the current delivery state in the same task timeline.", "review-request", "review-note", topics, ctx),
			buildThreadTemplate(projectID, "Handoff note", "Publish a handoff reply so the next operator can resume this task with the current state, pending work, and delivery context attached.", "handoff", "handoff-note", topics, ctx),
			buildAssetTemplate(projectID, "skill", "Publish supporting skill", "Publish a reusable capability asset that this task depends on or discovered during execution.", "<skill title>", "<how this skill supports the task>", topics, ctx, map[string]any{}),
			buildAssetTemplate(projectID, "markdown", "Publish delivery note", "Publish a durable markdown note that captures delivery context, findings, or operator handoff detail for this task.", "<delivery note title>", "<delivery note for the current task>", topics, ctx, map[string]any{
				"md.kind":       "delivery-note",
				"md.collection": "task-notes",
			}),
			buildAssetTemplate(projectID, "code", "Publish implementation asset", "Publish a code asset that captures the implementation surface needed to complete or support this task.", "<implementation asset title>", "<implementation summary for the current task>", topics, ctx, map[string]any{
				"code.kind": "implementation",
			}),
			buildReactionTemplate(projectID, "task", "Approve delivery", "Publish a positive delivery signal when this task is ready to proceed or merge.", "<approval signal>", "<why this task is approved>", "vote", 1, topics, ctx),
			buildReactionTemplate(projectID, "task", "Needs attention", "Publish a negative signal when this task needs more work, review, or clarification before continuing.", "<attention signal>", "<what still needs attention>", "vote", -1, topics, ctx),
		}
	case "idea":
		return []CoordPublishTemplate{
			buildTaskPromotionTemplate(projectID, topics),
			buildAssetTemplate(projectID, "markdown", "Publish framing note", "Publish a durable knowledge note that captures background, tradeoffs, and constraints behind the current idea.", "<framing note title>", "<context and constraints for this idea>", topics, ctx, map[string]any{
				"md.kind":       "framing-note",
				"md.collection": "idea-notes",
			}),
			buildAssetTemplate(projectID, "code", "Publish prototype asset", "Publish a prototype or experiment asset that helps validate the current idea before it becomes a committed delivery task.", "<prototype asset title>", "<prototype summary for this idea>", topics, ctx, map[string]any{
				"code.kind": "prototype",
			}),
			buildReactionTemplate(projectID, "idea", "Signal support", "Publish a lightweight support signal when the current idea should stay visible in the backlog.", "<support signal>", "<why this idea deserves support>", "vote", 1, topics, ctx),
		}
	case "markdown":
		return []CoordPublishTemplate{
			buildAssetTemplate(projectID, "code", "Publish implementation follow-up", "Publish a code asset that implements or validates the knowledge captured in this note.", "<implementation follow-up title>", "<implementation follow-up for this note>", topics, ctx, map[string]any{
				"code.kind": "implementation",
			}),
			buildTaskPromotionTemplate(projectID, topics),
		}
	case "code":
		return []CoordPublishTemplate{
			buildAssetTemplate(projectID, "markdown", "Publish implementation note", "Publish a durable knowledge note that explains how this code asset should be used, reviewed, or extended.", "<implementation note title>", "<implementation note for this code asset>", topics, ctx, map[string]any{
				"md.kind":       "implementation-note",
				"md.collection": "code-notes",
			}),
			buildAssetTemplate(projectID, "skill", "Publish reusable skill", "Publish a reusable skill asset when this code surface should become an explicit capability in the workspace.", "<reusable skill title>", "<reusable capability extracted from this code asset>", topics, ctx, map[string]any{}),
		}
	case "agent":
		return []CoordPublishTemplate{
			buildAgentAvailabilityTemplate(projectID, topics, ctx),
			buildAgentSkillTemplate(projectID, topics, ctx),
			buildAssetTemplate(projectID, "markdown", "Publish operator note", "Publish a durable operator note that documents how this agent should be used, coordinated, or handed off inside the workspace.", "<operator note title>", "<operator note for this agent>", topics, ctx, map[string]any{
				"md.kind":       "operator-note",
				"md.collection": "agent-notes",
			}),
		}
	default:
		return nil
	}
}

func publishWorkflowFor(coordType string) []CoordPublishWorkflowStep {
	switch coordType {
	case "task":
		return []CoordPublishWorkflowStep{
			{
				Title:          "Keep the task timeline current",
				Description:    "Use reply-based updates when execution changes hands, needs review, or hits blockers so the current task stays self-explanatory.",
				TemplateLabels: []string{"Progress update", "Blocker update", "Request review", "Handoff note"},
			},
			{
				Title:          "Capture linked delivery assets",
				Description:    "Publish supporting skills, notes, and implementation assets when delivery produces reusable artifacts outside the timeline.",
				TemplateLabels: []string{"Publish supporting skill", "Publish delivery note", "Publish implementation asset"},
			},
			{
				Title:          "Signal delivery readiness",
				Description:    "Use lightweight reactions after timeline and asset updates are in place to show whether the task is ready or still needs attention.",
				TemplateLabels: []string{"Approve delivery", "Needs attention"},
			},
		}
	case "idea":
		return []CoordPublishWorkflowStep{
			{
				Title:          "Frame the idea",
				Description:    "Capture rationale and backlog support before the idea turns into committed delivery work.",
				TemplateLabels: []string{"Publish framing note", "Signal support"},
			},
			{
				Title:          "Turn exploration into execution",
				Description:    "Publish a prototype or promote the idea into a task when the direction is ready for implementation.",
				TemplateLabels: []string{"Publish prototype asset", "Promote to task"},
			},
		}
	case "agent":
		return []CoordPublishWorkflowStep{
			{
				Title:          "Keep the operator profile current",
				Description:    "Refresh the profile first so availability and declared capabilities stay accurate for the rest of the workspace.",
				TemplateLabels: []string{"Availability update"},
			},
			{
				Title:          "Extract reusable operating knowledge",
				Description:    "Capture skills and operator notes when the agent profile should produce durable assets other operators can reuse.",
				TemplateLabels: []string{"Publish capability skill", "Publish operator note"},
			},
		}
	case "markdown":
		return []CoordPublishWorkflowStep{
			{
				Title:          "Move from note to delivery",
				Description:    "Turn durable knowledge into executable work by publishing follow-up code or promoting the note into a task.",
				TemplateLabels: []string{"Publish implementation follow-up", "Promote to task"},
			},
		}
	case "code":
		return []CoordPublishWorkflowStep{
			{
				Title:          "Document and extract reuse",
				Description:    "Pair the implementation surface with a note or reusable skill when the code should become a durable workspace asset.",
				TemplateLabels: []string{"Publish implementation note", "Publish reusable skill"},
			},
		}
	default:
		return nil
	}
}

func bindWorkflowTemplates(steps []CoordPublishWorkflowStep, templates []CoordPublishTemplate) []CoordPublishWorkflowStep {
	if len(steps) == 0 || len(templates) == 0 {
		return steps
	}
	byLabel := make(map[string]CoordPublishTemplate, len(templates))
	for _, template := range templates {
		byLabel[template.Label] = template
	}
	out := make([]CoordPublishWorkflowStep, 0, len(steps))
	for _, step := range steps {
		bound := CoordPublishWorkflowStep{
			Title:          step.Title,
			Description:    step.Description,
			TemplateLabels: append([]string(nil), step.TemplateLabels...),
			Templates:      make([]CoordPublishTemplate, 0, len(step.TemplateLabels)),
		}
		for _, label := range step.TemplateLabels {
			template, ok := byLabel[label]
			if !ok {
				continue
			}
			bound.Templates = append(bound.Templates, template)
		}
		out = append(out, bound)
	}
	return out
}

func buildThreadTemplate(projectID, label, description, role, resultType string, topics []string, ctx CoordPublishContext) CoordPublishTemplate {
	channel := publishChannelForCoordType(projectID, "task")
	identityFile := identityFileForThread()
	payload := map[string]any{
		"project":            strings.TrimSpace(projectID),
		"coord.type":         "thread",
		"thread.role":        role,
		"thread.result_type": resultType,
	}
	if taskID := extensionText(ctx.Extensions["task.id"]); taskID != "" {
		payload["thread.task_id"] = taskID
	} else {
		payload["thread.task_id"] = "<task-id>"
	}
	if len(topics) > 0 {
		payload["topics"] = topics
	}
	infoHash := placeholderOrValue(strings.TrimSpace(ctx.InfoHash), "<parent-infohash>")
	magnet := placeholderOrValue(strings.TrimSpace(ctx.Magnet), "<parent-magnet>")
	command := fmt.Sprintf(
		"go -C ./AiP2P run ./cmd/aip2p publish --store \"$HOME/.aip2p-sharing/aip2p/.aip2p\" --identity-file \"$HOME/.aip2p-sharing/identities/task-operator-01.json\" --kind reply --channel \"%s\" --title \"<title>\" --body \"<body>\" --reply-infohash \"%s\" --reply-magnet \"%s\" --extensions-json '%s'",
		channel,
		infoHash,
		magnet,
		mustCompactJSON(payload),
	)
	return CoordPublishTemplate{
		Category:              "Timeline",
		Label:                 label,
		Description:           description,
		ExecutionMode:         "Reply in timeline",
		TargetSummary:         "Replies to the current task asset and keeps the execution timeline attached to the same work item.",
		ReferenceTarget:       "Reply target: " + infoHash,
		PublisherRole:         "Task operator",
		IdentityFile:          identityFile,
		StorePath:             publishStorePath(),
		Channel:               channel,
		Kind:                  "reply",
		Prerequisites: []string{
			"Keep the current task infohash and magnet available so the reply stays attached to the right timeline.",
			"Use an operator identity that is allowed to post task updates on this node.",
		},
		StarterExtensionsJSON: mustPrettyJSON(payload),
		StarterCommand:        command,
	}
}

func buildTaskPromotionTemplate(projectID string, topics []string) CoordPublishTemplate {
	channel := publishChannelForCoordType(projectID, "task")
	identityFile := identityFileForGuide("task")
	payload := map[string]any{
		"project":                   strings.TrimSpace(projectID),
		"coord.type":                "task",
		"task.id":                   "<task-id>",
		"task.status":               "queued",
		"task.priority":             "<priority>",
		"task.parent":               "<upstream-work-item>",
		"task.required_assets":      []string{"<value>"},
		"task.expected_result_type": "<value>",
	}
	if len(topics) > 0 {
		payload["topics"] = topics
	}
	command := fmt.Sprintf(
		"go -C ./AiP2P run ./cmd/aip2p publish --store \"$HOME/.aip2p-sharing/aip2p/.aip2p\" --identity-file \"$HOME/.aip2p-sharing/identities/task-01.json\" --kind post --channel \"%s\" --title \"<task title>\" --body \"<delivery scope>\" --extensions-json '%s'",
		channel,
		mustCompactJSON(payload),
	)
	return CoordPublishTemplate{
		Category:              "Task follow-up",
		Label:                 "Promote to task",
		Description:           "Turn the current idea into an executable task with an initial status and stable task identifier.",
		ExecutionMode:         "Create asset",
		TargetSummary:         "Creates a new task asset as the execution follow-up to the current item.",
		ReferenceTarget:       "No current asset target. This creates a fresh task asset.",
		PublisherRole:         "Task publisher",
		IdentityFile:          identityFile,
		StorePath:             publishStorePath(),
		Channel:               channel,
		Kind:                  "post",
		Prerequisites: []string{
			"Choose a stable task id before publishing so later replies and follow-up assets can refer back to the same work item.",
			"Replace placeholder delivery metadata with the actual initial status, priority, and expected result type.",
		},
		StarterExtensionsJSON: mustPrettyJSON(payload),
		StarterCommand:        command,
	}
}

func buildAgentAvailabilityTemplate(projectID string, topics []string, ctx CoordPublishContext) CoordPublishTemplate {
	channel := publishChannelForCoordType(projectID, "agent")
	identityFile := identityFileForGuide("agent")
	payload := map[string]any{
		"project":            strings.TrimSpace(projectID),
		"coord.type":         "agent",
		"agent.id":           placeholderOrValue(extensionText(ctx.Extensions["agent.id"]), "<agent-id>"),
		"agent.availability": "<availability>",
		"agent.models":       []string{"<model>"},
		"agent.tools":        []string{"<tool>"},
		"agent.skills":       []string{"<skill>"},
	}
	if len(topics) > 0 {
		payload["topics"] = topics
	}
	command := fmt.Sprintf(
		"go -C ./AiP2P run ./cmd/aip2p publish --store \"$HOME/.aip2p-sharing/aip2p/.aip2p\" --identity-file \"$HOME/.aip2p-sharing/identities/agent-01.json\" --kind post --channel \"%s\" --title \"<profile update title>\" --body \"<availability update>\" --extensions-json '%s'",
		channel,
		mustCompactJSON(payload),
	)
	return CoordPublishTemplate{
		Category:              "Profile update",
		Label:                 "Availability update",
		Description:           "Publish a fresh agent profile update when availability, tools, models, or declared skill coverage changes.",
		ExecutionMode:         "Create asset",
		TargetSummary:         "Creates a new agent profile update asset in the agent channel.",
		ReferenceTarget:       "No current asset target. This creates a fresh agent profile asset.",
		PublisherRole:         "Agent profile owner",
		IdentityFile:          identityFile,
		StorePath:             publishStorePath(),
		Channel:               channel,
		Kind:                  "post",
		Prerequisites: []string{
			"Use the agent profile identity that should own this availability update.",
			"Replace placeholder availability, model, tool, and skill values before publishing.",
		},
		StarterExtensionsJSON: mustPrettyJSON(payload),
		StarterCommand:        command,
	}
}

func buildAssetTemplate(projectID, coordType, label, description, titlePlaceholder, bodyPlaceholder string, topics []string, ctx CoordPublishContext, overrides map[string]any) CoordPublishTemplate {
	spec, ok := CoordMetadataSpec(coordType)
	if !ok {
		return CoordPublishTemplate{}
	}
	channel := publishChannelForCoordType(projectID, coordType)
	identityFile := identityFileForGuide(coordType)
	payload := map[string]any{
		"project":    strings.TrimSpace(projectID),
		"coord.type": coordType,
	}
	for _, field := range spec.Fields {
		payload[field.Key] = starterValueForField(field.Key)
	}
	for key, value := range overrides {
		payload[key] = value
	}
	if len(topics) > 0 {
		payload["topics"] = topics
	}
	taskID := strings.TrimSpace(extensionText(ctx.Extensions["task.id"]))
	if taskID != "" {
		bodyPlaceholder = bodyPlaceholder + " (task: " + taskID + ")"
	}
	command := fmt.Sprintf(
		"go -C ./AiP2P run ./cmd/aip2p publish --store \"$HOME/.aip2p-sharing/aip2p/.aip2p\" --identity-file \"%s\" --kind post --channel \"%s\" --title \"%s\" --body \"%s\" --extensions-json '%s'",
		identityFile,
		channel,
		titlePlaceholder,
		bodyPlaceholder,
		mustCompactJSON(payload),
	)
	return CoordPublishTemplate{
		Category:              "Linked asset",
		Label:                 label,
		Description:           description,
		ExecutionMode:         "Create linked asset",
		TargetSummary:         "Creates a new linked " + coordType + " asset that can be discovered from the current workspace context.",
		ReferenceTarget:       "No direct reply target. This publishes a new " + coordType + " asset into its typed channel.",
		PublisherRole:         publisherRoleForGuide(coordType),
		IdentityFile:          identityFile,
		StorePath:             publishStorePath(),
		Channel:               channel,
		Kind:                  "post",
		Prerequisites: []string{
			"Replace the placeholder metadata and body text with the actual reusable asset context before publishing.",
		},
		StarterExtensionsJSON: mustPrettyJSON(payload),
		StarterCommand:        command,
	}
}

func buildReactionTemplate(projectID, parentCoordType, label, description, titlePlaceholder, bodyPlaceholder, reactionType string, value any, topics []string, ctx CoordPublishContext) CoordPublishTemplate {
	channel := publishChannelForCoordType(projectID, parentCoordType)
	identityFile := identityFileForReaction(parentCoordType)
	payload := map[string]any{
		"project":       strings.TrimSpace(projectID),
		"reaction_type": reactionType,
		"value":         value,
		"explanation":   bodyPlaceholder,
		"subject": map[string]any{
			"infohash": placeholderOrValue(strings.TrimSpace(ctx.InfoHash), "<subject-infohash>"),
		},
	}
	if magnet := strings.TrimSpace(ctx.Magnet); magnet != "" {
		subject, _ := payload["subject"].(map[string]any)
		subject["magnet"] = magnet
	}
	if len(topics) > 0 {
		payload["topics"] = topics
	}
	command := fmt.Sprintf(
		"go -C ./AiP2P run ./cmd/aip2p publish --store \"$HOME/.aip2p-sharing/aip2p/.aip2p\" --identity-file \"$HOME/.aip2p-sharing/identities/%s-operator-01.json\" --kind reaction --channel \"%s\" --title \"%s\" --body \"%s\" --extensions-json '%s'",
		parentCoordType,
		channel,
		titlePlaceholder,
		bodyPlaceholder,
		mustCompactJSON(payload),
	)
	return CoordPublishTemplate{
		Category:              "Signal",
		Label:                 label,
		Description:           description,
		ExecutionMode:         "React to asset",
		TargetSummary:         "Publishes a reaction against the current " + parentCoordType + " asset without creating a separate post in the timeline.",
		ReferenceTarget:       "Reaction subject: " + placeholderOrValue(strings.TrimSpace(ctx.InfoHash), "<subject-infohash>"),
		PublisherRole:         strings.Title(parentCoordType) + " operator",
		IdentityFile:          identityFile,
		StorePath:             publishStorePath(),
		Channel:               channel,
		Kind:                  "reaction",
		Prerequisites: []string{
			"Keep the current asset infohash available so the reaction lands on the intended subject.",
			"Use an operator identity that is meant to emit delivery or backlog signals for this asset.",
		},
		StarterExtensionsJSON: mustPrettyJSON(payload),
		StarterCommand:        command,
	}
}

func buildAgentSkillTemplate(projectID string, topics []string, ctx CoordPublishContext) CoordPublishTemplate {
	agentID := placeholderOrValue(extensionText(ctx.Extensions["agent.id"]), "<agent-id>")
	return buildAssetTemplate(projectID, "skill", "Publish capability skill", "Publish a reusable skill asset that formalizes one capability currently declared on this agent profile.", "<capability skill title>", "<capability this agent can provide> (agent: "+agentID+")", topics, ctx, map[string]any{
		"skill.id": "<skill-id>",
	})
}

func placeholderOrValue(value, placeholder string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return placeholder
	}
	return value
}

func mustCompactJSON(payload map[string]any) string {
	data, err := json.Marshal(payload)
	if err != nil {
		return "{}"
	}
	return string(data)
}

func mustPrettyJSON(payload map[string]any) string {
	data, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return "{}"
	}
	return string(data)
}
