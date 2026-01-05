// Package mcp provides Model Context Protocol server and client for MakeaTUI
package mcp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/makeatui/makeatui/pkg/agent"
	"github.com/makeatui/makeatui/pkg/ai"
	"github.com/makeatui/makeatui/pkg/templates"
)

// Server implements an MCP server for TUI design
type Server struct {
	sessions       map[string]*Session
	sessionMu      sync.RWMutex
	templateEngine *templates.TemplateEngine
	port           int
}

// Session represents an active design session
type Session struct {
	ID       string     `json:"id"`
	Name     string     `json:"name"`
	API      *agent.API `json:"-"`
	AIAgent  *ai.TUIAgent `json:"-"`
}

// NewServer creates a new MCP server
func NewServer(port int) *Server {
	return &Server{
		sessions:       make(map[string]*Session),
		templateEngine: templates.NewTemplateEngine(),
		port:           port,
	}
}

// Start starts the MCP server
func (s *Server) Start() error {
	mux := http.NewServeMux()

	// Session management
	mux.HandleFunc("/sessions", s.handleSessions)
	mux.HandleFunc("/sessions/", s.handleSession)

	// Tools
	mux.HandleFunc("/tools/add_component", s.handleAddComponent)
	mux.HandleFunc("/tools/move_component", s.handleMoveComponent)
	mux.HandleFunc("/tools/remove_component", s.handleRemoveComponent)
	mux.HandleFunc("/tools/set_text", s.handleSetText)
	mux.HandleFunc("/tools/generate", s.handleGenerate)
	mux.HandleFunc("/tools/export", s.handleExport)

	// Templates
	mux.HandleFunc("/templates", s.handleTemplates)
	mux.HandleFunc("/templates/apply", s.handleApplyTemplate)

	// Resources
	mux.HandleFunc("/resources/canvas", s.handleCanvasResource)
	mux.HandleFunc("/resources/components", s.handleComponentsResource)

	// Health check
	mux.HandleFunc("/health", s.handleHealth)

	fmt.Printf("🚀 MakeaTUI MCP Server running on port %d\n", s.port)
	return http.ListenAndServe(fmt.Sprintf(":%d", s.port), mux)
}

// CreateSession creates a new design session
func (s *Server) CreateSession(name string) *Session {
	s.sessionMu.Lock()
	defer s.sessionMu.Unlock()

	id := generateSessionID()
	session := &Session{
		ID:      id,
		Name:    name,
		API:     agent.NewAPI(name),
		AIAgent: ai.NewTUIAgent(),
	}

	s.sessions[id] = session
	return session
}

// GetSession retrieves a session by ID
func (s *Server) GetSession(id string) *Session {
	s.sessionMu.RLock()
	defer s.sessionMu.RUnlock()
	return s.sessions[id]
}

// HTTP Handlers

func (s *Server) handleSessions(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.sessionMu.RLock()
		sessions := make([]*Session, 0, len(s.sessions))
		for _, sess := range s.sessions {
			sessions = append(sessions, sess)
		}
		s.sessionMu.RUnlock()
		respondJSON(w, sessions)

	case http.MethodPost:
		var req struct {
			Name string `json:"name"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "invalid request")
			return
		}
		session := s.CreateSession(req.Name)
		respondJSON(w, session)

	default:
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (s *Server) handleSession(w http.ResponseWriter, r *http.Request) {
	// Extract session ID from path
	id := r.URL.Path[len("/sessions/"):]
	session := s.GetSession(id)
	if session == nil {
		respondError(w, http.StatusNotFound, "session not found")
		return
	}

	switch r.Method {
	case http.MethodGet:
		respondJSON(w, session)
	case http.MethodDelete:
		s.sessionMu.Lock()
		delete(s.sessions, id)
		s.sessionMu.Unlock()
		respondJSON(w, map[string]string{"status": "deleted"})
	default:
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (s *Server) handleAddComponent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req struct {
		SessionID string `json:"session_id"`
		Type      string `json:"type"`
		Name      string `json:"name"`
		Text      string `json:"text"`
		X         int    `json:"x"`
		Y         int    `json:"y"`
		Width     int    `json:"width"`
		Height    int    `json:"height"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request")
		return
	}

	session := s.GetSession(req.SessionID)
	if session == nil {
		respondError(w, http.StatusNotFound, "session not found")
		return
	}

	var id string
	switch req.Type {
	case "box":
		id = session.API.AddBox(req.Name, req.Text, req.X, req.Y, req.Width, req.Height)
	case "text":
		id = session.API.AddText(req.Name, req.Text, req.X, req.Y)
	case "button":
		id = session.API.AddButton(req.Name, req.Text, req.X, req.Y)
	default:
		respondError(w, http.StatusBadRequest, "unknown component type")
		return
	}

	respondJSON(w, map[string]string{"id": id})
}

