package service

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"os"
)

type (
	File struct {
		Name     string `json:"name"`
		Checksum string `json:"checksum"`
	}
	Service struct {
		RootPath string
	}
)

func NewService(rootPath string) *Service {
	return &Service{RootPath: rootPath}
}

func FileMD5(path string) string {
	h := md5.New()
	f, err := os.Open(path)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	_, err = io.Copy(h, f)
	if err != nil {
		log.Println(err)
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}
