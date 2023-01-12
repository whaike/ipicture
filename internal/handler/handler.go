package handler

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"ipicture/internal/hooks"
	"ipicture/internal/model"
	"ipicture/pkg"
	"os"
	"strings"
	"time"
)

type (
	Handler struct {
		FileCh   chan *File
		db       *model.IAV
		hookList []hooks.IHook
	}
	File struct {
		*MetaInfo
	}
	MetaInfo struct {
		Name        string
		Path        string
		Type        string
		Suffix      string
		Md5         string
		Lng         string
		Lat         string
		CreatedTime string
	}
)

func NewHandler(fileCh chan *File, db *model.IAV, hks []hooks.IHook) *Handler {
	return &Handler{
		FileCh:   fileCh,
		db:       db,
		hookList: hks,
	}
}

func (h *Handler) FileCheck() {
	for {
		select {
		case c := <-h.FileCh:
			err := c.Typed()
			if err != nil {
				continue
			}

			err = c.md5()
			if err != nil {
				continue
			}
			err = c.metaInfo()
			if err != nil {
				continue
			}
			err = h.hooks(c)
			if err != nil {
				continue
			}
			h.UpInsert(c)
		}
	}
}

func (h *Handler) UpInsert(fi *File) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("[UpInsert] 捕获异常: ", err)
		}
	}()
	pm := &model.IAVModel{
		Name:    fi.Name,
		Path:    fi.Path,
		Md5:     fi.Md5,
		Type:    fi.Type,
		Suffix:  fi.Suffix,
		ShootAt: fi.CreatedTime,
		Lng:     fi.Lng,
		Lat:     fi.Lat,
	}
	old, err := h.db.Query(pm)
	if err != nil {
		fmt.Println("查询失败", pm.Name, pm.Md5, err.Error())
		return
	}
	if old != nil && old.Md5 == pm.Md5 {
		if old.Path == pm.Path && old.Name == pm.Name {
			return
		} else {
			fi.deleteSelf()
		}
	} else {
		h.db.Insert(pm)
		fmt.Println("insert ", pm.Path)
	}
}

func (h *Handler) hooks(fi *File) error {
	for _, hook := range h.hookList {
		err := hook.Hook(fi)
		if err != nil {
			return err
		}
	}
	return nil
}

func (fi *File) Typed() error {
	var err error
	suffix := ""
	if !strings.HasPrefix(fi.Name, ".") {
		sp := strings.Split(fi.Name, ".")
		if len(sp) > 0 {
			suffix = sp[len(sp)-1]
		} else {
			err = fmt.Errorf("没有后缀")
			fmt.Println(fi.Path, err.Error())
			return err
		}
	}
	suffix = strings.ToLower(suffix)
	fi.Suffix = suffix
	switch suffix {
	case "jpg", "png", "jpeg", "gif":
		fi.Type = "image"
	case "mov", "mp4", "mkv":
		fi.Type = "movie"
	default:
		err = fmt.Errorf("无法处理的后缀")
		fmt.Println(fi.Path, err.Error())
		return err
	}
	return nil
}

// -----------------------------MetaInfo----------------------------------------------

func (mi *MetaInfo) md5() error {
	f, err := os.Open(mi.Path)
	if err != nil {
		fmt.Println("打开文件失败", mi.Path)
		return err
	}
	hash := md5.New()
	_, err = io.Copy(hash, f)
	mi.Md5 = hex.EncodeToString(hash.Sum(nil))
	fmt.Println("拷贝文件,计算MD5失败", mi.Path)
	return err
}

func (mi *MetaInfo) metaInfo() error {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("[metaInfo] 捕获异常: ", err)
		}
	}()
	metaTool := pkg.NewExifGo()
	mt := metaTool.MetaInfo(mi.Path)
	err := mi.getCreateTime(mt)
	if err != nil {
		return err
	}
	err = mi.getGpsLocations(mt)
	return err
}

func (mi *MetaInfo) getCreateTime(mp map[string]string) error {
	shootAt := ""
	if a, ok := mp["DateTimeDigitized"]; ok {
		shootAt = a
	} else if b, ok := mp["DateTimeOriginal"]; ok {
		shootAt = b
	} else if c, ok := mp["DateTime"]; ok {
		shootAt = c
	} else {
		shootAt = mi.guessTimeFromName(mi.Name)
	}
	if shootAt != "" && len(shootAt) == 16 {
		tm, err := time.Parse("2006:01:02 15:04", shootAt)
		if err != nil {
			fmt.Println("parse time error: ", mi.Path, shootAt)
			return err
		} else {
			shootAt = tm.Format("2006:01:02 15:04:05")
		}
	}
	mi.CreatedTime = shootAt
	return nil
}

func (mi *MetaInfo) getGpsLocations(mp map[string]string) error {
	if a, ok := mp["GPSLatitude"]; ok {
		mi.Lat = a
	}
	if b, ok := mp["GPSLongitude"]; ok {
		mi.Lng = b
	}
	return nil
}

func (mi *MetaInfo) guessTimeFromName(name string) string {
	ns := strings.Split(name, "_")
	if len(ns) < 6 {
		return ""
	}
	ts := fmt.Sprintf("%s:%s:%s %s:%s", ns[0], ns[1], ns[2], ns[3], ns[4])
	t, err := time.Parse("2006:01:02 15:04", ts)
	if err != nil {
		return ""
	}
	return t.Format("2006:01:02 15:04:05")
}

func (mi *MetaInfo) deleteSelf() {
	//err := os.Remove(mi.Path)
	//if err != nil {
	//	fmt.Println(fmt.Sprintf("delete %s error, %s", mi.Path, err.Error()))
	//} else {
	//	fmt.Println(fmt.Sprintf("delete %s,%s \n", mi.Md5, mi.Path))
	//}

	fmt.Println(fmt.Sprintf("delete %s,%s \n", mi.Md5, mi.Path))
}
