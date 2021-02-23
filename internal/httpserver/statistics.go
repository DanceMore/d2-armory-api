package httpserver

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/nokka/d2-armory-api/internal/domain"
)

// statisticsService encapsulates the business logic around statistics.
type statisticsService interface {
	// Gets the character statistics.
	GetCharacter(ctx context.Context, character string) (*domain.StatisticsRequest, error)

	// Parse parses a character binary.
	Parse(ctx context.Context, stats []domain.StatisticsRequest) error
}

// statisticsHandler is used to put parse statistics requests.
type statisticsHandler struct {
	encoder           *encoder
	statisticsService statisticsService
}

func (h statisticsHandler) Routes(router chi.Router) {
	router.Post("/", h.postStatistics)
	router.Get("/", h.getStatistics)
}

func (h statisticsHandler) postStatistics(w http.ResponseWriter, r *http.Request) {
	var stats []domain.StatisticsRequest
	if err := json.NewDecoder(r.Body).Decode(&stats); err != nil {
		h.encoder.Error(w, err)
		return
	}

	// Pass the request context in order to make use of cancellation for lower level work.
	err := h.statisticsService.Parse(r.Context(), stats)
	if err != nil {
		h.encoder.Error(w, err)
		return
	}

	h.encoder.Response(w, "statistics")
}

func (h statisticsHandler) getStatistics(w http.ResponseWriter, r *http.Request) {
	characterName := r.URL.Query().Get("character")

	//Pass the request context in order to make use of cancellation for lower level work.
	stats, err := h.statisticsService.GetCharacter(r.Context(), characterName)
	if err != nil {
		h.encoder.Error(w, err)
		return
	}

	h.encoder.Response(w, stats)
}

func newStatisticsHandler(encoder *encoder, statisticsService statisticsService) *statisticsHandler {
	return &statisticsHandler{
		encoder:           encoder,
		statisticsService: statisticsService,
	}
}