func (s *Server) handleMoveComponent(w http.ResponseWriter, r *http.Request) {
	var req struct {
		SessionID   string `json:"session_id"`
		ComponentID string `json:"component_id"`
		X           int    `json:"x"`
		Y           int    `json:"y"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request")
		return
	}

	session := s.GetSession(req.SessionID)
	if session == nil {
		respondError(w, http.StatusNotFound, "session not found")
		return
	}

	_ = session.API.Move(req.ComponentID, req.X, req.Y)
	respondJSON(w, map[string]string{"status": "moved"})
}

func (s *Server) handleRemoveComponent(w http.ResponseWriter, r *http.Request) {
	var req struct {
		SessionID   string `json:"session_id"`
		ComponentID string `json:"component_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request")
		return
	}

	session := s.GetSession(req.SessionID)
	if session == nil {
		respondError(w, http.StatusNotFound, "session not found")
		return
	}

	_ = session.API.Delete(req.ComponentID)
	respondJSON(w, map[string]string{"status": "removed"})
}

func (s *Server) handleSetText(w http.ResponseWriter, r *http.Request) {
	var req struct {
		SessionID   string `json:"session_id"`
		ComponentID string `json:"component_id"`
		Text        string `json:"text"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request")
		return
	}

	session := s.GetSession(req.SessionID)
	if session == nil {
		respondError(w, http.StatusNotFound, "session not found")
		return
	}

	_ = session.API.SetText(req.ComponentID, req.Text)
	respondJSON(w, map[string]string{"status": "updated"})
}

func (s *Server) handleGenerate(w http.ResponseWriter, r *http.Request) {
	var req struct {
		SessionID   string `json:"session_id"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request")
		return
	}

	session := s.GetSession(req.SessionID)
	if session == nil {
		respondError(w, http.StatusNotFound, "session not found")
		return
	}

	api, err := session.AIAgent.GenerateFromDescription(req.Description)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	session.API = api
	respondJSON(w, map[string]string{"status": "generated"})
}

func (s *Server) handleExport(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session_id")
	format := r.URL.Query().Get("format")

	session := s.GetSession(sessionID)
	if session == nil {
		respondError(w, http.StatusNotFound, "session not found")
		return
	}

	switch format {
	case "json":
		jsonStr, _ := session.API.ExportJSON()
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(jsonStr))
	default:
		code := session.API.Export()
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write([]byte(code))
	}
}

func (s *Server) handleTemplates(w http.ResponseWriter, r *http.Request) {
	templates := s.templateEngine.List()
	respondJSON(w, templates)
}

func (s *Server) handleApplyTemplate(w http.ResponseWriter, r *http.Request) {
	var req struct {
		SessionID    string `json:"session_id"`
		TemplateName string `json:"template_name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request")
		return
	}

	session := s.GetSession(req.SessionID)
	if session == nil {
		respondError(w, http.StatusNotFound, "session not found")
		return
	}

	api, err := s.templateEngine.Apply(req.TemplateName)
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}

	session.API = api
	respondJSON(w, map[string]string{"status": "applied"})
}

func (s *Server) handleCanvasResource(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session_id")
	session := s.GetSession(sessionID)
	if session == nil {
		respondError(w, http.StatusNotFound, "session not found")
		return
	}

	jsonStr, _ := session.API.ExportJSON()
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(jsonStr))
}

func (s *Server) handleComponentsResource(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session_id")
	session := s.GetSession(sessionID)
	if session == nil {
		respondError(w, http.StatusNotFound, "session not found")
		return
	}

	respondJSON(w, session.API.ListComponents())
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, map[string]string{"status": "healthy"})
}

// Helper functions
func respondJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
}

var sessionCounter int
var sessionMu sync.Mutex

func generateSessionID() string {
	sessionMu.Lock()
	defer sessionMu.Unlock()
	sessionCounter++
	return fmt.Sprintf("session_%d", sessionCounter)
}

