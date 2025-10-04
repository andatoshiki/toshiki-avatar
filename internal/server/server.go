package server


import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"strconv"
	"github.com/andatoshiki/toshiki-avatar/internal/avatar"
	"github.com/andatoshiki/toshiki-avatar/internal/api"
	apperr "github.com/andatoshiki/toshiki-avatar/internal/errors"
	encode "github.com/andatoshiki/toshiki-avatar/internal/encode"
)

//go:embed web/* web/**/*
var staticFiles embed.FS



// StaticHandler serves the embedded static frontend (Next.js export)
func StaticHandler() http.Handler {
	content, err := fs.Sub(staticFiles, "web")
	if err != nil {
		panic("failed to get embedded static files: " + err.Error())
	}
	return http.FileServer(http.FS(content))
}

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

	if len(s.AvatarService.Images) == 0 {
		http.Error(w, apperr.ErrAvatarNotFound.Error(), http.StatusNotFound)
		return
	}

	imgPath := s.AvatarService.PickImage(hash)

       sizeStr := r.URL.Query().Get("s")
       size, err := strconv.Atoi(sizeStr)
       if sizeStr == "" || err != nil || size <= 0 {
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
	       http.Error(w, "failed to decode or resize image", http.StatusInternalServerError)
	       return
       }

	switch s.ImgType {
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

func (s *Server) Start(port int) error {
	// randomHandler := api.NewRandomHandler(s.AvatarService, s.ImgType)
	http.Handle("/", StaticHandler())
	http.HandleFunc("/avatar/", s.AvatarHandler)
	// http.Handle("/random", randomHandler)
	http.HandleFunc("/healthz", api.HealthHandler)
	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("toshiki-anime avatar server running on %s\n", addr)
	return http.ListenAndServe(addr, nil)
}
