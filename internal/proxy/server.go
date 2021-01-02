/**
Created by: Joseph Han <joseph.bing.han@gmail.com>.
Created at: 2021-01-02 17:26
*/

package proxy

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
)

type Stream struct {
	data chan []byte
}

var _server *Stream

func init() {
	_server = &Stream{}
	_server.data = make(chan []byte)
	proxyAddr := ":8089"
	httpAddr := ":8090"
	proxy, errp := net.Listen("tcp", proxyAddr)
	if errp != nil {
		log.Fatalf("Listen fail on %v %v", proxyAddr, errp)
	}

	go func() {
		for {
			conn, errc := proxy.Accept()
			if errc != nil {
				log.Fatalf("Accept request fail %v", errc)
			}
			go func(con net.Conn) {
				defer con.Close()
				for {
					buffer := make([]byte, 1024)
					con.SetReadDeadline(time.Now().Add(30 * time.Second))
					_, errR := con.Read(buffer)
					if errR != nil {
						fmt.Errorf("Unable read input %v", errR)
						return
					}
					_server.data <- buffer
				}

			}(conn)
		}
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Handle Http Server Request")
		w.Header().Set("Content-Type", "video/x-flv")
		w.Header().Set("Transfer-Encoding", "chunked")
		w.Header().Set("Connection", "close")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(200)
		flusher := w.(http.Flusher)
		for {
			buf := <-_server.data
			w.Write(buf)
			flusher.Flush()
		}
	})
	go func() {
		errH := http.ListenAndServe(httpAddr, nil)
		if errH != nil {
			log.Fatalf("ListenAndServe: %v ", errH)
		}
	}()

}
