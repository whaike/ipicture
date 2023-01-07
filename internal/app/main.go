package main

import (
	"ipicture/internal/handler"
	"ipicture/internal/model"
	"ipicture/internal/picture"
)

func main() {
	path := "/Users/whaike/Documents/我的.txt/backup"
	picModel := model.NewPictureDB("./ipictures.db")
	fileCh := make(chan *handler.File)
	imgCh := make(chan *picture.Picture)
	hand := handler.NewHandler(fileCh, imgCh)
	wk := NewWK(hand, path, picModel)
	wk.Start()
}
