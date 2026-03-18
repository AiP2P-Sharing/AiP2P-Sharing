package newsplugin

import (
	"encoding/json"
	"fmt"
	"strings"
)

const humanArticlePreviewLimit = 200

func StandardArticleClasses() []string {
	return []string{"note", "json", "code", "link", "profile"}
}

type ArticleLink struct {
	Label        string
	Value        string
	CompactValue string
}

func ArticlePreviewForPost(post Post) string {
	if strings.TrimSpace(post.Summary) != "" {
		return summarize(post.Summary, humanArticlePreviewLimit)
	}
	if code := strings.TrimSpace(ArticleCodeForPost(post)); code != "" {
		return summarize(code, humanArticlePreviewLimit)
	}
	return summarize(ArticleBodyForPost(post), humanArticlePreviewLimit)
}

func StandardArticleClassForPost(post Post) string {
	ext := post.Message.Extensions
	for _, key := range []string{"article.class", "content.class"} {
		value := strings.ToLower(strings.TrimSpace(extensionText(ext[key])))
		switch value {
		case "note", "json", "code", "link", "profile":
			return value
		}
	}
	if ArticleIsJSONForPost(post) {
		return "json"
	}
	if ArticleIsCodeForPost(post) {
		return "code"
	}
	if PostCoordType(post) == "agent" || strings.TrimSpace(PublisherBioForPost(post)) != "" {
		return "profile"
	}
	if len(ArticleLinksForPost(post)) > 0 && strings.TrimSpace(ArticleBodyForPost(post)) == "" {
		return "link"
	}
	return "note"
}

func ArticleHasMoreForPost(post Post) bool {
	full := strings.TrimSpace(ArticleBodyForPost(post))
	if ArticleIsCodeForPost(post) {
		full = strings.TrimSpace(ArticleCodeForPost(post))
	}
	preview := strings.TrimSpace(ArticlePreviewForPost(post))
	return full != "" && full != preview
}

func ArticleBodyForPost(post Post) string {
	return strings.TrimSpace(post.Body)
}

func ArticleIsCodeForPost(post Post) bool {
	coordType := strings.ToLower(strings.TrimSpace(PostCoordType(post)))
	if coordType == "code" {
		return true
	}
	if ArticleIsJSONForPost(post) {
		return true
	}
	ext := post.Message.Extensions
	for _, key := range []string{"article.content_type", "content.type", "body.type", "render.mode", "article.format"} {
		if strings.EqualFold(strings.TrimSpace(extensionText(ext[key])), "code") {
			return true
		}
	}
	return strings.TrimSpace(extensionText(ext["code.language"])) != ""
}

func ArticleCodeLanguageForPost(post Post) string {
	if ArticleIsJSONForPost(post) {
		return "json"
	}
	value := strings.TrimSpace(extensionText(post.Message.Extensions["code.language"]))
	if value == "" {
		return "text"
	}
	return value
}

func ArticleIsJSONForPost(post Post) bool {
	ext := post.Message.Extensions
	for _, key := range []string{"article.content_type", "content.type", "body.type", "article.format"} {
		if strings.EqualFold(strings.TrimSpace(extensionText(ext[key])), "json") {
			return true
		}
	}
	body := strings.TrimSpace(post.Body)
	if body == "" {
		return false
	}
	if !(strings.HasPrefix(body, "{") || strings.HasPrefix(body, "[")) {
		return false
	}
	var payload any
	return json.Unmarshal([]byte(body), &payload) == nil
}

func ArticleCodeForPost(post Post) string {
	ext := post.Message.Extensions
	for _, key := range []string{"article.code", "code.snippet", "snippet", "body.code"} {
		if value := strings.TrimSpace(extensionText(ext[key])); value != "" {
			return value
		}
	}
	if ArticleIsJSONForPost(post) {
		var payload any
		if json.Unmarshal([]byte(post.Body), &payload) == nil {
			if pretty, err := json.MarshalIndent(payload, "", "  "); err == nil {
				return string(pretty)
			}
		}
	}
	return strings.TrimSpace(post.Body)
}

func PublisherBioForPost(post Post) string {
	ext := post.Message.Extensions
	for _, key := range []string{"publisher.bio", "author.bio", "profile.bio", "signature.bio", "agent.bio"} {
		if value := strings.TrimSpace(extensionText(ext[key])); value != "" {
			return value
		}
	}
	return ""
}

func ArticleLinksForPost(post Post) []ArticleLink {
	links := make([]ArticleLink, 0, 6)
	seen := map[string]struct{}{}
	add := func(label, value string) {
		value = strings.TrimSpace(value)
		label = strings.TrimSpace(label)
		if value == "" {
			return
		}
		if _, ok := seen[value]; ok {
			return
		}
		seen[value] = struct{}{}
		if label == "" {
			label = "Link"
		}
		links = append(links, ArticleLink{
			Label:        label,
			Value:        value,
			CompactValue: summarize(value, 64),
		})
	}

	add("Source", post.SourceURL)
	ext := post.Message.Extensions
	for _, item := range []struct {
		Key   string
		Label string
	}{
		{Key: "skill.repo", Label: "Repository"},
		{Key: "code.repo", Label: "Repository"},
		{Key: "md.source_repo", Label: "Repository"},
		{Key: "external.url", Label: "External"},
		{Key: "github", Label: "GitHub"},
		{Key: "youtube", Label: "YouTube"},
		{Key: "via", Label: "Via"},
	} {
		add(item.Label, extensionText(ext[item.Key]))
	}
	for _, key := range []string{"external.links", "article.links", "links"} {
		for _, value := range extensionStrings(ext[key]) {
			add("Link", value)
		}
	}
	return links
}

func ReplyPreview(reply Reply) string {
	return summarize(strings.TrimSpace(reply.Body), humanArticlePreviewLimit)
}

func ReplyHasMore(reply Reply) bool {
	full := strings.TrimSpace(reply.Body)
	return full != "" && full != ReplyPreview(reply)
}

func ExternalLinkText(value any) string {
	return fmt.Sprint(value)
}
