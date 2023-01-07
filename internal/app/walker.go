package main

import (
	"io/fs"
	"ipicture/internal/handler"
	"ipicture/internal/model"
	"ipicture/internal/picture"
	"path/filepath"
	"strings"
)

type WK struct {
	hand     *handler.Handler
	rootPath string
	picDB    *model.Picture
}

func NewWK(hand *handler.Handler, root string, picdb *model.Picture) *WK {
	return &WK{
		hand:     hand,
		rootPath: root,
		picDB:    picdb,
	}
}

func (w *WK) Start() {
	pics := picture.NewPictures(w.hand.ImgCh, w.picDB)
	go pics.Info()
	go w.hand.FileCheck()
	filepath.Walk(w.rootPath, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() && !strings.HasPrefix(info.Name(), ".") {
			//fmt.Println(path)
			w.hand.FileCh <- &handler.File{
				Name: info.Name(),
				Path: path,
			}

			//picModel.Insert(&model.PictureModel{
			//	Name:    info.Name(),
			//	Path:    path,
			//	Type:    "",
			//	Suffix:  suffix,
			//	Tags:    "",
			//	ShootAt: "",
			//	Lng:     "",
			//	Lat:     "",
			//})
		}
		return nil
	})
}
