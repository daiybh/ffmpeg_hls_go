package main

import (
	"ffmpeg_hls_go/internal/configs"
	"ffmpeg_hls_go/internal/logger"
	"ffmpeg_hls_go/internal/video"
	"ffmpeg_hls_go/internal/video/handles"

	"fmt"
	"net/http"
)

var ffmpegMgr *video.FFmpegMgr

func main() {
	config := configs.GetConfigInstance()
	log := logger.GetLoggerInstance()
	log.Info("")
	log.Info("##############Starting server...#####################")

	// Start the FFmpeg Manager
	ffmpegMgr = video.GetFFmpegMgr()
	ffmpegMgr.Start(config)

	// Setup HTTP routes

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/ch/", handles.PlayHandler)
	http.HandleFunc("/status/", handles.StatusHandler)

	// Start the HTTP server
	address := fmt.Sprintf(":%d", config.Server.Port)
	log.Println("Server is listening on port " + address + "...")
	log.Fatal(http.ListenAndServe(address, nil))
}
