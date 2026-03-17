package newsplugin

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

func FeedOptionsFromRequest(r *http.Request) FeedOptions {
	if r == nil {
		return FeedOptions{Now: time.Now()}
	}
	return FeedOptions{
		Channel:   strings.TrimSpace(r.URL.Query().Get("channel")),
		Topic:     strings.TrimSpace(r.URL.Query().Get("topic")),
		Source:    strings.TrimSpace(r.URL.Query().Get("source")),
		MetaKey:   strings.TrimSpace(r.URL.Query().Get("meta_key")),
		MetaValue: strings.TrimSpace(r.URL.Query().Get("meta_value")),
		Sort:      strings.TrimSpace(r.URL.Query().Get("sort")),
		Query:     strings.TrimSpace(r.URL.Query().Get("q")),
		Window:    CanonicalWindow(r.URL.Query().Get("window")),
		Page:      parsePositiveRequestInt(r.URL.Query().Get("page"), 1),
		PageSize:  parseRequestFeedPageSize(r.URL.Query().Get("page_size")),
		Now:       time.Now(),
	}
}

func CanonicalWindow(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "24h":
		return "24h"
	case "7d":
		return "7d"
	case "30d":
		return "30d"
	default:
		return ""
	}
}

func parsePositiveRequestInt(raw string, fallback int) int {
	value, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil || value < 1 {
		return fallback
	}
	return value
}

func parseRequestFeedPageSize(raw string) int {
	value := parsePositiveRequestInt(raw, 20)
	if value < 1 {
		return 20
	}
	if value > 200 {
		return 200
	}
	return value
}
