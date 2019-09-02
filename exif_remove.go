package exifremove

import (
	"bytes"
	"errors"
	"image"
	"fmt"

	"github.com/dsoprea/go-exif"
	"github.com/dsoprea/go-jpeg-image-structure"
	"github.com/dsoprea/go-png-image-structure"
)

const (
	JpegMediaType  = "jpeg"
	PngMediaType   = "png"
	OtherMediaType = "other"
	StartBytes     = 0
	EndBytes       = 0
)

type MediaContext struct {
	MediaType string
	RootIfd   *exif.Ifd
	RawExif   []byte
	Media     interface{}
}

func RemoveEXIF(data []byte) ([]byte, error) {
	jmp := jpegstructure.NewJpegMediaParser()
	pmp := pngstructure.NewPngMediaParser()
	mc := &MediaContext{
		MediaType: OtherMediaType,
		RootIfd:   nil,
		RawExif:   nil,
		Media:     nil,
	}
	filtered := []byte{}

	if jmp.LooksLikeFormat(data) {
		fmt.Printf("JPG\n")

		mc.MediaType = JpegMediaType
		sl, err := jmp.ParseBytes(data)
		if err != nil {
			return nil, err
		}
		mc.Media = sl

		if rootIfd, rawExif, err := sl.Exif(); err != nil {
			return nil, err
		} else {
			mc.RootIfd = rootIfd
			mc.RawExif = rawExif
		}

		if _, _, err := sl.FindExif(); err != nil {
			return nil, err
		} else {

			startExifBytes := StartBytes
			endExifBytes := EndBytes

			if bytes.Contains(data, mc.RawExif) {
				for i := 0; i < len(data)-len(mc.RawExif); i++ {
					if bytes.Compare(data[i:i+len(mc.RawExif)], mc.RawExif) == 0 {
						startExifBytes = i
						endExifBytes = i + len(mc.RawExif)
					}
				}
				fill := make([]byte, len(data[startExifBytes:endExifBytes]))
				copy(data[startExifBytes:endExifBytes], fill)
			}

			filtered = data

			_, _, err = image.Decode(bytes.NewReader(filtered))
			if err != nil {
				return nil, errors.New("EXIF removal corrupted " + err.Error())
			}

		}

	} else if pmp.LooksLikeFormat(data) {
		mc.MediaType = PngMediaType
		fmt.Printf("PNG\n")

		cs, err := pmp.ParseBytes(data)
		if err != nil {
			return nil, err
		}
		mc.Media = cs

		if rootIfd, rawExif, err := cs.Exif(); err != nil {
			return nil, err
		} else {
			mc.RootIfd = rootIfd
			mc.RawExif = rawExif
		}
		if _, err := cs.FindExif(); err != nil {
			return nil, err
		} else {

			startExifBytes := StartBytes
			endExifBytes := EndBytes
			print("wow\n")
			if bytes.Contains(data, mc.RawExif) {
				fmt.Printf("hi\n")
				for i := 0; i < len(data)-len(mc.RawExif); i++ {
					if bytes.Compare(data[i:i+len(mc.RawExif)], mc.RawExif) == 0 {
						startExifBytes = i
						endExifBytes = i + len(mc.RawExif)
					}
				}
				fill := make([]byte, len(data[startExifBytes:endExifBytes]))
				copy(data[startExifBytes:endExifBytes], fill)
			}

			filtered = data

			_, _, err = image.Decode(bytes.NewReader(filtered))
			if err != nil {
				return nil, errors.New("EXIF removal corrupted " + err.Error())
			}

		}
	}

	return filtered, nil
}
