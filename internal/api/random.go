package api

import (
	"encoding/json"
	"net/http"
	"github.com/andatoshiki/toshiki-avatar/internal/avatar"
	"strconv"
	apperr "github.com/andatoshiki/toshiki-avatar/internal/errors"
	encode "github.com/andatoshiki/toshiki-avatar/internal/encode"
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
    // Prevent caching so each refresh gets a new random avatar
    w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate")
    w.Header().Set("Pragma", "no-cache")
    w.Header().Set("Expires", "0")
    w.Header().Set("Surrogate-Control", "no-store")

	if len(h.AvatarService.Images) == 0 {
		http.Error(w, apperr.ErrAvatarNotFound.Error(), http.StatusNotFound)
		return
	}

	imgPath := h.AvatarService.RandomImage()
       sizeStr := r.URL.Query().Get("s")
       size, err := strconv.Atoi(sizeStr)
       if sizeStr == "" || err != nil || size <= 0 {
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
	       http.Error(w, "failed to decode or resize image", http.StatusInternalServerError)
	       return
       }

	switch h.ImgType {
	case "jpg", "jpeg":
		w.Header().Set("Content-Type", "image/jpeg")
		err := encode.EncodeJPEG(w, resized)
		if err != nil {
			http.Error(w, "failed to encode jpeg", http.StatusInternalServerError)
		}
	case "webp":
		w.Header().Set("Content-Type", "image/webp")
		err := encode.EncodeWebP(w, resized)
		if err != nil {
			http.Error(w, "failed to encode webp", http.StatusInternalServerError)
		}
	case "png":
		w.Header().Set("Content-Type", "image/png")
		err := encode.EncodePNG(w, resized)
		if err != nil {
			http.Error(w, "failed to encode png", http.StatusInternalServerError)
		}
	default:
		http.Error(w, apperr.ErrInvalidImageType.Error(), http.StatusBadRequest)
	}
}
