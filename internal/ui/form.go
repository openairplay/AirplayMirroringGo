/**
Created by: Joseph Han <joseph.bing.han@gmail.com>.
Created at: 2020-12-30 13:16
*/

package ui

import (
	"AirplayMirroringGo/internal/airplay"
	"AirplayMirroringGo/internal/screen"
	"bytes"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"image/jpeg"
	"log"
	"net"
	"time"
)

const FPS = 24

type MainForm struct {
	mainApp fyne.App
	window  fyne.Window

	btnOK     *widget.Button
	selScreen *widget.Select
	labChoose *widget.Label

	videoProvider  *screen.VideoProvider
	screens        []screen.Screen
	discovery      *airplay.Discovery
	startMirroring bool
	appleTV        net.IP
	client         *airplay.AppleTVClient
}

func NewMainForm() MainForm {
	app := app.New()
	window := app.NewWindow("AppleTV Mirroring")
	window.SetIcon(resourceMirroringPng)
	window.Resize(fyne.NewSize(280, 140))
	window.SetFixedSize(true)
	window.CenterOnScreen()
	main := MainForm{mainApp: app, window: window}
	return main
}

func (f *MainForm) Run() {
	log.Print("Main form starting")
	f.startMirroring = false

	v := screen.NewVideoProvider()
	f.videoProvider = v
	f.screens = v.GetScreens()

	var options []string
	for _, s := range f.screens {
		item := fmt.Sprintf("Screen %d : %s", s.Index+1, s.Bounds.String())
		options = append(options, item)
	}

	f.btnOK = widget.NewButton("Start Mirroring", f.onButtonClick)

	f.selScreen = widget.NewSelect(options, func(s string) {
	})
	f.selScreen.SetSelectedIndex(0) // default selected first

	f.labChoose = widget.NewLabel("Choose the screen to mirror")

	f.window.SetContent(container.NewVBox(
		f.labChoose,
		f.selScreen,
		layout.NewSpacer(),
		container.NewCenter(
			f.btnOK,
		),
		layout.NewSpacer(),
	))

	f.discovery = airplay.NewDiscovery()
	go f.searchAppleTV()

	c := airplay.NewAppleTVClient(f.appleTV)
	f.client = c

	f.window.ShowAndRun()
}

func (f *MainForm) onButtonClick() {
	log.Println("Click Button")
	if f.startMirroring {
		f.startMirroring = false
		f.btnOK.SetText("Start Mirroring")
		f.videoProvider.Stop()
		f.client.Stop()
		f.selScreen.Enable()
	} else {
		if f.appleTV == nil {
			go f.searchAppleTV()
			return
		}

		f.startMirroring = true
		index := f.selScreen.SelectedIndex()
		sc := f.screens[index]
		f.videoProvider.ChooseScreen(sc, FPS)
		f.videoProvider.Start()

		f.client.Start()
		go func() {
			frames := f.videoProvider.GetFrames()
			opt := &jpeg.Options{Quality: 100}
			for f.startMirroring {
				select {
				case frame := <-frames:
					buf := &bytes.Buffer{}
					err := jpeg.Encode(buf, frame, opt)
					if err != nil {
						log.Fatalf("JPEG encodeing fail %v", err)
						return
					}
					f.client.Stream <- buf.Bytes()
					buf.Reset()
				}
			}
		}()

		f.btnOK.SetText("Stop Mirroring")
		f.selScreen.Disable()
	}

}

func (f *MainForm) searchAppleTV() {
	defer func() {
		f.labChoose.SetText("Start Mirroring")
		f.btnOK.Enable()
	}()
	go func() {
		info := ""
		for f.appleTV == nil {
			info += "."

			f.btnOK.Disable()
			f.labChoose.SetText("Searching Apple TV" + info)

			if info == "........" {
				info = ""
			}
			time.Sleep(800 * time.Millisecond)
		}
	}()
	f.appleTV = f.discovery.GetAirPlayService().AddrIPv4
	if f.appleTV == nil {
		dialog.ShowInformation("Error", "Cannot find AppleTV, try again please!", f.window)
		f.appleTV = f.discovery.GetAirPlayService().AddrIPv4
		return
	}
	log.Printf("Found AppleTV at %v", f.appleTV)
}
