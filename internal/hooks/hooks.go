package hooks

import "ipicture/internal/handler"

type IHook interface {
	Hook(fi *handler.File) error
}

var hookList []IHook

func Register(hook IHook) {
	hookList = append(hookList, hook)
}

func NewHookList() []IHook {
	imgHook := NewImageHook()
	Register(imgHook)
	return hookList
}
