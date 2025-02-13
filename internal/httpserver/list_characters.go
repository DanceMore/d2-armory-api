// list_characters.go
package httpserver

import (
	"net/http"
	
	"github.com/go-chi/chi"
)

// listCharactersHandler is used to list available characters.
type listCharactersHandler struct{}

func newListCharactersHandler() *listCharactersHandler {
	return &listCharactersHandler{}
}

func (h *listCharactersHandler) Routes(router chi.Router) {
	router.Get("/", h.list)
}

func (h *listCharactersHandler) list(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement actual listing logic
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("List all characters"))
}
