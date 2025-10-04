package main


import (
	"bufio"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/andatoshiki/toshiki-avatar/internal/avatar"
	"github.com/andatoshiki/toshiki-avatar/internal/config"
	"github.com/andatoshiki/toshiki-avatar/internal/server"
	utils "github.com/andatoshiki/toshiki-avatar/internal/utils"
	apperr "github.com/andatoshiki/toshiki-avatar/internal/errors"
)



func main() {
	rand.Seed(time.Now().UnixNano())
	conf := config.ParseFlags()
	var images []string

	// CLI flag conflict: both -d and -l provided or both empty
	if (conf.Dir == "" && conf.ListFile == "") || (conf.Dir != "" && conf.ListFile != "") {
		log.Fatal(apperr.ErrFlagConflict)
	}

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
			       if conf.Dir != "" && !utils.IsPathWithinBase(conf.Dir, line) {
				       log.Printf("skipping path outside base dir: %s", line)
				       continue
			       }
			       if !utils.IsReadable(line) {
				       log.Printf("skipping unreadable file: %s", line)
				       continue
			       }
			       if !utils.SupportedImageExt(line) {
				       log.Printf("skipping unsupported file type: %s", line)
				       continue
			       }
			       images = append(images, utils.NormalizePath(line))
		       }
	       }
	       if err := scanner.Err(); err != nil {
		       log.Fatalf("error reading list file: %v", err)
	       }
       } else {
	       if conf.Dir == "" {
		       log.Fatal("-d (directory) cannot be empty. Please specify a valid directory or use -l for a list file.")
	       }
	       err := filepath.Walk(conf.Dir, func(path string, info os.FileInfo, err error) error {
		       if err == nil && !info.IsDir() {
			       if !utils.IsReadable(path) {
				       log.Printf("skipping unreadable file: %s", path)
				       return nil
			       }
			       if !utils.SupportedImageExt(path) {
				       log.Printf("skipping unsupported file type: %s", path)
				       return nil
			       }
			       images = append(images, utils.NormalizePath(path))
		       }
		       return nil
	       })
	       if err != nil {
		       log.Fatal(err)
	       }
       }
	if len(images) == 0 {
		log.Fatal(apperr.ErrNoAvatars)
	}
	avatarService := avatar.NewAvatarService(images)
	srv := server.NewServer(avatarService, conf.ImgType)
       if err := srv.Start(conf.Port); err != nil {
	       if os.IsExist(err) || (err != nil && err.Error() == "listen tcp :"+strconv.Itoa(conf.Port)+": bind: address already in use") {
		       log.Fatalf("port %d already in use", conf.Port)
	       }
	       log.Fatal(err)
       }
}

