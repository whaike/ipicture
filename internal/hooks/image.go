package hooks

import (
	"ipicture/internal/handler"
)

type Image struct {
}

func (i *Image) Hook(fi *handler.File) error {
	//TODO implement me
	return nil
}

func NewImageHook() *Image {
	return &Image{}
}
