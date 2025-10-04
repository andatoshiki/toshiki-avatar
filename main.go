package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"github.com/andatoshiki/toshiki-avatar/internal/avatar"
	"github.com/andatoshiki/toshiki-avatar/internal/config"
	"github.com/andatoshiki/toshiki-avatar/internal/server"
)



func main() {
	conf := config.ParseFlags()
	var images []string
	if conf.ListFile != "" {
		file, err := os.Open(conf.ListFile)
		if err != nil {
			log.Fatalf("failed to open list file: %v", err)
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if line != "" {
				images = append(images, line)
			}
		}
		if err := scanner.Err(); err != nil {
			log.Fatalf("error reading list file: %v", err)
		}
	} else {
		err := filepath.Walk(conf.Dir, func(path string, info os.FileInfo, err error) error {
			if err == nil && !info.IsDir() {
				images = append(images, path)
			}
			return nil
		})
		if err != nil {
			log.Fatal(err)
		}
	}
	if len(images) == 0 {
		log.Fatalf("no avatars found in %s or list file", conf.Dir)
	}
	avatarService := avatar.NewAvatarService(images)
	srv := server.NewServer(avatarService, conf.ImgType)
	if err := srv.Start(conf.Port); err != nil {
		log.Fatal(err)
	}
}

