package handler

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"ipicture/g"
	"ipicture/internal/model"
	"ipicture/pkg"
	"os"
	"strings"
	"time"
)

type (
	Handler struct {
		FileCh chan *File
		db     *model.IAV
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

func NewHandler(fileCh chan *File, db *model.IAV) *Handler {
	return &Handler{
		FileCh: fileCh,
		db:     db,
	}
}

func (h *Handler) FileCheck() {
	for {
		select {
		case c := <-h.FileCh:
			start := time.Now()
			err := c.TypeCheck()
			if err != nil {
				continue
			}
			a1 := time.Since(start).Milliseconds()

			err = c.md5()
			if err != nil {
				continue
			}
			a2 := time.Since(start).Milliseconds()

			err = c.metaInfo()
			if err != nil {
				continue
			}
			a3 := time.Since(start).Milliseconds()
			//err = h.hooks(c)
			//if err != nil {
			//	continue
			//}
			h.UpInsert(c)
			a4 := time.Since(start).Milliseconds()
			g.Logs.Debugf("data life: prepare[%dms],md5[%dms],exiftool[%dms],db[%dms]", a1, a2, a3, a4)
		}
	}
}

func (h *Handler) UpInsert(fi *File) {
	defer func() {
		if err := recover(); err != nil {
			g.Logs.Errorf("[UpInsert] 捕获异常: %s", err)
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
		g.Logs.Errorf("查询失败: %s, %s, %s", pm.Name, pm.Md5, err.Error())
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
		g.Logs.Infof("insert %s", pm.Path)
	}
}

func (fi *File) TypeCheck() error {
	var err error
	suffix := ""
	if !strings.HasPrefix(fi.Name, ".") {
		sp := strings.Split(fi.Name, ".")
		if len(sp) > 0 {
			suffix = sp[len(sp)-1]
		} else {
			err = fmt.Errorf("没有后缀")
			//fmt.Println(fi.Path, err.Error())
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
		//fmt.Println(fi.Path, err.Error())
		return err
	}
	return nil
}

// -----------------------------MetaInfo----------------------------------------------

func (mi *MetaInfo) md5() error {
	f, err := os.Open(mi.Path)
	if err != nil {
		g.Logs.Errorf("打开文件失败: %s", mi.Path)
		return err
	}
	hash := md5.New()
	_, err = io.Copy(hash, f)
	if err != nil {
		g.Logs.Errorf("拷贝文件,准备计算MD5失败: %s", mi.Path)
		return err
	}
	mi.Md5 = hex.EncodeToString(hash.Sum(nil))
	return nil
}

func (mi *MetaInfo) metaInfo() error {
	defer func() {
		if err := recover(); err != nil {
			g.Logs.Errorf("[metaInfo] 捕获异常: %s", err)
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
			g.Logs.Errorf("parse time error: %s, %s", mi.Path, shootAt)
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

	g.Logs.Debugf("delete %s,%s \n", mi.Md5, mi.Path)
}
