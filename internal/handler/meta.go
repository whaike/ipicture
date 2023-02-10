package handler

import (
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

func NewExifGo() *ExifGo {
	e := &ExifGo{
		Options: &ExifOptions{},
	}
	return e
}

func (e *ExifGo) MetaInfos(filenames ...string) []map[string]string {
	et, err := exiftool.NewExiftool(exiftool.CoordFormant("%+f"))
	if err != nil {
		g.Logs.Errorf("Error when intializing: %v", err)
		return nil
	}
	defer et.Close()

	start := time.Now()
	fileInfos := et.ExtractMetadata(filenames...)
	end := time.Since(start).Milliseconds()

	if len(fileInfos) == 0 {
		return nil
	}
	result := make([]map[string]string, len(filenames))

	for _, fileInfo := range fileInfos {
		rs := make(map[string]string)
		s := reflect.ValueOf(*e.Options)
		for i := 0; i < s.NumField(); i++ {
			v := s.Field(i).Interface().(bool)
			if v {
				key := s.Type().Field(i).Name
				res, err := fileInfo.GetString(key)
				if err != nil {
					continue
				}
				rs[key] = res
			}
		}
	}

	others := time.Since(start).Milliseconds()
	g.Logs.Debugf("MetaInfo %d files: exiftool.ExtractMetadata cost: %d ms, others: %d", len(filenames), end, others)
	return result
}
