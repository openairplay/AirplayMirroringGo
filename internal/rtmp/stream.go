/**
Created by: Joseph Han <joseph.bing.han@gmail.com>.
Created at: 2021-01-02 16:09
*/
package rtmp

import (
	"github.com/gwuhaolin/livego/configure"
	"github.com/gwuhaolin/livego/protocol/hls"
	"github.com/gwuhaolin/livego/protocol/httpflv"
	"github.com/gwuhaolin/livego/protocol/rtmp"
	"log"
	"net"
)

const ROOM = "mirroring"

var _stream *rtmp.RtmpStream
var _key string

func init() {
	hlsAddr := ":8089"
	rtmpAddr := ":8090"
	//apiAddr := ":8088"
	httpflvAddr:=":8087"

	_stream = rtmp.NewRtmpStream()

	hlsListen, err := net.Listen("tcp", hlsAddr)
	if err != nil {
		log.Fatal(err)
	}

	hlsServer := hls.NewServer()
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Fatalf("HLS server panic: %v", r)
			}
		}()
		log.Printf("HLS listen On %v", hlsAddr)
		hlsServer.Serve(hlsListen)
	}()

	flvListen, err := net.Listen("tcp", httpflvAddr)
	if err != nil {
		log.Fatal(err)
	}

	hdlServer := httpflv.NewServer(_stream)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Fatalf("HTTP-FLV server panic: %v", r)
			}
		}()
		log.Printf("HTTP-FLV listen On %v", httpflvAddr)
		hdlServer.Serve(flvListen)
	}()

	//opListen, err := net.Listen("tcp", apiAddr)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//opServer := api.NewServer(_stream, rtmpAddr)
	//go func() {
	//	defer func() {
	//		if r := recover(); r != nil {
	//			log.Fatalf("HTTP-API server panic: %v ", r)
	//		}
	//	}()
	//	log.Printf("HTTP-API listen on %v", apiAddr)
	//	opServer.Serve(opListen)
	//}()

	rtmpListen, err := net.Listen("tcp", rtmpAddr)
	if err != nil {
		log.Fatal(err)
	}

	var rtmpServer *rtmp.Server

	rtmpServer = rtmp.NewRtmpServer(_stream, hlsServer)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Fatalf("RTMP server panic: %v", r)
			}
		}()
		log.Printf("RTMP Listen On %v", rtmpAddr)
		rtmpServer.Serve(rtmpListen)
	}()

	initKey()

}

func initKey() {
	msg, err := configure.RoomKeys.GetKey(ROOM)
	if err != nil {
		log.Fatalf("Create API key fail %v", err)
	}
	_key = msg
}

func GetKey() string {
	return _key
}
