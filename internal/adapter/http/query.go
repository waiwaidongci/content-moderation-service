// Package implementation for policy-driven content moderation and human review.
package httpadapter

import (
	"github.com/ali/go-0821/content-moderation-service/internal/infrastructure/config"
	"net/http"
	"strconv"
	"strings"
)

type Query struct {
	ProjectID     string
	EnvironmentID string
	Page          int
	Size          int
	Tags          map[string]string
}

func ParseQuery(r *http.Request) Query {
	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	size, _ := strconv.Atoi(q.Get("size"))
	return Query{ProjectID: strings.TrimSpace(q.Get("project_id")), EnvironmentID: strings.TrimSpace(q.Get("environment_id")), Page: page, Size: size, Tags: config.ParseTags(q.Get("tags"))}
}
func boolQuery(r *http.Request, key string, def bool) bool {
	v := r.URL.Query().Get(key)
	if v == "" {
		return def
	}
	b, e := strconv.ParseBool(v)
	if e != nil {
		return def
	}
	return b
}
