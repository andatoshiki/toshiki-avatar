
package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"github.com/andatoshiki/toshiki-avatar/internal/avatar"
	"github.com/andatoshiki/toshiki-avatar/internal/api"
)

type Server struct {
	AvatarService *avatar.AvatarService
	ImgType       string
}

func NewServer(avatarService *avatar.AvatarService, imgType string) *Server {
	return &Server{
		AvatarService: avatarService,
		ImgType:       imgType,
	}
}

func (s *Server) AvatarHandler(w http.ResponseWriter, r *http.Request) {
	hash := r.URL.Path[len("/avatar/"):]
	if hash == "" {
		http.Error(w, "missing hash", http.StatusBadRequest)
		return
	}

	imgPath := s.AvatarService.PickImage(hash)

	sizeStr := r.URL.Query().Get("s")
	size, _ := strconv.Atoi(sizeStr)
	if size == 0 {
		size = 128
	}

	format := r.URL.Query().Get("format")
	if format == "json" || r.Header.Get("Accept") == "application/json" {
		url := fmt.Sprintf("http://%s/avatar/%s?s=%d", r.Host, hash, size)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"hash": hash,
			"url":  url,
			"path": imgPath,
			"size": size,
			"type": s.ImgType,
		})
		return
	}

	resized, err := avatar.ResizeImage(imgPath, size)
	if err != nil {
		http.Error(w, "failed to open image", http.StatusInternalServerError)
		return
	}

	switch s.ImgType {
	case "jpg", "jpeg":
		w.Header().Set("Content-Type", "image/jpeg")
		avatar.EncodeJPEG(w, resized)
	default:
		w.Header().Set("Content-Type", "image/png")
		avatar.EncodePNG(w, resized)
	}
}

func (s *Server) Start(port int) error {
	randomHandler := api.NewRandomHandler(s.AvatarService, s.ImgType)
	http.HandleFunc("/avatar/", s.AvatarHandler)
	http.Handle("/random", randomHandler)
	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("Anime avatar server running on %s\n", addr)
	return http.ListenAndServe(addr, nil)
}
