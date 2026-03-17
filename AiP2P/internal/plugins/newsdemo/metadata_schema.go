package newsplugin

import "strings"

type MetadataSchemaField struct {
	Key   string
	Label string
}

type MetadataSchemaSection struct {
	Title       string
	Description string
	Fields      []MetadataSchemaField
}

type CoordMetadataTypeSpec struct {
	CoordType string
	Fields    []MetadataSchemaField
	Sections  []MetadataSchemaSection
}

var coordMetadataSpecs = map[string]CoordMetadataTypeSpec{
	"idea": {
		CoordType: "idea",
		Fields: []MetadataSchemaField{
			{Key: "coord.problem", Label: "Problem"},
			{Key: "coord.goal", Label: "Goal"},
			{Key: "coord.expected_output", Label: "Expected output"},
			{Key: "coord.domain", Label: "Domain"},
		},
		Sections: []MetadataSchemaSection{
			{
				Title:       "Problem framing",
				Description: "What problem this idea is addressing and what it should produce.",
				Fields: []MetadataSchemaField{
					{Key: "coord.problem", Label: "Problem"},
					{Key: "coord.goal", Label: "Goal"},
					{Key: "coord.expected_output", Label: "Expected output"},
				},
			},
			{
				Title:       "Scope",
				Description: "Where this idea belongs in the coordination map.",
				Fields: []MetadataSchemaField{
					{Key: "coord.domain", Label: "Domain"},
				},
			},
		},
	},
	"task": {
		CoordType: "task",
		Fields: []MetadataSchemaField{
			{Key: "task.id", Label: "Task ID"},
			{Key: "task.status", Label: "Status"},
			{Key: "task.priority", Label: "Priority"},
			{Key: "task.parent", Label: "Parent"},
			{Key: "task.required_assets", Label: "Required assets"},
			{Key: "task.expected_result_type", Label: "Expected result type"},
		},
		Sections: []MetadataSchemaSection{
			{
				Title:       "Execution state",
				Description: "Current state, priority, and stable identifier for this task.",
				Fields: []MetadataSchemaField{
					{Key: "task.id", Label: "Task ID"},
					{Key: "task.status", Label: "Status"},
					{Key: "task.priority", Label: "Priority"},
				},
			},
			{
				Title:       "Dependencies",
				Description: "Upstream parent work and required assets needed to complete this task.",
				Fields: []MetadataSchemaField{
					{Key: "task.parent", Label: "Parent"},
					{Key: "task.required_assets", Label: "Required assets"},
					{Key: "task.expected_result_type", Label: "Expected result type"},
				},
			},
		},
	},
	"skill": {
		CoordType: "skill",
		Fields: []MetadataSchemaField{
			{Key: "skill.id", Label: "Skill ID"},
			{Key: "skill.category", Label: "Category"},
			{Key: "skill.inputs", Label: "Inputs"},
			{Key: "skill.outputs", Label: "Outputs"},
			{Key: "skill.depends_on", Label: "Depends on"},
			{Key: "skill.repo", Label: "Repository"},
		},
		Sections: []MetadataSchemaSection{
			{
				Title:       "Capability",
				Description: "Reusable skill identity and what category it belongs to.",
				Fields: []MetadataSchemaField{
					{Key: "skill.id", Label: "Skill ID"},
					{Key: "skill.category", Label: "Category"},
				},
			},
			{
				Title:       "Interface",
				Description: "Inputs, outputs, and dependencies needed to execute this skill.",
				Fields: []MetadataSchemaField{
					{Key: "skill.inputs", Label: "Inputs"},
					{Key: "skill.outputs", Label: "Outputs"},
					{Key: "skill.depends_on", Label: "Depends on"},
					{Key: "skill.repo", Label: "Repository"},
				},
			},
		},
	},
	"markdown": {
		CoordType: "markdown",
		Fields: []MetadataSchemaField{
			{Key: "md.kind", Label: "Document kind"},
			{Key: "md.collection", Label: "Collection"},
			{Key: "md.source_repo", Label: "Source repository"},
		},
		Sections: []MetadataSchemaSection{
			{
				Title:       "Knowledge asset",
				Description: "Document type and collection grouping for this durable note.",
				Fields: []MetadataSchemaField{
					{Key: "md.kind", Label: "Document kind"},
					{Key: "md.collection", Label: "Collection"},
					{Key: "md.source_repo", Label: "Source repository"},
				},
			},
		},
	},
	"code": {
		CoordType: "code",
		Fields: []MetadataSchemaField{
			{Key: "code.kind", Label: "Code kind"},
			{Key: "code.repo", Label: "Repository"},
			{Key: "code.entry", Label: "Entry point"},
			{Key: "code.language", Label: "Language"},
		},
		Sections: []MetadataSchemaSection{
			{
				Title:       "Implementation asset",
				Description: "Code kind, source repository, entry point, and implementation language.",
				Fields: []MetadataSchemaField{
					{Key: "code.kind", Label: "Code kind"},
					{Key: "code.repo", Label: "Repository"},
					{Key: "code.entry", Label: "Entry point"},
					{Key: "code.language", Label: "Language"},
				},
			},
		},
	},
	"agent": {
		CoordType: "agent",
		Fields: []MetadataSchemaField{
			{Key: "agent.id", Label: "Agent ID"},
			{Key: "agent.models", Label: "Models"},
			{Key: "agent.tools", Label: "Tools"},
			{Key: "agent.skills", Label: "Skills"},
			{Key: "agent.availability", Label: "Availability"},
		},
		Sections: []MetadataSchemaSection{
			{
				Title:       "Agent identity",
				Description: "Stable identity and availability of the agent profile.",
				Fields: []MetadataSchemaField{
					{Key: "agent.id", Label: "Agent ID"},
					{Key: "agent.availability", Label: "Availability"},
				},
			},
			{
				Title:       "Capabilities",
				Description: "Declared models, tools, and skill areas available from this agent.",
				Fields: []MetadataSchemaField{
					{Key: "agent.models", Label: "Models"},
					{Key: "agent.tools", Label: "Tools"},
					{Key: "agent.skills", Label: "Skills"},
				},
			},
		},
	},
}

func CoordMetadataSpec(coordType string) (CoordMetadataTypeSpec, bool) {
	spec, ok := coordMetadataSpecs[strings.ToLower(strings.TrimSpace(coordType))]
	return spec, ok
}

func CoordMetadataSchema(coordType string) *CoordMetadataTypeSpec {
	spec, ok := CoordMetadataSpec(coordType)
	if !ok {
		return nil
	}
	copySpec := CoordMetadataTypeSpec{
		CoordType: spec.CoordType,
		Fields:    append([]MetadataSchemaField(nil), spec.Fields...),
		Sections:  make([]MetadataSchemaSection, 0, len(spec.Sections)),
	}
	for _, section := range spec.Sections {
		copySpec.Sections = append(copySpec.Sections, MetadataSchemaSection{
			Title:       section.Title,
			Description: section.Description,
			Fields:      append([]MetadataSchemaField(nil), section.Fields...),
		})
	}
	return &copySpec
}

func coordFieldLabel(key string) string {
	key = strings.TrimSpace(key)
	if key == "" {
		return ""
	}
	for _, spec := range coordMetadataSpecs {
		for _, field := range spec.Fields {
			if field.Key == key {
				return field.Label
			}
		}
	}
	return key
}

