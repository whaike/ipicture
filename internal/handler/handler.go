package handler

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"ipicture/g"
	"ipicture/internal/model"
	"os"
	"strings"
	"time"
)

type (
	Handler struct {
		FileCh       chan *File
		metasCh      chan File
		db           *model.IAV
		tmp          int  // 缓存长度
		delDuplicate bool // 是否删除重复文件
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

func NewHandler(fileCh chan *File, db *model.IAV, delDuplicate bool) *Handler {
	return &Handler{
		FileCh:       fileCh,
		metasCh:      make(chan File),
		db:           db,
		tmp:          100,
		delDuplicate: delDuplicate,
	}
}

func (h *Handler) FileCheck() {
	for {
		select {
		case c := <-h.FileCh:
			err := c.TypeCheck()
			if err != nil {
				continue
			}
			err = c.md5()
			if err != nil {
				continue
			}
			h.metasCh <- *c
		}
	}
}

func (h *Handler) MetaAndSave() {
	fs := make([]File, 0)
	work := time.Now()
	for {
		select {
		case c := <-h.metasCh:
			if c.MetaInfo == nil {
				continue
			}
			fs = append(fs, c)
			if len(fs) >= h.tmp || time.Since(work).Seconds() > 1 {
				go h.Do(fs)
				fs = make([]File, 0)
				work = time.Now()
			}
		}
	}
}

func (h *Handler) Do(fsIn []File) {
	fileops := make([]File, 0)
	for _, in := range fsIn {
		if in.MetaInfo != nil {
			fileops = append(fileops, in)
		}
	}
	start := time.Now()
	defer func() {
		if err := recover(); err != nil {
			g.Logs.Errorf("Do 函数崩溃，当前文件数量: %d, err=%v", len(fileops), err)
		}
	}()
	fs := make([]string, 0)
	for _, v := range fileops {
		fs = append(fs, v.MetaInfo.Path)
	}
	et := NewExifGo()
	mp := et.MetaInfos(fs...)
	a1 := time.Since(start).Milliseconds()
	for i, f := range fileops {
		err := f.metaInfo2(mp[i]) // 补充元数据信息
		if err != nil {
			continue
		}
		h.UpInsert(&f)
	}
	a2 := time.Since(start).Milliseconds()
	g.Logs.Debugf("本轮基本信息查询结束,获取[%d 个]文件信息, 其中元数据耗时[%d ms], 其他数据及存储耗时%d ms]", len(fileops), a1, a2-a1)
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
			if h.delDuplicate {
				fi.deleteSelf()
			} else {
				g.Logs.Infof("duplicated file: %s, %s", fi.Md5, fi.Path)
			}
		}
	} else {
		h.db.Insert(pm)
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

// 补充其他元数据信息
func (mi *MetaInfo) metaInfo2(mt map[string]string) error {
	defer func() {
		if err := recover(); err != nil {
			g.Logs.Errorf("[metaInfo2] 捕获异常: %s", err)
		}
	}()
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
	err := os.Remove(mi.Path)
	if err != nil {
		g.Logs.Errorf("delete duplicate %s error, %s", mi.Path, err.Error())
	} else {
		g.Logs.Infof("delete duplicate %s, %s", mi.Md5, mi.Path)
	}
}
