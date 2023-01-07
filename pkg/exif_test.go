package pkg

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestExifGo_MetaInfo(t *testing.T) {
	filename := "/Users/whaike/Documents/我的.txt/backup/20220110/屏幕快照/2021_05_19_17_37_IMG_1950.jpg"

	E := NewExifGo()
	//E := NewExifGo(WithFilters(&ExifOptions{
	//	DateTimeDigitized: true,
	//	DateTimeOriginal:  true,
	//	GPSLongitude:      true,
	//	GPSLatitude:       true,
	//	UserComment:       true,
	//}))
	res := E.MetaInfo(filename)
	if res != nil {
		s, _ := json.MarshalIndent(res, " ", " ")
		fmt.Println(string(s))
	}
}
