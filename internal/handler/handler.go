package handler

import (
	"fmt"
	"ipicture/internal/picture"
	"strings"
)

type File struct {
	Name   string
	Path   string
	Type   string
	Suffix string
}

type Handler struct {
	FileCh chan *File
	ImgCh  chan *picture.Picture
}

func NewHandler(fileCh chan *File, imgCh chan *picture.Picture) *Handler {
	return &Handler{
		FileCh: fileCh,
		ImgCh:  imgCh,
	}
}

func (h *Handler) FileCheck() {
	for {
		select {
		case c := <-h.FileCh:
			c.Typed()
			switch c.Type {
			case "image":
				h.ImgCh <- &picture.Picture{
					Name:   c.Name,
					Path:   c.Path,
					Type:   c.Type,
					Suffix: c.Suffix,
				}
			default:
				fmt.Println("无法操作的数据类型", c)
			}

		}
	}
}

func (f *File) Typed() {
	suffix := ""
	if !strings.HasPrefix(f.Name, ".") {
		sp := strings.Split(f.Name, ".")
		if len(sp) > 0 {
			suffix = sp[len(sp)-1]
		}
	}
	suffix = strings.ToLower(suffix)
	f.Suffix = suffix
	switch suffix {
	case "jpg", "png", "jpeg":
		f.Type = "image"
	}
}
