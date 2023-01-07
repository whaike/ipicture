package picture

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
	"github.com/rwcarlsen/goexif/tiff"
	"io"
	"ipicture/internal/model"
	"os"
	"strings"
	"time"
)

type Picture struct {
	Name      string
	Path      string
	Type      string
	Suffix    string
	Md5       string
	Lng       string
	Lat       string
	PhotoTime string
}

type Pictures struct {
	db    *model.Picture
	imgCh chan *Picture
}

func NewPictures(imgCh chan *Picture, db *model.Picture) *Pictures {
	return &Pictures{
		imgCh: imgCh,
		db:    db,
	}
}

func (p *Pictures) Info() {
	for {
		select {
		case i := <-p.imgCh:
			i.md5()
			i.exifInfos()
			//if i.Lng != "" && i.Lat != "" {
			//	fmt.Println(i.Name, i.Lng, i.Lat)
			//}
			pm := &model.PictureModel{
				Name:    i.Name,
				Path:    i.Path,
				Md5:     i.Md5,
				Type:    i.Type,
				Suffix:  i.Suffix,
				ShootAt: i.PhotoTime,
				Lng:     i.Lng,
				Lat:     i.Lat,
			}
			old, err := p.db.Query(pm)
			if err != nil {
				fmt.Println("查询失败", pm.Name, pm.Md5, err.Error())
				continue
			}
			if old != nil && old.Md5 == pm.Md5 {
				if old.Path == pm.Path && old.Name == pm.Name {
					continue
				} else {
					i.deleteSelf()
				}
			} else {
				p.db.Insert(pm)
				fmt.Println("insert ", pm.Path)
			}
		}
	}
}

func (pi *Picture) deleteSelf() {
	fmt.Println(fmt.Sprintf("delete %s \n", pi.Path))
}

func (pi *Picture) md5() {
	f, err := os.Open(pi.Path)
	if err != nil {
		fmt.Println("打开文件失败", pi.Path)
		return
	}
	hash := md5.New()
	_, _ = io.Copy(hash, f)
	pi.Md5 = hex.EncodeToString(hash.Sum(nil))
}

func (pi *Picture) exifInfos() {
	f, err := os.Open(pi.Path)
	if err != nil {
		fmt.Println("打开文件失败", pi.Path)
		return
	}
	exif.RegisterParsers(mknote.All...)
	x, err := exif.Decode(f)
	if err != nil {
		fmt.Println("未获取到meta信息", pi.Path)
		return
	}
	//camModel, err := x.Get(exif.Model) // normally, don't ignore errors!
	//if err != nil {
	//	fmt.Println(err.Error())
	//}
	//fmt.Println(camModel.StringVal())
	//
	//focal, _ := x.Get(exif.FocalLength)
	//numer, denom, _ := focal.Rat2(0) // retrieve first (only) rat. value
	//fmt.Printf("%v/%v", numer, denom)
	//
	//// Two convenience functions exist for date/time taken and GPS coords:
	tm, err := x.DateTime()
	if err != nil {
		tm2, err := DateTimeDigitized(x)
		if err == nil {
			pi.PhotoTime = tm2.String()
		} else {
			// 通过名字猜创建日期
			pi.PhotoTime = GuessTimeFromName(pi.Name)
		}
	} else {
		pi.PhotoTime = tm.String()
	}

	//fmt.Println("Taken: ", pi.PhotoTime)

	lat, long, _ := x.LatLong()
	if long == lat && lat == 0 {
		return
	}
	pi.Lng = fmt.Sprintf("%f", long)
	pi.Lat = fmt.Sprintf("%f", lat)

}

func DateTimeDigitized(x *exif.Exif) (time.Time, error) {
	var dt time.Time
	tag, err := x.Get(exif.DateTimeDigitized)
	if err != nil {
		return dt, err
	}
	if tag.Format() != tiff.StringVal {
		return dt, errors.New("DateTime[DateTimeDigitized] not in string format")
	}
	exifTimeLayout := "2006:01:02 15:04:05"
	dateStr := strings.TrimRight(string(tag.Val), "\x00")
	timeZone := time.Local
	if tz, _ := x.TimeZone(); tz != nil {
		timeZone = tz
	}
	return time.ParseInLocation(exifTimeLayout, dateStr, timeZone)
}

func GuessTimeFromName(name string) string {
	ns := strings.Split(name, "_")
	if len(ns) < 6 {
		return ""
	}
	ts := fmt.Sprintf("%s:%s:%s %s:%s", ns[0], ns[1], ns[2], ns[3], ns[4])
	t, err := time.Parse("2006:01:02 15:04", ts)
	if err != nil {
		return ""
	}
	return t.Format("2006:01:02 15:04")
}
