package main

import (
	"ffmpeg_hls_go/internal/configs"
	"ffmpeg_hls_go/internal/handles"
	"ffmpeg_hls_go/internal/logger"
	"ffmpeg_hls_go/internal/token"
	"ffmpeg_hls_go/internal/video"
	"fmt"
	"net/http"
	"os"
)

var (
	ffmpegMgr *video.FFmpegMgr
	gitHash   string
	buildTime string
	goVersion string
)

func main() {
	config := configs.GetConfigInstance()
	log := logger.GetLoggerInstance()

	pid := os.Getpid()
	log.Info("")
	log.Infof("##############Starting server [%d]...#####################", pid)
	log.Infof("##############GitHash: %s, BuildTime: %s, GoVersion: %s....#####################", gitHash, buildTime, goVersion)

	token.Init()
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
