/**
Created by: Joseph Han <joseph.bing.han@gmail.com>.
Created at: 2020-12-30 09:24
*/

package main

import (
	"AirplayMirroringGo/internal/ui"
	_ "AirplayMirroringGo/internal/ui"
)

func main() {
	//dc := airplay.NewDiscovery()
	//log.Printf("Search AppleTV-airplay: %v", dc.GetAirPlayService())
	//log.Printf("Search AppleTV-RAOP: %v", dc.GetRemoteAudioOutputProtocolService())
	ui.AppFormRun()
}
