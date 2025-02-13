// list_characters.go
package httpserver

import (
	"net/http"
	"fmt"
	"os"
	"sort"
	"path/filepath"
	"time"

	"github.com/go-chi/chi"
)

// listCharactersHandler is used to list available characters.
type listCharactersHandler struct {
	encoder          *encoder
	d2sPath          string
	characterService characterService
}

// Add this new struct to hold file info
type CharacterFileInfo struct {
	Name          string    `json:"name"`
	LastModified  time.Time `json:"last_modified"`
}


func newListCharactersHandler() *listCharactersHandler {
	return &listCharactersHandler{}
}

func (h *listCharactersHandler) Routes(router chi.Router) {
	router.Get("/", h.list)
}

func (h *listCharactersHandler) list(w http.ResponseWriter, r *http.Request) {
    var d2sPath            = env.String("D2S_PATH", "")

    // Read directory entries
    files, err := os.ReadDir(d2sPath)
    if err != nil {
        h.encoder.Error(w, fmt.Errorf("failed to read directory: %v", err))
        return
    }

    var characterFiles []CharacterFileInfo

    for _, file := range files {
        // Skip directories and only process character files (no extension needed)
        if !file.IsDir() {
            fullPath := filepath.Join(d2sPath, file.Name())

            // Get full file info including modification time
            fileInfo, err := os.Stat(fullPath)
            if err != nil {
                continue // Skip files that can't be accessed
            }

            characterFiles = append(characterFiles, CharacterFileInfo{
                Name:          file.Name(),
                LastModified:  fileInfo.ModTime().UTC(), // Get UTC time for consistency
            })
        }
    }

    // Sort by last modified time descending (most recent first)
    sort.SliceStable(characterFiles, func(i, j int) bool {
        return characterFiles[i].LastModified.After(characterFiles[j].LastModified)
    })

    h.encoder.Response(w, struct {
        Characters []CharacterFileInfo `json:"characters"`
    }{
        Characters: characterFiles,
    })
}
