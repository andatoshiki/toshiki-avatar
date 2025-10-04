package api
package api

import (
	"encoding/json"
	"net/http"
	"github.com/andatoshiki/toshiki-avatar/internal/avatar"
	"strconv"
)

type RandomHandler struct {
	AvatarService *avatar.AvatarService
	ImgType       string
}

func NewRandomHandler(avatarService *avatar.AvatarService, imgType string) *RandomHandler {
	return &RandomHandler{
		AvatarService: avatarService,
		ImgType:       imgType,
	}
}

func (h *RandomHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	imgPath := h.AvatarService.RandomImage()
	sizeStr := r.URL.Query().Get("s")
	size, _ := strconv.Atoi(sizeStr)
	if size == 0 {
		size = 128
	}

	format := r.URL.Query().Get("format")
	if format == "json" || r.Header.Get("Accept") == "application/json" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"path": imgPath,
			"size": size,
			"type": h.ImgType,
		})
		return
	}

	resized, err := avatar.ResizeImage(imgPath, size)
	if err != nil {
		http.Error(w, "failed to open image", http.StatusInternalServerError)
		return
	}

	switch h.ImgType {
	case "jpg", "jpeg":
		w.Header().Set("Content-Type", "image/jpeg")
		avatar.EncodeJPEG(w, resized)
	default:
		w.Header().Set("Content-Type", "image/png")
		avatar.EncodePNG(w, resized)
	}
}
