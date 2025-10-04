
package config

import "flag"

type Config struct {
	Port     int
	ImgType  string
	Dir      string
	ListFile string
}

func ParseFlags() *Config {
	port := flag.Int("p", 8080, "Port to run the server on")
	imgType := flag.String("t", "png", "Output type: png, jpg")
	dir := flag.String("d", "./avatars", "Directory with avatar images")
	listFile := flag.String("l", "", "File with list of image paths or URLs (one per line)")
	flag.Parse()

	return &Config{
		Port:     *port,
		ImgType:  *imgType,
		Dir:      *dir,
		ListFile: *listFile,
	}
}
