package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/istvzsig/eve-gate-watch/internal/intel"
	"github.com/istvzsig/eve-gate-watch/internal/systems"
)

type Handler struct {
	intelService   *intel.Service
	systemsService *systems.Service
}

func NewHandler(
	intelService *intel.Service,
	systemsService *systems.Service,
) *Handler {
	return &Handler{
		intelService:   intelService,
		systemsService: systemsService,
	}
}

func (h *Handler) HandleDScan(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	systemName := r.URL.Query().Get("system")

	result, err := h.intelService.DScan(systemName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) HandleRecentPVP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	systemName := r.URL.Query().Get("system")

	result, err := h.intelService.RecentPVP(systemName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) HandleSystems(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query().Get("q")

	results, _ := h.systemsService.Search(r.Context(), query)

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(results); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
