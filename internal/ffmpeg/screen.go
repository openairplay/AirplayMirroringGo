/**
Created by: Joseph Han <joseph.bing.han@gmail.com>.
Created at: 2021-01-02 10:15
*/

package ffmpeg

import (
	"AirplayMirroringGo/internal/screen"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
)

type Ffmpeg struct {
	chanExit chan bool
	cmd      *exec.Cmd
}

var _fmpeg *Ffmpeg

func init() {
	cmd, err := exec.Command("ffmpeg", "-version").Output()
	if err != nil {
		log.Fatalf("ffmpeg get version fail: %v", err)
	}
	log.Printf("Get ffmpeg version: %v", string(cmd))
	_fmpeg = &Ffmpeg{chanExit: make(chan bool)}
}

func FFmpeg() *Ffmpeg {
	return _fmpeg
}

func (f *Ffmpeg) Start(sc screen.Screen, fps int) {
	size := fmt.Sprintf("%dx%d", sc.Bounds.Dx(), sc.Bounds.Dy())
	log.Print(size)
	scr := fmt.Sprintf(":0.0+%d,%d", sc.Bounds.Min.X-sc.Min.X, sc.Bounds.Min.Y-sc.Min.Y)
	log.Printf("Min:%v,Bounds:%v", sc.Min, sc.Bounds)
	//rtmpAddr := fmt.Sprintf("rtmp://localhost:8090/live/%s", key)
	outFile := "tcp://localhost:8089/"
	os.Remove(outFile)
	go func() {
		switch runtime.GOOS {
		case "linux":
			f.cmd = exec.Command("ffmpeg", "-video_size", size, "-framerate", strconv.Itoa(fps),
				"-f", "x11grab", "-i", scr, "-c:v", "libx264", "-async", "1", "-vsync", "1",
				"-qp", "0", "-preset", "ultrafast", "-f", "flv", outFile)
			log.Print(f.cmd.String())
		case "windows":
			f.cmd = exec.Command("ffmpeg", "-video_size", size, "-framerate", strconv.Itoa(fps),
				"-offset_x", strconv.Itoa(sc.Bounds.Min.X), "-offset_y", strconv.Itoa(20), "-f", "gdigrab",
				"-i", "desktop", "-c:v", "libx264", "-qp", "0", "-preset", "ultrafast")
		case "darwin":
			f.cmd = exec.Command("ffmpeg", "-video_size", size, "-framerate", strconv.Itoa(fps), "-f",
				"avfoundation", "-i", "1:0", "-c:v", "libx264", "-qp", "0", "-preset",
				"ultrafast")
		}
		f.cmd.Stdout = os.Stdout
		f.cmd.Stderr = os.Stderr
		if err := f.cmd.Start(); err != nil {
			log.Fatalf("Start mirroring using ffmpeg fail %v", err)
		}

		f.cmd.Wait()

	}()
	log.Print("Started ffmpeg")

}

func (f *Ffmpeg) Stop() {
	f.cmd.Process.Kill()
}
