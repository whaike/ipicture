package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/pyroscope-io/pyroscope/pkg/agent/profiler"
	"io/fs"
	"ipicture/g"
	"ipicture/internal/etc"
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
	g.Logs.Info("start walking")
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

	g.Logs.Info("stop walking")

	// order pictures by created time and move them to other place
	// 1、首先对所有图片/视频排序并取得时间范围
	// 2、在时间范围内按年+月新建文件夹
	// 3、将所有图片移入
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	path := flag.String("path", ".", "your file path")
	conFile := flag.String("config", "etc/config.yaml", "config file. default etc/config.yaml")
	level := flag.String("level", "info", "log level, default info")
	flag.Parse()

	c := etc.LoadConfig(*conFile)
	if *path != "." {
		c.Path = *path
	}
	if *level != "info" {
		c.ZapLog.Level = *level
	}
	g.InitLog(&c.ZapLog)

	if c.PyroscopeEnable {
		if c.PyroscopeAddr == "" {
			panic(errors.New("PyroscopeAddr 为空"))
		}
		profiler.Start(profiler.Config{
			ApplicationName: "ipicture.golang.app",
			ServerAddress:   c.PyroscopeAddr,
		})
	}

	iavModel := model.NewIAVSModel("./ipictures.db")
	fileCh := make(chan *handler.File)
	hand := handler.NewHandler(fileCh, iavModel)
	wk := newWK(hand, *path)
	wk.Start()
}
