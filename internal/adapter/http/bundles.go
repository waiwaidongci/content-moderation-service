package httpadapter

import (
	"github.com/ali/go-0821/content-moderation-service/internal/domain"
	"net/http"
)

func (s *Server) moderationWorkspaces(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(405)
		return
	}
	var v domain.ModerationWorkspace
	if err := decode(r, &v); err != nil {
		writeErr(w, err)
		return
	}
	if err := s.bundles.CreateModerationWorkspace(r.Context(), v); err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, 201, v)
}
func (s *Server) moderationChannels(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		o, e := s.bundles.ListModerationChannels(r.Context(), r.URL.Query().Get("workspace_id"))
		if e != nil {
			writeErr(w, e)
			return
		}
		writeJSON(w, 200, o)
		return
	}
	var v domain.ModerationChannel
	if err := decode(r, &v); err != nil {
		writeErr(w, err)
		return
	}
	if err := s.bundles.CreateModerationChannel(r.Context(), v); err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, 201, v)
}
func (s *Server) ruleBundles(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		o, e := s.bundles.ListRuleBundles(r.Context(), r.URL.Query().Get("workspace_id"), r.URL.Query().Get("channel_id"))
		if e != nil {
			writeErr(w, e)
			return
		}
		writeJSON(w, 200, o)
		return
	}
	var v domain.RuleBundle
	if err := decode(r, &v); err != nil {
		writeErr(w, err)
		return
	}
	if err := s.bundles.CreateRuleBundle(r.Context(), v); err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, 201, v)
}
