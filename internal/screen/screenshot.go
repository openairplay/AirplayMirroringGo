/**
Created by: Joseph Han <joseph.bing.han@gmail.com>.
Created at: 2020-12-30 10:50
*/

package screen

import (
	"github.com/kbinani/screenshot"
	"go.uber.org/atomic"
	"image"
	"log"
	"time"
)

type Screen struct {
	Index  int
	Bounds image.Rectangle
	Min    image.Point
}

type VideoProvider struct {
	fps    int
	screen Screen
	frames chan *image.RGBA
	stop   atomic.Bool
}

func NewVideoProvider() *VideoProvider {
	v := VideoProvider{
		frames: make(chan *image.RGBA),
	}
	v.stop.Store(true)
	return &v
}

//获取所有可用屏幕
func (p *VideoProvider) GetScreens() []Screen {
	numScreens := screenshot.NumActiveDisplays()
	screens := make([]Screen, numScreens)
	var min image.Point
	for i := 0; i < numScreens; i++ {
		screens[i] = Screen{
			Index:  i,
			Bounds: screenshot.GetDisplayBounds(i),
		}
		if screens[i].Bounds.Min.X < min.X {
			min.X = screens[i].Bounds.Min.X
		}
		if screens[i].Bounds.Min.Y < min.Y {
			min.Y = screens[i].Bounds.Min.Y
		}
	}
	for i := 0; i < numScreens; i++ {
		screens[i].Min = min
	}

	return screens
}

//选择要投评的屏幕以及帧数
func (p *VideoProvider) ChooseScreen(screen Screen, fps int) {
	p.screen = screen
	p.fps = fps
}

//获取视频流的通道
func (p *VideoProvider) GetFrames() <-chan *image.RGBA {
	return p.frames
}

//开始录制视频
func (p *VideoProvider) Start() {
	p.stop.Store(false)
	dur := time.Duration(1000/p.fps) * time.Millisecond
	log.Printf("Duration %v", dur)
	go func() {
		for !p.stop.Load() {
			lastTime := time.Now()
			if img, err := screenshot.CaptureRect(p.screen.Bounds); err != nil {
				log.Fatalf("Screen capture fail: %v", err)
				return
			} else {
				p.frames <- img
				sub := time.Now().Sub(lastTime)
				sleepTime := dur - sub
				if sleepTime > 0 {
					log.Printf("Sleep %v", sleepTime)
					time.Sleep(sleepTime)
				}

			}
		}
	}()
}

//停止录制视频流
func (p VideoProvider) Stop() {
	p.stop.Store(true)
}
