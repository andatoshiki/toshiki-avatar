package config

import (
	"flag"
	"fmt"
	"github.com/andatoshiki/toshiki-avatar/internal"
	"os"
	"runtime/debug"
)

type Config struct {
	Port     int
	ImgType  string
	Dir      string
	ListFile string
}

func ParseFlags() *Config {
	port := flag.Int("p", 8080, "Port to run the server on")
	imgType := flag.String("t", "png", "Output type: png, jpg, webp")
	dir := flag.String("d", "./avatars", "Directory with avatar images")
	listFile := flag.String("l", "", "File with list of image paths or URLs (one per line)")
	help := flag.Bool("h", false, "Show help")

	flag.Usage = func() {
		version := internal.Version
		if info, ok := debug.ReadBuildInfo(); ok {
			for _, s := range info.Settings {
				if s.Key == "vcs.revision" && len(s.Value) >= 7 {
					version = s.Value[:7]
				}
				if s.Key == "vcs.tag" && s.Value != "" {
					version = s.Value
				}
			}
		}
		fmt.Fprintf(os.Stderr, "toshiki-avatar v%s - A modern, simpler yet better drop-in self-hosted lightweight gravatar replacemenet avatar api server in single file.\n\n", version)
		// OSC 8 hyperlink for clickable author link in supporting terminals
		fmt.Fprintf(os.Stderr, "Author: \033]8;;https://toshiki.dev\033\\Anda Toshiki\033]8;;\033\\ <andatoshiki@proton.me>\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.VisitAll(func(f *flag.Flag) {
			fmt.Fprintf(os.Stderr, "  -%s %s\n      %s (default: %s)\n", f.Name, f.DefValue, f.Usage, f.DefValue)
		})
	}

	flag.Parse()

	// Support --help as well, but don't show it in the flag list
	for _, arg := range os.Args[1:] {
		if arg == "--help" {
			flag.Usage()
			os.Exit(0)
		}
	}
	if *help {
		flag.Usage()
		os.Exit(0)
	}

	return &Config{
		Port:     *port,
		ImgType:  *imgType,
		Dir:      *dir,
		ListFile: *listFile,
	}
}
