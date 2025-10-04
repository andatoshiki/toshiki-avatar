package avatar

import (
	"math/big"
	"sort"
	"github.com/disintegration/imaging"
	"image"
	"math/rand"
	"time"
)

// RandomImage returns a random image path from the available images
func (a *AvatarService) RandomImage() string {
	rand.Seed(time.Now().UnixNano())
	idx := rand.Intn(len(a.Images))
	return a.Images[idx]
}

type AvatarService struct {
	Images  []string
}

// NewAvatarService creates a new AvatarService and sorts the images deterministically
func NewAvatarService(images []string) *AvatarService {
	sort.Strings(images)
	return &AvatarService{Images: images}
}

// PickImage deterministically maps a hash -> avatar index
func (a *AvatarService) PickImage(hash string) string {
	n := new(big.Int)
	n.SetString(hash, 16)
	idx := new(big.Int).Mod(n, big.NewInt(int64(len(a.Images)))).Int64()
	return a.Images[idx]
}

// ResizeImage resizes the image at the given path to the specified size
func ResizeImage(path string, size int) (image.Image, error) {
	src, err := imaging.Open(path)
	if err != nil {
		return nil, err
	}
	resized := imaging.Resize(src, size, size, imaging.Lanczos)
	return resized, nil
}


