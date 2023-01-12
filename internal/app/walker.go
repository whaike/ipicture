package main

import (
	"io/fs"
	"ipicture/internal/handler"
	"path/filepath"
	"strings"
)

type WK struct {
	hand     *handler.Handler
	rootPath string
}

func NewWK(hand *handler.Handler, root string) *WK {
	return &WK{
		hand:     hand,
		rootPath: root,
	}
}

func (w *WK) Start() {
	go w.hand.FileCheck()
	filepath.Walk(w.rootPath, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() && !strings.HasPrefix(info.Name(), ".") {
			//fmt.Println(path)
			w.hand.FileCh <- &handler.File{
				MetaInfo: &handler.MetaInfo{
					Name: info.Name(),
					Path: path,
				},
			}
		}
		return nil
	})

	// order pictures by created time and move them to other place
	// 1、首先对所有图片/视频排序并取得时间范围
	// 2、在时间范围内按年+月新建文件夹
	// 3、将所有图片移入
}
