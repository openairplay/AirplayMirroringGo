/**
Created by: Joseph Han <joseph.bing.han@gmail.com>.
Created at: 2020-12-30 10:01
*/

package airplay

import (
	"github.com/oleksandr/bonjour"
	"log"
	"time"
)

//goland:noinspection ALL
const (
	DNS_SD_TYPE_AP   = "_airplay._tcp"
	DNS_SD_TYPE_RAOP = "_raop._tcp"
)

type Discovery struct {
	resolver *bonjour.Resolver
	chanAP   chan *bonjour.ServiceEntry
	chanRAOP chan *bonjour.ServiceEntry
}

func NewDiscovery() *Discovery {
	dc := Discovery{}
	resolver, err := bonjour.NewResolver(nil)
	if err != nil {
		log.Fatalf("Failed to initialize resolver: %v", err)
	}
	dc.resolver = resolver
	dc.chanAP = make(chan *bonjour.ServiceEntry)
	dc.chanRAOP = make(chan *bonjour.ServiceEntry)
	return &dc
}

func (d *Discovery) GetAirPlayService() *bonjour.ServiceEntry {
	go func() {
		err := d.resolver.Browse(DNS_SD_TYPE_AP, "local.", d.chanAP)
		if err != nil {
			log.Fatalf("Fail to browse AppleTV: %v", err)
		}
	}()
	select {
	case ret := <-d.chanAP:
		return ret
	case <-time.After(time.Minute):
		return nil
	}
}

func (d *Discovery) GetRemoteAudioOutputProtocolService() *bonjour.ServiceEntry {
	go func() {
		err := d.resolver.Browse(DNS_SD_TYPE_RAOP, "local.", d.chanRAOP)
		if err != nil {
			log.Fatalf("Fail to browse AppleTV: %v", err)
		}
	}()
	select {
	case ret := <-d.chanRAOP:
		return ret
	case <-time.After(time.Minute):
		return nil
	}
}
