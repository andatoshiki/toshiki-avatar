package main

import (
	"flag"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"

	"github.com/disintegration/imaging"
	"github.com/kolesa-team/go-webp/encoder"
	"github.com/kolesa-team/go-webp/webp"
)

var images []string
var imgType string

// pickImage deterministically maps a hash -> avatar index
func pickImage(hash string) string {
	n := new(big.Int)
	n.SetString(hash, 16)
	idx := new(big.Int).Mod(n, big.NewInt(int64(len(images)))).Int64()
	return images[idx]
}

func main() {
	// CLI flags
	port := flag.Int("p", 8080, "Port to run the server on")
	imgTypeFlag := flag.String("t", "png", "Output type: png, jpg, webp")
	dir := flag.String("d", "./avatars", "Directory with avatar images")
	flag.Parse()

	imgType = *imgTypeFlag

	// load images from dir
	err := filepath.Walk(*dir, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() {
			images = append(images, path)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	if len(images) == 0 {
		log.Fatalf("no avatars found in %s", *dir)
	}

	// ensure deterministic order
	sort.Strings(images)

	// Main avatar handler
	http.HandleFunc("/avatar/", func(w http.ResponseWriter, r *http.Request) {
		hash := r.URL.Path[len("/avatar/"):]
		if hash == "" {
			http.Error(w, "missing hash", http.StatusBadRequest)
			return
		}

		imgPath := pickImage(hash)

		// ?s=size param
		sizeStr := r.URL.Query().Get("s")
		size, _ := strconv.Atoi(sizeStr)
		if size == 0 {
			size = 128
		}

		// JSON response mode
		format := r.URL.Query().Get("format")
		if format == "json" || r.Header.Get("Accept") == "application/json" {
			url := fmt.Sprintf("http://%s/avatar/%s?s=%d", r.Host, hash, size)
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{
  "hash": "%s",
  "url": "%s",
  "path": "%s",
  "size": %d,
  "type": "%s"
}`, hash, url, imgPath, size, imgType)
			return
		}

		// Serve image
		src, err := imaging.Open(imgPath)
		if err != nil {
			http.Error(w, "failed to open image", http.StatusInternalServerError)
			return
		}
		resized := imaging.Resize(src, size, size, imaging.Lanczos)

		switch imgType {
		case "jpg", "jpeg":
			w.Header().Set("Content-Type", "image/jpeg")
			imaging.Encode(w, resized, imaging.JPEG)
		case "webp":
			w.Header().Set("Content-Type", "image/webp")
			opts, _ := encoder.NewLossyEncoderOptions(encoder.PresetDefault, 75)
			if err := webp.Encode(w, resized, opts); err != nil {
				http.Error(w, "failed to encode webp", http.StatusInternalServerError)
			}
		default:
			w.Header().Set("Content-Type", "image/png")
			imaging.Encode(w, resized, imaging.PNG)
		}
	})

	addr := fmt.Sprintf(":%d", *port)
	fmt.Printf("Anime avatar server running on %s, dir=%s, type=%s\n", addr, *dir, imgType)
	log.Fatal(http.ListenAndServe(addr, nil))
}

