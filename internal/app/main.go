package main

import (
	"ipicture/internal/handler"
	"ipicture/internal/hooks"
	"ipicture/internal/model"
)

func main() {
	path := "/Users/whaike/Documents/我的.txt/backup"
	iavModel := model.NewIAVSModel("./ipictures.db")
	fileCh := make(chan *handler.File)
	hookList := hooks.NewHookList()
	hand := handler.NewHandler(fileCh, iavModel, hookList)
	wk := NewWK(hand, path)
	wk.Start()
}
