package httpapi

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/iszsig/eve-gate-watch/internal/dscan"
	"github.com/iszsig/eve-gate-watch/internal/intel"
	"github.com/iszsig/eve-gate-watch/internal/ship"
)

type Handler struct {
	intelService *intel.Service
}

type DScanRecord struct {
	SystemName string      `json:"systemName"`
	Ships      []ship.Ship `json:"ships"`
}

type StatusResponse struct {
	Status string `json:"status"`
}

func NewHandler(intelService *intel.Service) *Handler {
	return &Handler{
		intelService: intelService,
	}
}

func (h *Handler) HandleDScan(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	record := &DScanRecord{}

	if err := json.NewDecoder(r.Body).Decode(record); err != nil {
		http.Error(w, "failed to decode record", http.StatusInternalServerError)
		return
	}

	ships := make([]ship.Ship, 0, len(record.Ships))

	for _, s := range record.Ships {
		ship := ship.Ship{
			Name: s.Name,
		}
		ships = append(ships, ship)
	}

	input := dscan.DScanInput{
		SystemName: record.SystemName,
		Ships:      ships,
	}

	if err := h.intelService.RecordDscan(
		r.Context(),
		input,
	); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(
		StatusResponse{
			Status: "created",
		}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) HandleIntel(w http.ResponseWriter, r *http.Request) {
	systemName := r.URL.Query().Get("system")

	report, err := h.intelService.RecentIntel(
		r.Context(),
		systemName,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(report); err != nil {
		log.Printf("encode intel response: %v", err)
	}
}
