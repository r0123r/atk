// Copyright 2018 visualfc. All rights reserved.

package tk

import (
	"errors"
	"fmt"
	"image"
	_ "image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/visualfc/atk/tk/interp"
)

type Image struct {
	id    string
	photo *interp.Photo
}

func (i *Image) Id() string {
	return i.id
}

type ImageOpt struct {
	key   string
	value interface{}
}

func ImageOptId(id string) *ImageOpt {
	return &ImageOpt{"id", id}
}

func ImageOptGamma(gamma float64) *ImageOpt {
	return &ImageOpt{"gamma", gamma}
}

func LoadImage(file string, options ...*ImageOpt) (*Image, error) {
	if file == "" {
		return nil, os.ErrInvalid
	}
	var fileImage image.Image
	if filepath.Ext(file) == ".gif" {
		options = append(options, &ImageOpt{"file", file})
	} else {
		file, err := os.Open(file)
		if err != nil {
			return nil, err
		}
		im, _, err := image.Decode(file)
		file.Close()
		if err != nil {
			return nil, err
		}
		fileImage = im
	}
	im := NewImage(options...)
	if im == nil {
		return nil, errors.New("NewImage failed")
	}
	if fileImage != nil {
		im.SetImage(fileImage)
	}
	return im, nil
}

func NewImage(options ...*ImageOpt) *Image {
	var iid string
	var optList []string
	for _, opt := range options {
		if opt == nil {
			continue
		}
		if opt.key == "id" {
			if v, ok := opt.value.(string); ok {
				iid = v
			}
			continue
		}
		optList = append(optList, fmt.Sprintf("-%v {%v}", opt.key, opt.value))
	}
	if iid == "" {
		iid = MakeImageId()
	}
	script := fmt.Sprintf("image create photo %v", iid)
	if len(optList) > 0 {
		script += " " + strings.Join(optList, " ")
	}
	err := eval(script)
	if err != nil {
		return nil
	}
	photo := interp.FindPhoto(mainInterp, iid)
	if photo == nil {
		return nil
	}
	return &Image{iid, photo}
}

func (i *Image) IsValid() bool {
	return i.id != "" && i.photo != nil
}

func (i *Image) SetImage(img image.Image) *Image {
	err := i.photo.PutImage(img)
	if err != nil {
		dumpError(err)
	}
	return i
}

func (i *Image) SetZoomedImage(img image.Image, zoomX, zoomY, subsampleX, subsampleY int) *Image {
	err := i.photo.PutZoomedImage(img, zoomX, zoomY, subsampleX, subsampleY)
	if err != nil {
		dumpError(err)
	}
	return i
}

func (i *Image) ToImage() image.Image {
	return i.photo.ToImage()
}

func (i *Image) Blank() *Image {
	i.photo.Blank()
	return i
}

func (i *Image) SizeN() (width int, height int) {
	return i.photo.Size()
}

func (i *Image) Size() Size {
	w, h := i.SizeN()
	return Size{w, h}
}

func (i *Image) SetSizeN(width int, height int) *Image {
	err := i.photo.SetSize(width, height)
	if err != nil {
		dumpError(err)
	}
	return i
}

func (i *Image) SetSize(sz Size) *Image {
	return i.SetSizeN(sz.Width, sz.Height)
}

func (i *Image) Gamma() float64 {
	v, _ := evalAsFloat64(fmt.Sprintf("%v cget -gamma", i.id))
	return v
}

func (i *Image) SetGamma(v float64) *Image {
	eval(fmt.Sprintf("%v configure -gamma {%v}", i.id, v))
	return i
}

func parserImageResult(id string, err error) *Image {
	if err != nil {
		return nil
	}
	photo := interp.FindPhoto(mainInterp, id)
	if photo == nil {
		return nil
	}
	return &Image{id, photo}
}
