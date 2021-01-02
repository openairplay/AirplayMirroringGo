/**
Created by: Joseph Han <joseph.bing.han@gmail.com>.
Created at: 2020-12-30 13:16
*/

package ui

import (
	"AirplayMirroringGo/internal/airplay"
	"AirplayMirroringGo/internal/ffmpeg"
	_ "AirplayMirroringGo/internal/proxy"
	"AirplayMirroringGo/internal/screen"
	"fmt"
	_ "github.com/ying32/govcl/pkgs/winappres"
	"github.com/ying32/govcl/vcl"
	"github.com/ying32/govcl/vcl/types"
	"github.com/ying32/govcl/vcl/types/colors"
	"log"
	"net"
	"time"
)

const FPS = 15

type MainForm struct {
	*vcl.TForm
	LabChoose      *vcl.TLabel
	BtnOK          *vcl.TButton
	CombList       *vcl.TComboBox
	Tray           *vcl.TTrayIcon
	videoProvider  *screen.VideoProvider
	screens        []screen.Screen
	discovery      *airplay.Discovery
	startMirroring bool
	appleTV        net.IP
	client         *airplay.AppleTVClient
	ffmpeg         *ffmpeg.Ffmpeg
}

var (
	mainForm *MainForm
)

func AppFormRun() {
	vcl.Application.Initialize()
	vcl.Application.SetMainFormOnTaskBar(true)
	vcl.Application.CreateForm(&mainForm)
	vcl.Application.Run()
}

func (f *MainForm) OnFormCreate(sender vcl.IObject) {
	log.Print("Main form starting")
	f.startMirroring = false

	f.SetCaption("AppleTV Mirroring")
	f.EnabledMaximize(false)
	f.SetWidth(260)
	f.SetHeight(180)
	f.SetPosition(types.PoScreenCenter)
	f.SetBorderStyle(types.BsDialog)
	f.Icon().LoadFromFile("mirroring.ico")

	lab := vcl.NewLabel(f)
	lab.SetParent(f)
	lab.SetBounds(30, 10, 88, 20)
	ft := vcl.NewFont()
	ft.SetColor(colors.ClDarkblue)
	ft.SetSize(12)
	lab.SetFont(ft)
	lab.SetCaption("Choose the screen to mirror")
	f.LabChoose = lab

	v := screen.NewVideoProvider()
	f.videoProvider = v
	f.screens = v.GetScreens()
	cmb := vcl.NewComboBox(f)
	cmb.SetParent(f)
	cmb.SetBounds(20, 48, 220, 28)
	for i, s := range f.screens {
		item := fmt.Sprintf("Screen %d : %s", s.Index+1, s.Bounds.String())
		cmb.AddItem(item, nil)
		if i == 0 {
			cmb.SetSelText(item)
		}
	}
	f.CombList = cmb

	btn := vcl.NewButton(f)
	btn.SetParent(f)
	btn.SetBounds(80, 100, 110, 28)
	btn.SetCaption("Start Mirroring")
	btn.SetOnClick(f.OnButtonClick)
	f.BtnOK = btn

	f.discovery = airplay.NewDiscovery()
	//go f.searchAppleTV()

	c := airplay.NewAppleTVClient(f.appleTV)
	f.client = c

	f.ffmpeg = ffmpeg.FFmpeg()

}

func (f *MainForm) OnButtonClick(sender vcl.IObject) {
	if f.startMirroring {
		f.startMirroring = false
		f.ffmpeg.Stop()
		f.BtnOK.SetCaption("Start Mirroring")
	} else {
		//if f.appleTV == nil {
		//	go f.searchAppleTV()
		//	return
		//}

		f.startMirroring = true
		index := f.CombList.ItemIndex()
		sc := f.screens[index]
		f.ffmpeg.Start(sc, FPS)
		f.BtnOK.SetCaption("Stop Mirroring")
	}
	f.CombList.SetEnabled(!f.startMirroring)
}

func (f *MainForm) searchAppleTV() {
	defer vcl.ThreadSync(func() {
		f.LabChoose.SetCaption("Start Mirroring")
		f.BtnOK.SetEnabled(true)
	})
	go func() {
		info := ""
		for f.appleTV == nil {
			info += "."
			vcl.ThreadSync(func() {
				f.BtnOK.SetEnabled(false)
				f.LabChoose.SetCaption("Searching Apple TV" + info)

			})
			if info == "........" {
				info = ""
			}
			time.Sleep(800 * time.Millisecond)
		}
	}()
	f.appleTV = f.discovery.GetAirPlayService().AddrIPv4
	if f.appleTV == nil {
		vcl.ThreadSync(func() {
			vcl.ShowMessage("Cannot find AppleTV, try again please!")
		})
		f.appleTV = f.discovery.GetAirPlayService().AddrIPv4
		return
	}
	log.Printf("Found AppleTV at %v", f.appleTV)
}
