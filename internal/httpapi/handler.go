package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/iszsig/eve-gate-watch/internal/intel"
)

type Handler struct {
	intelService *intel.Service
}

func NewHandler(intelService *intel.Service) *Handler {
	return &Handler{
		intelService: intelService,
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
