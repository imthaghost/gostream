package main

import (
	"fmt"
	"net"
	"path"
	"runtime"
	"time"

	"github.com/imthaghost/gostream/configure"
	"github.com/imthaghost/gostream/protocol/api"
	"github.com/imthaghost/gostream/protocol/hls"
	"github.com/imthaghost/gostream/protocol/httpflv"
	"github.com/imthaghost/gostream/protocol/rtmp"

	log "github.com/sirupsen/logrus"
)

var VERSION = "master"

func startHls() *hls.Server {
	hlsAddr := configure.Config.GetString("hls_addr")
	hlsListen, err := net.Listen("tcp", hlsAddr)
	if err != nil {
		log.Fatal(err)
	}

	hlsServer := hls.NewServer()
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Error("HLS server panic: ", r)
			}
		}()
		log.Info("HLS listen On ", hlsAddr)
		hlsServer.Serve(hlsListen)
	}()
	return hlsServer
}

var rtmpAddr string

func startRtmp(stream *rtmp.RtmpStream, hlsServer *hls.Server) {
	rtmpAddr = configure.Config.GetString("rtmp_addr")

	rtmpListen, err := net.Listen("tcp", rtmpAddr)
	if err != nil {
		log.Fatal(err)
	}

	var rtmpServer *rtmp.Server

	if hlsServer == nil {
		rtmpServer = rtmp.NewRtmpServer(stream, nil)
		log.Info("HLS server disable....")
	} else {
		rtmpServer = rtmp.NewRtmpServer(stream, hlsServer)
		log.Info("HLS server enable....")
	}

	defer func() {
		if r := recover(); r != nil {
			log.Error("RTMP server panic: ", r)
		}
	}()
	log.Info("RTMP Listen On ", rtmpAddr)
	rtmpServer.Serve(rtmpListen)
}

func startHTTPFlv(stream *rtmp.RtmpStream) {
	httpflvAddr := configure.Config.GetString("httpflv_addr")

	flvListen, err := net.Listen("tcp", httpflvAddr)
	if err != nil {
		log.Fatal(err)
	}

	hdlServer := httpflv.NewServer(stream)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Error("HTTP-FLV server panic: ", r)
			}
		}()
		log.Info("HTTP-FLV listen On ", httpflvAddr)
		hdlServer.Serve(flvListen)
	}()
}

func startAPI(stream *rtmp.RtmpStream) {
	apiAddr := configure.Config.GetString("api_addr")

	if apiAddr != "" {
		opListen, err := net.Listen("tcp", apiAddr)
		if err != nil {
			log.Fatal(err)
		}
		opServer := api.NewServer(stream, rtmpAddr)
		go func() {
			defer func() {
				if r := recover(); r != nil {
					log.Error("HTTP-API server panic: ", r)
				}
			}()
			log.Info("HTTP-API listen On ", apiAddr)
			opServer.Serve(opListen)
		}()
	}
}

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)
			return fmt.Sprintf("%s()", f.Function), fmt.Sprintf(" %s:%d", filename, f.Line)
		},
	})
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Error("gostream panic: ", r)
			time.Sleep(1 * time.Second)
		}
	}()

	log.Infof(`
	________  ________  ________  _________  ________  _______   ________  _____ ______      
	|\   ____\|\   __  \|\   ____\|\___   ___\\   __  \|\  ___ \ |\   __  \|\   _ \  _   \    
	\ \  \___|\ \  \|\  \ \  \___|\|___ \  \_\ \  \|\  \ \   __/|\ \  \|\  \ \  \\\__\ \  \   
	 \ \  \  __\ \  \\\  \ \_____  \   \ \  \ \ \   _  _\ \  \_|/_\ \   __  \ \  \\|__| \  \  
	  \ \  \|\  \ \  \\\  \|____|\  \   \ \  \ \ \  \\  \\ \  \_|\ \ \  \ \  \ \  \    \ \  \ 
	   \ \_______\ \_______\____\_\  \   \ \__\ \ \__\\ _\\ \_______\ \__\ \__\ \__\    \ \__\
		\|_______|\|_______|\_________\   \|__|  \|__|\|__|\|_______|\|__|\|__|\|__|     \|__|
						   \|_________|                                                       
												 
        version: %s
	`, VERSION)

	stream := rtmp.NewRtmpStream()
	hlsServer := startHls()
	startHTTPFlv(stream)
	startAPI(stream)

	startRtmp(stream, hlsServer)
}