/**
Created by: Joseph Han <joseph.bing.han@gmail.com>.
Created at: 2020-12-30 12:02
*/

package encoder

import (
	"bytes"
	"github.com/gen2brain/x264-go"
	"image"
	"log"
	"math"
)

type H264Encoder struct {
	buffer   *bytes.Buffer
	encoder  *x264.Encoder
	realSize image.Point
}

func NewH264Encoder(size image.Point, frameRate int) (*H264Encoder, error) {
	buffer := bytes.NewBuffer(make([]byte, 0))
	realSize := findBestSizeForH264Profile(size)
	log.Printf("Best size is %v", realSize)
	opts := x264.Options{
		Width:     realSize.X,
		Height:    realSize.Y,
		FrameRate: frameRate,
		Tune:      "zerolatency",
		Preset:    "ultrafast",
		Profile:   "baseline",
		LogLevel:  x264.LogError,
	}
	encoder, err := x264.NewEncoder(buffer, &opts)
	if err != nil {
		return nil, err
	}
	return &H264Encoder{
		buffer:   buffer,
		encoder:  encoder,
		realSize: realSize,
	}, nil
}

func findBestSizeForH264Profile(constraints image.Point) image.Point {
	sizes := []image.Point{
		{2048, 1024},
		{1920, 1080},
		{1280, 720},
		{720, 576},
		{720, 480},
	}
	minRatioDiff := math.MaxFloat64
	var minRatioSize image.Point
	for _, size := range sizes {
		if size == constraints {
			return size
		}
		lowerRes := size.X < constraints.X && size.Y < constraints.Y
		hRatio := float64(constraints.X) / float64(size.X)
		vRatio := float64(constraints.Y) / float64(size.Y)
		ratioDiff := math.Abs(hRatio - vRatio)
		if lowerRes && (ratioDiff) < 0.0001 {
			return size
		} else if ratioDiff < minRatioDiff {
			minRatioDiff = ratioDiff
			minRatioSize = size
		}
	}
	return minRatioSize

}

func (e *H264Encoder) Encode(frame *image.RGBA) ([]byte, error) {
	err := e.encoder.Encode(frame)
	if err != nil {
		return nil, err
	}
	err = e.encoder.Flush()
	if err != nil {
		return nil, err
	}
	payload := e.buffer.Bytes()
	e.buffer.Reset()
	return payload, nil
}
