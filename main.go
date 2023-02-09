package main

import (
	"flag"
	"github.com/pyroscope-io/pyroscope/pkg/agent/profiler"
	"io/fs"
	"ipicture/internal/handler"
	"ipicture/internal/model"
	"path/filepath"
	"strings"
)

type WK struct {
	hand     *handler.Handler
	rootPath string
}

func newWK(hand *handler.Handler, root string) *WK {
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

func main() {
	path := flag.String("path", ".", "your file path")
	c_pyroscope := flag.Bool("pyroscope_open", false, "pyroscope closed by default")
	pserver := flag.String("pyroscope_addr", "http://127.0.0.1:4040", "pyroscope server address")
	if *c_pyroscope {
		profiler.Start(profiler.Config{
			ApplicationName: "ipicture.golang.app",
			ServerAddress:   *pserver,
		})
	}
	flag.Parse()
	iavModel := model.NewIAVSModel("./ipictures.db")
	fileCh := make(chan *handler.File)
	hand := handler.NewHandler(fileCh, iavModel)
	wk := newWK(hand, *path)
	wk.Start()
}
