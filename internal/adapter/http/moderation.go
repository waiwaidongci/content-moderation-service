// Package implementation for policy-driven content moderation and human review.
package httpadapter

import (
	"github.com/ali/go-0821/content-moderation-service/internal/domain"
	"net/http"
	"strings"
)

func (s *Server) moderationPolicies(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(405)
		return
	}
	var p domain.ModerationPolicy
	if err := decode(r, &p); err != nil {
		writeErr(w, err)
		return
	}
	o, err := s.moderation.PutPolicy(r.Context(), p)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, 201, o)
}
func (s *Server) moderationDecision(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 4 {
		w.WriteHeader(404)
		return
	}
	id := parts[3]
	if len(parts) > 4 && parts[4] == "explain" {
		o, e := s.moderation.Explain(id)
		if e != nil {
			writeErr(w, e)
			return
		}
		writeJSON(w, 200, map[string]any{"sample_id": id, "hits": o})
		return
	}
	o, e := s.moderation.Decision(id)
	if e != nil {
		writeErr(w, e)
		return
	}
	writeJSON(w, 200, o)
}
func (s *Server) moderationAppeals(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var in struct {
			SampleID string `json:"sample_id"`
			Reason   string `json:"reason"`
		}
		if e := decode(r, &in); e != nil {
			writeErr(w, e)
			return
		}
		o, e := s.moderation.Appeal(r.Context(), in.SampleID, in.Reason)
		if e != nil {
			writeErr(w, e)
			return
		}
		writeJSON(w, 201, o)
		return
	}
	if r.Method == "PATCH" {
		var in struct {
			ID     string `json:"id"`
			Status string `json:"status"`
		}
		if e := decode(r, &in); e != nil {
			writeErr(w, e)
			return
		}
		o, e := s.moderation.ResolveAppeal(in.ID, in.Status)
		if e != nil {
			writeErr(w, e)
			return
		}
		writeJSON(w, 200, o)
		return
	}
	w.WriteHeader(405)
}
func (s *Server) moderationDictionaries(w http.ResponseWriter, r *http.Request) {
	var in struct {
		ID    string   `json:"id"`
		Terms []string `json:"terms"`
	}
	if e := decode(r, &in); e != nil {
		writeErr(w, e)
		return
	}
	if e := s.moderation.AddDictionary(in.ID, in.Terms); e != nil {
		writeErr(w, e)
		return
	}
	writeJSON(w, 201, map[string]any{"id": in.ID, "terms": len(in.Terms)})
}
func (s *Server) publishModerationPolicy(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(405)
		return
	}
	var in struct {
		ID string `json:"id"`
	}
	if err := decode(r, &in); err != nil {
		writeErr(w, err)
		return
	}
	o, err := s.moderation.Publish(r.Context(), in.ID)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, 200, o)
}
func (s *Server) moderateContent(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(405)
		return
	}
	var in domain.Sample
	if err := decode(r, &in); err != nil {
		writeErr(w, err)
		return
	}
	o, err := s.moderation.Check(r.Context(), in)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, 200, o)
}
func (s *Server) claimReviewTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(405)
		return
	}
	var in struct {
		Reviewer string `json:"reviewer"`
	}
	if err := decode(r, &in); err != nil {
		writeErr(w, err)
		return
	}
	o, err := s.moderation.Claim(r.Context(), in.Reviewer)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, 200, o)
}
func (s *Server) submitReviewTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(405)
		return
	}
	var in struct {
		TaskID   string          `json:"task_id"`
		Reviewer string          `json:"reviewer"`
		Decision domain.Decision `json:"decision"`
	}
	if err := decode(r, &in); err != nil {
		writeErr(w, err)
		return
	}
	o, err := s.moderation.Submit(r.Context(), in.TaskID, in.Reviewer, in.Decision)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, 200, o)
}
