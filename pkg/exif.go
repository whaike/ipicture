package pkg

import (
	"fmt"
	"github.com/barasher/go-exiftool"
	"ipicture/g"
	"reflect"
	"time"
)

type ExifGo struct {
	Options *ExifOptions
}

type ExifOptions struct {
	// Primary EXIF fields
	ImageWidth                 bool
	ImageLength                bool
	BitsPerSample              bool
	Compression                bool
	PhotometricInterpretation  bool
	Orientation                bool
	SamplesPerPixel            bool
	PlanarConfiguration        bool
	YCbCrSubSampling           bool
	YCbCrPositioning           bool
	XResolution                bool
	YResolution                bool
	ResolutionUnit             bool
	DateTime                   bool
	ImageDescription           bool
	Make                       bool
	Model                      bool
	Software                   bool
	Artist                     bool
	Copyright                  bool
	ExifIFDPointer             bool
	GPSInfoIFDPointer          bool
	InteroperabilityIFDPointer bool
	ExifVersion                bool
	FlashpixVersion            bool
	ColorSpace                 bool
	ComponentsConfiguration    bool
	CompressedBitsPerPixel     bool
	PixelXDimension            bool
	PixelYDimension            bool
	MakerNote                  bool
	UserComment                bool
	RelatedSoundFile           bool
	DateTimeOriginal           bool
	DateTimeDigitized          bool
	SubSecTime                 bool
	SubSecTimeOriginal         bool
	SubSecTimeDigitized        bool
	ImageUniqueID              bool
	ExposureTime               bool
	FNumber                    bool
	ExposureProgram            bool
	SpectralSensitivity        bool
	ISOSpeedRatings            bool
	OECF                       bool
	ShutterSpeedValue          bool
	ApertureValue              bool
	BrightnessValue            bool
	ExposureBiasValue          bool
	MaxApertureValue           bool
	SubjectDistance            bool
	MeteringMode               bool
	LightSource                bool
	Flash                      bool
	FocalLength                bool
	SubjectArea                bool
	FlashEnergy                bool
	SpatialFrequencyResponse   bool
	FocalPlaneXResolution      bool
	FocalPlaneYResolution      bool
	FocalPlaneResolutionUnit   bool
	SubjectLocation            bool
	ExposureIndex              bool
	SensingMethod              bool
	FileSource                 bool
	SceneType                  bool
	CFAPattern                 bool
	CustomRendered             bool
	ExposureMode               bool
	WhiteBalance               bool
	DigitalZoomRatio           bool
	FocalLengthIn35mmFilm      bool
	SceneCaptureType           bool
	GainControl                bool
	Contrast                   bool
	Saturation                 bool
	Sharpness                  bool
	DeviceSettingDescription   bool
	SubjectDistanceRange       bool
	LensMake                   bool
	LensModel                  bool

	// Windows-specific tags
	XPTitle    bool
	XPComment  bool
	XPAuthor   bool
	XPKeywords bool
	XPSubject  bool

	// thumbnail fields
	ThumbJPEGInterchangeFormat       bool
	ThumbJPEGInterchangeFormatLength bool

	// GPS fields
	GPSVersionID        bool
	GPSLatitudeRef      bool
	GPSLatitude         bool
	GPSLongitudeRef     bool
	GPSLongitude        bool
	GPSAltitudeRef      bool
	GPSAltitude         bool
	GPSTimeStamp        bool
	GPSSatelites        bool
	GPSStatus           bool
	GPSMeasureMode      bool
	GPSDOP              bool
	GPSSpeedRef         bool
	GPSSpeed            bool
	GPSTrackRef         bool
	GPSTrack            bool
	GPSImgDirectionRef  bool
	GPSImgDirection     bool
	GPSMapDatum         bool
	GPSDestLatitudeRef  bool
	GPSDestLatitude     bool
	GPSDestLongitudeRef bool
	GPSDestLongitude    bool
	GPSDestBearingRef   bool
	GPSDestBearing      bool
	GPSDestDistanceRef  bool
	GPSDestDistance     bool
	GPSProcessingMethod bool
	GPSAreaInformation  bool
	GPSDateStamp        bool
	GPSDifferential     bool

	// interoperability fields
	InteroperabilityIndex bool
}

type Option func(exifgo *ExifGo)

func WithFilters(filter *ExifOptions) Option {
	return func(e *ExifGo) {
		// if filter is nil, will init all options true
		if filter == nil {
			o := &ExifOptions{}
			s := reflect.ValueOf(o).Elem()
			n := s.NumField()
			for i := 0; i < n; i++ {
				s.Field(i).SetBool(true)
			}
			e.Options = o
		} else {
			e.Options = filter
		}
	}
}

func NewExifGo(options ...Option) *ExifGo {
	e := &ExifGo{}
	if len(options) == 0 {
		options = append(options, WithFilters(nil))
	}
	for _, opt := range options {
		opt(e)
	}

	return e
}

func (e *ExifGo) MetaInfo(filename string) map[string]string {
	et, err := exiftool.NewExiftool(exiftool.CoordFormant("%+f"))
	if err != nil {
		fmt.Printf("Error when intializing: %v\n", err)
		return nil
	}
	defer et.Close()

	start := time.Now()
	fileInfos := et.ExtractMetadata(filename)
	end := time.Since(start).Milliseconds()

	if len(fileInfos) == 0 {
		return nil
	}
	result := make(map[string]string)

	s := reflect.ValueOf(*e.Options)
	for i := 0; i < s.NumField(); i++ {
		v := s.Field(i).Interface().(bool)
		if v {
			key := s.Type().Field(i).Name
			res, err := fileInfos[0].GetString(key)
			if err != nil {
				continue
			}
			result[key] = res
		}
	}
	others := time.Since(start).Milliseconds()
	g.Logs.Debugf("MetaInfo: exiftool.ExtractMetadata cost: %d ms, others: %d", end, others)
	return result
}
